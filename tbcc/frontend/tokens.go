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
	TokTypePlus
	TokTypeMinus
	TokTypeAsterisk
	TokTypeSlash
	TokTypePercent
	TokTypeMinusMinus
	TokTypeAmpersand
	TokTypePipe
	TokTypeCaret
	TokTypeLessLess
	TokTypeGreaterGreater
	TokTypeExclMark
	TokTypeAmperAmper
	TokTypePipePipe
	TokTypeEqEq
	TokTypeExclMarkEq
	TokTypeGt
	TokTypeGtEq
	TokTypeLt
	TokTypeLtEq
	TokTypeEq
	TokTypePlusEq
	TokTypeMinusEq
	TokTypeAsteriskEq
	TokTypeSlashEq
	TokTypePercentEq
)

var tokenTypeToRegexStr = map[TokenType]string{
	TokTypeIdentifier:     "[a-zA-Z_]\\w*\\b",
	TokTypeIntConstant:    "[0-9]+\\b",
	TokTypeLeftParen:      "\\(",
	TokTypeRightParen:     "\\)",
	TokTypeLeftBrace:      "{",
	TokTypeRightBrace:     "}",
	TokTypeSemicolon:      ";",
	TokTypeTilde:          "~",
	TokTypePlus:           "\\+",
	TokTypeMinus:          "-",
	TokTypeAsterisk:       "\\*",
	TokTypeSlash:          "/",
	TokTypePercent:        "%",
	TokTypeMinusMinus:     "--",
	TokTypeAmpersand:      "&",
	TokTypePipe:           "\\|",
	TokTypeCaret:          "\\^",
	TokTypeLessLess:       "<<",
	TokTypeGreaterGreater: ">>",
	TokTypeExclMark:       "!",
	TokTypeAmperAmper:     "&&",
	TokTypePipePipe:       "\\|\\|",
	TokTypeEqEq:           "==",
	TokTypeExclMarkEq:     "!=",
	TokTypeGt:             ">",
	TokTypeGtEq:           ">=",
	TokTypeLt:             "<",
	TokTypeLtEq:           "<=",
	TokTypeEq:             "=",
	TokTypePlusEq:         "\\+=",
	TokTypeMinusEq:        "-=",
	TokTypeAsteriskEq:     "\\*=",
	TokTypeSlashEq:        "/=",
	TokTypePercentEq:      "%=",
}

var strToKeyword = map[string]TokenType{
	"int":    TokTypeInt,
	"void":   TokTypeVoid,
	"return": TokTypeReturn,
}

type Associativity int

const (
	AssocLeft Associativity = iota
	AssocRight
)

type PrefInfo struct {
	Level int
	Assoc Associativity
}

var binOpPreference = map[TokenType]PrefInfo{
	TokTypeAsterisk:       {50, AssocLeft},
	TokTypeSlash:          {50, AssocLeft},
	TokTypePercent:        {50, AssocLeft},
	TokTypePlus:           {45, AssocLeft},
	TokTypeMinus:          {45, AssocLeft},
	TokTypeLessLess:       {40, AssocLeft},
	TokTypeGreaterGreater: {40, AssocLeft},
	TokTypeLt:             {35, AssocLeft},
	TokTypeLtEq:           {35, AssocLeft},
	TokTypeGt:             {35, AssocLeft},
	TokTypeGtEq:           {35, AssocLeft},
	TokTypeEqEq:           {30, AssocLeft},
	TokTypeExclMarkEq:     {30, AssocLeft},
	TokTypeAmpersand:      {26, AssocLeft},
	TokTypeCaret:          {25, AssocLeft},
	TokTypePipe:           {20, AssocLeft},
	TokTypeAmperAmper:     {10, AssocLeft},
	TokTypePipePipe:       {5, AssocLeft},
	TokTypeEq:             {1, AssocRight},
	TokTypePlusEq:         {1, AssocRight},
	TokTypeMinusEq:        {1, AssocRight},
	TokTypeAsteriskEq:     {1, AssocRight},
	TokTypeSlashEq:        {1, AssocRight},
	TokTypePercentEq:      {1, AssocRight},
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
