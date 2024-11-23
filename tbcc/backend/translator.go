package backend

import (
	"fmt"
	"github.com/thomasbollmeier/writing-a-c-compiler/tbcc/tacky"
	"math"
	"slices"
)

var argRegisters = []string{RegDI, RegSI, RegDX, RegCX, RegR8, RegR9}

type Translator struct{}

func NewTranslator() *Translator {
	return &Translator{}
}

func (t *Translator) Translate(program *tacky.Program) *Program {
	var funcDefs []FunctionDef
	for _, fun := range program.Funs {
		funcDefs = append(funcDefs, *t.translateFunctionDef(fun))
	}
	prog := NewProgram(funcDefs)
	prog, stackSizes := NewPseudoRegReplacer().Replace(prog)
	prog = NewInstructionAdapter(stackSizes).Adapt(prog)
	return prog
}

func (t *Translator) translateFunctionDef(fun tacky.Function) *FunctionDef {
	var instructions []Instruction
	name := fun.Ident

	numArgRegisters := len(argRegisters)
	numParams := len(fun.Parameters)

	if numParams <= numArgRegisters {
		for i, param := range fun.Parameters {
			regName := argRegisters[i]
			instructions = append(instructions, NewMov(NewRegister(regName), NewPseudoReg(param)))
		}
	} else {
		for i := 0; i < numArgRegisters; i++ {
			regName := argRegisters[i]
			param := fun.Parameters[i]
			instructions = append(instructions, NewMov(NewRegister(regName), NewPseudoReg(param)))
		}
		for i, param := range fun.Parameters[numArgRegisters:] {
			offset := 8 + (i+1)*8
			instructions = append(instructions, NewMov(NewStack(offset), NewPseudoReg(param)))
		}
	}

	instructions = append(instructions, t.translateAllInstructions(fun.Body)...)

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

	instrType := instruction.GetType()

	switch instrType {
	case tacky.TacReturn:
		ret := instruction.(*tacky.Return)
		operand := t.translateOperand(ret.Val)
		result = append(result,
			NewMov(operand, NewRegister(RegAX)),
			NewReturn())
		return result
	case tacky.TacUnary:
		unary := instruction.(*tacky.Unary)
		src := t.translateOperand(unary.Src)
		dst := t.translateOperand(unary.Dst)
		if unary.Op.GetType() != tacky.TacNot {
			op := t.translateUnaryOperator(unary.Op)
			result = append(result,
				NewMov(src, dst),
				NewUnary(op, dst))
		} else {
			result = append(result,
				NewCmp(NewImmediate(0), src),
				NewMov(NewImmediate(0), dst),
				NewSetCC(CcEq, dst))
		}
		return result
	case tacky.TacBinary:
		binary := instruction.(*tacky.Binary)
		src1 := t.translateOperand(binary.Src1)
		src2 := t.translateOperand(binary.Src2)
		dst := t.translateOperand(binary.Dst)
		switch binary.Op.GetType() {
		case tacky.TacAdd, tacky.TacSub, tacky.TacMul,
			tacky.TacBitAnd, tacky.TacBitOr, tacky.TacBitXor,
			tacky.TacBitShiftLeft, tacky.TacBitShiftRight:
			op := t.translateBinaryOperator(binary.Op)
			result = append(result,
				NewMov(src1, dst),
				NewBinary(op, src2, dst))
			return result
		case tacky.TacDiv:
			return t.createIDivInstructions(true, src1, src2, dst)
		case tacky.TacRemainder:
			return t.createIDivInstructions(false, src1, src2, dst)
		case tacky.TacEq, tacky.TacNotEq,
			tacky.TacGt, tacky.TacGtEq,
			tacky.TacLt, tacky.TacLtEq:
			return t.translateRelation(binary)
		default:
			panic("unsupported binary operator")
		}
	case tacky.TacJump:
		jump := instruction.(*tacky.Jump)
		return []Instruction{NewJump(jump.Target)}
	case tacky.TacJumpIfZero:
		jumpIfZero := instruction.(*tacky.JumpIfZero)
		cond := t.translateOperand(jumpIfZero.Condition)
		return []Instruction{
			NewCmp(NewImmediate(0), cond),
			NewJumpCC(CcEq, jumpIfZero.Target),
		}
	case tacky.TacJumpIfNotZero:
		jumpIfZero := instruction.(*tacky.JumpIfNotZero)
		cond := t.translateOperand(jumpIfZero.Condition)
		return []Instruction{
			NewCmp(NewImmediate(0), cond),
			NewJumpCC(CcNotEq, jumpIfZero.Target),
		}
	case tacky.TacCopy:
		cp := instruction.(*tacky.Copy)
		src := t.translateOperand(cp.Src)
		dst := t.translateOperand(cp.Dst)
		return []Instruction{NewMov(src, dst)}
	case tacky.TacLabel:
		label := instruction.(*tacky.Label)
		return []Instruction{NewLabel(label.Name)}
	case tacky.TacFunCall:
		funCall := instruction.(*tacky.FunctionCall)
		return t.translateFunctionCall(funCall)
	default:
		panic("unsupported instruction type")
	}
}

