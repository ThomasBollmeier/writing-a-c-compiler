package frontend

import (
	"errors"
	"fmt"
)

type typeChecker struct {
	envs      *Environments
	errorList []error
}

func newTypeChecker(envs *Environments) *typeChecker {
	return &typeChecker{
		envs:      envs,
		errorList: []error{},
	}
}

func (tc *typeChecker) check(p *Program) []error {
	tc.errorList = make([]error, 0)
	p.Accept(tc)
	return tc.errorList
}

func (tc *typeChecker) VisitProgram(p *Program) {
	for _, f := range p.Declarations {
		f.Accept(tc)
	}
}

func (tc *typeChecker) VisitFunction(f *Function) {
	entry, definingEnv := tc.envs.Get(f.Name)

	if entry == nil || entry.category != idCatFunction && definingEnv != tc.envs.Block {
		tc.envs.set(f.Name, f.Name, linkExternal, idCatFunction,
			&FuncInfo{
				NumParams: len(f.Params),
				IsDefined: f.Body != nil,
			})
	} else {
		for {
			if entry.category != idCatFunction {
				tc.errorList = append(tc.errorList,
					errors.New(fmt.Sprintf("%s defined as a non-function", f.Name)))
				break
			}
			fnInfo := entry.typeInfo.(*FuncInfo)
			if fnInfo.NumParams != len(f.Params) {
				tc.errorList = append(tc.errorList,
					errors.New(fmt.Sprintf("%s is already declared with different signature", f.Name)))
				break
			}
			if fnInfo.IsDefined && f.Body != nil {
				tc.errorList = append(tc.errorList,
					errors.New(fmt.Sprintf("%s is already defined", f.Name)))
				break
			}
			if f.Body != nil {
				fnInfo.IsDefined = true
			}
			break
		}
	}

	if f.Body != nil {
		tc.envs.beginBlock()

		for _, param := range f.Params {
			tc.envs.set(param.Name, param.Name, linkNone, idCatParameter, &IntInfo{})
		}

		f.Body.Accept(tc)

		tc.envs.endBlock()
	}
}

func (tc *typeChecker) VisitVarDecl(v *VarDecl) {
	tc.envs.set(v.Name, v.Name, linkNone, idCatVariable, &IntInfo{})

	if v.InitValue != nil {
		v.InitValue.Accept(tc)
	}
}

func (tc *typeChecker) VisitReturn(r *ReturnStmt) {
	if r.Expression != nil {
		r.Expression.Accept(tc)
	}
}

func (tc *typeChecker) VisitExprStmt(e *ExpressionStmt) {
	e.Expression.Accept(tc)
}

func (tc *typeChecker) VisitIfStmt(i *IfStmt) {
	i.Condition.Accept(tc)
	i.Consequent.Accept(tc)
	if i.Alternate != nil {
		i.Alternate.Accept(tc)
	}
}

func (tc *typeChecker) VisitBlockStmt(b *BlockStmt) {
	tc.envs.beginBlock()
	for _, item := range b.Items {
		item.Accept(tc)
	}
	tc.envs.endBlock()
}

func (tc *typeChecker) VisitGotoStmt(*GotoStmt) {}

func (tc *typeChecker) VisitLabelStmt(*LabelStmt) {}

func (tc *typeChecker) VisitDoWhileStmt(d *DoWhileStmt) {
	d.Condition.Accept(tc)
	d.Body.Accept(tc)
}

func (tc *typeChecker) VisitWhileStmt(w *WhileStmt) {
	w.Condition.Accept(tc)
	w.Body.Accept(tc)
}

func (tc *typeChecker) VisitForStmt(f *ForStmt) {
	f.InitStmt.Accept(tc)
	if f.Condition != nil {
		f.Condition.Accept(tc)
	}
	if f.Post != nil {
		f.Post.Accept(tc)
	}
	f.Body.Accept(tc)
}

func (tc *typeChecker) VisitBreakStmt(*BreakStmt) {}

func (tc *typeChecker) VisitContinueStmt(*ContinueStmt) {}

func (tc *typeChecker) VisitSwitchStmt(s *SwitchStmt) {
	s.Expr.Accept(tc)
	s.Body.Accept(tc)
}

func (tc *typeChecker) VisitCaseStmt(c *CaseStmt) {
	if c.Value != nil {
		c.Value.Accept(tc)
	}
}

func (tc *typeChecker) VisitNullStmt() {}

func (tc *typeChecker) VisitInteger(*IntegerLiteral) {}

func (tc *typeChecker) VisitVariable(v *Variable) {
	entry, _ := tc.envs.Get(v.Name)
	if entry != nil && entry.category != idCatVariable && entry.category != idCatParameter {
		tc.errorList = append(tc.errorList, errors.New(fmt.Sprintf("%s defined as a non-variable", v.Name)))
	}
}

func (tc *typeChecker) VisitFunctionCall(f *FunctionCall) {
	entry, _ := tc.envs.Get(f.Callee)
	if entry != nil {
		if entry.category != idCatFunction {
			tc.errorList = append(tc.errorList, errors.New(fmt.Sprintf("%s is not a function", f.Callee)))
		} else {
			fnInfo := entry.typeInfo.(*FuncInfo)
			if len(f.Args) != fnInfo.NumParams {
				tc.errorList = append(tc.errorList,
					errors.New(fmt.Sprintf("%s: #arguments <> #params (%d <> %d)",
						f.Callee, len(f.Args), fnInfo.NumParams)))
			}
		}
	}
	for _, arg := range f.Args {
		arg.Accept(tc)
	}
}

func (tc *typeChecker) VisitUnary(u *UnaryExpression) {
	u.Right.Accept(tc)
}

func (tc *typeChecker) VisitPostfixIncDec(p *PostfixIncDec) {
	p.Operand.Accept(tc)
}

func (tc *typeChecker) VisitBinary(b *BinaryExpression) {
	b.Left.Accept(tc)
	b.Right.Accept(tc)
}

func (tc *typeChecker) VisitConditional(c *Conditional) {
	c.Condition.Accept(tc)
	c.Consequent.Accept(tc)
	c.Alternate.Accept(tc)
}
