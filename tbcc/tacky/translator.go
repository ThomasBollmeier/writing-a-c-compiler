package tacky

import (
	"fmt"
	"github.com/thomasbollmeier/writing-a-c-compiler/tbcc/frontend"
)

type Translator struct {
	nextCounter   uint
	labelCounters map[string]uint
}

func NewTranslator() *Translator {
	return &Translator{
		0,
		make(map[string]uint),
	}
}

func (t *Translator) Translate(program *frontend.Program) *Program {
	return &Program{t.translateFunction(&program.Func)}
}

func (t *Translator) translateFunction(f *frontend.Function) Function {
	return Function{
		f.Name,
		t.translateBody(f.Body),
	}
}

func (t *Translator) translateBody(body []frontend.BodyItem) []Instruction {
	var ret []Instruction

	for _, item := range body {
		switch item.GetType() {
		case frontend.AstReturn:
			retStmt := item.(*frontend.ReturnStmt)
			val, instructions := t.translateExpr(retStmt.Expression)
			ret = append(ret, instructions...)
			ret = append(ret, &Return{val})
		default:
			panic("unsupported statement type")
		}
	}

	return ret
}

func (t *Translator) translateExpr(expr frontend.Expression) (Value, []Instruction) {
	switch expr.GetType() {
	case frontend.AstInteger:
		val := expr.(*frontend.IntegerLiteral).Value
		return &IntConstant{val}, nil
	case frontend.AstUnary:
		unary := expr.(*frontend.UnaryExpression)
		unaryOp := t.getUnaryOp(unary.Operator)
		src, instructions := t.translateExpr(unary.Right)
		dst := &Var{t.createVarName()}
		instructions = append(instructions, &Unary{unaryOp, src, dst})
		return dst, instructions
	case frontend.AstBinary:
		binary := expr.(*frontend.BinaryExpression)
		binaryOp := t.getBinaryOp(binary.Operator)
		binaryOpType := binaryOp.GetType()
		if binaryOpType != TacAnd && binaryOpType != TacOr {
			src1, instructions := t.translateExpr(binary.Left)
			src2, instructions2 := t.translateExpr(binary.Right)
			instructions = append(instructions, instructions2...)
			dst := &Var{t.createVarName()}
			instructions = append(instructions, &Binary{binaryOp, src1, src2, dst})
			return dst, instructions
		} else {
			return t.translateExprWithShortCircuit(binaryOp, binary.Left, binary.Right)
		}
	default:
		panic("unsupported expression type")
	}
}

func (t *Translator) translateExprWithShortCircuit(
	op BinaryOp,
	left, right frontend.Expression) (Value, []Instruction) {

	var instructions []Instruction

	varResult := &Var{t.createVarName()}
	valLeft, instructionsLeft := t.translateExpr(left)
	varLeft := &Var{t.createVarName()}
	valRight, instructionsRight := t.translateExpr(right)
	varRight := &Var{t.createVarName()}
	labelEnd := t.createLabelName("end")
	labelFalse := t.createLabelName("false")
	labelTrue := t.createLabelName("true")

	switch op.GetType() {
	case TacAnd:
		instructions = append(
			instructionsLeft,
			&Copy{valLeft, varLeft},
			&JumpIfZero{varLeft, labelFalse})
		instructions = append(instructions,
			instructionsRight...)
		instructions = append(instructions,
			&Copy{valRight, varRight},
			&JumpIfZero{varRight, labelFalse},
			&Copy{&IntConstant{1}, varResult},
			&Jump{labelEnd},
			&Label{labelFalse},
			&Copy{&IntConstant{0}, varResult},
			&Label{labelEnd},
		)
	case TacOr:
		instructions = append(
			instructionsLeft,
			&Copy{valLeft, varLeft},
			&JumpIfNotZero{varLeft, labelTrue})
		instructions = append(instructions,
			instructionsRight...)
		instructions = append(instructions,
			&Copy{valRight, varRight},
			&JumpIfNotZero{varRight, labelTrue},
			&Copy{&IntConstant{0}, varResult},
			&Jump{labelEnd},
			&Label{labelTrue},
			&Copy{&IntConstant{1}, varResult},
			&Label{labelEnd},
		)
	default:
		panic("unsupported logical operator")
	}

	return varResult, instructions
}

func (t *Translator) createVarName() string {
	varName := fmt.Sprintf("tmp.%d", t.nextCounter)
	t.nextCounter++
	return varName
}

func (t *Translator) createLabelName(prefix string) string {
	current, ok := t.labelCounters[prefix]
	if !ok {
		current = 0
	}
	ret := fmt.Sprintf("%s%d", prefix, current)
	current++
	t.labelCounters[prefix] = current
	return ret
}

func (t *Translator) getUnaryOp(op string) UnaryOp {
	switch op {
	case "-":
		return &Negate{}
	case "~":
		return &Complement{}
	case "!":
		return &Not{}
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
	case "&":
		return &BitAnd{}
	case "|":
		return &BitOr{}
	case "^":
		return &BitXor{}
	case "<<":
		return &BitShiftLeft{}
	case ">>":
		return &BitShiftRight{}
	case "==":
		return &Equal{}
	case "!=":
		return &NotEqual{}
	case ">":
		return &Greater{}
	case ">=":
		return &GreaterEq{}
	case "<":
		return &Less{}
	case "<=":
		return &LessEq{}
	case "&&":
		return &And{}
	case "||":
		return &Or{}
	default:
		panic("unsupported operator: " + op)
	}
}
