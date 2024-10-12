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

func codeToAsm(code string) *Program {
	tokens, _ := frontend.Tokenize(code)
	ast, _ := frontend.NewParser(tokens).ParseProgram()
	tackyAst := tacky.NewTranslator().Translate(ast)
	return NewTranslator().Translate(tackyAst)
}
