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
	body, err := p.parseStatement()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(TokTypeRightBrace)
	if err != nil {
		return nil, err
	}

	return &Function{name, body}, nil
}

func (p *Parser) parseStatement() (Statement, error) {
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

		binOpPref, ok := binOpPreference[token.tokenType]
		if !ok || binOpPref < minPrecedence {
			return ret, nil
		}

		binOpToken := token
		_, _ = p.consume()

		right, err := p.parseExpression(binOpPref + 1)
		if err != nil {
			return nil, err
		}

		ret = &BinaryExpression{
			binOpToken.lexeme,
			ret,
			right,
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
	case TokTypeMinus, TokTypeTilde:
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
