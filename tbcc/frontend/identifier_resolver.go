package frontend

import (
	"errors"
	"fmt"
)

type idResolverResult struct {
	ast AST
	err error
}

type identifierResolver struct {
	nameCreator     NameCreator
	envs            *Environments
	labelMap        map[string]string
	functionNesting int
	result          idResolverResult
}

func newIdentifierResolver(nameCreator NameCreator) *identifierResolver {
	return &identifierResolver{
		nameCreator:     nameCreator,
		envs:            NewEnvironments(),
		labelMap:        make(map[string]string),
		functionNesting: 0,
		result:          idResolverResult{nil, nil},
	}
}

func (ir *identifierResolver) resolve(program *Program) (*Program, error) {
	ast, err := ir.evalAst(program)
	if err != nil {
		return nil, err
	}
	return ast.(*Program), nil
}

func (ir *identifierResolver) VisitProgram(p *Program) {
	var newDeclarations []Declaration

	for _, decl := range p.Declarations {
		ast, err := ir.evalAst(decl)
		if err != nil {
			return
		}
		newDeclarations = append(newDeclarations, ast)
	}

	ir.setResult(&Program{newDeclarations}, nil)
}

func (ir *identifierResolver) VisitFunction(f *Function) {
	var newBody *BlockStmt
	var newParams []Parameter

	if !allParamsUnique(f.Params) {
		ir.setResult(nil,
			errors.New(fmt.Sprintf("parameters of function %s must be unique", f.Name)))
		return
	}

	if f.Body != nil && ir.functionNesting > 0 {
		ir.setResult(nil,
			errors.New(fmt.Sprintf("function %s must not be defined within another function", f.Name)))
		return
	}

	entry, env := ir.envs.Get(f.Name)
	if env != nil {
		if ir.envs.Block == env && entry.category != idCatFunction {
			ir.setResult(nil, errors.New(fmt.Sprintf("%s is already defined", f.Name)))
			return
		}
	}
	ir.envs.set(f.Name, f.Name, linkExternal, idCatFunction, nil)

	if f.Body != nil {
		ir.envs.beginBlock()
		ir.functionNesting++

		for _, param := range f.Params {
			uniqueName := ir.nameCreator.VarName()
			ir.envs.set(param.Name, uniqueName, linkNone, idCatParameter, &IntInfo{})
			newParams = append(newParams, Parameter{
				Name: uniqueName,
				TyId: param.TyId,
			})
		}

		ast, err := ir.evalAst(f.Body)
		if err != nil {
			ir.envs.endBlock()
			return
		}
		newBody = ast.(*BlockStmt)

		ir.functionNesting--
		ir.envs.endBlock()

	} else {
		newParams = f.Params
		newBody = nil
	}

	ir.setResult(&Function{
		Name:   f.Name,
		Params: newParams,
		Body:   newBody,
	}, nil)
}

func allParamsUnique(params []Parameter) bool {
	paramSet := make(map[string]bool)
	for _, param := range params {
		_, ok := paramSet[param.Name]
		if ok {
			return false
		}
		paramSet[param.Name] = true
	}
	return true
}

func (ir *identifierResolver) VisitVarDecl(v *VarDecl) {

	entry, definingEnv := ir.envs.Get(v.Name)
	alreadyDefined := false
	if definingEnv != nil {
		if definingEnv == ir.envs.Block {
			alreadyDefined = true
		} else if definingEnv == ir.envs.Block.getParent() && entry.category == idCatParameter {
			alreadyDefined = true
		}
	}
	if alreadyDefined {
		ir.setResult(nil, errors.New(fmt.Sprintf("variable %s already defined", v.Name)))
		return
	}

	uniqueName := ir.nameCreator.VarName()
	ir.envs.set(v.Name, uniqueName, linkNone, idCatVariable, &IntInfo{})

	var newInitValue AST
	var err error

	if v.InitValue != nil {
		newInitValue, err = ir.evalAst(v.InitValue)
		if err != nil {
			return
		}
	} else {
		newInitValue = nil
	}

	ir.setResult(&VarDecl{uniqueName, newInitValue, v.StorageClass}, nil)
}

