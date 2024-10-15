package backend

import (
	"fmt"
)

type CodeGenerator struct {
	code string
}

func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{}
}

func (cg *CodeGenerator) GenerateCode(program Program) string {
	cg.code = ""
	program.Accept(cg)
	cg.writeln(".section .note.GNU-stack,\"\",@progbits")
	return cg.code
}

func (cg *CodeGenerator) VisitProgram(p *Program) {
	p.FuncDef.Accept(cg)
}

func (cg *CodeGenerator) VisitFunctionDef(f *FunctionDef) {
	cg.writeln("\t.globl " + f.Name)
	cg.writeln(f.Name + ":")
	cg.writeln("\tpushq %rbp")
	cg.writeln("\tmovq %rsp, %rbp")
	for _, instr := range f.Instructions {
		instr.Accept(cg)
	}
}

func (cg *CodeGenerator) VisitMov(m *Mov) {
	cg.write("\tmovl ")
	m.Src.Accept(cg)
	cg.write(", ")
	m.Dst.Accept(cg)
	cg.writeln("")
}

func (cg *CodeGenerator) VisitUnary(u *Unary) {
	cg.write("\t")
	u.Op.Accept(cg)
	cg.write(" ")
	u.Operand.Accept(cg)
	cg.writeln("")
}

func (cg *CodeGenerator) VisitBinary(b *Binary) {
	cg.write("\t")
	b.Op.Accept(cg)
	cg.write(" ")
	b.Operand1.Accept(cg)
	cg.write(", ")
	b.Operand2.Accept(cg)
	cg.writeln("")
}

func (cg *CodeGenerator) VisitIDiv(i *IDiv) {
	cg.write("\tidivl ")
	i.Operand.Accept(cg)
	cg.writeln("")
}

func (cg *CodeGenerator) VisitCdq(*Cdq) {
	cg.writeln("\tcdq")
}

func (cg *CodeGenerator) VisitAllocStack(a *AllocStack) {
	cg.writeln(fmt.Sprintf("\tsubq $%d, %%rsp", a.N))
}

func (cg *CodeGenerator) VisitReturn() {
	cg.writeln("\tmovq %rbp, %rsp")
	cg.writeln("\tpopq %rbp")
	cg.writeln("\tret")
}

func (cg *CodeGenerator) VisitNeg(*Neg) {
	cg.write("negl")
}

func (cg *CodeGenerator) VisitNot(*Not) {
	cg.write("notl")
}

func (cg *CodeGenerator) VisitAdd(*Add) {
	cg.write("addl")
}

func (cg *CodeGenerator) VisitSub(*Sub) {
	cg.write("subl")
}

func (cg *CodeGenerator) VisitMul(*Mul) {
	cg.write("imull")
}

func (cg *CodeGenerator) VisitImmediate(i *Immediate) {
	cg.write(fmt.Sprintf("$%d", i.Value))
}

func (cg *CodeGenerator) VisitRegister(r *Register) {
	switch r.Name {
	case RegAX:
		cg.write("%eax")
	case RegDX:
		cg.write("%edx")
	case RegR10:
		cg.write("%r10d")
	case RegR11:
		cg.write("%r11d")
	default:
		panic("unknown register name")
	}
}

func (cg *CodeGenerator) VisitPseudoReg(*PseudoReg) {
	panic("this should not be called")
}

func (cg *CodeGenerator) VisitStack(s *Stack) {
	cg.write(fmt.Sprintf("%d(%%rbp)", s.N))
}

func (cg *CodeGenerator) write(text string) {
	cg.code += text
}

func (cg *CodeGenerator) writeln(text string) {
	cg.write(text + "\n")
}
