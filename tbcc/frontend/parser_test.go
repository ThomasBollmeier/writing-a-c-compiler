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

func TestParser_ParseBodyItems(t *testing.T) {
	code := `
int main(void) {
	int answer = 42;
	40 + 2;
	int a;
	int b;
	int c = b = a = 7 * 6;
	return answer;
}`
	runParserWithCode(t, code, false)
}

func TestParser_ParseUnary(t *testing.T) {
	code := `
int main(void) {
	return ~(-42);
}`
	runParserWithCode(t, code, false)
}

func TestParser_ParseBinary(t *testing.T) {
	code := `
int main(void) {
	return 1 + 2 * 3 - 4;
}`
	runParserWithCode(t, code, false)
}

func TestParser_ParseShiftLeft(t *testing.T) {
	code := `
int main(void) {
	return 21 << 1;
}`
	runParserWithCode(t, code, false)
}

func TestParser_ParseBitwiseOr(t *testing.T) {
	code := `
int main(void) {
	return 1 | 2;
}`
	runParserWithCode(t, code, false)
}

func TestParser_ParseBitwisePreference(t *testing.T) {
	code := `int main(void) {
		return 80 >> 2 | 1 ^ 5 & 7 << 1;
	}`
	runParserWithCode(t, code, false)
}

func TestParser_ParseBitwiseShiftPreference(t *testing.T) {
	code := `int main(void) {
		return 40 << 4 + 12 >> 1;
	}`
	runParserWithCode(t, code, false)
}

func TestParser_ParseCompoundAssignment(t *testing.T) {
	code := `int main(void) {
		int a = 21;
		a *= 1 + 1;
		return a;
	}`
	runParserWithCode(t, code, false)
}

func TestParser_ParseCompoundAssignmentBinary(t *testing.T) {
	code := `int main(void) {
		int a = 21;
		a &= 1 + 1;
		return a;
	}`
	runParserWithCode(t, code, false)
}

func TestParser_ParsePrefixIncrement(t *testing.T) {
	code := `int main(void) {
		int a = 41;
		return ++a;
	}`
	runParserWithCode(t, code, false)
}

func TestParser_ParseIncrements(t *testing.T) {
	code := `int main(void) {
		int a = 0;
		int b = 0;
		a++;
		++a;
		++a;
		b--;
		--b;
		return (a == 3 && b == -2);
	}`
	runParserWithCode(t, code, false)
}

func TestParser_ParsePostfixIncrement(t *testing.T) {
	code := `int main(void) {
		int a = 42;
		return a++;
	}`
	runParserWithCode(t, code, false)
}

func TestParser_ParseIfStatement(t *testing.T) {
	code := `int main(void) {
		int ok = 1;
		if (ok)
			return 42;
	}`

	runParserWithCode(t, code, false)
}

func TestParser_ParseIfElseStatement(t *testing.T) {
	code := `int main(void) {
		int ok = 1;
		if (ok)
			return 42;
		else
			return 23;
	}`

	runParserWithCode(t, code, false)
}

func TestParser_ParseConditional(t *testing.T) {
	code := `int main(void) {
    	int a = 2;
    	int b = 1;
    	a > b ? a = 1 : a;
    	return a;
	}`

	runParserWithCode(t, code, false)
}

func TestParser_ParseNestedTernary(t *testing.T) {
	code := `int main(void) {
 	   	int a = 1;
		int b = 2;
    	int flag = 0;
    	
		return a > b ? 5 : flag ? 6 : 7;
	}`

	runParserWithCode(t, code, false)
}

func TestParser_ParseGoto(t *testing.T) {
	code := `int main(void) {
 	   	int a = 42;
		goto end;

		a = 23;
		
	end:
		return a;
	}`

	runParserWithCode(t, code, false)
}

func TestParser_ParseLabelMultiple(t *testing.T) {
	code := `int main(void) {
 	   	int a = 42;
		goto end;
	end:
		a = 23;
		
	end:
		return a;
	}`

	runParserWithCode(t, code, true)
}

func TestParser_ParseLabelAtDecl(t *testing.T) {
	code := `int main(void) {
 	forbidden_before_decl:
 	   	int a = 42;
		goto end;
		a = 23;
		
	end:
		return a;
	}`

	runParserWithCode(t, code, true)
}

func TestParser_ParseProgramFail(t *testing.T) {
	code := `
int main(void) {
	return 42;
}
foo`
	runParserWithCode(t, code, true)
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

func runParserWithCode(t *testing.T, code string, expectError bool) {
	tokens, err := Tokenize(code)
	if err != nil {
		t.Errorf("Tokenize() error = %v", err)
	}
	parser := NewParser(tokens)

	program, err := parser.ParseProgram()
	if err == nil {
		program, err = AnalyzeSemantics(program, NewNameCreator())
	}

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
