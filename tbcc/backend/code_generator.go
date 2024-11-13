package backend

import (
	"fmt"
)

type CodeGenerator struct {
	code             string
	use1ByteRegister bool
}

func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{
		code:             "",
		use1ByteRegister: false,
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
	cg.use1ByteRegister = true
	cg.write(fmt.Sprintf("\tset%s ", cg.getCondInstrSuffix(s.CondCode)))
	s.Op.Accept(cg)
	cg.writeln("")
	cg.use1ByteRegister = false
}

func (cg *CodeGenerator) VisitLabel(l *Label) {
	cg.writeln(fmt.Sprintf(".L%s:", l.Identifier))
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
	if !cg.use1ByteRegister {
		switch r.Name {
		case RegAX:
			cg.write("%eax")
		case RegCX:
			cg.write("%ecx")
		case RegDX:
			cg.write("%edx")
		case RegR10:
			cg.write("%r10d")
		case RegR11:
			cg.write("%r11d")
		default:
			panic("unknown register name")
		}
	} else {
		switch r.Name {
		case RegAX:
			cg.write("%al")
		case RegCX:
			cg.write("%cl")
		case RegDX:
			cg.write("%dl")
		case RegR10:
			cg.write("%r10b")
		case RegR11:
			cg.write("%r11b")
		default:
			panic("unknown register name")
		}
	}
}

func (cg *CodeGenerator) VisitPseudoReg(*PseudoReg) {
	panic("this should not be called")
}

func (cg *CodeGenerator) VisitStack(s *Stack) {
	cg.write(fmt.Sprintf("%d(%%rbp)", s.N))
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
