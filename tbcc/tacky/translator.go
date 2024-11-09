package tacky

import (
	"fmt"
	"github.com/thomasbollmeier/writing-a-c-compiler/tbcc/frontend"
)

type Translator struct {
	nameCreator  frontend.NameCreator
	switchValues []Value
}

func NewTranslator(nameCreator frontend.NameCreator) *Translator {
	return &Translator{nameCreator, make([]Value, 0)}
}

func (t *Translator) Translate(program *frontend.Program) *Program {
	return &Program{t.translateFunction(&program.Func)}
}

func (t *Translator) translateFunction(f *frontend.Function) Function {

	bodyInstructions := t.translateBlock(&f.Body)
	bodyInstructions = append(bodyInstructions, &Return{&IntConstant{0}})

	return Function{
		f.Name,
		bodyInstructions,
	}
}

func (t *Translator) translateBlock(b *frontend.BlockStmt) []Instruction {
	var ret []Instruction

	for _, item := range b.Items {
		ret = append(ret, t.translateBodyItem(item)...)
	}

	return ret
}

func (t *Translator) translateBodyItem(item frontend.BodyItem) []Instruction {
	switch item.GetType() {
	case frontend.AstVarDecl:
		var ret []Instruction
		varDecl := item.(*frontend.VarDecl)
		if varDecl.InitValue != nil {
			val, instructions := t.translateExpr(varDecl.InitValue)
			ret = append(ret, instructions...)
			ret = append(ret, &Copy{val, &Var{varDecl.Name}})
		}
		return ret
	default:
		return t.translateStatement(item)
	}
}

func (t *Translator) translateStatement(stmt frontend.Statement) []Instruction {
	var ret []Instruction
	var val Value

	switch stmt.GetType() {
	case frontend.AstReturn:
		retStmt := stmt.(*frontend.ReturnStmt)
		val, ret = t.translateExpr(retStmt.Expression)
		ret = append(ret, &Return{val})
	case frontend.AstExprStmt:
		exprStmt := stmt.(*frontend.ExpressionStmt)
		_, ret = t.translateExpr(exprStmt.Expression)
	case frontend.AstIfStmt:
		ret = t.translateIfStmt(stmt.(*frontend.IfStmt))
	case frontend.AstBlockStmt:
		ret = t.translateBlock(stmt.(*frontend.BlockStmt))
	case frontend.AstGotoStmt:
		gotoStmt := stmt.(*frontend.GotoStmt)
		ret = []Instruction{&Jump{gotoStmt.Target}}
	case frontend.AstLabelStmt:
		labelStmt := stmt.(*frontend.LabelStmt)
		ret = []Instruction{&Label{labelStmt.Name}}
	case frontend.AstDoWhileStmt:
		ret = t.translateDoWhileStmt(stmt.(*frontend.DoWhileStmt))
	case frontend.AstWhileStmt:
		ret = t.translateWhileStmt(stmt.(*frontend.WhileStmt))
	case frontend.AstForStmt:
		ret = t.translateForStmt(stmt.(*frontend.ForStmt))
	case frontend.AstSwitchStmt:
		ret = t.translateSwitchStmt(stmt.(*frontend.SwitchStmt))
	case frontend.AstCaseStmt:
		ret = t.translateCaseStmt(stmt.(*frontend.CaseStmt))
	case frontend.AstContinueStmt:
		continueStmt := stmt.(*frontend.ContinueStmt)
		ret = []Instruction{&Jump{t.loopLabelContinue(continueStmt.Label)}}
	case frontend.AstBreakStmt:
		breakStmt := stmt.(*frontend.BreakStmt)
		ret = []Instruction{&Jump{t.loopLabelBreak(breakStmt.Label)}}
	case frontend.AstNullStmt:
	default:
		panic("unsupported statement type")
	}

	return ret
}

