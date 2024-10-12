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
	p.FuncDef.Accept(ap)
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

func (ap *AsmPrinter) VisitAllocStack(a *AllocStack) {
	text := fmt.Sprintf("AllocStack(%d)", a.N)
	ap.println(text)
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
