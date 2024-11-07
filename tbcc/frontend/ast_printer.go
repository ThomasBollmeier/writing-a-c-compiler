package frontend

import "fmt"

type AstPrinter struct {
	offset          int
	delta           int
	suppressPadding bool
}

func NewAstPrinter(delta int) *AstPrinter {
	return &AstPrinter{0, delta, false}
}

func (ap *AstPrinter) VisitProgram(p *Program) {
	ap.println("Program(")
	ap.indent()
	p.Func.Accept(ap)
	ap.dedent()
	ap.println(")")
}

func (ap *AstPrinter) VisitFunction(f *Function) {
	ap.println("Function(")
	ap.indent()
	ap.println("name=\"" + f.Name + "\"")
	ap.print("body=")
	ap.suppressPadding = true
	f.Body.Accept(ap)
	ap.dedent()
	ap.println(")")
}

func (ap *AstPrinter) VisitVarDecl(v *VarDecl) {
	ap.println("VarDeclaration(")
	ap.indent()
	ap.println("name=\"" + v.Name + "\"")
	if v.InitValue != nil {
		ap.print("initValue=")
		ap.suppressPadding = true
		v.InitValue.Accept(ap)
	}
	ap.dedent()
	ap.println(")")
}

func (ap *AstPrinter) VisitReturn(r *ReturnStmt) {
	ap.println("Return(")
	ap.indent()
	r.Expression.Accept(ap)
	ap.dedent()
	ap.println(")")
}

func (ap *AstPrinter) VisitExprStmt(e *ExpressionStmt) {
	ap.println("ExpressionStatement(")
	ap.indent()
	ap.print("expression=")
	ap.suppressPadding = true
	e.Expression.Accept(ap)
	ap.dedent()
	ap.println(")")
}

func (ap *AstPrinter) VisitIfStmt(i *IfStmt) {
	ap.println("IfStatement(")
	ap.indent()
	ap.print("condition=")
	ap.suppressPadding = true
	i.Condition.Accept(ap)
	ap.print("then=")
	ap.suppressPadding = true
	i.Consequent.Accept(ap)
	if i.Alternate != nil {
		ap.print("else=")
		ap.suppressPadding = true
		i.Alternate.Accept(ap)
	}
	ap.dedent()
	ap.println(")")
}

func (ap *AstPrinter) VisitBlockStmt(b *BlockStmt) {
	ap.println("BlockStmt(")
	ap.indent()
	for _, item := range b.Items {
		item.Accept(ap)
	}
	ap.dedent()
	ap.println(")")
}

func (ap *AstPrinter) VisitGotoStmt(g *GotoStmt) {
	ap.println("GotoStatement(")
	ap.indent()
	ap.println("target=" + g.Target)
	ap.dedent()
	ap.println(")")
}

func (ap *AstPrinter) VisitLabelStmt(l *LabelStmt) {
	ap.println("LabelStatement(")
	ap.indent()
	ap.println("name=" + l.Name)
	ap.dedent()
	ap.println(")")
}

func (ap *AstPrinter) VisitDoWhileStmt(d *DoWhileStmt) {
	ap.println("DoWhileStatement(")
	ap.indent()
	if d.Label != "" {
		ap.println(fmt.Sprintf("label=%s", d.Label))
	}
	ap.print("body=")
	ap.suppressPadding = true
	d.Body.Accept(ap)
	ap.print("condition=")
	ap.suppressPadding = true
	d.Condition.Accept(ap)
	ap.dedent()
	ap.println(")")
}

func (ap *AstPrinter) VisitWhileStmt(w *WhileStmt) {
	ap.println("WhileStatement(")
	ap.indent()
	if w.Label != "" {
		ap.println(fmt.Sprintf("label=%s", w.Label))
	}
	ap.print("condition=")
	ap.suppressPadding = true
	w.Condition.Accept(ap)
	ap.print("body=")
	ap.suppressPadding = true
	w.Body.Accept(ap)
	ap.dedent()
	ap.println(")")
}

