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

func (cg *CodeGenerator) write(text string) {
	cg.code += text
}

func (cg *CodeGenerator) writeln(text string) {
	cg.write(text + "\n")
}

func (cg *CodeGenerator) VisitProgram(p *Program) {
	p.FuncDef.Accept(cg)
}

func (cg *CodeGenerator) VisitFunctionDef(f *FunctionDef) {
	cg.writeln("\t.globl " + f.Name)
	cg.writeln(f.Name + ":")
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

func (cg *CodeGenerator) VisitReturn() {
	cg.writeln("\tret")
}

func (cg *CodeGenerator) VisitImmediate(i *Immediate) {
	text := fmt.Sprintf("$%d", i.Value)
	cg.write(text)
}

func (cg *CodeGenerator) VisitRegister() {
	cg.write("%eax")
}
