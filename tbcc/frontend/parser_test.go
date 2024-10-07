package frontend

import (
	"testing"
)

func createParser() *Parser {
	var pos Position
	tokens := []Token{
		*NewToken(TokTypeReturn, "return", pos),
		*NewToken(TokTypeIntConstant, "0", pos),
		*NewToken(TokTypeSemicolon, ";", pos),
	}
	return NewParser(tokens)
}

func TestNewParser(t *testing.T) {

	got := createParser()
	if got == nil {
		t.Errorf("NewParser() returned nil")
		return
	}

	if got.currIdx != 0 {
		t.Errorf("NewParser().currIdx = %d, want 0", got.currIdx)
	}

	if got.maxIdx != 2 {
		t.Errorf("NewParser().maxIdx = %d, want 2", got.maxIdx)
	}

}

func TestParser_ParseProgram(t *testing.T) {
	code := `
int main(void) {
	return 42;
}`
	runParserWithCode(t, code, false)
}

func TestParser_ParseProgramFail(t *testing.T) {
	code := `
int main(void) {
	return 42;
}
foo`
	runParserWithCode(t, code, true)
}

func runParserWithCode(t *testing.T, code string, expectError bool) {
	tokens, err := Tokenize(code)
	if err != nil {
		t.Errorf("Tokenize() error = %v", err)
	}
	parser := NewParser(tokens)

	program, err := parser.ParseProgram()

	if !expectError {
		if err != nil {
			t.Errorf("ParseProgram() error = %v", err)
		}

		if program.GetType() != AstProgram {
			t.Errorf("program.GetType() = %v, want %v", program.GetType(), AstProgram)
		}

		program.Accept(NewAstPrinter(4))

	} else {
		if err == nil {
			t.Errorf("ParseProgram() should have returned an error")
		}
	}
}

func TestParser_consume(t *testing.T) {

	p := createParser()

	_, err := p.consume(TokTypeReturn)
	if err != nil {
		t.Errorf("Parser.consume(TokTypeReturn) = %v, want nil", err)
	}

	_, err = p.consume(TokTypeIntConstant)
	if err != nil {
		t.Errorf("Parser.consume(TokTypeIntConstant) = %v, want error", err)
	}

	_, err = p.consume(TokTypeSemicolon)
	if err != nil {
		t.Errorf("Parser.consume(TokTypeSemicolon) = %v, want error", err)
	}

	token, err := p.consume()
	if err == nil {
		t.Errorf("Parser.consume() = %v, want error", token)
	}
}
