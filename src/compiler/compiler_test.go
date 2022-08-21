package compiler

import (
	"mnk/src/code"
	"testing"
)

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{{
		input:             "1 + 2",
		expectedConstants: []interface{}{1, 2},
		expectedInstructions: []code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpAdd),
		}},
	}

	runCompilerTests(t, tests)
}
