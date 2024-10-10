package tacky

import (
	"fmt"
	"github.com/thomasbollmeier/writing-a-c-compiler/tbcc/frontend"
	"testing"
)

func TestEmitter_Emit(t *testing.T) {
	code := `
int main(void) {
	return ~(-42);
}`
	ast := parse(code)
	emitter := NewEmitter()
	program := emitter.Emit(ast)

	fmt.Println(program)
}

func parse(code string) *frontend.Program {
	tokens, _ := frontend.Tokenize(code)
	parser := frontend.NewParser(tokens)
	ret, _ := parser.ParseProgram()
	return ret
}
