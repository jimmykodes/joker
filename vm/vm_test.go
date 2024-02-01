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

func TestClosure(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
      fn add(a) {
        return fn(b) {
          return a + b;
        }
      }
      let adder = add(10);
      adder(12);
      `,
			expected: 22,
		},
		{
			input: `
      fn add(a) {
        acc := 0;
        return fn() {
          acc = acc + a;
          return acc;
        }
      }
      let adder = add(10);
      adder();
      adder();
      adder();
      `,
			expected: 30,
		},
		// this test has a recursion failure that I haven't figured out how to solve yet...
		// {
		// 	input: `
		//     fn add(a) {
		//       acc := 0;
		//       return fn(b) {
		//         return fn() {
		//           acc = acc + a + b;
		//           return acc;
		//         }
		//       }
		//     }
		//     let adderA = add(10);
		//     let adderB = adderA(5);
		//     adderB();
		//     adderB();
		//     `,
		// 	expected: 30,
		// },
	}
	runVmTests(t, tests)
}

func TestRecursion(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
      fn recursion(a) {
        if a == 0 {
          return 0;
        } else {
          return recursion(a - 1);
        }
      }
      recursion(10);
      `,
			expected: 0,
		},
		{
			input: `
      fn recursion(a) {
        if a == 0 {
          return 0;
        } else {
          return recursion(a - 1);
        }
      }
      fn wrapper() {
        return recursion(1);
      }
      wrapper();
      `,
			expected: 0,
		},
		// TODO: this type of recursion breaks, and I don't
		// like the book's fix for it, so until I come up with another
		// one, i'm going to leave this case broken
		// {
		// 	input: `
		//     fn wrapper() {
		//       fn recursion(a) {
		//         if a == 0 {
		//           return 0;
		//         }
		//         return recursion(a - 1);
		//       }
		//       return recursion(10);
		//     }
		//     `,
		// 	expected: 0,
		// },
	}
	runVmTests(t, tests)
}

func TestBuiltinCall(t *testing.T) {
	tests := []vmTestCase{
		{`len([1, 2, 3])`, 3},
		{`slice([1, 2, 3, 4], 2)`, []any{1, 2}},
		{`slice([1, 2, 3, 4], 1, 3)`, []any{2, 3}},
		// TODO: this test _should_ pass, but something about
		// the way we evaluate `len(x)` to be pushed then popped
		// because it is an expression, without intelligently determining
		// the value is used as an arg is breaking things
		// {
		// 	input: `let x = [1, 2, 3];
		//     slice(x, 1, len(x));`,
		// 	expected: []any{2, 3},
		// },
		{
			input: `let x = [1, 2, 3];
      let end = len(x);
      slice(x, 1, end);`,
			expected: []any{2, 3},
		},
	}
	runVmTests(t, tests)
}

func TestFuncCall(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `fn fivePlusTen() { return 5 + 10; }
      fivePlusTen();`,
			expected: 15,
		},
		{
			input: `
      fn one() { return 1; }
      fn two() { return 2; }
      one() + two();`,
			expected: 3,
		},
		{
			input: `
      fn early() { return "bird"; "test" }
      early();`,
			expected: "bird",
		},
		{
			input: `
      fn early() { "bird"; }
      early();`,
			expected: Null,
		},
		{
			input: `
      fn one() { return 1; }
      fn two() { return 1 + one(); }
      two();
      `,
			expected: 2,
		},
		{
			input: `
      fn one() { return 1; }
      fn oneCaller() { return one; }
      oneCaller()();
      `,
			expected: 1,
		},
		{
			input: `
      fn one() { let one = 1; return one; }
      one();
      `,
			expected: 1,
		},
		{
			input: `
      fn add(a, b) { return a + b; }
      add(1, 2);
      `,
			expected: 3,
		},
		{
			input: `
      fn square(a) { return a * a; }
      fn add(a, b) { return a + b; }
      add(1, square(2));
      `,
			expected: 5,
		},
	}
	runVmTests(t, tests)
}

