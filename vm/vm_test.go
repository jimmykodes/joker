package vm

import (
	"fmt"
	"testing"

	"github.com/jimmykodes/joker/ast"
	"github.com/jimmykodes/joker/compiler"
	"github.com/jimmykodes/joker/lexer"
	"github.com/jimmykodes/joker/object"
	"github.com/jimmykodes/joker/parser"
)

type vmTestCase struct {
	input    string
	expected any
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"-1", -1},
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"4 - 3", 1},
		{"-4 - 3", -7},
		{"6 / 3", 2},
		{"3 * 4", 12},
		{"15 % 7", 1},
	}
	runVmTests(t, tests)
}

func TestFloatArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"-1.0", -1.0},
		{"1.0", 1.0},
		{"2.5", 2.5},
		{"1.0 + 2.5", 3.5},
		{"2 + 2.5", 4.5},
		{"4.0 - 3", 1.0},
		{"6.0 / 4.0", 1.5},
		{"3.0 * 4", 12.0},
		{"-3.0 * 4", -12.0},
	}
	runVmTests(t, tests)
}

func TestStringArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{`"test"`, "test"},
		{`"taco"`, "taco"},
		{`"test" + "taco"`, "testtaco"},
	}
	runVmTests(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"false", false},
		{"true", true},

		{"1 < 2", true},
		{"1 < 1", false},
		{"2 < 1", false},

		{"2 > 1", true},
		{"1 > 2", false},
		{"1 > 1", false},

		{"1 <= 2", true},
		{"1 <= 1", true},
		{"2 <= 1", false},

		{"1 >= 2", false},
		{"1 >= 1", true},
		{"2 >= 1", true},

		{"1 == 2", false},
		{"1 == 1", true},
		{"2 == 1", false},

		{"1 != 2", true},
		{"1 != 1", false},
		{"2 != 1", true},

		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},

		{"(2 > 1) == true", true},
		{"(2 > 1) == false", false},

		{"(1 <= 1) == true", true},
		{"(1 <= 1) == false", false},

		{"(1 >= 1) == true", true},
		{"(1 >= 1) == false", false},

		{"(1 == 1) == true", true},
		{"(1 == 1) == false", false},

		{"(1 != 2) == true", true},
		{"(1 != 2) == false", false},

		{"!false", true},
		{"!true", false},
		{"!!false", false},
		{"!!true", true},
		{"!(2 != 1)", false},
		{"!(2 == 1)", true},
		{"!((1 < 2) == true)", false},
		{"!((1 < 2) == false)", true},
	}
	runVmTests(t, tests)
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program := parse(tt.input)

			comp := compiler.New()
			if err := comp.Compile(program); err != nil {
				t.Errorf("compiler error: %s", err)
				return
			}

			vm := New(comp.Bytecode())
			if err := vm.Run(); err != nil {
				t.Errorf("vm error: %s", err)
				return
			}

			stackElem := vm.LastPoppedStackElem()
			testExpectedObject(t, tt.expected, stackElem)
		})
	}
}

func testExpectedObject(t *testing.T, want any, got object.Object) {
	t.Helper()
	switch want := want.(type) {
	case int:
		if err := testIntegerObject(int64(want), got); err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	case float64:
		if err := testFloatObject(want, got); err != nil {
			t.Errorf("testFloatObject failed: %s", err)
		}
	case string:
		if err := testStringObject(want, got); err != nil {
			t.Errorf("testStringObject failed: %s", err)
		}
	case bool:
		if err := testBoolObject(want, got); err != nil {
			t.Errorf("testBoolObject failed: %s", err)
		}
	default:
		t.Errorf("missing test for type: %T", want)

	}
}

func testIntegerObject(want int64, got object.Object) error {
	result, ok := got.(*object.Integer)
	if !ok {
		return fmt.Errorf("object not an integer. got %T (%v)", got, got)
	}
	if result.Value != want {
		return fmt.Errorf("incorrect value: got %d - want %d", result.Value, want)
	}
	return nil
}

func testStringObject(want string, got object.Object) error {
	result, ok := got.(*object.String)
	if !ok {
		return fmt.Errorf("object not a string. got %T (%v)", got, got)
	}
	if result.Value != want {
		return fmt.Errorf("incorrect value: got %s - want %s", result.Value, want)
	}
	return nil
}

func testFloatObject(want float64, got object.Object) error {
	result, ok := got.(*object.Float)
	if !ok {
		return fmt.Errorf("object not a float. got %T (%v)", got, got)
	}
	if result.Value != want {
		return fmt.Errorf("incorrect value: got %f - want %f", result.Value, want)
	}
	return nil
}

func testBoolObject(want bool, got object.Object) error {
	if want {
		if got != object.True {
			return fmt.Errorf("invalid value: got %v - want %v", got, want)
		}
	} else {
		if got != object.False {
			return fmt.Errorf("invalid value: got %v - want %v", got, want)
		}
	}
	return nil
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}
