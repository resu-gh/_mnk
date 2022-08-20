package vm

import "testing"

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 2}, // FIXME
	}

	runVmTests(t, tests)
}
