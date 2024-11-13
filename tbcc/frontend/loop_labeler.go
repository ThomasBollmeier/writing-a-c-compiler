package frontend

import (
	"errors"
	"fmt"
)

type labelContext uint

const (
	labelCtxLoop labelContext = iota
	labelCtxSwitch
)

type labelInfo struct {
	name        string
	ctx         labelContext
	switchInfo_ *switchInfo
}

type switchInfo struct {
	nextCaseIdx uint
	prevCase    *CaseStmt
	cases       map[string]bool
	caseLabels  []string
}

type loopLabeler struct {
	nameCreator NameCreator
	labelStack  []labelInfo
	err         error
}

func newLoopLabeler(nameCreator NameCreator) *loopLabeler {
	return &loopLabeler{nameCreator: nameCreator}
}

func (ll *loopLabeler) addLabels(p *Program) error {
	ll.labelStack = make([]labelInfo, 0)
	ll.err = nil
	p.Accept(ll)
	return ll.err
}

func (ll *loopLabeler) pushNewLabel(ctx labelContext) labelInfo {
	var prefix string
	var switchInfo_ *switchInfo = nil

	if ctx == labelCtxLoop {
		prefix = "loop"
	} else {
		prefix = "switch"
		switchInfo_ = &switchInfo{
			nextCaseIdx: 0,
			prevCase:    nil,
			cases:       make(map[string]bool),
			caseLabels:  make([]string, 0),
		}
	}
	label := ll.nameCreator.LabelName(prefix)
	ret := labelInfo{
		name:        label,
		ctx:         ctx,
		switchInfo_: switchInfo_,
	}
	ll.labelStack = append(ll.labelStack, ret)
	return ret
}

func (ll *loopLabeler) peekLabel() *labelInfo {
	size := len(ll.labelStack)
	if size > 0 {
		return &ll.labelStack[size-1]
	} else {
		return nil
	}
}

func (ll *loopLabeler) getLoopIdx() int {
	size := len(ll.labelStack)
	for i := size - 1; i >= 0; i-- {
		if ll.labelStack[i].ctx == labelCtxLoop {
			return i
		}
	}
	return -1
}

func (ll *loopLabeler) getSwitchIdx() int {
	size := len(ll.labelStack)
	for i := size - 1; i >= 0; i-- {
		if ll.labelStack[i].ctx == labelCtxSwitch {
			return i
		}
	}
	return -1
}

func (ll *loopLabeler) popLabel() *labelInfo {
	size := len(ll.labelStack)
	if size > 0 {
		ret := &ll.labelStack[size-1]
		ll.labelStack = ll.labelStack[:size-1]
		return ret
	} else {
		return nil
	}
}

func (ll *loopLabeler) VisitProgram(p *Program) {
	for _, fun := range p.Functions {
		fun.Accept(ll)
	}
}

func (ll *loopLabeler) VisitFunction(f *Function) {
	if f.Body != nil {
		f.Body.Accept(ll)
	}
}

func (ll *loopLabeler) VisitVarDecl(*VarDecl) {}

func (ll *loopLabeler) VisitReturn(*ReturnStmt) {}

func (ll *loopLabeler) VisitExprStmt(*ExpressionStmt) {}

func (ll *loopLabeler) VisitIfStmt(i *IfStmt) {
	i.Consequent.Accept(ll)
	if ll.err != nil {
		return
	}
	if i.Alternate != nil {
		i.Alternate.Accept(ll)
		if ll.err != nil {
			return
		}
	}
}

func (ll *loopLabeler) VisitBlockStmt(b *BlockStmt) {
	for _, item := range b.Items {
		item.Accept(ll)
		if ll.err != nil {
			return
		}
	}
}

func (ll *loopLabeler) VisitGotoStmt(*GotoStmt) {}

func (ll *loopLabeler) VisitLabelStmt(*LabelStmt) {}

func (ll *loopLabeler) VisitDoWhileStmt(d *DoWhileStmt) {
	d.Label = ll.pushNewLabel(labelCtxLoop).name
	d.Body.Accept(ll)
	ll.popLabel()
}

func (ll *loopLabeler) VisitWhileStmt(w *WhileStmt) {
	w.Label = ll.pushNewLabel(labelCtxLoop).name
	w.Body.Accept(ll)
	ll.popLabel()
}

func (ll *loopLabeler) VisitForStmt(f *ForStmt) {
	f.Label = ll.pushNewLabel(labelCtxLoop).name
	f.Body.Accept(ll)
	ll.popLabel()
}

func (ll *loopLabeler) VisitBreakStmt(b *BreakStmt) {
	lInfo := ll.peekLabel()
	if lInfo == nil {
		ll.err = errors.New("break statement outside of loop/switch")
		return
	}
	b.Label = lInfo.name
}

func (ll *loopLabeler) VisitContinueStmt(c *ContinueStmt) {
	loopIdx := ll.getLoopIdx()
	if loopIdx == -1 {
		ll.err = errors.New("continue statement outside of loop")
		return
	}
	c.Label = ll.labelStack[loopIdx].name
}

func (ll *loopLabeler) VisitSwitchStmt(s *SwitchStmt) {
	s.Label = ll.pushNewLabel(labelCtxSwitch).name
	s.Body.Accept(ll)
	lInfo := ll.popLabel()
	if lInfo != nil {
		caseLabels := lInfo.switchInfo_.caseLabels
		if len(caseLabels) > 0 {
			s.FirstCaseLabel = caseLabels[0]
		}
	}
}

func (ll *loopLabeler) VisitCaseStmt(c *CaseStmt) {
	switchIdx := ll.getSwitchIdx()
	if switchIdx == -1 {
		ll.err = errors.New("case/default statement outside of switch")
		return
	}
	switchData := ll.labelStack[switchIdx]

	var caseValueStr string
	if c.Value != nil {
		caseValueStr = fmt.Sprintf("%d", c.Value.(*IntegerLiteral).Value)
	} else {
		caseValueStr = "default"
	}
	_, ok := switchData.switchInfo_.cases[caseValueStr]
	if ok {
		ll.err = errors.New("there is already a case clause for value " + caseValueStr)
		return
	}

	label := fmt.Sprintf("%s.case.%d", switchData.name, switchData.switchInfo_.nextCaseIdx)
	switchData.switchInfo_.cases[caseValueStr] = true
	switchData.switchInfo_.caseLabels = append(switchData.switchInfo_.caseLabels, label)

	if switchData.switchInfo_.prevCase != nil {
		switchData.switchInfo_.prevCase.NextCaseLabel = label
		c.PrevCaseLabel = switchData.switchInfo_.prevCase.Label
	}
	c.Label = label
	c.NextCaseLabel = fmt.Sprintf("%s.break", switchData.name)

	switchData.switchInfo_.prevCase = c
	switchData.switchInfo_.nextCaseIdx++
	ll.labelStack[switchIdx] = switchData
}

func (ll *loopLabeler) VisitNullStmt() {}

func (ll *loopLabeler) VisitInteger(*IntegerLiteral) {}

func (ll *loopLabeler) VisitVariable(*Variable) {}

func (ll *loopLabeler) VisitFunctionCall(*FunctionCall) {}

func (ll *loopLabeler) VisitUnary(*UnaryExpression) {}

func (ll *loopLabeler) VisitPostfixIncDec(*PostfixIncDec) {}

func (ll *loopLabeler) VisitBinary(*BinaryExpression) {}

func (ll *loopLabeler) VisitConditional(*Conditional) {}