func (t *Translator) translateFunctionCall(funCall *tacky.FunctionCall) []Instruction {
	var ret []Instruction
	var registerArgs []tacky.Value
	var stackArgs []tacky.Value
	var stackPadding int

	numRegs := len(argRegisters)
	numArgs := len(funCall.Args)

	if numArgs <= numRegs {
		registerArgs = funCall.Args[:numArgs]
	} else {
		registerArgs = funCall.Args[:numRegs]
		stackArgs = funCall.Args[numRegs:]
		slices.Reverse(stackArgs)
	}

	if len(stackArgs)%2 != 0 {
		stackPadding = 8
		ret = append(ret, NewAllocStack(stackPadding))
	}

	// Fill registers with call arguments
	for i, argValue := range registerArgs {
		arg := t.translateOperand(argValue)
		ret = append(ret, NewMov(arg, NewRegister(argRegisters[i])))
	}

	ax := NewRegister(RegAX)

	var immediate *Immediate
	var register *Register
	var ok bool

	// Push remaining args onto the stack
	for _, argValue := range stackArgs {
		arg := t.translateOperand(argValue)
		if immediate, ok = arg.(*Immediate); ok {
			ret = append(ret, NewPush(immediate))
		} else if register, ok = arg.(*Register); ok {
			ret = append(ret, NewPush(register))
		} else {
			ret = append(ret, NewMov(arg, ax), NewPush(ax))
		}
	}

	ret = append(ret, NewCall(funCall.Name))

	// adjust stack pointer
	bytesToRemove := 8*len(stackArgs) + stackPadding
	if bytesToRemove > 0 {
		ret = append(ret, NewDeAllocStack(bytesToRemove))
	}

	// Set result
	dst := t.translateOperand(funCall.Dst)
	ret = append(ret, NewMov(ax, dst))

	return ret
}

func (t *Translator) translateRelation(binary *tacky.Binary) []Instruction {
	src1 := t.translateOperand(binary.Src1)
	src2 := t.translateOperand(binary.Src2)
	dst := t.translateOperand(binary.Dst)
	result := []Instruction{
		NewCmp(src2, src1), // order of operands switched!
		NewMov(NewImmediate(0), dst),
	}
	var conditionCode ConditionCode
	switch binary.Op.GetType() {
	case tacky.TacEq:
		conditionCode = CcEq
	case tacky.TacNotEq:
		conditionCode = CcNotEq
	case tacky.TacGt:
		conditionCode = CcGt
	case tacky.TacGtEq:
		conditionCode = CcGtEq
	case tacky.TacLt:
		conditionCode = CcLt
	case tacky.TacLtEq:
		conditionCode = CcLtEq
	default:
		panic(fmt.Sprintf("unsupported relation type: %v", binary.Op.GetType()))
	}

	result = append(result, NewSetCC(conditionCode, dst))

	return result
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
		panic(fmt.Sprintf("unsupported operator type: %v", op.GetType()))
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
	case tacky.TacBitAnd:
		return NewBitAnd()
	case tacky.TacBitOr:
		return NewBitOr()
	case tacky.TacBitXor:
		return NewBitXor()
	case tacky.TacBitShiftLeft:
		return NewBitShiftLeft()
	case tacky.TacBitShiftRight:
		return NewBitShiftRight()
	default:
		panic("unsupported operator type")
	}
}

