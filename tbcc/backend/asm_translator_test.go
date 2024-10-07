package backend

import (
	"github.com/thomasbollmeier/writing-a-c-compiler/tbcc/frontend"
	"testing"
)

func TestAsmTranslator_Translate(t *testing.T) {
	code := `
int main(void) {
	return 42;
}`
	tokens, _ := frontend.Tokenize(code)
	parser := frontend.NewParser(tokens)
	program, _ := parser.ParseProgram()

	translator := NewAsmTranslator()
	asm_program := translator.Translate(program)

	asm_program.Accept(NewAsmPrinter(4))
}
