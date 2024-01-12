package compiler

import (
	"fmt"
	"testing"

	"github.com/jimmykodes/joker/ast"
	"github.com/jimmykodes/joker/code"
	"github.com/jimmykodes/joker/lexer"
	"github.com/jimmykodes/joker/object"
	"github.com/jimmykodes/joker/parser"
)

type compilerTestCase struct {
	input                string
	expectedConstants    []any
	expectedInstructions []code.Instructions
}

func TestDictLiterals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "{}",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpMap, 0),
				code.Instruction(code.OpPop),
			},
		},
		{
			input:             "{1: 12}",
			expectedConstants: []any{1, 12},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpConstant, 1),
				code.Instruction(code.OpMap, 1),
				code.Instruction(code.OpPop),
			},
		},
		{
			input:             `{"test": 12, "thing": 44}`,
			expectedConstants: []any{"test", 12, "thing", 44},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpConstant, 1),
				code.Instruction(code.OpConstant, 2),
				code.Instruction(code.OpConstant, 3),
				code.Instruction(code.OpMap, 2),
				code.Instruction(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestArrayLiterals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "[]",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpArray, 0),
				code.Instruction(code.OpPop),
			},
		},
		{
			input:             "[1, 2, 3]",
			expectedConstants: []any{1, 2, 3},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpConstant, 1),
				code.Instruction(code.OpConstant, 2),
				code.Instruction(code.OpArray, 3),
				code.Instruction(code.OpPop),
			},
		},
		{
			input:             "[1+2, 3+4, 5+6]",
			expectedConstants: []any{1, 2, 3, 4, 5, 6},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpConstant, 1),
				code.Instruction(code.OpAdd),
				code.Instruction(code.OpConstant, 2),
				code.Instruction(code.OpConstant, 3),
				code.Instruction(code.OpAdd),
				code.Instruction(code.OpConstant, 4),
				code.Instruction(code.OpConstant, 5),
				code.Instruction(code.OpAdd),
				code.Instruction(code.OpArray, 3),
				code.Instruction(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
      let one = 1;
      let two = 2;
      `,
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpSetGlobal, 0),
				code.Instruction(code.OpConstant, 1),
				code.Instruction(code.OpSetGlobal, 1),
			},
		},
		{
			input: `
      let one = 1;
      one;
      `,
			expectedConstants: []any{1},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpSetGlobal, 0),
				code.Instruction(code.OpGetGlobal, 0),
				code.Instruction(code.OpPop),
			},
		},
		{
			input: `
      let one = 1;
      let two = one;
      two;
      `,
			expectedConstants: []any{1},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpSetGlobal, 0),
				code.Instruction(code.OpGetGlobal, 0),
				code.Instruction(code.OpSetGlobal, 1),
				code.Instruction(code.OpGetGlobal, 1),
				code.Instruction(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestGlobalDefineStatements(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
      one := 1;
      two := 2;
      `,
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpSetGlobal, 0),
				code.Instruction(code.OpConstant, 1),
				code.Instruction(code.OpSetGlobal, 1),
			},
		},
		{
			input: `
      one := 1;
      one;
      `,
			expectedConstants: []any{1},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpSetGlobal, 0),
				code.Instruction(code.OpGetGlobal, 0),
				code.Instruction(code.OpPop),
			},
		},
		{
			input: `
      one := 1;
      two := one;
      two;
      `,
			expectedConstants: []any{1},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpSetGlobal, 0),
				code.Instruction(code.OpGetGlobal, 0),
				code.Instruction(code.OpSetGlobal, 1),
				code.Instruction(code.OpGetGlobal, 1),
				code.Instruction(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "if true { 10 }; 3333;",
			expectedConstants: []any{10, 3333},
			expectedInstructions: []code.Instructions{
				// 0000
				code.Instruction(code.OpTrue),
				// 0001
				code.Instruction(code.OpJumpNotTruthy, 10),
				// 0004
				code.Instruction(code.OpConstant, 0),
				// 0007
				code.Instruction(code.OpJump, 11),
				// 0010
				code.Instruction(code.OpNull),
				// 0011
				code.Instruction(code.OpPop),
				// 0012
				code.Instruction(code.OpConstant, 1),
				// 0015
				code.Instruction(code.OpPop),
			},
		},
		{
			input:             "if true { 10 } else { 12 }; 3333;",
			expectedConstants: []any{10, 12, 3333},
			expectedInstructions: []code.Instructions{
				// 0000
				code.Instruction(code.OpTrue),
				// 0001
				code.Instruction(code.OpJumpNotTruthy, 10),
				// 0004
				code.Instruction(code.OpConstant, 0),
				// 0007
				code.Instruction(code.OpJump, 13),
				// 0010
				code.Instruction(code.OpConstant, 1),
				// 0011
				code.Instruction(code.OpPop),
				// 0012
				code.Instruction(code.OpConstant, 2),
				// 0013
				code.Instruction(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestStringArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             `"test" + "thing"`,
			expectedConstants: []any{"test", "thing"},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpConstant, 1),
				code.Instruction(code.OpAdd),
				code.Instruction(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestFloatArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "1.0 + 2.0",
			expectedConstants: []any{1.0, 2.0},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpConstant, 1),
				code.Instruction(code.OpAdd),
				code.Instruction(code.OpPop),
			},
		},
		{
			input:             "1.0 - 2.0",
			expectedConstants: []any{1.0, 2.0},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpConstant, 1),
				code.Instruction(code.OpSub),
				code.Instruction(code.OpPop),
			},
		},
		{
			input:             "1.0 * 2.0",
			expectedConstants: []any{1.0, 2.0},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpConstant, 1),
				code.Instruction(code.OpMult),
				code.Instruction(code.OpPop),
			},
		},
		{
			input:             "1.0 / 2.0",
			expectedConstants: []any{1.0, 2.0},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpConstant, 1),
				code.Instruction(code.OpDiv),
				code.Instruction(code.OpPop),
			},
		},
		{
			input:             "1.0 % 2.0",
			expectedConstants: []any{1.0, 2.0},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpConstant, 1),
				code.Instruction(code.OpMod),
				code.Instruction(code.OpPop),
			},
		},
		{
			input:             "-1.0 - 2.0",
			expectedConstants: []any{1.0, 2.0},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpMinus),
				code.Instruction(code.OpConstant, 1),
				code.Instruction(code.OpSub),
				code.Instruction(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "1 + 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpConstant, 1),
				code.Instruction(code.OpAdd),
				code.Instruction(code.OpPop),
			},
		},
		{
			input:             "1 - 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpConstant, 1),
				code.Instruction(code.OpSub),
				code.Instruction(code.OpPop),
			},
		},
		{
			input:             "1 * 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpConstant, 1),
				code.Instruction(code.OpMult),
				code.Instruction(code.OpPop),
			},
		},
		{
			input:             "1 / 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpConstant, 1),
				code.Instruction(code.OpDiv),
				code.Instruction(code.OpPop),
			},
		},
		{
			input:             "1 % 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpConstant, 1),
				code.Instruction(code.OpMod),
				code.Instruction(code.OpPop),
			},
		},
		{
			input:             "-1 - 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpMinus),
				code.Instruction(code.OpConstant, 1),
				code.Instruction(code.OpSub),
				code.Instruction(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "true",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpTrue),
				code.Instruction(code.OpPop),
			},
		},
		{
			input:             "false",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpFalse),
				code.Instruction(code.OpPop),
			},
		},
		{
			input:             "!true",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpTrue),
				code.Instruction(code.OpBang),
				code.Instruction(code.OpPop),
			},
		},
		{
			input:             "!false",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpFalse),
				code.Instruction(code.OpBang),
				code.Instruction(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestComparisonExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "1 == 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpConstant, 1),
				code.Instruction(code.OpEQ),
				code.Instruction(code.OpPop),
			},
		},
		{
			input:             "1 != 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpConstant, 1),
				code.Instruction(code.OpNEQ),
				code.Instruction(code.OpPop),
			},
		},
		{
			input:             "1 < 2",
			expectedConstants: []any{2, 1},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpConstant, 1),
				code.Instruction(code.OpGT),
				code.Instruction(code.OpPop),
			},
		},
		{
			input:             "1 > 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpConstant, 1),
				code.Instruction(code.OpGT),
				code.Instruction(code.OpPop),
			},
		},
		{
			input:             "1 <= 2",
			expectedConstants: []any{2, 1},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpConstant, 1),
				code.Instruction(code.OpGTE),
				code.Instruction(code.OpPop),
			},
		},
		{
			input:             "1 >= 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpConstant, 1),
				code.Instruction(code.OpGTE),
				code.Instruction(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()
	for _, tt := range tests {
		program := parse(tt.input)
		compiler := New()
		err := compiler.Compile(program)
		if err != nil {
			t.Errorf("compiler error: %s", err)
			continue
		}
		bytecode := compiler.Bytecode()
		err = testInstructions(tt.expectedInstructions, bytecode.Instructions)
		if err != nil {
			t.Errorf("test instructions failed: %s", err)
			continue
		}

		err = testConstants(tt.expectedConstants, bytecode.Constants)
		if err != nil {
			t.Errorf("test constants failed: %s", err)
			continue
		}
	}
}

func testInstructions(want []code.Instructions, got code.Instructions) error {
	joined := concatInstructions(want)
	if len(got) != len(joined) {
		return fmt.Errorf("wrong lengths:\ngot %q\nwant %q", got, joined)
	}
	for i, ins := range joined {
		if ins != got[i] {
			return fmt.Errorf("mismatched instruction at %d:\ngot  %q\nwant %q", i, got, joined)
		}
	}
	return nil
}

func testConstants(want []any, got []object.Object) error {
	if len(want) != len(got) {
		return fmt.Errorf("invalid number of constants: got %d - want %d", len(got), len(want))
	}
	for i, constant := range want {
		switch constant := constant.(type) {
		case int:
			if err := testIntegerObject(int64(constant), got[i]); err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s", i, err)
			}
		case float64:
			if err := testFloatObject(constant, got[i]); err != nil {
				return fmt.Errorf("constant %d - testFloatObject failed: %s", i, err)
			}
		case string:
			if err := testStringObject(constant, got[i]); err != nil {
				return fmt.Errorf("constant %d - testStringObject failed: %s", i, err)
			}
		default:
			return fmt.Errorf("missing test for constant: %T", constant)
		}
	}
	return nil
}

func concatInstructions(in []code.Instructions) code.Instructions {
	var out code.Instructions
	for _, inst := range in {
		out = append(out, inst...)
	}
	return out
}

func testIntegerObject(want int64, got object.Object) error {
	res, ok := got.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not an integer. got %T (%v)", got, got)
	}
	if res.Value != want {
		return fmt.Errorf("values not equal: got %d - want %d", res.Value, want)
	}
	return nil
}

func testFloatObject(want float64, got object.Object) error {
	res, ok := got.(*object.Float)
	if !ok {
		return fmt.Errorf("object is not a float. got %T (%v)", got, got)
	}
	if res.Value != want {
		return fmt.Errorf("values not equal: got %f - want %f", res.Value, want)
	}
	return nil
}

func testStringObject(want string, got object.Object) error {
	res, ok := got.(*object.String)
	if !ok {
		return fmt.Errorf("object is not a string. got %T (%v)", got, got)
	}
	if res.Value != want {
		return fmt.Errorf("values not equal: got %s - want %s", res.Value, want)
	}
	return nil
}
