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

func TestStringArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             `"test" + "thing"`,
			expectedConstants: []any{"test", "thing"},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
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
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
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
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
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
		t.Run(tt.input, func(t *testing.T) {
			program := parse(tt.input)
			compiler := New()
			err := compiler.Compile(program)
			if err != nil {
				t.Fatalf("compiler error: %s", err)
			}
			bytecode := compiler.Bytecode()
			err = testInstructions(tt.expectedInstructions, bytecode.Instructions)
			if err != nil {
				t.Fatalf("test instructions failed: %s", err)
			}

			err = testConstants(tt.expectedConstants, bytecode.Constants)
			if err != nil {
				t.Fatalf("test constants failed: %s", err)
			}
		})
	}
}

func testInstructions(want []code.Instructions, got code.Instructions) error {
	joined := concatInstructions(want)
	if len(got) != len(joined) {
		return fmt.Errorf("wrong lengths:\ngot %q\nwant %q", got, joined)
	}
	for i, ins := range joined {
		if ins != got[i] {
			return fmt.Errorf("mismatched instruction at %d: got %q - want %q", i, got[i], ins)
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
