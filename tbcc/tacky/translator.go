package tacky

import (
	"github.com/thomasbollmeier/writing-a-c-compiler/tbcc/frontend"
)

type Translator struct {
	nameCreator frontend.NameCreator
}

func NewTranslator(nameCreator frontend.NameCreator) *Translator {
	return &Translator{nameCreator}
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
		case frontend.AstVarDecl:
			varDecl := item.(*frontend.VarDecl)
			if varDecl.InitValue != nil {
				val, instructions := t.translateExpr(varDecl.InitValue)
				ret = append(ret, instructions...)
				ret = append(ret, &Copy{val, &Var{varDecl.Name}})
			}
		case frontend.AstReturn:
			retStmt := item.(*frontend.ReturnStmt)
			val, instructions := t.translateExpr(retStmt.Expression)
			ret = append(ret, instructions...)
			ret = append(ret, &Return{val})
		case frontend.AstExprStmt:
			exprStmt := item.(*frontend.ExpressionStmt)
			_, instructions := t.translateExpr(exprStmt.Expression)
			ret = append(ret, instructions...)
		case frontend.AstNullStmt:
		default:
			panic("unsupported statement type")
		}
	}

	ret = append(ret, &Return{&IntConstant{0}})

	return ret
}

func (t *Translator) translateExpr(expr frontend.Expression) (Value, []Instruction) {
	switch expr.GetType() {
	case frontend.AstInteger:
		val := expr.(*frontend.IntegerLiteral).Value
		return &IntConstant{val}, nil
	case frontend.AstVariable:
		variable := expr.(*frontend.Variable)
		return &Var{variable.Name}, nil
	case frontend.AstUnary:
		unary := expr.(*frontend.UnaryExpression)
		unaryOp := t.getUnaryOp(unary.Operator)
		src, instructions := t.translateExpr(unary.Right)
		dst := &Var{t.createVarName()}
		instructions = append(instructions, &Unary{unaryOp, src, dst})
		return dst, instructions
	case frontend.AstBinary:
		binary := expr.(*frontend.BinaryExpression)
		if binary.Operator == "=" {
			return t.translateAssignment(binary)
		} else {
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
		}
	default:
		panic("unsupported expression type")
	}
}

func (t *Translator) translateAssignment(assignment *frontend.BinaryExpression) (Value, []Instruction) {
	rhsValue, instructions := t.translateExpr(assignment.Right)
	v := assignment.Left.(*frontend.Variable)
	variable := &Var{v.Name}
	instructions = append(instructions, &Copy{rhsValue, variable})
	return variable, instructions
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

func (t *Translator) createVarName() string {
	return t.nameCreator.VarName()
}

func (t *Translator) createLabelName(prefix string) string {
	return t.nameCreator.LabelName(prefix)
}
