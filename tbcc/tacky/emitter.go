package tacky

import (
	"fmt"
	"github.com/thomasbollmeier/writing-a-c-compiler/tbcc/frontend"
)

type Emitter struct {
	nextCounter uint
}

func NewEmitter() *Emitter {
	return &Emitter{0}
}

func (e *Emitter) Emit(progam *frontend.Program) *Program {
	return &Program{e.emitFunction(&progam.Func)}
}

func (e *Emitter) emitFunction(f *frontend.Function) Function {
	return Function{
		f.Name,
		e.emitBody(f.Body),
	}
}

func (e *Emitter) emitBody(body frontend.Statement) []Instruction {
	switch body.GetType() {
	case frontend.AstReturn:
		retStmt := body.(*frontend.ReturnStmt)
		val, instructions := e.emitExpr(retStmt.Expression)
		instructions = append(instructions, &Return{val})
		return instructions
	default:
		panic("unsupported statement type")
	}
}

func (e *Emitter) emitExpr(expr frontend.Expression) (Value, []Instruction) {
	switch expr.GetType() {
	case frontend.AstInteger:
		val := expr.(*frontend.IntegerLiteral).Value
		return &IntConstant{val}, nil
	case frontend.AstUnary:
		unary := expr.(*frontend.UnaryExpression)
		src, instructions := e.emitExpr(unary.Right)
		dst := &Var{e.createVarName()}
		unaryOp := e.getUnaryOp(unary.Operator)
		instructions = append(instructions, &Unary{unaryOp, src, dst})
		return dst, instructions
	default:
		panic("unsupported expression type")
	}
}

func (e *Emitter) createVarName() string {
	varName := fmt.Sprintf("tmp.%d", e.nextCounter)
	e.nextCounter++
	return varName
}

func (e *Emitter) getUnaryOp(op string) UnaryOp {
	switch op {
	case "-":
		return &Negate{}
	case "~":
		return &Complement{}
	default:
		panic("unsupported operator")
	}
}