func (t *Translator) translateCaseStmt(stmt *frontend.CaseStmt) []Instruction {
	var ret []Instruction

	if stmt.PrevCaseLabel != "" {
		ret = append(ret, &Jump{stmt.Label + ".fallthrough"})
	}

	ret = append(ret, &Label{stmt.Label})

	if stmt.Value != nil {
		caseVal, caseValInstructions := t.translateExpr(stmt.Value)
		ret = append(ret, caseValInstructions...)

		switchNestingLevel := len(t.switchValues)
		selectVar := t.switchValues[switchNestingLevel-1]
		caseVar := &Var{t.createVarName()}
		resultVar := &Var{t.createVarName()}
		ret = append(ret,
			&Copy{caseVal, caseVar},
			&Binary{
				&Sub{},
				caseVar,
				selectVar,
				resultVar,
			},
			&JumpIfNotZero{resultVar, stmt.NextCaseLabel},
		)

	}

	if stmt.PrevCaseLabel != "" {
		ret = append(ret, &Label{stmt.Label + ".fallthrough"})
	}

	return ret
}

func (t *Translator) translateSwitchStmt(stmt *frontend.SwitchStmt) []Instruction {
	var ret []Instruction

	breakLabel := t.loopLabelBreak(stmt.Label)

	selectVal, selInstructions := t.translateExpr(stmt.Expr)
	ret = append(ret, selInstructions...)

	if stmt.FirstCaseLabel == "" {
		return ret
	}

	selectVar := &Var{t.createVarName()}

	// Push current selection var to stack to make it available
	// for case statements
	t.switchValues = append(t.switchValues, selectVar)

	ret = append(ret, &Copy{selectVal, selectVar})
	switch stmt.Body.GetType() {
	case frontend.AstBlockStmt:
		{
			var bodyInstructions []Instruction
			endOfVarDecls := false
			for _, item := range stmt.Body.(*frontend.BlockStmt).Items {
				switch item.GetType() {
				case frontend.AstVarDecl:
				default:
					if !endOfVarDecls {
						bodyInstructions = append(bodyInstructions,
							&Jump{stmt.FirstCaseLabel})
					}
					endOfVarDecls = true
				}
				bodyInstructions = append(bodyInstructions, t.translateBodyItem(item)...)
			}
			ret = append(ret, bodyInstructions...)
		}
	default:
		ret = append(ret, t.translateStatement(stmt.Body)...)
	}
	ret = append(ret, &Label{breakLabel})

	// Pop current selection var from stack
	size := len(t.switchValues)
	t.switchValues = t.switchValues[0 : size-1]

	return ret
}

func (t *Translator) translateForStmt(stmt *frontend.ForStmt) []Instruction {

	var ret []Instruction
	startLabel := t.loopLabelStart(stmt.Label)
	continueLabel := t.loopLabelContinue(stmt.Label)
	breakLabel := t.loopLabelBreak(stmt.Label)

	ret = append(ret, t.translateBodyItem(stmt.InitStmt)...)
	ret = append(ret, &Label{startLabel})

	if stmt.Condition != nil {
		condVal, condInstructions := t.translateExpr(stmt.Condition)
		condResult := &Var{t.createVarName()}
		ret = append(ret, condInstructions...)
		ret = append(ret,
			&Copy{condVal, condResult},
			&JumpIfZero{condResult, breakLabel},
		)
	}
	ret = append(ret, t.translateStatement(stmt.Body)...)
	ret = append(ret, &Label{continueLabel})

	if stmt.Post != nil {
		_, postInstructions := t.translateExpr(stmt.Post)
		ret = append(ret, postInstructions...)
	}

	ret = append(ret,
		&Jump{startLabel},
		&Label{breakLabel},
	)

	return ret
}

func (t *Translator) translateWhileStmt(stmt *frontend.WhileStmt) []Instruction {
	continueLabel := t.loopLabelContinue(stmt.Label)
	breakLabel := t.loopLabelBreak(stmt.Label)

	ret := []Instruction{&Label{continueLabel}}
	condVal, condInstructions := t.translateExpr(stmt.Condition)
	condResult := &Var{t.createVarName()}
	ret = append(ret, condInstructions...)
	ret = append(ret,
		&Copy{condVal, condResult},
		&JumpIfZero{condResult, breakLabel},
	)
	ret = append(ret, t.translateStatement(stmt.Body)...)
	ret = append(ret,
		&Jump{continueLabel},
		&Label{breakLabel},
	)

	return ret
}

