package frontend

type TokenType int

const (
	TokTypeUnknown TokenType = iota
	TokTypeIdentifier
	TokTypeIntConstant
	TokTypeInt
	TokTypeVoid
	TokTypeReturn
	TokTypeLeftParen
	TokTypeRightParen
	TokTypeLeftBrace
	TokTypeRightBrace
	TokTypeSemicolon
	TokTypeTilde
	TokTypeMinus
	TokTypeMinusMinus
)

var tokenTypeToRegexStr = map[TokenType]string{
	TokTypeIdentifier:  "[a-zA-Z_]\\w*\\b",
	TokTypeIntConstant: "[0-9]+\\b",
	TokTypeLeftParen:   "\\(",
	TokTypeRightParen:  "\\)",
	TokTypeLeftBrace:   "{",
	TokTypeRightBrace:  "}",
	TokTypeSemicolon:   ";",
	TokTypeTilde:       "~",
	TokTypeMinus:       "-",
	TokTypeMinusMinus:  "--",
}

var strToKeyword = map[string]TokenType{
	"int":    TokTypeInt,
	"void":   TokTypeVoid,
	"return": TokTypeReturn,
}

type Position struct {
	Line, Col int
}

func (p Position) Advance(ch rune) Position {
	if string(ch) != "\n" {
		return Position{p.Line, p.Col + 1}
	} else {
		return Position{p.Line + 1, 1}
	}
}

type Token struct {
	tokenType TokenType
	lexeme    string
	position  Position
}

func NewToken(tokenType TokenType, lexeme string, position Position) *Token {
	return &Token{
		tokenType: tokenType,
		lexeme:    lexeme,
		position:  position,
	}
}
