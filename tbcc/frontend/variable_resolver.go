package frontend

import (
	"errors"
	"fmt"
)

type varResolverResult struct {
	ast AST
	err error
}

type variableResolver struct {
	nameCreator NameCreator
	env         *environment
	labelMap    map[string]string
	result      varResolverResult
}

func newVariableResolver(nameCreator NameCreator) *variableResolver {
	return &variableResolver{
		nameCreator: nameCreator,
		env:         newEnvironment(nil),
		labelMap:    make(map[string]string),
		result:      varResolverResult{nil, nil},
	}
}

func (vr *variableResolver) resolve(program *Program) (*Program, error) {
	ast, err := vr.evalAst(program)
	if err != nil {
		return nil, err
	}
	return ast.(*Program), nil
}

func (vr *variableResolver) VisitProgram(p *Program) {
	var newFunctions []Function

	for _, fun := range p.Functions {
		ast, err := vr.evalAst(&fun)
		if err != nil {
			return
		}
		newFunctions = append(newFunctions, *ast.(*Function))
	}

	vr.setResult(&Program{newFunctions}, nil)
}

func (vr *variableResolver) VisitFunction(f *Function) {
	var newBody *BlockStmt
	var newParams []Parameter

	if f.Body != nil {

		vr.env = newEnvironment(vr.env)

		for _, param := range f.Params {
			uniqueName := vr.nameCreator.VarName()
			vr.env.set(param.Name, uniqueName)
			newParams = append(newParams, Parameter{
				Name: uniqueName,
				TyId: param.TyId,
			})
		}

		ast, err := vr.evalAst(f.Body)
		if err != nil {
			vr.env = vr.env.getParent()
			return
		}
		newBody = ast.(*BlockStmt)

		vr.env = vr.env.getParent()

	} else {
		newParams = f.Params
		newBody = nil
	}

	vr.setResult(&Function{
		Name:   f.Name,
		Params: newParams,
		Body:   newBody,
	}, nil)
}

func (vr *variableResolver) VisitVarDecl(v *VarDecl) {
	if vr.env.isSet(v.Name) {
		vr.setResult(nil, errors.New(fmt.Sprintf("variable %s already defined", v.Name)))
		return
	}
	uniqueName := vr.nameCreator.VarName()
	vr.env.set(v.Name, uniqueName)

	var newInitValue AST
	var err error

	if v.InitValue != nil {
		newInitValue, err = vr.evalAst(v.InitValue)
		if err != nil {
			return
		}
	} else {
		newInitValue = nil
	}

	vr.setResult(&VarDecl{uniqueName, newInitValue}, nil)
}

func (vr *variableResolver) VisitReturn(r *ReturnStmt) {
	newExpr, err := vr.evalAst(r.Expression)
	if err != nil {
		return
	}
	vr.setResult(&ReturnStmt{newExpr}, nil)
}

func (vr *variableResolver) VisitExprStmt(e *ExpressionStmt) {
	newExpr, err := vr.evalAst(e.Expression)
	if err != nil {
		return
	}
	vr.setResult(&ExpressionStmt{newExpr}, nil)
}

func (vr *variableResolver) VisitIfStmt(i *IfStmt) {
	newCondition, err := vr.evalAst(i.Condition)
	if err != nil {
		return
	}
	newConsequent, err := vr.evalAst(i.Consequent)
	if err != nil {
		return
	}
	var newAlternate Statement = nil
	if i.Alternate != nil {
		newAlternate, err = vr.evalAst(i.Alternate)
		if err != nil {
			return
		}
	}
	vr.setResult(&IfStmt{
		newCondition,
		newConsequent,
		newAlternate,
	}, nil)
}

func (vr *variableResolver) VisitBlockStmt(b *BlockStmt) {
	defer func() {
		vr.env = vr.env.getParent()
	}()

	var newItems []BodyItem

	vr.env = newEnvironment(vr.env)

	for _, item := range b.Items {
		newItem, err := vr.evalAst(item)
		if err != nil {
			return
		}
		newItems = append(newItems, newItem)
	}

	vr.setResult(&BlockStmt{newItems}, nil)
}

func (vr *variableResolver) VisitGotoStmt(g *GotoStmt) {
	uniqueTarget, ok := vr.labelMap[g.Target]
	if !ok {
		uniqueTarget = vr.nameCreator.LabelName(g.Target)
		vr.labelMap[g.Target] = uniqueTarget
	}

	vr.setResult(&GotoStmt{uniqueTarget}, nil)
}

func (vr *variableResolver) VisitLabelStmt(l *LabelStmt) {
	uniqueName, ok := vr.labelMap[l.Name]
	if !ok {
		uniqueName = vr.nameCreator.LabelName(l.Name)
		vr.labelMap[l.Name] = uniqueName
	}
	vr.setResult(&LabelStmt{uniqueName}, nil)
}

func (vr *variableResolver) VisitDoWhileStmt(d *DoWhileStmt) {
	newCondition, err := vr.evalAst(d.Condition)
	if err != nil {
		return
	}
	newBody, err := vr.evalAst(d.Body)
	if err != nil {
		return
	}
	vr.setResult(&DoWhileStmt{
		Condition: newCondition,
		Body:      newBody,
		Label:     d.Label,
	}, nil)
}

