package backend

import (
	"github.com/thomasbollmeier/writing-a-c-compiler/tbcc/tacky"
)

type Translator struct{}

func NewTranslator() *Translator {
	return &Translator{}
}

func (t *Translator) Translate(program *tacky.Program) *Program {
	functionDef := t.translateFunctionDef(program.Fun)
	prog := NewProgram(*functionDef)
	prog, stackSize := NewPseudoRegReplacer().Replace(prog)
	prog = NewInstructionAdapter(stackSize).Adapt(prog)
	return prog
}

func (t *Translator) translateFunctionDef(fun tacky.Function) *FunctionDef {
	name := fun.Ident
	instructions := t.translateAllInstructions(fun.Body)
	return NewFunctionDef(name, instructions)
}

func (t *Translator) translateAllInstructions(instructions []tacky.Instruction) []Instruction {
	var result []Instruction
	for _, instruction := range instructions {
		result = append(result, t.translateInstructions(instruction)...)
	}
	return result
}

func (t *Translator) translateInstructions(instruction tacky.Instruction) []Instruction {
	var result []Instruction
	switch instruction.GetType() {
	case tacky.TacReturn:
		ret := instruction.(*tacky.Return)
		operand := t.translateOperand(ret.Val)
		result = append(result,
			NewMov(operand, NewRegister(RegAX)),
			NewReturn())
		return result
	case tacky.TacUnary:
		unary := instruction.(*tacky.Unary)
		op := t.translateUnaryOperator(unary.Op)
		src := t.translateOperand(unary.Src)
		dst := t.translateOperand(unary.Dst)
		result = append(result,
			NewMov(src, dst),
			NewUnary(op, dst))
		return result
	case tacky.TacBinary:
		binary := instruction.(*tacky.Binary)
		src1 := t.translateOperand(binary.Src1)
		src2 := t.translateOperand(binary.Src2)
		dst := t.translateOperand(binary.Dst)
		switch binary.Op.GetType() {
		case tacky.TacAdd, tacky.TacSub, tacky.TacMul:
			op := t.translateBinaryOperator(binary.Op)
			result = append(result,
				NewMov(src1, dst),
				NewBinary(op, src2, dst))
			return result
		case tacky.TacDiv:
			return t.createIDivInstructions(true, src1, src2, dst)
		case tacky.TacRemainder:
			return t.createIDivInstructions(false, src1, src2, dst)
		default:
			panic("unsupported binary operator")
		}
	default:
		panic("unsupported instruction type")
	}
}

func (t *Translator) createIDivInstructions(calcQuotient bool, src1, src2, dst Operand) []Instruction {
	var result []Instruction
	result = append(result, NewMov(src1, NewRegister(RegAX)))
	result = append(result, NewCdq())
	result = append(result, NewIDiv(src2))
	if calcQuotient {
		result = append(result, NewMov(NewRegister(RegAX), dst))
	} else {
		result = append(result, NewMov(NewRegister(RegDX), dst))
	}
	return result
}

func (t *Translator) translateOperand(value tacky.Value) Operand {
	switch value.GetType() {
	case tacky.TacIntConstant:
		intLiteral := value.(*tacky.IntConstant)
		return NewImmediate(intLiteral.Val)
	case tacky.TacVar:
		variable := value.(*tacky.Var)
		return NewPseudoReg(variable.Ident)
	default:
		panic("unsupported value type")
	}
}

func (t *Translator) translateUnaryOperator(op tacky.UnaryOp) UnaryOp {
	switch op.GetType() {
	case tacky.TacComplement:
		return NewNot()
	case tacky.TacNegate:
		return NewNeg()
	default:
		panic("unsupported operator type")
	}
}

func (t *Translator) translateBinaryOperator(op tacky.BinaryOp) BinaryOp {
	switch op.GetType() {
	case tacky.TacAdd:
		return NewAdd()
	case tacky.TacSub:
		return NewSub()
	case tacky.TacMul:
		return NewMul()
	default:
		panic("unsupported operator type")
	}
}

