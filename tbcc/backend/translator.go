package backend

import "github.com/thomasbollmeier/writing-a-c-compiler/tbcc/tacky"

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
			NewMov(operand, NewRegister("AX")),
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
	default:
		panic("Unsupported instruction type")
	}
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
		panic("Unsupported value type")
	}
}

func (t *Translator) translateUnaryOperator(op tacky.UnaryOp) UnaryOp {
	switch op.GetType() {
	case tacky.TacComplement:
		return NewNot()
	case tacky.TacNegate:
		return NewNeg()
	default:
		panic("Unsupported operator type")
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
		r10 := NewRegister("R10")
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