func (t *Translator) translateDoWhileStmt(stmt *frontend.DoWhileStmt) []Instruction {
	startLabel := t.loopLabelStart(stmt.Label)
	ret := []Instruction{&Label{startLabel}}
	ret = append(ret, t.translateStatement(stmt.Body)...)
	ret = append(ret, &Label{t.loopLabelContinue(stmt.Label)})
	condVal, condInstructions := t.translateExpr(stmt.Condition)
	condResult := &Var{t.createVarName()}
	ret = append(ret, condInstructions...)
	ret = append(ret,
		&Copy{condVal, condResult},
		&JumpIfNotZero{condResult, startLabel},
		&Label{t.loopLabelBreak(stmt.Label)},
	)

	return ret
}

func (t *Translator) loopLabelStart(loopLabel string) string {
	return fmt.Sprintf("%s.start", loopLabel)
}

func (t *Translator) loopLabelEnd(loopLabel string) string {
	return fmt.Sprintf("%s.end", loopLabel)
}

func (t *Translator) loopLabelBreak(loopLabel string) string {
	return fmt.Sprintf("%s.break", loopLabel)
}

func (t *Translator) loopLabelContinue(loopLabel string) string {
	return fmt.Sprintf("%s.continue", loopLabel)
}

func (t *Translator) translateIfStmt(ifStmt *frontend.IfStmt) []Instruction {
	var ret []Instruction
	var condValue Value

	if ifStmt.Alternate == nil {
		condValue, ret = t.translateExpr(ifStmt.Condition)
		endLabelName := t.createLabelName("end")
		ret = append(ret, &JumpIfZero{condValue, endLabelName})
		ret = append(ret, t.translateStatement(ifStmt.Consequent)...)
		ret = append(ret, &Label{endLabelName})
	} else {
		condValue, ret = t.translateExpr(ifStmt.Condition)
		endLabelName := t.createLabelName("end")
		elseLabelName := t.createLabelName("else")
		ret = append(ret, &JumpIfZero{condValue, elseLabelName})
		ret = append(ret, t.translateStatement(ifStmt.Consequent)...)
		ret = append(ret, &Jump{endLabelName})
		ret = append(ret, &Label{elseLabelName})
		ret = append(ret, t.translateStatement(ifStmt.Alternate)...)
		ret = append(ret, &Label{endLabelName})
	}

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
	case frontend.AstPostfixIncDec:
		return t.translatePostfixIncDec(expr.(*frontend.PostfixIncDec))
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
	case frontend.AstConditional:
		conditional := expr.(*frontend.Conditional)
		return t.translateConditional(conditional)
	default:
		panic("unsupported expression type")
	}
}

func (t *Translator) translateConditional(conditional *frontend.Conditional) (Value, []Instruction) {
	resultValue := &Var{t.createVarName()}

	condValue, instructions := t.translateExpr(conditional.Condition)
	endLabelName := t.createLabelName("end")
	elseLabelName := t.createLabelName("else")

	instructions = append(instructions, &JumpIfZero{condValue, elseLabelName})
	consValue, consInstructions := t.translateExpr(conditional.Consequent)
	instructions = append(instructions, consInstructions...)
	instructions = append(instructions, &Copy{consValue, resultValue})
	instructions = append(instructions, &Jump{endLabelName})
	instructions = append(instructions, &Label{elseLabelName})
	altValue, altInstructions := t.translateExpr(conditional.Alternate)
	instructions = append(instructions, altInstructions...)
	instructions = append(instructions, &Copy{altValue, resultValue})
	instructions = append(instructions, &Label{endLabelName})

	return resultValue, instructions
}

func (t *Translator) translatePostfixIncDec(postfixIncDec *frontend.PostfixIncDec) (Value, []Instruction) {
	resultValue := &Var{t.createVarName()}
	value := &Var{postfixIncDec.Operand.Name}

	var binOp BinaryOp
	if postfixIncDec.Operator == "++" {
		binOp = &Add{}
	} else {
		binOp = &Sub{}
	}

	instructions := []Instruction{&Copy{value, resultValue}}
	instructions = append(instructions, &Binary{
		binOp,
		value,
		&IntConstant{1},
		value,
	})

	return resultValue, instructions
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
