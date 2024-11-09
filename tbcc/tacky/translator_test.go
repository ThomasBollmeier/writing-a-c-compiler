package tacky

import (
	"fmt"
	"github.com/thomasbollmeier/writing-a-c-compiler/tbcc/frontend"
	"testing"
)

func TestTranslator_Translate(t *testing.T) {
	code := `
int main(void) {
	return ~(-42);
}`
	program := translate(code)
	fmt.Println(program)

}

func TestTranslator_Translate_ShiftLeft(t *testing.T) {
	code := `
int main(void) {
	return 21 << 1;
}`
	program := translate(code)
	fmt.Println(program)
}

func TestTranslator_Translate_VarDecl(t *testing.T) {
	code := `
int main(void) {
	int a = 42;
	int b;
	41 + 1;
	b = 22 + 1;
	return b;
}`
	program := translate(code)
	fmt.Println(program)
}

func TestTranslator_Translate_PostfixInc(t *testing.T) {
	code := `
int main(void) {
	int a = 42;
	return a++;
}`
	program := translate(code)
	fmt.Println(program)
}

func TestTranslator_Translate_Switch(t *testing.T) {
	code := `int main(void) {
		int a = 4;
		int b = 9;
		int c = 0;
		switch (a ? b : 7) {
			case 9:
				c = 2;
			case 1:
				c = c + 4;
		}
		return c;
	}`

	program := translate(code)
	fmt.Println(program)
}

func TestTranslator_TranslateSwitchWithNestedCase(t *testing.T) {
	code := `int main(void) {
		int answer = 42;
		switch(1) {
			int i = 1;
			if (0) {
				answer = 0;	
				case 1: 
				i = 2; break;
			}
		}
	
		return answer == 42;
	}`

	program := translate(code)

	program.Accept(NewAstPrinter(2))
}

func translate(code string) *Program {
	nameCreator := frontend.NewNameCreator()
	tokens, _ := frontend.Tokenize(code)
	parser := frontend.NewParser(tokens)
	ast, _ := parser.ParseProgram()
	ast, _ = frontend.AnalyzeSemantics(ast, nameCreator)
	translator := NewTranslator(nameCreator)
	return translator.Translate(ast)
}
