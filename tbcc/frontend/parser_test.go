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

func TestParser_ParseBlock(t *testing.T) {
	code := `int main(void) {
    	int a = 42;
		{	
			int b = a + 1;
			int a = 41;
			return b;
		}
    	return a;
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

func TestParser_ParseForStmt(t *testing.T) {
	code := `int main(void) {
		int sum = 0;
		int i = 42;
		int counter;
		for (int i = 0; i <= 10; i = i + 1) {
			counter = i;
			if (i % 2 == 0)
				continue;
			sum = sum + 1;
		}
	
		return sum == 5 && counter == 10;
	}`

	runParserWithCode(t, code, false)
}

func TestParser_ParseForLoopWithoutCond(t *testing.T) {
	code := `int main(void) {
		int sum = 0;
		int i = 0;
		int counter;
		for (;;) {
			i++;
			sum = sum + i;
			if (i == 5)
				break;
		}
	
		return sum == 15;
	}`

	runParserWithCode(t, code, false)
}

func TestParser_ParseMultipleCont(t *testing.T) {
	code := `int main(void) {
		int x = 10;
		int y = 0;
		int z = 0;
		do {
			z = z + 1;
			if (x <= 0)
				continue;
			x = x - 1;
			if (y >= 10)
				continue;
			y = y + 1;
		} while (z != 50);
		return z == 50 && x == 0 && y == 10;
	}`

	runParserWithCode(t, code, false)
}

func TestParser_ParseSwitchStatement(t *testing.T) {
	code := `int main(void) {
		int count = 37;
		int iterations = (count + 4) / 5;
		switch (count % 5) {
			case 0:
				do {
					count = count - 1;
					case 4:
						count = count - 1;
					case 3:
						count = count - 1;
					case 2:
						count = count - 1;
					case 1:
						count = count - 1;
				} while ((iterations = iterations - 1) > 0);
		}
		return (count == 0 && iterations == 0);
	}`

	runParserWithCode(t, code, false)
}

func TestParser_ParseSwitchWithContinue(t *testing.T) {
	code := `int main(void) {
		int sum = 0;
		for (int i = 0; i < 10; i = i + 1) {
			switch(i % 2) {
				case 0: continue;
				default: sum = sum + 1;
			}
		}
		return sum;
	}`

	runParserWithCode(t, code, false)
}

func TestParser_ParseSwitchWithDecl(t *testing.T) {
	code := `int main(void) {
		int a = 3;
		int b = 0;
		switch(a) {
			int a = (b = 5);
			while (1);
			int answer = 42;
		case 3:
			a = 4;
			b = b + a;
		}
	
		return a == 3 && b == 4;
	}`

	runParserWithCode(t, code, false)
}

func TestParser_ParseSwitchNestedCase(t *testing.T) {
	code := `int main(void) {
		int answer = 0;
		switch(3) {
		case 0: return 0;
		case 1:
			if (0) {
				case 3: answer = 42; break;
			}
		}
	
		return answer == 42;
	}`

	runParserWithCode(t, code, false)
}

func TestParser_ParseProgramWithFunctions(t *testing.T) {
	code := `
	int init(void);

	int add(int a, int b) {
		int init(void); 
		init();
		return a + b;
	}

	int main(void) {
		int a = add(41, 1);
		return a;
	}`

	runParserWithCode(t, code, false)
}

func TestParser_NestedFuncDef(t *testing.T) {
	code := `int main(void) {
		int foo(void) {
    	    return 1;
		}
    	return foo();
	}`

	runParserWithCode(t, code, true)
}

func TestParser_Recursion(t *testing.T) {
	code := `int fib(int n) {
		if (n == 0 || n == 1) {
			return n;
		} else {
			return fib(n - 1) + fib(n - 2);
		}
	}
	
	int main(void) {
		int n = 6;
		return fib(n);
	}`

	runParserWithCode(t, code, false)
}

func TestParser_FunctionDecl(t *testing.T) {
	code := `int main(void) {
		int foo = 3;
		int bar = 4;
		if (foo + bar > 0) {
			int foo(void);
			bar = foo();
		}
		return foo + bar;
	}
	
	int foo(void) {
		return 8;
	}`

	runParserWithCode(t, code, false)
}

func TestParser_ParameterShadowsFunction(t *testing.T) {
	code := `int a(void) {
			return 1;
		}
		
		int b(int a) {
			return a;
		}
		
		int main(void) {
			return a() + b(2);
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
		program, _, err = AnalyzeSemantics(program, NewNameCreator())
	}

	if !expectError {
		if err != nil {
			t.Errorf("ParseProgram() error = %v", err)
			return
		}

		if program.GetType() != AstProgram {
			t.Errorf("program.GetType() = %v, want %v", program.GetType(), AstProgram)
			return
		}

		program.Accept(NewAstPrinter(4))

	} else {
		if err == nil {
			t.Errorf("ParseProgram() should have returned an error")
		}
	}
}
