package tacky

import "fmt"

type AstPrinter struct {
	offset          int
	delta           int
	suppressPadding bool
}

func NewAstPrinter(delta int) *AstPrinter {
	return &AstPrinter{0, delta, false}
}

func (ap *AstPrinter) visitProgram(p *Program) {
	ap.println("Program(")
	ap.indent()
	p.Fun.Accept(ap)
	ap.dedent()
	ap.println(")")
}

func (ap *AstPrinter) visitFunction(f *Function) {
	ap.println("Function(")
	ap.indent()
	ap.println("name=" + f.Ident)
	ap.println("body=[")
	ap.indent()
	for _, inst := range f.Body {
		inst.Accept(ap)
	}
	ap.dedent()
	ap.println("]")
	ap.dedent()
	ap.println(")")
}

func (ap *AstPrinter) visitReturn(r *Return) {
	ap.println("Return(")
	ap.indent()
	r.Val.Accept(ap)
	ap.println("")
	ap.dedent()
	ap.println(")")
}

func (ap *AstPrinter) visitUnary(u *Unary) {
	ap.println("Unary(")
	ap.indent()
	ap.print("operator=")
	ap.suppressPadding = true
	u.Op.Accept(ap)
	ap.println("")
	ap.print("src=")
	ap.suppressPadding = true
	u.Src.Accept(ap)
	ap.println("")
	ap.print("dst=")
	ap.suppressPadding = true
	u.Dst.Accept(ap)
	ap.println("")
	ap.dedent()
	ap.println(")")
}

func (ap *AstPrinter) visitBinary(b *Binary) {
	ap.println("Binary(")
	ap.indent()
	ap.print("operator=")
	ap.suppressPadding = true
	b.Op.Accept(ap)
	ap.println("")
	ap.print("src1=")
	ap.suppressPadding = true
	b.Src1.Accept(ap)
	ap.println("")
	ap.print("src2=")
	ap.suppressPadding = true
	b.Src2.Accept(ap)
	ap.println("")
	ap.print("dst=")
	ap.suppressPadding = true
	b.Dst.Accept(ap)
	ap.println("")
	ap.dedent()
	ap.println(")")
}

func (ap *AstPrinter) visitCopy(c *Copy) {
	ap.println("Copy (")
	ap.indent()
	ap.print("src=")
	ap.suppressPadding = true
	c.Src.Accept(ap)
	ap.println("")
	ap.print("dst=")
	ap.suppressPadding = true
	c.Dst.Accept(ap)
	ap.println("")
	ap.dedent()
	ap.println(")")
}

func (ap *AstPrinter) visitJump(j *Jump) {
	ap.println("Jump(")
	ap.indent()
	ap.println("target=" + j.Target)
	ap.dedent()
	ap.println(")")
}

func (ap *AstPrinter) visitJumpIfZero(j *JumpIfZero) {
	ap.println("JumpIfZero(")
	ap.indent()
	ap.print("condition=")
	ap.suppressPadding = true
	j.Condition.Accept(ap)
	ap.println("")
	ap.println("target=" + j.Target)
	ap.dedent()
	ap.println(")")
}

func (ap *AstPrinter) visitJumpIfNotZero(j *JumpIfNotZero) {
	ap.println("JumpIfNotZero(")
	ap.indent()
	ap.print("condition=")
	ap.suppressPadding = true
	j.Condition.Accept(ap)
	ap.println("")
	ap.println("target=" + j.Target)
	ap.dedent()
	ap.println(")")
}

func (ap *AstPrinter) visitLabel(l *Label) {
	ap.println("Label(name=" + l.Name + ")")
}

func (ap *AstPrinter) visitIntConstant(i *IntConstant) {
	ap.print(fmt.Sprintf("IntConstant(%d)", i.Val))
}

func (ap *AstPrinter) visitVar(v *Var) {
	ap.print("Var(" + v.Ident + ")")
}

func (ap *AstPrinter) visitComplement() {
	ap.print("~")
}

func (ap *AstPrinter) visitNegate() {
	ap.print("-")
}

func (ap *AstPrinter) visitNot() {
	ap.print("!")
}

func (ap *AstPrinter) visitAdd() {
	ap.print("+")
}

func (ap *AstPrinter) visitSub() {
	ap.print("-")
}

func (ap *AstPrinter) visitMul() {
	ap.print("*")
}

func (ap *AstPrinter) visitDiv() {
	ap.print("/")
}

func (ap *AstPrinter) visitRemainder() {
	ap.print("%")
}

func (ap *AstPrinter) visitBitAnd() {
	ap.print("&")
}

func (ap *AstPrinter) visitBitOr() {
	ap.print("|")
}

func (ap *AstPrinter) visitBitXor() {
	ap.print("^")
}

func (ap *AstPrinter) visitBitShiftLeft() {
	ap.print("<<")
}

func (ap *AstPrinter) visitBitShiftRight() {
	ap.print(">>")
}

func (ap *AstPrinter) visitAnd() {
	ap.print("&&")
}

func (ap *AstPrinter) visitOr() {
	ap.print("||")
}

func (ap *AstPrinter) visitEqual() {
	ap.print("==")
}

func (ap *AstPrinter) visitNotEqual() {
	ap.print("!=")
}

func (ap *AstPrinter) visitGreater() {
	ap.print(">")
}

func (ap *AstPrinter) visitGreaterEq() {
	ap.print(">=")
}

func (ap *AstPrinter) visitLess() {
	ap.print("<")
}

func (ap *AstPrinter) visitLessEq() {
	ap.print("<=")
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