func (ir *identifierResolver) VisitReturn(r *ReturnStmt) {
	newExpr, err := ir.evalAst(r.Expression)
	if err != nil {
		return
	}
	ir.setResult(&ReturnStmt{newExpr}, nil)
}

func (ir *identifierResolver) VisitExprStmt(e *ExpressionStmt) {
	newExpr, err := ir.evalAst(e.Expression)
	if err != nil {
		return
	}
	ir.setResult(&ExpressionStmt{newExpr}, nil)
}

func (ir *identifierResolver) VisitIfStmt(i *IfStmt) {
	newCondition, err := ir.evalAst(i.Condition)
	if err != nil {
		return
	}
	newConsequent, err := ir.evalAst(i.Consequent)
	if err != nil {
		return
	}
	var newAlternate Statement = nil
	if i.Alternate != nil {
		newAlternate, err = ir.evalAst(i.Alternate)
		if err != nil {
			return
		}
	}
	ir.setResult(&IfStmt{
		newCondition,
		newConsequent,
		newAlternate,
	}, nil)
}

func (ir *identifierResolver) VisitBlockStmt(b *BlockStmt) {
	defer func() {
		ir.envs.endBlock()
	}()

	var newItems []BodyItem

	ir.envs.beginBlock()

	for _, item := range b.Items {
		newItem, err := ir.evalAst(item)
		if err != nil {
			return
		}
		newItems = append(newItems, newItem)
	}

	ir.setResult(&BlockStmt{newItems}, nil)
}

func (ir *identifierResolver) VisitGotoStmt(g *GotoStmt) {
	uniqueTarget, ok := ir.labelMap[g.Target]
	if !ok {
		uniqueTarget = ir.nameCreator.LabelName(g.Target)
		ir.labelMap[g.Target] = uniqueTarget
	}

	ir.setResult(&GotoStmt{uniqueTarget}, nil)
}

func (ir *identifierResolver) VisitLabelStmt(l *LabelStmt) {
	uniqueName, ok := ir.labelMap[l.Name]
	if !ok {
		uniqueName = ir.nameCreator.LabelName(l.Name)
		ir.labelMap[l.Name] = uniqueName
	}
	ir.setResult(&LabelStmt{uniqueName}, nil)
}

func (ir *identifierResolver) VisitDoWhileStmt(d *DoWhileStmt) {
	newCondition, err := ir.evalAst(d.Condition)
	if err != nil {
		return
	}
	newBody, err := ir.evalAst(d.Body)
	if err != nil {
		return
	}
	ir.setResult(&DoWhileStmt{
		Condition: newCondition,
		Body:      newBody,
		Label:     d.Label,
	}, nil)
}

func (ir *identifierResolver) VisitWhileStmt(w *WhileStmt) {
	newCondition, err := ir.evalAst(w.Condition)
	if err != nil {
		return
	}
	newBody, err := ir.evalAst(w.Body)
	if err != nil {
		return
	}
	ir.setResult(&WhileStmt{
		Condition: newCondition,
		Body:      newBody,
		Label:     w.Label,
	}, nil)
}

func (ir *identifierResolver) VisitForStmt(f *ForStmt) {

	var newCondition Expression
	var newPost Expression

	ir.envs.beginBlock()

	defer func() {
		ir.envs.endBlock()
	}()

	newInitStmt, err := ir.evalAst(f.InitStmt)
	if err != nil {
		return
	}

	if f.Condition != nil {
		newCondition, err = ir.evalAst(f.Condition)
		if err != nil {
			return
		}
	}

	if f.Post != nil {
		newPost, err = ir.evalAst(f.Post)
		if err != nil {
			return
		}
	}

	newBody, err := ir.evalAst(f.Body)
	if err != nil {
		return
	}
	ir.setResult(&ForStmt{
		InitStmt:  newInitStmt,
		Condition: newCondition,
		Post:      newPost,
		Body:      newBody,
		Label:     f.Label,
	}, nil)

}

