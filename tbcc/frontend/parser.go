package frontend

import (
	"errors"
	"fmt"
	"strconv"
)

type Parser struct {
	tokens  []Token
	currIdx int
	maxIdx  int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		currIdx: 0,
		maxIdx:  len(tokens) - 1,
	}
}

func (p *Parser) ParseProgram() (*Program, error) {
	var fs []Function

	for !p.endOfInput() {
		f, err := p.parseFunction()
		if err != nil {
			return nil, err
		}
		fs = append(fs, *f)
	}

	return &Program{fs}, nil
}

func (p *Parser) parseFunction() (*Function, error) {

	_, err := p.consume(TokTypeInt) // only integer types are supported for now
	if err != nil {
		return nil, err
	}

	token, err := p.consume(TokTypeIdentifier)
	if err != nil {
		return nil, err
	}
	name := token.lexeme

	_, err = p.consume(TokTypeLeftParen)
	if err != nil {
		return nil, err
	}

	var params []Parameter
	var param *Parameter
	token, err = p.peek()
	if err != nil {
		return nil, err
	}
	if token.tokenType != TokTypeVoid {
	paramsLoop:
		for {
			param, err = p.parseParameter()
			if err != nil {
				return nil, err
			}
			params = append(params, *param)
			token, err = p.peek()
			if err != nil {
				return nil, err
			}

			switch token.tokenType {
			case TokTypeRightParen:
				break paramsLoop
			case TokTypeComma:
				_, _ = p.consume(TokTypeComma)
			default:
				return nil, errors.New("expected comma or parenthesis")
			}
		}
	} else {
		_, _ = p.consume(TokTypeVoid)
	}

	_, err = p.consume(TokTypeRightParen)
	if err != nil {
		return nil, err
	}

	token, err = p.peek()
	if err != nil {
		return nil, err
	}

	var body *BlockStmt
	if token.tokenType != TokTypeSemicolon {
		body, err = p.parseBlockStmt()
		if err != nil {
			return nil, err
		}
	} else {
		_, _ = p.consume(TokTypeSemicolon)
	}

	return &Function{
		Name:   name,
		Params: params,
		Body:   body,
	}, nil
}

func (p *Parser) parseParameter() (*Parameter, error) {
	_, err := p.consume(TokTypeInt) // only integer types are supported for now
	if err != nil {
		return nil, err
	}
	identifier, err := p.consume(TokTypeIdentifier)
	if err != nil {
		return nil, err
	}

	return &Parameter{
		Name: identifier.lexeme,
		TyId: TypeInt,
	}, nil
}

