package frontend

import (
	"errors"
	"regexp"
	"unicode"
)

func Tokenize(code string) ([]Token, error) {
	pos := Position{1, 1}
	remaining := code
	tokens := make([]Token, 0)

	for remaining != "" {
		remaining, pos = skipWhitespace(remaining, pos)
		if remaining == "" {
			break
		}
		tokenPtr, err := maxMunch(remaining, pos)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, *tokenPtr)
		pos = updatePosition(tokenPtr, pos)
		remaining = remaining[len(tokenPtr.lexeme):]
	}

	return tokens, nil
}

func updatePosition(tokenPtr *Token, pos Position) Position {
	newPos := pos
	for _, ch := range tokenPtr.lexeme {
		newPos = newPos.Advance(ch)
	}
	return newPos
}

func maxMunch(code string, pos Position) (*Token, error) {
	maxTokenType := TokTypeUnknown
	maxLen := 0
	maxLexeme := ""

	for tokenType, regexStr := range tokenTypeToRegexStr {
		regexStr = "^" + regexStr // match must be at the beginning
		regex := regexp.MustCompile(regexStr)
		match := regex.Find([]byte(code))
		if match == nil {
			continue
		}
		lexeme := string(match)
		if len(lexeme) > maxLen {
			maxTokenType = tokenType
			maxLen = len(lexeme)
			maxLexeme = lexeme
		}
	}

	if maxTokenType != TokTypeUnknown {
		tokenType := adaptTokenType(maxTokenType, maxLexeme)
		return NewToken(tokenType, maxLexeme, pos), nil
	} else {
		return nil, errors.New("code matches no token")
	}
}

func adaptTokenType(tokenType TokenType, lexeme string) TokenType {
	if tokenType != TokTypeIdentifier {
		return tokenType
	}
	keywordType, ok := strToKeyword[lexeme]
	if !ok {
		return tokenType
	} else {
		return keywordType
	}
}

func skipWhitespace(code string, startPos Position) (string, Position) {
	pos := startPos
	for i, ch := range code {
		if !unicode.IsSpace(ch) {
			return code[i:], pos
		}
		pos = pos.Advance(ch)
	}
	return "", pos
}