func (vr *variableResolver) VisitWhileStmt(w *WhileStmt) {
	newCondition, err := vr.evalAst(w.Condition)
	if err != nil {
		return
	}
	newBody, err := vr.evalAst(w.Body)
	if err != nil {
		return
	}
	vr.setResult(&WhileStmt{
		Condition: newCondition,
		Body:      newBody,
		Label:     w.Label,
	}, nil)
}

func (vr *variableResolver) VisitForStmt(f *ForStmt) {

	var newCondition Expression
	var newPost Expression

	vr.env = newEnvironment(vr.env)

	defer func() {
		vr.env = vr.env.getParent()
	}()

	newInitStmt, err := vr.evalAst(f.InitStmt)
	if err != nil {
		return
	}

	if f.Condition != nil {
		newCondition, err = vr.evalAst(f.Condition)
		if err != nil {
			return
		}
	}

	if f.Post != nil {
		newPost, err = vr.evalAst(f.Post)
		if err != nil {
			return
		}
	}

	newBody, err := vr.evalAst(f.Body)
	if err != nil {
		return
	}
	vr.setResult(&ForStmt{
		InitStmt:  newInitStmt,
		Condition: newCondition,
		Post:      newPost,
		Body:      newBody,
		Label:     f.Label,
	}, nil)

}

func (vr *variableResolver) VisitBreakStmt(b *BreakStmt) {
	vr.setResult(b, nil)
}

func (vr *variableResolver) VisitContinueStmt(c *ContinueStmt) {
	vr.setResult(c, nil)
}

func (vr *variableResolver) VisitSwitchStmt(s *SwitchStmt) {
	newExpr, err := vr.evalAst(s.Expr)
	if err != nil {
		return
	}

	var newBody Statement
	newBody, err = vr.evalAst(s.Body)
	if err != nil {
		return
	}

	vr.setResult(&SwitchStmt{
		Expr:           newExpr,
		Body:           newBody,
		Label:          s.Label,
		FirstCaseLabel: s.FirstCaseLabel,
	}, nil)
}

func (vr *variableResolver) VisitCaseStmt(c *CaseStmt) {
	vr.setResult(c, nil)
}

func (vr *variableResolver) VisitNullStmt() {
	vr.setResult(&NullStmt{}, nil)
}

func (vr *variableResolver) VisitInteger(i *IntegerLiteral) {
	vr.setResult(i, nil)
}

func (vr *variableResolver) VisitVariable(v *Variable) {
	uniqueName, err := vr.env.lookup(v.Name)
	if err != nil {
		vr.setResult(nil, err)
		return
	}
	vr.setResult(&Variable{uniqueName}, nil)
}

func (vr *variableResolver) VisitFunctionCall(f *FunctionCall) {
	var newArgs []Expression

	for _, arg := range f.Args {
		newArg, err := vr.evalAst(arg)
		if err != nil {
			return
		}
		newArgs = append(newArgs, newArg)
	}

	vr.setResult(&FunctionCall{f.Callee, newArgs}, nil)
}

func (vr *variableResolver) VisitUnary(u *UnaryExpression) {
	newRight, err := vr.evalAst(u.Right)
	if err != nil {
		return
	}
	vr.setResult(&UnaryExpression{
		Operator: u.Operator,
		Right:    newRight,
	}, nil)
}

func (vr *variableResolver) VisitPostfixIncDec(p *PostfixIncDec) {
	newOperand, err := vr.evalAst(&p.Operand)
	if err != nil {
		return
	}
	vr.setResult(&PostfixIncDec{
		Operator: p.Operator,
		Operand:  *newOperand.(*Variable),
	}, nil)
}

func (vr *variableResolver) VisitBinary(b *BinaryExpression) {
	var newLeft Expression
	var newRight Expression
	var err error

	newLeft, err = vr.evalAst(b.Left)
	if err != nil {
		return
	}

	newRight, err = vr.evalAst(b.Right)
	if err != nil {
		return
	}

	// For assignment check if left expression is LVALUE
	if b.Operator == "=" && newLeft.GetType() != AstVariable {
		vr.setResult(nil, errors.New("invalid lvalue"))
		return
	}

	vr.setResult(&BinaryExpression{
		Operator: b.Operator,
		Left:     newLeft,
		Right:    newRight,
	}, nil)
}

func (vr *variableResolver) VisitConditional(cond *Conditional) {
	newCond, err := vr.evalAst(cond.Condition)
	if err != nil {
		return
	}
	newConsequent, err := vr.evalAst(cond.Consequent)
	if err != nil {
		return
	}
	newAlternate, err := vr.evalAst(cond.Alternate)
	if err != nil {
		return
	}
	vr.setResult(&Conditional{
		newCond,
		newConsequent,
		newAlternate,
	}, nil)
}

func (vr *variableResolver) evalAst(ast AST) (AST, error) {
	ast.Accept(vr)
	return vr.result.ast, vr.result.err
}

func (vr *variableResolver) setResult(ast AST, err error) {
	vr.result.ast = ast
	vr.result.err = err
}
