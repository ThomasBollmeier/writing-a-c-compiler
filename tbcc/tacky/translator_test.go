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

func translate(code string) *Program {
	nameCreator := frontend.NewNameCreator()
	tokens, _ := frontend.Tokenize(code)
	parser := frontend.NewParser(tokens)
	ast, _ := parser.ParseProgram()
	ast, _ = frontend.AnalyzeSemantics(ast, nameCreator)
	translator := NewTranslator(nameCreator)
	return translator.Translate(ast)
}
