package tacky

import (
	"fmt"
	"github.com/thomasbollmeier/writing-a-c-compiler/tbcc/frontend"
)

type Translator struct {
	nextCounter uint
}

func NewTranslator() *Translator {
	return &Translator{0}
}

func (t *Translator) Translate(progam *frontend.Program) *Program {
	return &Program{t.translateFunction(&progam.Func)}
}

func (t *Translator) translateFunction(f *frontend.Function) Function {
	return Function{
		f.Name,
		t.translateBody(f.Body),
	}
}

func (t *Translator) translateBody(body frontend.Statement) []Instruction {
	switch body.GetType() {
	case frontend.AstReturn:
		retStmt := body.(*frontend.ReturnStmt)
		val, instructions := t.translateExpr(retStmt.Expression)
		instructions = append(instructions, &Return{val})
		return instructions
	default:
		panic("unsupported statement type")
	}
}

func (t *Translator) translateExpr(expr frontend.Expression) (Value, []Instruction) {
	switch expr.GetType() {
	case frontend.AstInteger:
		val := expr.(*frontend.IntegerLiteral).Value
		return &IntConstant{val}, nil
	case frontend.AstUnary:
		unary := expr.(*frontend.UnaryExpression)
		src, instructions := t.translateExpr(unary.Right)
		dst := &Var{t.createVarName()}
		unaryOp := t.getUnaryOp(unary.Operator)
		instructions = append(instructions, &Unary{unaryOp, src, dst})
		return dst, instructions
	case frontend.AstBinary:
		binary := expr.(*frontend.BinaryExpression)
		src1, instructions := t.translateExpr(binary.Left)
		src2, instructions2 := t.translateExpr(binary.Right)
		instructions = append(instructions, instructions2...)
		dst := &Var{t.createVarName()}
		binaryOp := t.getBinaryOp(binary.Operator)
		instructions = append(instructions, &Binary{binaryOp, src1, src2, dst})
		return dst, instructions
	default:
		panic("unsupported expression type")
	}
}

func (t *Translator) createVarName() string {
	varName := fmt.Sprintf("tmp.%d", t.nextCounter)
	t.nextCounter++
	return varName
}

func (t *Translator) getUnaryOp(op string) UnaryOp {
	switch op {
	case "-":
		return &Negate{}
	case "~":
		return &Complement{}
	default:
		panic("unsupported operator")
	}
}

func (t *Translator) getBinaryOp(op string) BinaryOp {
	switch op {
	case "+":
		return &Add{}
	case "-":
		return &Sub{}
	case "*":
		return &Mul{}
	case "/":
		return &Div{}
	case "%":
		return &Remainder{}
	default:
		panic("unsupported operator")
	}
}
