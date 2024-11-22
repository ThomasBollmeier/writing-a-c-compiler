package backend

import (
	"fmt"
)

type AsmPrinter struct {
	offset          int
	delta           int
	suppressPadding bool
}

func NewAsmPrinter(delta int) *AsmPrinter {
	return &AsmPrinter{0, delta, false}
}

func (ap *AsmPrinter) VisitProgram(p *Program) {
	ap.println("Program(")
	ap.indent()
	for _, funcDef := range p.FuncDefs {
		funcDef.Accept(ap)
	}
	ap.dedent()
	ap.println(")")
}

func (ap *AsmPrinter) VisitFunctionDef(f *FunctionDef) {
	ap.println("FunctionDef(")
	ap.indent()
	ap.println("name=\"" + f.Name + "\"")
	ap.println("instructions=[")
	ap.indent()
	for _, inst := range f.Instructions {
		inst.Accept(ap)
	}
	ap.dedent()
	ap.println("]")
	ap.dedent()
	ap.println(")")
}

func (ap *AsmPrinter) VisitMov(m *Mov) {
	ap.println("Mov(")
	ap.indent()
	ap.print("src=")
	ap.suppressPadding = true
	m.Src.Accept(ap)
	ap.print("dst=")
	ap.suppressPadding = true
	m.Dst.Accept(ap)
	ap.dedent()
	ap.println(")")
}

func (ap *AsmPrinter) VisitUnary(u *Unary) {
	ap.println("Unary(")
	ap.indent()
	ap.print("op=")
	ap.suppressPadding = true
	u.Op.Accept(ap)
	ap.print("operand=")
	ap.suppressPadding = true
	u.Operand.Accept(ap)
	ap.dedent()
	ap.println(")")
}

func (ap *AsmPrinter) VisitBinary(b *Binary) {
	ap.println("Binary(")
	ap.indent()
	ap.print("op=")
	ap.suppressPadding = true
	b.Op.Accept(ap)
	ap.print("operand1=")
	ap.suppressPadding = true
	b.Operand1.Accept(ap)
	ap.print("operand2=")
	ap.suppressPadding = true
	b.Operand2.Accept(ap)
	ap.dedent()
	ap.println(")")
}

func (ap *AsmPrinter) VisitCmp(c *Cmp) {
	ap.println("Cmp(")
	ap.indent()
	ap.print("left=")
	ap.suppressPadding = true
	c.Left.Accept(ap)
	ap.print("right=")
	ap.suppressPadding = true
	c.Right.Accept(ap)
	ap.dedent()
	ap.println(")")
}

func (ap *AsmPrinter) VisitIDiv(i *IDiv) {
	ap.println("IDiv(")
	ap.indent()
	ap.print("operand=")
	ap.suppressPadding = true
	i.Operand.Accept(ap)
	ap.dedent()
	ap.println(")")
}

func (ap *AsmPrinter) VisitCdq(*Cdq) {
	ap.println("Cdq")
}

func (ap *AsmPrinter) VisitJump(j *Jump) {
	ap.println(fmt.Sprintf("Jump(target=\"%s\")", j.Identifier))
}

func (ap *AsmPrinter) VisitJumpCC(j *JumpCC) {
	ap.println("JumpCC(")
	ap.indent()
	ap.println(fmt.Sprintf("conditionCode=%d", j.CondCode))
	ap.println(fmt.Sprintf("target=\"%s\")", j.Identifier))
	ap.dedent()
	ap.println(")")
}

func (ap *AsmPrinter) VisitSetCC(s *SetCC) {
	ap.println("SetCC(")
	ap.indent()
	ap.println(fmt.Sprintf("conditionCode=%d", s.CondCode))
	ap.print("operand=")
	ap.suppressPadding = true
	s.Op.Accept(ap)
	ap.dedent()
	ap.println(")")
}

func (ap *AsmPrinter) VisitLabel(l *Label) {
	ap.println(fmt.Sprintf("Label(name=\"%s\")", l.Identifier))
}

func (ap *AsmPrinter) VisitAllocStack(a *AllocStack) {
	ap.println(fmt.Sprintf("AllocStack(%d)", a.N))
}

func (ap *AsmPrinter) VisitDeAllocStack(d *DeAllocStack) {
	ap.println(fmt.Sprintf("DeAllocStack(%d)", d.N))
}

func (ap *AsmPrinter) VisitPush(p *Push) {
	ap.println("Push(")
	ap.indent()
	p.Op.Accept(ap)
	ap.dedent()
	ap.println(")")
}

func (ap *AsmPrinter) VisitCall(c *Call) {
	ap.println(fmt.Sprintf("Call(%s)", c.Identifier))
}

func (ap *AsmPrinter) VisitReturn() {
	ap.println("Ret")
}

func (ap *AsmPrinter) VisitNeg(*Neg) {
	ap.println("Neg")
}

func (ap *AsmPrinter) VisitNot(*Not) {
	ap.println("Not")
}

func (ap *AsmPrinter) VisitAdd(*Add) {
	ap.println("Add")
}

func (ap *AsmPrinter) VisitSub(*Sub) {
	ap.println("Sub")
}

func (ap *AsmPrinter) VisitMul(*Mul) {
	ap.println("Mul")
}

func (ap *AsmPrinter) VisitBitOp(op BinaryOp) {
	switch op.GetType() {
	case AsmBitAnd:
		ap.println("BitAnd")
	case AsmBitOr:
		ap.println("BitOr")
	case AsmBitXor:
		ap.println("BitXor")
	case AsmBitShiftLeft:
		ap.println("BitShiftLeft")
	case AsmBitShiftRight:
		ap.println("BitShiftRight")
	default:
		panic("unknown bit operator")
	}
}

func (ap *AsmPrinter) VisitImmediate(i *Immediate) {
	text := fmt.Sprintf("Immediate(%d)", i.Value)
	ap.println(text)
}

func (ap *AsmPrinter) VisitRegister(r *Register) {
	ap.println("Register(" + r.Name + ")")
}

func (ap *AsmPrinter) VisitPseudoReg(p *PseudoReg) {
	text := fmt.Sprintf("PseudoReg(%s)", p.Ident)
	ap.println(text)
}

func (ap *AsmPrinter) VisitStack(s *Stack) {
	text := fmt.Sprintf("Stack(%d)", s.N)
	ap.println(text)
}

func (ap *AsmPrinter) indent() {
	ap.offset += ap.delta
}

func (ap *AsmPrinter) dedent() {
	ap.offset -= ap.delta
}

func (ap *AsmPrinter) print(text string) {
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

func (ap *AsmPrinter) println(text string) {
	ap.print(text + "\n")
}
