package backend

import (
	"fmt"
	"github.com/thomasbollmeier/writing-a-c-compiler/tbcc/frontend"
	"github.com/thomasbollmeier/writing-a-c-compiler/tbcc/tacky"
	"testing"
)

func TestCodeGenerator_GenerateCode(t *testing.T) {
	code := `
int main(void) {
	return ~(-42);
}`
	asmProgram := codeToAsm(code)
	asm := NewCodeGenerator().GenerateCode(*asmProgram)

	fmt.Print(asm)
}

func TestCodeGenerator_GenerateCode_Associativity(t *testing.T) {
	code := `
int main(void) {
    return (3 / 2 * 4) + (5 - 4 + 3);
}`
	asmProgram := codeToAsm(code)
	asm := NewCodeGenerator().GenerateCode(*asmProgram)

	fmt.Print(asm)
}

func TestCodeGenerator_GenerateCode_BitAnd(t *testing.T) {
	code := `
int main(void) {
    return 3 & 5;
}`
	asmProgram := codeToAsm(code)
	asm := NewCodeGenerator().GenerateCode(*asmProgram)

	fmt.Print(asm)
}

func TestCodeGenerator_GenerateCode_AndFalse(t *testing.T) {
	code := `int main(void) {
    	return (10 && 0) + (0 && 4) + (0 && 0);
	}`

	asmProgram := codeToAsm(code)
	asm := NewCodeGenerator().GenerateCode(*asmProgram)

	fmt.Print(asm)
}

func codeToAsm(code string) *Program {
	tokens, _ := frontend.Tokenize(code)
	ast, _ := frontend.NewParser(tokens).ParseProgram()
	nameCreator := frontend.NewNameCreator()
	tackyAst := tacky.NewTranslator(nameCreator).Translate(ast)
	return NewTranslator().Translate(tackyAst)
}