type varOffsetsPerFunc map[string]int

type PseudoRegReplacer struct {
	varSizeByte  int
	currFunction string
	varOffsets   map[string]varOffsetsPerFunc
	result       any
}

func NewPseudoRegReplacer() *PseudoRegReplacer {
	return &PseudoRegReplacer{}
}

type VarSizesPerFunc map[string]int

func (pr *PseudoRegReplacer) Replace(p *Program) (*Program, VarSizesPerFunc) {
	pr.initialize()
	prog := pr.eval(p).(*Program)
	sizes := make(VarSizesPerFunc)
	for fn, offsets := range pr.varOffsets {
		sizes[fn] = len(offsets) * pr.varSizeByte
	}
	return prog, sizes
}

func (pr *PseudoRegReplacer) initialize() {
	pr.varSizeByte = 4
	pr.varOffsets = make(map[string]varOffsetsPerFunc)
	pr.result = nil
}

func (pr *PseudoRegReplacer) VisitProgram(p *Program) {
	var newFuncDefs []FunctionDef
	for _, fun := range p.FuncDefs {
		newFuncDef := pr.eval(&fun).(*FunctionDef)
		newFuncDefs = append(newFuncDefs, *newFuncDef)
	}
	pr.result = &Program{newFuncDefs}
}

