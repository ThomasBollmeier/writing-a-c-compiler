package frontend

import (
	"errors"
	"fmt"
)

type labelChecker struct {
	gotoStmts  map[string]bool
	labelStmts map[string]error
	caseErrors map[string]error
}

func newLabelChecker() *labelChecker {
	return &labelChecker{}
}

func (lc *labelChecker) check(program *Program) error {

	lc.gotoStmts = map[string]bool{}
	lc.labelStmts = map[string]error{}
	lc.caseErrors = map[string]error{}

	program.Accept(lc)

	for target := range lc.gotoStmts {
		_, ok := lc.labelStmts[target]
		if !ok {
			return errors.New("target " + target + " does not exist")
		}
	}
	for _, err := range lc.labelStmts {
		if err != nil {
			return err
		}
	}
	for _, err := range lc.caseErrors {
		if err != nil {
			return err
		}
	}
	return nil
}

func (lc *labelChecker) VisitProgram(p *Program) {
	p.Func.Accept(lc)
}

func (lc *labelChecker) VisitFunction(f *Function) {
	f.Body.Accept(lc)
}

func (lc *labelChecker) VisitVarDecl(*VarDecl) {}

func (lc *labelChecker) VisitReturn(*ReturnStmt) {}

func (lc *labelChecker) VisitExprStmt(*ExpressionStmt) {}

func (lc *labelChecker) VisitIfStmt(i *IfStmt) {
	i.Consequent.Accept(lc)
	if i.Alternate != nil {
		i.Alternate.Accept(lc)
	}
}

func (lc *labelChecker) VisitBlockStmt(b *BlockStmt) {
	var labelName string
	var caseName string
	var caseLabel string

	for _, item := range b.Items {
		item.Accept(lc)
		if item.GetType() == AstLabelStmt {
			labelName = item.(*LabelStmt).Name
		} else if item.GetType() == AstCaseStmt {
			caseStmt := item.(*CaseStmt)
			if caseStmt.Value != nil {
				caseName = fmt.Sprintf("case %d", caseStmt.Value)
			} else {
				caseName = "default"
			}
			caseLabel = caseStmt.Label
		} else if labelName != "" {
			if item.GetType() == AstVarDecl {
				lc.labelStmts[labelName] = errors.New("label " + labelName +
					" is not allowed before a variable declaration")
			}
			labelName = ""
		} else if caseLabel != "" {
			if item.GetType() == AstVarDecl {
				lc.caseErrors[caseLabel] = errors.New(caseName +
					" is not allowed before a variable declaration")
			}
			caseName = ""
			caseLabel = ""
		}
	}

	if labelName != "" {
		lc.labelStmts[labelName] = errors.New("label " + labelName + " is not before any statement")
	}
	if caseName != "" {
		lc.caseErrors[caseLabel] = errors.New(caseName + " is not before any statement")
	}
}

func (lc *labelChecker) VisitGotoStmt(g *GotoStmt) {
	_, ok := lc.gotoStmts[g.Target]
	if !ok {
		lc.gotoStmts[g.Target] = true
	}
}

func (lc *labelChecker) VisitLabelStmt(l *LabelStmt) {
	_, ok := lc.labelStmts[l.Name]
	if !ok {
		lc.labelStmts[l.Name] = nil
	} else {
		lc.labelStmts[l.Name] = errors.New("label " + l.Name + " already exists")
	}
}

func (lc *labelChecker) VisitDoWhileStmt(d *DoWhileStmt) {
	d.Body.Accept(lc)
}

func (lc *labelChecker) VisitWhileStmt(w *WhileStmt) {
	w.Body.Accept(lc)
}

func (lc *labelChecker) VisitForStmt(f *ForStmt) {
	f.Body.Accept(lc)
}

func (lc *labelChecker) VisitBreakStmt(*BreakStmt) {}

func (lc *labelChecker) VisitContinueStmt(*ContinueStmt) {}

func (lc *labelChecker) VisitSwitchStmt(s *SwitchStmt) {
	s.Body.Accept(lc)
}

func (lc *labelChecker) VisitCaseStmt(*CaseStmt) {}

func (lc *labelChecker) VisitNullStmt() {}

func (lc *labelChecker) VisitInteger(*IntegerLiteral) {}

func (lc *labelChecker) VisitVariable(*Variable) {}

func (lc *labelChecker) VisitUnary(*UnaryExpression) {}

func (lc *labelChecker) VisitPostfixIncDec(*PostfixIncDec) {}

func (lc *labelChecker) VisitBinary(*BinaryExpression) {}

func (lc *labelChecker) VisitConditional(*Conditional) {}
