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
)

var tokenTypeToRegexStr = map[TokenType]string{
	TokTypeIdentifier:  "[a-zA-Z_]\\w*\\b",
	TokTypeIntConstant: "[0-9]+\\b",
	TokTypeLeftParen:   "\\(",
	TokTypeRightParen:  "\\)",
	TokTypeLeftBrace:   "{",
	TokTypeRightBrace:  "}",
	TokTypeSemicolon:   ";",
}

var strToKeyword = map[string]TokenType{
	"int":    TokTypeInt,
	"void":   TokTypeVoid,
	"return": TokTypeReturn,
}

type Position struct {
	Line, Col int
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