func (p *Parser) parseBlockStmt() (*BlockStmt, error) {
	_, err := p.consume(TokTypeLeftBrace)
	if err != nil {
		return nil, err
	}

	var items []BodyItem
	var item BodyItem
	var token *Token

	for {
		token, err = p.peek()
		if err != nil {
			return nil, err
		}
		if token.tokenType == TokTypeRightBrace {
			break
		}
		item, err = p.parseBodyItem()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	_, err = p.consume(TokTypeRightBrace)
	if err != nil {
		return nil, err
	}

	return &BlockStmt{items}, nil
}

func (p *Parser) parseBodyItem() (BodyItem, error) {
	tokens := p.peekN(3)

	if len(tokens) == 3 {
		if tokens[0].tokenType == TokTypeInt && tokens[1].tokenType == TokTypeIdentifier {
			if tokens[2].tokenType == TokTypeLeftParen {
				return p.parseFunction()
			} else {
				return p.parseVarDeclaration()
			}
		} else {
			return p.parseStatement()
		}
	} else {
		return p.parseStatement()
	}
}

func (p *Parser) parseVarDeclaration() (*VarDecl, error) {
	var ret *VarDecl

	_, err := p.consume(TokTypeInt)
	if err != nil {
		return nil, err
	}
	ident, err := p.consume(TokTypeIdentifier)
	if err != nil {
		return nil, err
	}

	token, err := p.peek()
	if err != nil {
		return nil, err
	}

	var initValue Expression

	switch token.tokenType {
	case TokTypeEq:
		_, _ = p.consume()
		initValue, err = p.parseExpression(0)
		if err != nil {
			return nil, err
		}
		ret = &VarDecl{
			Name:      ident.lexeme,
			InitValue: initValue,
		}
	case TokTypeSemicolon:
		ret = &VarDecl{
			Name:      ident.lexeme,
			InitValue: nil,
		}
	default:
		return nil, errors.New("unexpected token at var declaration: " + token.lexeme)
	}

	_, err = p.consume(TokTypeSemicolon)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (p *Parser) parseStatement() (Statement, error) {

	token, err := p.peek()
	if err != nil {
		return nil, err
	}

	switch token.tokenType {
	case TokTypeReturn:
		return p.parseReturnStmt()
	case TokTypeIf:
		return p.parseIfStmt()
	case TokTypeDo:
		return p.parseDoWhileStmt()
	case TokTypeWhile:
		return p.parseWhileStmt()
	case TokTypeFor:
		return p.parseForStmt()
	case TokTypeSwitch:
		return p.parseSwitchStmt()
	case TokTypeCase, TokTypeDefault:
		return p.parseCaseStmt()
	case TokTypeBreak, TokTypeContinue:
		_, _ = p.consume()
		_, err = p.consume(TokTypeSemicolon)
		if err != nil {
			return nil, err
		}
		if token.tokenType == TokTypeBreak {
			return &BreakStmt{}, nil
		} else {
			return &ContinueStmt{}, nil
		}
	case TokTypeLeftBrace:
		return p.parseBlockStmt()
	case TokTypeSemicolon:
		_, _ = p.consume()
		return &NullStmt{}, nil
	case TokTypeGoto:
		return p.parseGotoStmt()
	case TokTypeIdentifier:
		nextTokens := p.peekN(2)
		if len(nextTokens) < 2 {
			return nil, errors.New("expected a colon or semicolon")
		}
		nextNext := nextTokens[1]
		switch nextNext.tokenType {
		case TokTypeColon:
			name := token.lexeme
			_, _ = p.consume()
			_, _ = p.consume()
			return &LabelStmt{Name: name}, nil
		default:
			return p.parseExprStmt()
		}
	default:
		return p.parseExprStmt()
	}
}

func (p *Parser) parseCaseStmt() (*CaseStmt, error) {
	var value Expression = nil

	token, err := p.consume(TokTypeCase, TokTypeDefault)
	if err != nil {
		return nil, err
	}

	if token.tokenType == TokTypeCase {
		value, err = p.parseExpression(0)
		if err != nil {
			return nil, err
		}
		if value.GetType() != AstInteger {
			return nil, errors.New("expected integer as case value")
		}
	}

	_, err = p.consume(TokTypeColon)
	if err != nil {
		return nil, err
	}

	return &CaseStmt{value, "", "", ""}, nil
}

func (p *Parser) parseSwitchStmt() (*SwitchStmt, error) {
	_, err := p.consume(TokTypeSwitch)
	if err != nil {
		return nil, err
	}
	_, err = p.consume(TokTypeLeftParen)
	if err != nil {
		return nil, err
	}
	expr, err := p.parseExpression(0)
	if err != nil {
		return nil, err
	}

	_, err = p.consume(TokTypeRightParen)
	if err != nil {
		return nil, err
	}

	body, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	if body.GetType() == AstBlockStmt {
		body = p.hoistingSwitchBlockVars(body.(*BlockStmt))
	}

	return &SwitchStmt{
		Expr:           expr,
		Body:           body,
		Label:          "",
		FirstCaseLabel: "",
	}, nil
}

func (p *Parser) hoistingSwitchBlockVars(block *BlockStmt) *BlockStmt {
	var newItems []BodyItem
	var varDecls []BodyItem

	hoisting := true

	for _, item := range block.Items {
		switch item.GetType() {
		case AstVarDecl:
			if hoisting {
				varDecl := item.(*VarDecl)
				varDecls = append(varDecls, &VarDecl{Name: varDecl.Name})
				continue
			}
		case AstCaseStmt:
			hoisting = false
		default:
		}
		newItems = append(newItems, item)
	}

	newItems = append(varDecls, newItems...)

	return &BlockStmt{Items: newItems}
}

func (p *Parser) parseForStmt() (*ForStmt, error) {
	_, err := p.consume(TokTypeFor)
	if err != nil {
		return nil, err
	}
	_, err = p.consume(TokTypeLeftParen)
	if err != nil {
		return nil, err
	}

	initStmt, err := p.parseBodyItem()
	if err != nil {
		return nil, err
	}
	switch initStmt.GetType() {
	case AstVarDecl, AstExprStmt, AstNullStmt:
		break
	default:
		return nil, errors.New("init statement must be one of: varDecl, exprStmt or nullStmt")
	}

	var condition Expression = nil
	token, err := p.peek()
	if err != nil {
		return nil, err
	}
	if token.tokenType != TokTypeSemicolon {
		condition, err = p.parseExpression(0)
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(TokTypeSemicolon)
	if err != nil {
		return nil, err
	}

	var post Expression = nil
	token, err = p.peek()
	if err != nil {
		return nil, err
	}
	if token.tokenType != TokTypeRightParen {
		post, err = p.parseExpression(0)
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(TokTypeRightParen)
	if err != nil {
		return nil, err
	}

	body, err := p.parseStatement()

	if err != nil {
		return nil, err
	}

	return &ForStmt{
		initStmt,
		condition,
		post,
		body,
		"",
	}, nil
}

func (p *Parser) parseWhileStmt() (*WhileStmt, error) {
	_, err := p.consume(TokTypeWhile)
	if err != nil {
		return nil, err
	}
	_, err = p.consume(TokTypeLeftParen)
	if err != nil {
		return nil, err
	}
	condition, err := p.parseExpression(0)
	if err != nil {
		return nil, err
	}
	_, err = p.consume(TokTypeRightParen)
	if err != nil {
		return nil, err
	}
	body, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	return &WhileStmt{
		condition,
		body,
		"",
	}, nil
}

func (p *Parser) parseDoWhileStmt() (*DoWhileStmt, error) {
	_, err := p.consume(TokTypeDo)
	if err != nil {
		return nil, err
	}
	body, err := p.parseStatement()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(TokTypeWhile)
	if err != nil {
		return nil, err
	}
	_, err = p.consume(TokTypeLeftParen)
	if err != nil {
		return nil, err
	}
	condition, err := p.parseExpression(0)
	if err != nil {
		return nil, err
	}
	_, err = p.consume(TokTypeRightParen)
	if err != nil {
		return nil, err
	}
	_, err = p.consume(TokTypeSemicolon)
	if err != nil {
		return nil, err
	}
	return &DoWhileStmt{
		condition,
		body,
		"",
	}, nil
}

func (p *Parser) parseGotoStmt() (*GotoStmt, error) {
	_, err := p.consume(TokTypeGoto)
	if err != nil {
		return nil, err
	}
	target, err := p.consume(TokTypeIdentifier)
	if err != nil {
		return nil, err
	}
	_, err = p.consume(TokTypeSemicolon)
	if err != nil {
		return nil, err
	}
	return &GotoStmt{target.lexeme}, nil
}

func (p *Parser) parseIfStmt() (*IfStmt, error) {
	_, _ = p.consume(TokTypeIf)
	_, err := p.consume(TokTypeLeftParen)
	if err != nil {
		return nil, err
	}
	condition, err := p.parseExpression(0)
	if err != nil {
		return nil, err
	}
	_, err = p.consume(TokTypeRightParen)
	if err != nil {
		return nil, err
	}
	consequent, err := p.parseStatement()
	if err != nil {
		return nil, err
	}
	var alternate Statement
	nextToken, err := p.peek()
	if err == nil && nextToken.tokenType == TokTypeElse {
		_, _ = p.consume(TokTypeElse)
		alternate, err = p.parseStatement()
		if err != nil {
			return nil, err
		}
	}

	return &IfStmt{condition, consequent, alternate}, nil
}

func (p *Parser) parseExprStmt() (*ExpressionStmt, error) {
	expr, err := p.parseExpression(0)
	if err != nil {
		return nil, err
	}
	_, err = p.consume(TokTypeSemicolon)
	if err != nil {
		return nil, err
	}
	return &ExpressionStmt{expr}, nil
}

func (p *Parser) parseReturnStmt() (*ReturnStmt, error) {
	_, err := p.consume(TokTypeReturn)
	if err != nil {
		return nil, err
	}
	expr, err := p.parseExpression(0)
	if err != nil {
		return nil, err
	}
	_, err = p.consume(TokTypeSemicolon)
	if err != nil {
		return nil, err
	}
	return &ReturnStmt{expr}, nil
}

func (p *Parser) parseExpression(minPrecedence int) (Expression, error) {
	ret, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	for {

		token, err := p.peek()
		if err != nil {
			return ret, nil
		}

		prefInfo, ok := binOpPreference[token.tokenType]
		if !ok {
			return ret, nil
		}
		if prefInfo.Level < minPrecedence {
			return ret, nil
		}

		var nextPref int
		if prefInfo.Assoc == AssocLeft {
			nextPref = prefInfo.Level + 1
		} else {
			nextPref = prefInfo.Level
		}

		if token.tokenType == TokTypeQuestionMark {
			ret, err = p.parseConditional(ret, nextPref)
			if err != nil {
				return nil, err
			}
			continue
		}

		binOpToken := token
		_, _ = p.consume()

		right, err := p.parseExpression(nextPref)
		if err != nil {
			return nil, err
		}

		switch binOpToken.lexeme {
		case "+=", "-=", "*=", "/=", "%=",
			"&=", "|=", "^=", "<<=", ">>=":
			// Compound assignment => expand it:
			op := binOpToken.lexeme[0:1]
			ret = &BinaryExpression{
				"=",
				ret,
				&BinaryExpression{
					op,
					ret,
					right,
				},
			}
		default:
			ret = &BinaryExpression{
				binOpToken.lexeme,
				ret,
				right,
			}
		}
	}
}

func (p *Parser) parseConditional(condition Expression, minPrecedence int) (Expression, error) {
	_, _ = p.consume(TokTypeQuestionMark)
	consequent, err := p.parseExpression(0)
	if err != nil {
		return nil, err
	}
	_, _ = p.consume(TokTypeColon)
	alternate, err := p.parseExpression(minPrecedence)
	if err != nil {
		return nil, err
	}
	return &Conditional{
		condition,
		consequent,
		alternate,
	}, nil
}

func (p *Parser) parseFactor() (Expression, error) {

	var ret Expression = nil

	token, err := p.peek()
	if err != nil {
		return nil, err
	}

	switch token.tokenType {
	case TokTypeIntConstant:
		intLiteral, _ := p.consume()
		value, err := strconv.ParseInt(intLiteral.lexeme, 10, 64)
		if err != nil {
			return nil, err
		}
		ret = &IntegerLiteral{int(value)}
	case TokTypeIdentifier:
		ident, _ := p.consume()
		nextToken, err := p.peek()
		if err == nil && nextToken.tokenType == TokTypeLeftParen {
			args, err := p.parseArguments()
			if err != nil {
				return nil, err
			}
			ret = &FunctionCall{
				Callee: ident.lexeme,
				Args:   args,
			}
		} else {
			ret = &Variable{ident.lexeme}
		}
	case TokTypeMinus, TokTypeTilde, TokTypeExclMark:
		_, _ = p.consume()
		operator := token.lexeme
		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		ret = &UnaryExpression{operator, right}
	case TokTypePlusPlus, TokTypeMinusMinus:
		_, _ = p.consume()
		var operator string
		if token.tokenType == TokTypePlusPlus {
			operator = "+"
		} else {
			operator = "-"
		}
		factor, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		lvalue, ok := factor.(*Variable)
		if !ok {
			return nil, errors.New("unexpected token after increment or decrement")
		}

		ret = &BinaryExpression{
			"=",
			lvalue,
			&BinaryExpression{
				operator,
				lvalue,
				&IntegerLiteral{1},
			},
		}
	case TokTypeLeftParen:
		_, _ = p.consume()
		expr, err := p.parseExpression(0)
		if err != nil {
			return nil, err
		}
		_, err = p.consume(TokTypeRightParen)
		if err != nil {
			return nil, err
		}
		ret = expr
	default:
		return nil, errors.New("unexpected token: " + token.lexeme)
	}

	if lvalue, ok := ret.(*Variable); ok {
		nextToken, err := p.peek()
		if err == nil &&
			(nextToken.tokenType == TokTypePlusPlus || nextToken.tokenType == TokTypeMinusMinus) {
			_, _ = p.consume()
			ret = &PostfixIncDec{
				nextToken.lexeme,
				*lvalue,
			}
		}
	}

	return ret, nil
}

func (p *Parser) parseArguments() ([]Expression, error) {
	var args []Expression
	var arg Expression
	var token *Token
	var err error

	_, err = p.consume(TokTypeLeftParen)
	if err != nil {
		return nil, err
	}

	token, err = p.peek()
	if err != nil {
		return nil, err
	}
	if token.tokenType == TokTypeRightParen {
		_, _ = p.consume()
		return args, nil
	}

argsLoop:
	for {
		arg, err = p.parseExpression(0)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)

		token, err = p.peek()
		if err != nil {
			return nil, err
		}
		switch token.tokenType {
		case TokTypeRightParen:
			_, _ = p.consume()
			break argsLoop
		case TokTypeComma:
			_, _ = p.consume()
		default:
			return nil, errors.New("unexpected token: " + token.lexeme)
		}
	}

	return args, nil
}

func (p *Parser) consume(expected ...TokenType) (*Token, error) {
	if p.currIdx > p.maxIdx {
		return nil, errors.New("no more tokens")
	}
	if len(expected) == 0 {
		ret := p.tokens[p.currIdx]
		p.currIdx++
		return &ret, nil
	} else {
		ret := p.tokens[p.currIdx]
		for _, expType := range expected {
			if ret.tokenType == expType {
				p.currIdx++
				return &ret, nil
			}
		}
		message := fmt.Sprintf("token '%s' has unexpected token type", ret.lexeme)
		return nil, errors.New(message)
	}
}

func (p *Parser) peek() (*Token, error) {
	if p.currIdx > p.maxIdx {
		return nil, errors.New("no more tokens")
	}
	return &p.tokens[p.currIdx], nil
}

func (p *Parser) peekN(n int) []Token {
	if p.currIdx > p.maxIdx {
		return nil
	}
	lastIdx := p.currIdx + n - 1
	lastIdx = min(lastIdx, p.maxIdx)
	return p.tokens[p.currIdx : lastIdx+1]
}

func (p *Parser) endOfInput() bool {
	return p.currIdx > p.maxIdx
}