type PseudoRegReplacer struct {
	varSizeByte int
	numVars     int
	varOffsets  map[string]int
	result      any
}

func NewPseudoRegReplacer() *PseudoRegReplacer {
	return &PseudoRegReplacer{}
}

func (pr *PseudoRegReplacer) Replace(p *Program) (*Program, int) {
	pr.initialize()
	prog := pr.eval(p).(*Program)
	totalVarSize := pr.numVars * pr.varSizeByte
	return prog, totalVarSize
}

func (pr *PseudoRegReplacer) initialize() {
	pr.varSizeByte = 4
	pr.numVars = 0
	pr.varOffsets = make(map[string]int)
	pr.result = nil
}

func (pr *PseudoRegReplacer) VisitProgram(p *Program) {
	newFuncDef := pr.eval(&p.FuncDef).(*FunctionDef)
	pr.result = &Program{*newFuncDef}
}

func (pr *PseudoRegReplacer) VisitFunctionDef(f *FunctionDef) {
	var instructions []Instruction
	name := f.Name
	for _, instruction := range f.Instructions {
		instructions = append(instructions, pr.eval(instruction).(Instruction))
	}
	pr.result = &FunctionDef{name, instructions}
}

func (pr *PseudoRegReplacer) VisitMov(m *Mov) {
	src := pr.eval(m.Src).(Operand)
	dst := pr.eval(m.Dst).(Operand)
	pr.result = &Mov{src, dst}
}

func (pr *PseudoRegReplacer) VisitUnary(u *Unary) {
	operand := pr.eval(u.Operand).(Operand)
	pr.result = &Unary{u.Op, operand}
}

func (pr *PseudoRegReplacer) VisitBinary(b *Binary) {
	operand1 := pr.eval(b.Operand1).(Operand)
	operand2 := pr.eval(b.Operand2).(Operand)
	pr.result = &Binary{b.Op, operand1, operand2}
}

func (pr *PseudoRegReplacer) VisitIDiv(i *IDiv) {
	operand := pr.eval(i.Operand).(Operand)
	pr.result = &IDiv{operand}
}

func (pr *PseudoRegReplacer) VisitCdq(c *Cdq) {
	pr.result = c
}

func (pr *PseudoRegReplacer) VisitAllocStack(a *AllocStack) {
	pr.result = a
}

func (pr *PseudoRegReplacer) VisitReturn() {
	pr.result = &Return{}
}

func (pr *PseudoRegReplacer) VisitNeg(n *Neg) {
	pr.result = n
}

func (pr *PseudoRegReplacer) VisitNot(n *Not) {
	pr.result = n
}

func (pr *PseudoRegReplacer) VisitAdd(a *Add) {
	pr.result = a
}

func (pr *PseudoRegReplacer) VisitSub(s *Sub) {
	pr.result = s
}

func (pr *PseudoRegReplacer) VisitMul(m *Mul) {
	pr.result = m
}

func (pr *PseudoRegReplacer) VisitImmediate(i *Immediate) {
	pr.result = i
}

func (pr *PseudoRegReplacer) VisitRegister(r *Register) {
	pr.result = r
}

func (pr *PseudoRegReplacer) VisitPseudoReg(p *PseudoReg) {
	offset, ok := pr.varOffsets[p.Ident]
	if !ok {
		pr.numVars++
		offset = -pr.numVars * pr.varSizeByte
		pr.varOffsets[p.Ident] = offset
	}
	pr.result = NewStack(offset)
}

func (pr *PseudoRegReplacer) VisitStack(s *Stack) {
	pr.result = s
}

func (pr *PseudoRegReplacer) eval(ast AST) any {
	ast.Accept(pr)
	return pr.result
}

type InstructionAdapter struct {
	stackSize int
	result    any
}