func (ir *identifierResolver) VisitBreakStmt(b *BreakStmt) {
	ir.setResult(b, nil)
}

func (ir *identifierResolver) VisitContinueStmt(c *ContinueStmt) {
	ir.setResult(c, nil)
}

func (ir *identifierResolver) VisitSwitchStmt(s *SwitchStmt) {
	newExpr, err := ir.evalAst(s.Expr)
	if err != nil {
		return
	}

	var newBody Statement
	newBody, err = ir.evalAst(s.Body)
	if err != nil {
		return
	}

	ir.setResult(&SwitchStmt{
		Expr:           newExpr,
		Body:           newBody,
		Label:          s.Label,
		FirstCaseLabel: s.FirstCaseLabel,
	}, nil)
}

func (ir *identifierResolver) VisitCaseStmt(c *CaseStmt) {
	ir.setResult(c, nil)
}

func (ir *identifierResolver) VisitNullStmt() {
	ir.setResult(&NullStmt{}, nil)
}

func (ir *identifierResolver) VisitInteger(i *IntegerLiteral) {
	ir.setResult(i, nil)
}

func (ir *identifierResolver) VisitVariable(v *Variable) {
	uniqueName, err := ir.envs.Lookup(v.Name)
	if err != nil {
		ir.setResult(nil, err)
		return
	}
	ir.setResult(&Variable{uniqueName}, nil)
}

func (ir *identifierResolver) VisitFunctionCall(f *FunctionCall) {
	var newArgs []Expression

	entry, definingEnv := ir.envs.Get(f.Callee)
	if definingEnv == nil || entry.category != idCatFunction {
		ir.setResult(nil,
			errors.New(fmt.Sprintf("%s is not a function", f.Callee)))
		return
	}

	for _, arg := range f.Args {
		newArg, err := ir.evalAst(arg)
		if err != nil {
			return
		}
		newArgs = append(newArgs, newArg)
	}

	ir.setResult(&FunctionCall{f.Callee, newArgs}, nil)
}

func (ir *identifierResolver) VisitUnary(u *UnaryExpression) {
	newRight, err := ir.evalAst(u.Right)
	if err != nil {
		return
	}
	ir.setResult(&UnaryExpression{
		Operator: u.Operator,
		Right:    newRight,
	}, nil)
}

func (ir *identifierResolver) VisitPostfixIncDec(p *PostfixIncDec) {
	newOperand, err := ir.evalAst(&p.Operand)
	if err != nil {
		return
	}
	ir.setResult(&PostfixIncDec{
		Operator: p.Operator,
		Operand:  *newOperand.(*Variable),
	}, nil)
}

func (ir *identifierResolver) VisitBinary(b *BinaryExpression) {
	var newLeft Expression
	var newRight Expression
	var err error

	newLeft, err = ir.evalAst(b.Left)
	if err != nil {
		return
	}

	newRight, err = ir.evalAst(b.Right)
	if err != nil {
		return
	}

	// For assignment check if left expression is LVALUE
	if b.Operator == "=" && newLeft.GetType() != AstVariable {
		ir.setResult(nil, errors.New("invalid lvalue"))
		return
	}

	ir.setResult(&BinaryExpression{
		Operator: b.Operator,
		Left:     newLeft,
		Right:    newRight,
	}, nil)
}

func (ir *identifierResolver) VisitConditional(cond *Conditional) {
	newCond, err := ir.evalAst(cond.Condition)
	if err != nil {
		return
	}
	newConsequent, err := ir.evalAst(cond.Consequent)
	if err != nil {
		return
	}
	newAlternate, err := ir.evalAst(cond.Alternate)
	if err != nil {
		return
	}
	ir.setResult(&Conditional{
		newCond,
		newConsequent,
		newAlternate,
	}, nil)
}

func (ir *identifierResolver) evalAst(ast AST) (AST, error) {
	ast.Accept(ir)
	return ir.result.ast, ir.result.err
}

func (ir *identifierResolver) setResult(ast AST, err error) {
	ir.result.ast = ast
	ir.result.err = err
}