func (pr *PseudoRegReplacer) VisitFunctionDef(f *FunctionDef) {
	var instructions []Instruction
	pr.currFunction = f.Name
	pr.varOffsets[f.Name] = make(varOffsetsPerFunc)
	for _, instruction := range f.Instructions {
		instructions = append(instructions, pr.eval(instruction).(Instruction))
	}
	pr.result = &FunctionDef{pr.currFunction, instructions}
	pr.currFunction = ""
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

func (pr *PseudoRegReplacer) VisitCmp(c *Cmp) {
	left := pr.eval(c.Left).(Operand)
	right := pr.eval(c.Right).(Operand)
	pr.result = NewCmp(left, right)
}

func (pr *PseudoRegReplacer) VisitIDiv(i *IDiv) {
	operand := pr.eval(i.Operand).(Operand)
	pr.result = &IDiv{operand}
}

func (pr *PseudoRegReplacer) VisitCdq(c *Cdq) {
	pr.result = c
}

func (pr *PseudoRegReplacer) VisitJump(j *Jump) {
	pr.result = j
}

func (pr *PseudoRegReplacer) VisitJumpCC(j *JumpCC) {
	pr.result = j
}

func (pr *PseudoRegReplacer) VisitSetCC(s *SetCC) {
	op := pr.eval(s.Op).(Operand)
	pr.result = NewSetCC(s.CondCode, op)
}

func (pr *PseudoRegReplacer) VisitLabel(l *Label) {
	pr.result = l
}

func (pr *PseudoRegReplacer) VisitAllocStack(a *AllocStack) {
	pr.result = a
}

func (pr *PseudoRegReplacer) VisitDeAllocStack(d *DeAllocStack) {
	pr.result = d
}

func (pr *PseudoRegReplacer) VisitPush(p *Push) {
	op := pr.eval(p.Op).(Operand)
	pr.result = NewPush(op)
}

func (pr *PseudoRegReplacer) VisitCall(c *Call) {
	pr.result = c
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

func (pr *PseudoRegReplacer) VisitBitOp(op BinaryOp) {
	pr.result = op
}

func (pr *PseudoRegReplacer) VisitImmediate(i *Immediate) {
	pr.result = i
}

func (pr *PseudoRegReplacer) VisitRegister(r *Register) {
	pr.result = r
}

func (pr *PseudoRegReplacer) VisitPseudoReg(p *PseudoReg) {
	varOffsets := pr.varOffsets[pr.currFunction]
	offset, ok := varOffsets[p.Ident]
	if !ok {
		numVars := len(varOffsets)
		numVars++
		offset = -numVars * pr.varSizeByte
		pr.varOffsets[pr.currFunction][p.Ident] = offset
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
	stackSizes VarSizesPerFunc
	result     any
}

func NewInstructionAdapter(stackSizes VarSizesPerFunc) *InstructionAdapter {
	return &InstructionAdapter{stackSizes, nil}
}

func (ia *InstructionAdapter) Adapt(program *Program) *Program {
	return ia.eval(program).(*Program)
}

func (ia *InstructionAdapter) VisitProgram(p *Program) {
	var newFuncDefs []FunctionDef
	for _, fun := range p.FuncDefs {
		newFuncDef := ia.eval(&fun).(*FunctionDef)
		newFuncDefs = append(newFuncDefs, *newFuncDef)
	}
	ia.result = &Program{newFuncDefs}

}

func (ia *InstructionAdapter) VisitFunctionDef(f *FunctionDef) {
	stackSize := ia.stackSizes[f.Name]
	// round up to next multiple of 16 for stack alignment:
	stackSize = int(math.Ceil(float64(stackSize)/16.0) * 16.0)
	newInstructions := []Instruction{NewAllocStack(stackSize)}
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
	case AsmAdd, AsmSub, AsmBitAnd, AsmBitOr, AsmBitXor:
		if b.Operand1.GetType() == AsmStack && b.Operand2.GetType() == AsmStack {
			r10 := NewRegister(RegR10)
			ia.result = []Instruction{
				NewMov(b.Operand1, r10),
				NewBinary(b.Op, r10, b.Operand2),
			}
		} else if b.Operand2.GetType() == AsmImmediate {
			r11 := NewRegister(RegR11)
			ia.result = []Instruction{
				NewMov(b.Operand2, r11),
				NewBinary(b.Op, b.Operand1, r11),
			}
		} else {
			ia.result = []Instruction{b}
		}
	case AsmBitShiftLeft, AsmBitShiftRight:
		if b.Operand1.GetType() == AsmStack {
			cx := NewRegister(RegCX)
			ia.result = []Instruction{
				NewMov(b.Operand1, cx),
				NewBinary(b.Op, cx, b.Operand2),
				NewMov(cx, b.Operand1),
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

func (ia *InstructionAdapter) VisitCmp(c *Cmp) {
	leftType := c.Left.GetType()
	rightType := c.Right.GetType()

	if leftType == AsmStack && rightType == AsmStack {
		r10 := NewRegister(RegR10)
		ia.result = []Instruction{
			NewMov(c.Left, r10),
			NewCmp(r10, c.Right),
		}
	} else if rightType == AsmImmediate {
		r11 := NewRegister(RegR11)
		ia.result = []Instruction{
			NewMov(c.Right, r11),
			NewCmp(c.Left, r11),
		}
	} else {
		ia.result = []Instruction{c}
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

func (ia *InstructionAdapter) VisitJump(j *Jump) {
	ia.result = []Instruction{j}
}

func (ia *InstructionAdapter) VisitJumpCC(j *JumpCC) {
	ia.result = []Instruction{j}
}

func (ia *InstructionAdapter) VisitSetCC(s *SetCC) {
	ia.result = []Instruction{s}
}

func (ia *InstructionAdapter) VisitLabel(l *Label) {
	ia.result = []Instruction{l}
}

func (ia *InstructionAdapter) VisitAllocStack(a *AllocStack) {
	ia.result = []Instruction{a}
}

func (ia *InstructionAdapter) VisitDeAllocStack(d *DeAllocStack) {
	ia.result = []Instruction{d}
}

func (ia *InstructionAdapter) VisitPush(p *Push) {
	ia.result = []Instruction{p}
}

func (ia *InstructionAdapter) VisitCall(c *Call) {
	ia.result = []Instruction{c}
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

func (ia *InstructionAdapter) VisitBitOp(op BinaryOp) {
	ia.result = op
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
