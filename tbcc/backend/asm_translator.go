package backend

import (
	"github.com/thomasbollmeier/writing-a-c-compiler/tbcc/frontend"
)

type AsmTranslator struct {
	result any
}

func NewAsmTranslator() *AsmTranslator {
	return &AsmTranslator{nil}
}

func (a *AsmTranslator) Translate(program *frontend.Program) *Program {
	program.Accept(a)
	return a.result.(*Program)
}

func (a *AsmTranslator) VisitProgram(p *frontend.Program) {
	p.Func.Accept(a)
	funcDef := a.result.(*FunctionDef)
	a.result = NewProgram(*funcDef)
}

func (a *AsmTranslator) VisitFunction(f *frontend.Function) {
	name := f.Name
	f.Body.Accept(a)
	instructions := a.result.([]Instruction)
	a.result = NewFunctionDef(name, instructions)
}

func (a *AsmTranslator) VisitReturn(r *frontend.ReturnStmt) {
	var instructions []Instruction
	r.Expression.Accept(a)
	instructions = append(instructions, NewMov(a.result.(Operand), NewRegister("AX")))
	instructions = append(instructions, NewReturn())
	a.result = instructions
}

func (a *AsmTranslator) VisitInteger(i *frontend.IntegerLiteral) {
	a.result = &Immediate{i.Value}
}

func (a *AsmTranslator) VisitIdentifier(*frontend.Identifier) {
	panic("not implemented")
}

func (a *AsmTranslator) VisitUnary(*frontend.UnaryExpression) {
	panic("not implemented")
}
