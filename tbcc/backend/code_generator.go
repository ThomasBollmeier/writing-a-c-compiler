package backend

import (
	"fmt"
	"github.com/thomasbollmeier/writing-a-c-compiler/tbcc/frontend"
)

type regByteMode uint

const (
	regByteMode1 regByteMode = iota
	regByteMode4
	regByteMode8
)

var registerNames = map[string]map[regByteMode]string{
	RegAX: {
		regByteMode8: "%rax",
		regByteMode4: "%eax",
		regByteMode1: "%al"},
	RegCX: {
		regByteMode8: "%rcx",
		regByteMode4: "%ecx",
		regByteMode1: "%cl"},
	RegDX: {
		regByteMode8: "%rdx",
		regByteMode4: "%edx",
		regByteMode1: "%dl"},
	RegDI: {
		regByteMode8: "%rdi",
		regByteMode4: "%edi",
		regByteMode1: "%dil"},
	RegSI: {
		regByteMode8: "%rsi",
		regByteMode4: "%esi",
		regByteMode1: "%sil"},
	RegR8: {
		regByteMode8: "%r8",
		regByteMode4: "%r8d",
		regByteMode1: "%r8b"},
	RegR9: {
		regByteMode8: "%r9",
		regByteMode4: "%r9d",
		regByteMode1: "%r9b"},
	RegR10: {
		regByteMode8: "%r10",
		regByteMode4: "%r10d",
		regByteMode1: "%r10b"},
	RegR11: {
		regByteMode8: "%r11",
		regByteMode4: "%r11d",
		regByteMode1: "%r11b"},
}

type CodeGenerator struct {
	code   string
	rbmode regByteMode
	env    *frontend.Environment
}

func NewCodeGenerator(env *frontend.Environment) *CodeGenerator {
	return &CodeGenerator{
		code:   "",
		rbmode: regByteMode4,
		env:    env,
	}
}

func (cg *CodeGenerator) GenerateCode(program Program) string {
	cg.code = ""
	program.Accept(cg)
	cg.writeln(".section .note.GNU-stack,\"\",@progbits")
	return cg.code
}

func (cg *CodeGenerator) VisitProgram(p *Program) {
	for _, funcDef := range p.FuncDefs {
		funcDef.Accept(cg)
	}
}

func (cg *CodeGenerator) VisitFunctionDef(f *FunctionDef) {
	funcName := cg.getFunctionName(f.Name)
	cg.writeln("\t.globl " + funcName)
	cg.writeln(funcName + ":")
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

func (cg *CodeGenerator) VisitCmp(c *Cmp) {
	cg.write("\tcmpl ")
	c.Left.Accept(cg)
	cg.write(", ")
	c.Right.Accept(cg)
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

func (cg *CodeGenerator) VisitJump(j *Jump) {
	cg.writeln(fmt.Sprintf("\tjmp .L%s", j.Identifier))
}

func (cg *CodeGenerator) VisitJumpCC(j *JumpCC) {
	cg.writeln(fmt.Sprintf("\tj%s .L%s",
		cg.getCondInstrSuffix(j.CondCode),
		j.Identifier))
}

func (cg *CodeGenerator) VisitSetCC(s *SetCC) {
	savedRegByteMode := cg.rbmode
	cg.rbmode = regByteMode1
	cg.write(fmt.Sprintf("\tset%s ", cg.getCondInstrSuffix(s.CondCode)))
	s.Op.Accept(cg)
	cg.writeln("")
	cg.rbmode = savedRegByteMode
}

func (cg *CodeGenerator) VisitLabel(l *Label) {
	cg.writeln(fmt.Sprintf(".L%s:", l.Identifier))
}

func (cg *CodeGenerator) VisitAllocStack(a *AllocStack) {
	cg.writeln(fmt.Sprintf("\tsubq $%d, %%rsp", a.N))
}

func (cg *CodeGenerator) VisitDeAllocStack(d *DeAllocStack) {
	cg.writeln(fmt.Sprintf("\taddq $%d, %%rsp", d.N))
}

func (cg *CodeGenerator) VisitPush(p *Push) {
	savedRegByteMode := cg.rbmode
	cg.rbmode = regByteMode8
	cg.write("\tpushq ")
	p.Op.Accept(cg)
	cg.writeln("")
	cg.rbmode = savedRegByteMode
}

func (cg *CodeGenerator) VisitCall(c *Call) {
	funcName := cg.getFunctionName(c.Identifier)
	cg.writeln(fmt.Sprintf("\tcall %s", funcName))
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

func (cg *CodeGenerator) VisitBitOp(op BinaryOp) {
	switch op.GetType() {
	case AsmBitAnd:
		cg.write("and")
	case AsmBitOr:
		cg.write("or")
	case AsmBitXor:
		cg.write("xor")
	case AsmBitShiftLeft:
		cg.write("shl")
	case AsmBitShiftRight:
		cg.write("shr")
	default:
		panic(fmt.Sprintf("unknown op type: %v", op.GetType()))
	}
}

func (cg *CodeGenerator) VisitImmediate(i *Immediate) {
	cg.write(fmt.Sprintf("$%d", i.Value))
}

func (cg *CodeGenerator) VisitRegister(r *Register) {
	cg.write(registerNames[r.Name][cg.rbmode])
}

func (cg *CodeGenerator) VisitPseudoReg(*PseudoReg) {
	panic("this should not be called")
}

func (cg *CodeGenerator) VisitStack(s *Stack) {
	cg.write(fmt.Sprintf("%d(%%rbp)", s.N))
}

func (cg *CodeGenerator) getFunctionName(funcName string) string {
	if cg.isOwnFunction(funcName) {
		return funcName
	}
	return funcName + "@PLT" // Linux specific: Procedure Linkage Table
}

func (cg *CodeGenerator) isOwnFunction(funcName string) bool {

	entry, _ := cg.env.Get(funcName)
	if entry == nil {
		return false
	}

	typeInfo := entry.GetTypeInfo()
	if typeInfo.GetTypeId() != frontend.TypeFunc {
		return false
	}

	funcInfo := typeInfo.(*frontend.FuncInfo)

	return funcInfo.IsDefined
}

func (cg *CodeGenerator) getCondInstrSuffix(conditionCode ConditionCode) string {
	switch conditionCode {
	case CcEq:
		return "e"
	case CcNotEq:
		return "ne"
	case CcLt:
		return "l"
	case CcLtEq:
		return "le"
	case CcGt:
		return "g"
	case CcGtEq:
		return "ge"
	default:
		panic(fmt.Sprintf("unknown condition code: %v", conditionCode))
	}
}

func (cg *CodeGenerator) write(text string) {
	cg.code += text
}

func (cg *CodeGenerator) writeln(text string) {
	cg.write(text + "\n")
}
