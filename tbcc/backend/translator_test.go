package backend

import (
	"github.com/thomasbollmeier/writing-a-c-compiler/tbcc/frontend"
	"github.com/thomasbollmeier/writing-a-c-compiler/tbcc/tacky"
	"testing"
)

func TestTranslator_Translate(t1 *testing.T) {
	code := `
int main(void) {
	return ~(-42);
}`
	tokens, _ := frontend.Tokenize(code)
	program, _ := frontend.NewParser(tokens).ParseProgram()
	tackyProgram := tacky.NewTranslator().Translate(program)

	translator := NewTranslator()
	asmProgram := translator.Translate(tackyProgram)

	asmProgram.Accept(NewAsmPrinter(4))
}

func TestTranslator_TranslateBinOp(t1 *testing.T) {
	code := `
int main(void) {
	return 7 + 5 * 7;
}`
	tokens, _ := frontend.Tokenize(code)
	program, _ := frontend.NewParser(tokens).ParseProgram()
	tackyProgram := tacky.NewTranslator().Translate(program)

	translator := NewTranslator()
	asmProgram := translator.Translate(tackyProgram)

	asmProgram.Accept(NewAsmPrinter(4))
}
