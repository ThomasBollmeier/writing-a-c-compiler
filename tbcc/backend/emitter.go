package backend

import "github.com/thomasbollmeier/writing-a-c-compiler/tbcc/tacky"

type Emitter struct{}

func NewEmitter() *Emitter {
	return &Emitter{}
}

func (e *Emitter) Emit(program *tacky.Program) *Program {
	functionDef := e.emitFunctionDef(program.Fun)
	return &Program{*functionDef}
}

func (e *Emitter) emitFunctionDef(fun tacky.Function) *FunctionDef {
	name := fun.Ident
	instructions := e.emitAllInstructions(fun.Body)
	return &FunctionDef{name, instructions}
}

func (e *Emitter) emitAllInstructions(instructions []tacky.Instruction) []Instruction {
	var result []Instruction
	for _, instruction := range instructions {
		result = append(result, e.emitInstructions(instruction)...)
	}
	return result
}

func (e *Emitter) emitInstructions(instruction tacky.Instruction) []Instruction {
	var result []Instruction
	switch instruction.GetType() {
	case tacky.TacReturn:
		ret := instruction.(*tacky.Return)
		operand := e.emitOperand(ret.Val)
		result = append(result, NewMov(operand, NewRegister("AX")))
		result = append(result, NewReturn())
		return result
	case tacky.TacUnary:
		unary := instruction.(*tacky.Unary)
		op := e.emitUnaryOperator(unary.Op)
		src := e.emitOperand(unary.Src)
		dst := e.emitOperand(unary.Dst)
		result = append(result, NewMov(src, dst), NewUnary(op, dst))
		return result
	default:
		panic("Unsupported instruction type")
	}
}

func (e *Emitter) emitOperand(value tacky.Value) Operand {
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

func (e *Emitter) emitUnaryOperator(op tacky.UnaryOp) UnaryOp {
	switch op.GetType() {
	case tacky.TacComplement:
		return NewNot()
	case tacky.TacNegate:
		return NewNeg()
	default:
		panic("Unsupported operator type")
	}
}
