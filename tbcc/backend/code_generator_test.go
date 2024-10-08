package backend

import (
	"fmt"
	"testing"
)

func TestCodeGenerator_GenerateCode(t *testing.T) {
	prog := NewProgram(
		FunctionDef{
			"main",
			[]Instruction{
				NewMov(NewImmediate(42), NewRegister()),
				NewReturn(),
			},
		})

	asm := NewCodeGenerator().GenerateCode(*prog)

	fmt.Print(asm)
}