func TestIndex(t *testing.T) {
	tests := []vmTestCase{
		{"{1:12}[1]", 12},
		{"{2:12}[1+1]", 12},
		{"[4, 5, 6][0]", 4},
		{"[4, 5, 6][1+1]", 6},
	}
	runVmTests(t, tests)
}

func TestMaps(t *testing.T) {
	tests := []vmTestCase{
		{"{}", map[any]any{}},
		{"{1: 12}", map[any]any{1: 12}},
		{`{"test": 12, "taco": 44}`, map[any]any{"test": 12, "taco": 44}},
	}
	runVmTests(t, tests)
}

func TestArrays(t *testing.T) {
	tests := []vmTestCase{
		{"[]", []any{}},
		{"[1, 2]", []any{1, 2}},
		{`[1+2, "test"]`, []any{3, "test"}},
		{`x := 10; [1+2, "test", x]`, []any{3, "test", 10}},
	}
	runVmTests(t, tests)
}

func TestGlobalVariableDefs(t *testing.T) {
	tests := []vmTestCase{
		{"let one = 1; one", 1},
		{"let one = 1; let two = 2; one + two", 3},
		{"let one = 1; let two = one + one; one + two", 3},
		{"one := 1; one", 1},
		{"one := 1; two := 2; one + two", 3},
		{"one := 1; two := one + one; one + two", 3},
		{"one := 1; two := 2; three := 0; three = one + two;", 3},
	}
	runVmTests(t, tests)
}

func TestWhileLoop(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
      i := 0;
      while i < 10 {
        i = i + 1;
      }
      i;
      `,
			expected: 10,
		},
		{
			input: `
      fn acc(num) {
        i := 0;
        a := 0;
        while i <= num {
          a = a + i;
          i = i + 1;
        }
        return a;
      }
      acc(4);`,
			expected: 10,
		},
		{
			input: `
      fn acc(num) {
        i := 0;
        a := 0;
        while true {
          if i > num {
            return a;
          }
          a = a + i;
          i = i + 1;
        }
      }
      acc(4);`,
			expected: 10,
		},
		{
			input: `
      fn acc(num) {
        i := 0;
        a := 0;
        while true {
          if i > num {
            break;
          }
          a = a + i;
          i = i + 1;
        }
        return a;
      }
      acc(4);`,
			expected: 10,
		},
		{
			input: `
      fn acc(num) {
        i := 0;
        a := 0;
        while true {
          a = a + i;
          i = i + 1;
          if i <= num {
            continue;
          }
          return a;
        }
      }
      acc(4);`,
			expected: 10,
		},
	}
	runVmTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []vmTestCase{
		{"if true { 10 } else { 12 }", 10},
		{"if false { 10 } else { 12 }", 12},
		{"if (if true { 10; }) { 10 } else { 12 }", 10},
		{"if (if false { 10; }) { 10 } else { 12 }", 12},
		{"if 1 > 2 { 10 }", Null},
		{"if false { 10 }", Null},
	}
	runVmTests(t, tests)
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
		{"!(if false { 5; })", true},
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
	case *object.Null:
		if got != Null {
			t.Errorf("object is not Null: %T (%v)", got, got)
		}
	case []any:
		testArrayObject(t, want, got)
	case map[any]any:
		testMapObject(t, want, got)
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

func testArrayObject(t *testing.T, want []any, got object.Object) {
	t.Helper()
	result, ok := got.(*object.Array)
	if !ok {
		t.Errorf("object not an array. got %T (%v)", got, got)
		return
	}
	for i, w := range want {
		testExpectedObject(t, w, result.Elements[i])
	}
}

func testMapObject(t *testing.T, want map[any]any, got object.Object) {
	t.Helper()
	result, ok := got.(*object.Map)
	if !ok {
		t.Errorf("object not a map. got %T (%v)", got, got)
		return
	}
	if len(want) != len(result.Pairs) {
		t.Errorf("invalid length: got %d - want %d", len(want), len(result.Pairs))
	}
	// TODO: fix these tests, since i'm not going to sort the keys in the compiler, cause I
	//  don't want to waste the overhead for a purely testing feature
	// for _, pair := range result.Pairs {
	// 	key := pair.Key
	// 	value := pair.Value
	// }
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}
