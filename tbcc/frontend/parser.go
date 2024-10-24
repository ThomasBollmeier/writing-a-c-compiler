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
	f, err := p.parseFunction()
	if err != nil {
		return nil, err
	}
	if !p.endOfInput() {
		return nil, errors.New("expected end of input")
	}
	return &Program{*f}, nil
}

func (p *Parser) parseFunction() (*Function, error) {
	_, err := p.consume(TokTypeInt)
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
	_, err = p.consume(TokTypeVoid)
	if err != nil {
		return nil, err
	}
	_, err = p.consume(TokTypeRightParen)
	if err != nil {
		return nil, err
	}
	_, err = p.consume(TokTypeLeftBrace)
	if err != nil {
		return nil, err
	}

	var body []BodyItem
	var bodyItem BodyItem

	for {
		token, err = p.peek()
		if err != nil {
			return nil, err
		}
		if token.tokenType == TokTypeRightBrace {
			break
		}
		bodyItem, err = p.parseBodyItem()
		if err != nil {
			return nil, err
		}
		body = append(body, bodyItem)
	}

	_, err = p.consume(TokTypeRightBrace)
	if err != nil {
		return nil, err
	}

	return &Function{name, body}, nil
}

func (p *Parser) parseBodyItem() (BodyItem, error) {
	token, err := p.peek()
	if err != nil {
		return nil, err
	}
	if token.tokenType == TokTypeInt {
		return p.parseVarDeclaration()
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
	case TokTypeSemicolon:
		_, _ = p.consume()
		return &NullStmt{}, nil
	default:
		return p.parseExprStmt()
	}
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

		binOpToken := token
		_, _ = p.consume()

		right, err := p.parseExpression(nextPref)
		if err != nil {
			return nil, err
		}

		switch binOpToken.lexeme {
		case "+=", "-=", "*=", "/=", "%=":
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

func (p *Parser) parseFactor() (Expression, error) {
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
		return &IntegerLiteral{int(value)}, nil
	case TokTypeIdentifier:
		ident, _ := p.consume()
		return &Variable{Name: ident.lexeme}, nil
	case TokTypeMinus, TokTypeTilde, TokTypeExclMark:
		_, _ = p.consume()
		operator := token.lexeme
		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		return &UnaryExpression{operator, right}, nil
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
		return expr, nil
	default:
		return nil, errors.New("unexpected token: " + token.lexeme)
	}
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