func (ap *AstPrinter) VisitForStmt(f *ForStmt) {
	ap.println("ForStatement(")
	ap.indent()
	if f.Label != "" {
		ap.println(fmt.Sprintf("label=%s", f.Label))
	}
	ap.print("init=")
	ap.suppressPadding = true
	f.InitStmt.Accept(ap)
	if f.Condition != nil {
		ap.print("condition=")
		ap.suppressPadding = true
		f.Condition.Accept(ap)
	}
	if f.Post != nil {
		ap.print("post=")
		ap.suppressPadding = true
		f.Post.Accept(ap)
	}
	ap.print("body=")
	ap.suppressPadding = true
	f.Body.Accept(ap)
	ap.dedent()
	ap.println(")")
}

func (ap *AstPrinter) VisitBreakStmt(b *BreakStmt) {
	ap.println(fmt.Sprintf("BreakStmt(%s)", b.Label))
}

func (ap *AstPrinter) VisitContinueStmt(c *ContinueStmt) {
	ap.println(fmt.Sprintf("ContinueStmt(%s)", c.Label))
}

func (ap *AstPrinter) VisitSwitchStmt(s *SwitchStmt) {
	ap.println("SwitchStmt(")
	ap.indent()
	ap.print("expression=")
	ap.suppressPadding = true
	s.Expr.Accept(ap)
	ap.println("caseBlocks=[")
	ap.indent()
	for _, block := range s.CaseBlocks {
		ap.println("Case(")
		ap.indent()
		if block.Value != nil {
			ap.print("select=")
			ap.suppressPadding = true
			block.Value.Accept(ap)
		} else {
			ap.println("select=default")
		}
		ap.println("statements=[")
		ap.indent()
		for _, stmt := range block.Body {
			stmt.Accept(ap)
		}
		ap.dedent()
		ap.println("]")
		ap.dedent()
		ap.println(")")
	}
	ap.dedent()
	ap.println("]")
	ap.dedent()
	ap.println(")")
}

func (ap *AstPrinter) VisitNullStmt() {
	ap.println("NullStatement()")
}

func (ap *AstPrinter) VisitInteger(i *IntegerLiteral) {
	text := fmt.Sprintf("Constant(%d)", i.Value)
	ap.println(text)
}

func (ap *AstPrinter) VisitVariable(v *Variable) {
	text := fmt.Sprintf("Variable(%s)", v.Name)
	ap.println(text)
}

func (ap *AstPrinter) VisitUnary(unary *UnaryExpression) {
	ap.println("Unary(")
	ap.indent()
	ap.print("operator=\"" + unary.Operator + "\"\n")
	ap.print("right=")
	ap.suppressPadding = true
	unary.Right.Accept(ap)
	ap.dedent()
	ap.println(")")
}

func (ap *AstPrinter) VisitPostfixIncDec(p *PostfixIncDec) {
	ap.println("PostfixIncDec(")
	ap.indent()
	ap.print("operator=\"" + p.Operator + "\"\n")
	ap.print("operand=")
	ap.suppressPadding = true
	p.Operand.Accept(ap)
	ap.dedent()
	ap.println(")")
}

func (ap *AstPrinter) VisitBinary(binary *BinaryExpression) {
	ap.println("Binary(")
	ap.indent()
	ap.print("operator=\"" + binary.Operator + "\"\n")
	ap.print("left=")
	ap.suppressPadding = true
	binary.Left.Accept(ap)
	ap.print("right=")
	ap.suppressPadding = true
	binary.Right.Accept(ap)
	ap.dedent()
	ap.println(")")
}

func (ap *AstPrinter) VisitConditional(cond *Conditional) {
	ap.println("Conditional(")
	ap.indent()
	ap.print("condition=")
	ap.suppressPadding = true
	cond.Condition.Accept(ap)
	ap.print("then=")
	ap.suppressPadding = true
	cond.Consequent.Accept(ap)
	ap.print("else=")
	ap.suppressPadding = true
	cond.Alternate.Accept(ap)
	ap.dedent()
	ap.println(")")
}

func (ap *AstPrinter) indent() {
	ap.offset += ap.delta
}

func (ap *AstPrinter) dedent() {
	ap.offset -= ap.delta
}

func (ap *AstPrinter) print(text string) {
	if ap.suppressPadding {
		fmt.Print(text)
		ap.suppressPadding = false
		return
	}
	leftPadding := ""
	for i := 0; i < ap.offset; i++ {
		leftPadding += " "
	}
	fmt.Printf("%s%s", leftPadding, text)
}

func (ap *AstPrinter) println(text string) {
	ap.print(text + "\n")
}
