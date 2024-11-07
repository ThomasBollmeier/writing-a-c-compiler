package frontend

import "errors"

type loopLabeler struct {
	nameCreator NameCreator
	labelStack  []string
	err         error
}

func newLoopLabeler(nameCreator NameCreator) *loopLabeler {
	return &loopLabeler{nameCreator: nameCreator}
}

func (ll *loopLabeler) addLabels(p *Program) error {
	ll.labelStack = make([]string, 0)
	ll.err = nil

	p.Accept(ll)

	return ll.err
}

func (ll *loopLabeler) pushNewLabel() string {
	label := ll.nameCreator.LabelName("loop")
	ll.labelStack = append(ll.labelStack, label)
	return label
}

func (ll *loopLabeler) peekLabel() string {
	size := len(ll.labelStack)
	if size > 0 {
		return ll.labelStack[size-1]
	} else {
		return ""
	}
}

func (ll *loopLabeler) popLabel() string {
	size := len(ll.labelStack)
	if size > 0 {
		label := ll.labelStack[size-1]
		ll.labelStack = ll.labelStack[:size-1]
		return label
	} else {
		return ""
	}
}

func (ll *loopLabeler) VisitProgram(p *Program) {
	p.Func.Accept(ll)
}

func (ll *loopLabeler) VisitFunction(f *Function) {
	f.Body.Accept(ll)
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
	d.Label = ll.pushNewLabel()
	d.Body.Accept(ll)
	ll.popLabel()
}

func (ll *loopLabeler) VisitWhileStmt(w *WhileStmt) {
	w.Label = ll.pushNewLabel()
	w.Body.Accept(ll)
	ll.popLabel()
}

func (ll *loopLabeler) VisitForStmt(f *ForStmt) {
	f.Label = ll.pushNewLabel()
	f.Body.Accept(ll)
	ll.popLabel()
}

func (ll *loopLabeler) VisitBreakStmt(b *BreakStmt) {
	label := ll.peekLabel()
	if label == "" {
		ll.err = errors.New("break statement outside of loop/switch")
		return
	}
	b.Label = label
}

func (ll *loopLabeler) VisitContinueStmt(c *ContinueStmt) {
	label := ll.peekLabel()
	if label == "" {
		ll.err = errors.New("continue statement outside of loop/switch")
		return
	}
	c.Label = label
}

func (ll *loopLabeler) VisitSwitchStmt(*SwitchStmt) {}

func (ll *loopLabeler) VisitNullStmt() {}

func (ll *loopLabeler) VisitInteger(*IntegerLiteral) {}

func (ll *loopLabeler) VisitVariable(*Variable) {}

func (ll *loopLabeler) VisitUnary(*UnaryExpression) {}

func (ll *loopLabeler) VisitPostfixIncDec(*PostfixIncDec) {}

func (ll *loopLabeler) VisitBinary(*BinaryExpression) {}

func (ll *loopLabeler) VisitConditional(*Conditional) {}