func NewInstructionAdapter(stackSize int) *InstructionAdapter {
	return &InstructionAdapter{stackSize, nil}
}

func (ia *InstructionAdapter) Adapt(program *Program) *Program {
	return ia.eval(program).(*Program)
}

func (ia *InstructionAdapter) VisitProgram(p *Program) {
	newFuncDef := ia.eval(&p.FuncDef).(*FunctionDef)
	ia.result = &Program{*newFuncDef}
}

func (ia *InstructionAdapter) VisitFunctionDef(f *FunctionDef) {
	newInstructions := []Instruction{NewAllocStack(ia.stackSize)}
	for _, instruction := range f.Instructions {
		newInstructions = append(newInstructions, ia.eval(instruction).([]Instruction)...)
	}
	ia.result = &FunctionDef{f.Name, newInstructions}
}

func (ia *InstructionAdapter) VisitMov(m *Mov) {
	if m.Src.GetType() == AsmStack && m.Dst.GetType() == AsmStack {
		r10 := NewRegister(RegR10)
		ia.result = []Instruction{
			&Mov{m.Src, r10},
			&Mov{r10, m.Dst},
		}
	} else {
		ia.result = []Instruction{m}
	}
}

func (ia *InstructionAdapter) VisitUnary(u *Unary) {
	ia.result = []Instruction{u}
}

func (ia *InstructionAdapter) VisitBinary(b *Binary) {
	switch b.Op.GetType() {
	case AsmAdd, AsmSub:
		if b.Operand1.GetType() == AsmStack && b.Operand2.GetType() == AsmStack {
			r10 := NewRegister(RegR10)
			ia.result = []Instruction{
				NewMov(b.Operand1, r10),
				NewBinary(b.Op, r10, b.Operand2),
			}
		} else {
			ia.result = []Instruction{b}
		}
	case AsmMul:
		if b.Operand2.GetType() == AsmStack {
			r11 := NewRegister(RegR11)
			ia.result = []Instruction{
				NewMov(b.Operand2, r11),
				NewBinary(b.Op, b.Operand1, r11),
				NewMov(r11, b.Operand2),
			}
		} else {
			ia.result = []Instruction{b}
		}
	default:
		panic("unsupported binary operator")
	}
}

func (ia *InstructionAdapter) VisitIDiv(i *IDiv) {
	if i.Operand.GetType() == AsmImmediate {
		r10 := NewRegister(RegR10)
		ia.result = []Instruction{
			NewMov(i.Operand, r10),
			NewIDiv(r10),
		}
	} else {
		ia.result = []Instruction{i}
	}
}

func (ia *InstructionAdapter) VisitCdq(c *Cdq) {
	ia.result = []Instruction{c}
}

func (ia *InstructionAdapter) VisitAllocStack(a *AllocStack) {
	ia.result = []Instruction{a}
}

func (ia *InstructionAdapter) VisitReturn() {
	ia.result = []Instruction{&Return{}}
}

func (ia *InstructionAdapter) VisitNeg(n *Neg) {
	ia.result = n
}

func (ia *InstructionAdapter) VisitNot(n *Not) {
	ia.result = n
}

func (ia *InstructionAdapter) VisitAdd(a *Add) {
	ia.result = a
}

func (ia *InstructionAdapter) VisitSub(s *Sub) {
	ia.result = s
}

func (ia *InstructionAdapter) VisitMul(m *Mul) {
	ia.result = m
}

func (ia *InstructionAdapter) VisitImmediate(i *Immediate) {
	ia.result = i
}

func (ia *InstructionAdapter) VisitRegister(r *Register) {
	ia.result = r
}

func (ia *InstructionAdapter) VisitPseudoReg(p *PseudoReg) {
	ia.result = p
}

func (ia *InstructionAdapter) VisitStack(s *Stack) {
	ia.result = s
}

func (ia *InstructionAdapter) eval(ast AST) any {
	ast.Accept(ia)
	return ia.result
}
