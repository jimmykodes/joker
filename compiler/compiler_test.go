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

func TestClosures(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
      fn(a) {
        return fn(b) {
          return a + b;
        };
      }
      `,
			expectedConstants: []any{
				[]code.Instructions{
					code.Instruction(code.OpGetFree, 0),
					code.Instruction(code.OpGetLocal, 0),
					code.Instruction(code.OpAdd),
					code.Instruction(code.OpReturn),
				},
				[]code.Instructions{
					code.Instruction(code.OpGetLocal, 0),
					code.Instruction(code.OpClosure, 0, 1),
					code.Instruction(code.OpReturn),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpClosure, 1, 0),
				code.Instruction(code.OpPop),
			},
		},
		{
			input: `
      fn(a) {
        return fn(b) {
          return fn(c) {
            return a + b + c;
          }
        }
      }
      `,
			expectedConstants: []any{
				[]code.Instructions{
					code.Instruction(code.OpGetFree, 0),
					code.Instruction(code.OpGetFree, 1),
					code.Instruction(code.OpAdd),
					code.Instruction(code.OpGetLocal, 0),
					code.Instruction(code.OpAdd),
					code.Instruction(code.OpReturn),
				},
				[]code.Instructions{
					code.Instruction(code.OpGetFree, 0),
					code.Instruction(code.OpGetLocal, 0),
					code.Instruction(code.OpClosure, 0, 2),
					code.Instruction(code.OpReturn),
				},
				[]code.Instructions{
					code.Instruction(code.OpGetLocal, 0),
					code.Instruction(code.OpClosure, 1, 1),
					code.Instruction(code.OpReturn),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpClosure, 2, 0),
				code.Instruction(code.OpPop),
			},
		},
		{
			input: `
      let g = 55;
      fn() {
        a := 66;
        return fn() {
          b := 77;
          return fn() {
            c := 88;
            return g + a + b + c;
          }
        }
      }
      `,
			expectedConstants: []any{
				55, 66, 77, 88,
				[]code.Instructions{
					code.Instruction(code.OpConstant, 3),
					code.Instruction(code.OpSetLocal, 0),
					code.Instruction(code.OpGetGlobal, 0),
					code.Instruction(code.OpGetFree, 0),
					code.Instruction(code.OpAdd),
					code.Instruction(code.OpGetFree, 1),
					code.Instruction(code.OpAdd),
					code.Instruction(code.OpGetLocal, 0),
					code.Instruction(code.OpAdd),
					code.Instruction(code.OpReturn),
				},
				[]code.Instructions{
					code.Instruction(code.OpConstant, 2),
					code.Instruction(code.OpSetLocal, 0),
					code.Instruction(code.OpGetFree, 0),
					code.Instruction(code.OpGetLocal, 0),
					code.Instruction(code.OpClosure, 4, 2),
					code.Instruction(code.OpReturn),
				},
				[]code.Instructions{
					code.Instruction(code.OpConstant, 1),
					code.Instruction(code.OpSetLocal, 0),
					code.Instruction(code.OpGetLocal, 0),
					code.Instruction(code.OpClosure, 5, 1),
					code.Instruction(code.OpReturn),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpSetGlobal, 0),
				code.Instruction(code.OpClosure, 6, 0),
				code.Instruction(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestVarScopes(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
      let num = 255;
      fn() { return num; }`,
			expectedConstants: []any{
				255,
				[]code.Instructions{
					code.Instruction(code.OpGetGlobal, 0),
					code.Instruction(code.OpReturn),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpSetGlobal, 0),
				code.Instruction(code.OpClosure, 1, 0),
				code.Instruction(code.OpPop),
			},
		},
		{
			input: `
      fn() {
        let num = 255;
        return num; 
      }`,
			expectedConstants: []any{
				255,
				[]code.Instructions{
					code.Instruction(code.OpConstant, 0),
					code.Instruction(code.OpSetLocal, 0),
					code.Instruction(code.OpGetLocal, 0),
					code.Instruction(code.OpReturn),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpClosure, 1, 0),
				code.Instruction(code.OpPop),
			},
		},
		{
			input: `
      fn() {
        let a = 10;
        let b = 5;
        return a + b; 
      }`,
			expectedConstants: []any{
				10,
				5,
				[]code.Instructions{
					code.Instruction(code.OpConstant, 0),
					code.Instruction(code.OpSetLocal, 0),
					code.Instruction(code.OpConstant, 1),
					code.Instruction(code.OpSetLocal, 1),
					code.Instruction(code.OpGetLocal, 0),
					code.Instruction(code.OpGetLocal, 1),
					code.Instruction(code.OpAdd),
					code.Instruction(code.OpReturn),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpClosure, 2, 0),
				code.Instruction(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestFunctionCalls(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
      fn onePlusTwo() { return 1 + 2; }
      onePlusTwo();
      `,
			expectedConstants: []any{
				1,
				2,
				[]code.Instructions{
					code.Instruction(code.OpConstant, 0),
					code.Instruction(code.OpConstant, 1),
					code.Instruction(code.OpAdd),
					code.Instruction(code.OpReturn),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpClosure, 2, 0),
				code.Instruction(code.OpSetGlobal, 0),
				code.Instruction(code.OpGetGlobal, 0),
				code.Instruction(code.OpCall),
				code.Instruction(code.OpPop),
			},
		},
		{
			input: `
      let onePlusTwo = fn() { return 1 + 2; };
      onePlusTwo();
      `,
			expectedConstants: []any{
				1,
				2,
				[]code.Instructions{
					code.Instruction(code.OpConstant, 0),
					code.Instruction(code.OpConstant, 1),
					code.Instruction(code.OpAdd),
					code.Instruction(code.OpReturn),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpClosure, 2, 0),
				code.Instruction(code.OpSetGlobal, 0),
				code.Instruction(code.OpGetGlobal, 0),
				code.Instruction(code.OpCall),
				code.Instruction(code.OpPop),
			},
		},
		{
			input: `fn() { return 1 + 2; }()`,
			expectedConstants: []any{
				1,
				2,
				[]code.Instructions{
					code.Instruction(code.OpConstant, 0),
					code.Instruction(code.OpConstant, 1),
					code.Instruction(code.OpAdd),
					code.Instruction(code.OpReturn),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpClosure, 2, 0),
				code.Instruction(code.OpCall),
				code.Instruction(code.OpPop),
			},
		},
		{
			input: `fn add(a, b) { return a + b; }
		    add(12, 13);`,
			expectedConstants: []any{
				[]code.Instructions{
					code.Instruction(code.OpGetLocal, 0),
					code.Instruction(code.OpGetLocal, 1),
					code.Instruction(code.OpAdd),
					code.Instruction(code.OpReturn),
				},
				12,
				13,
			},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpClosure, 0, 0),
				code.Instruction(code.OpSetGlobal, 0),
				code.Instruction(code.OpGetGlobal, 0),
				code.Instruction(code.OpConstant, 1),
				code.Instruction(code.OpConstant, 2),
				code.Instruction(code.OpCall, 2),
				code.Instruction(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestFunctions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: "fn() { return 5 + 10 }",
			expectedConstants: []any{
				5,
				10,
				[]code.Instructions{
					code.Instruction(code.OpConstant, 0),
					code.Instruction(code.OpConstant, 1),
					code.Instruction(code.OpAdd),
					code.Instruction(code.OpReturn),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpClosure, 2, 0),
				code.Instruction(code.OpPop),
			},
		},
		{
			input: "fn() { 5 + 10 }",
			expectedConstants: []any{
				5,
				10,
				[]code.Instructions{
					code.Instruction(code.OpConstant, 0),
					code.Instruction(code.OpConstant, 1),
					code.Instruction(code.OpAdd),
					code.Instruction(code.OpPop),
					code.Instruction(code.OpNull),
					code.Instruction(code.OpReturn),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpClosure, 2, 0),
				code.Instruction(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestIndexExpression(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "[1, 2, 3][2]",
			expectedConstants: []any{1, 2, 3, 2},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpConstant, 1),
				code.Instruction(code.OpConstant, 2),
				code.Instruction(code.OpArray, 3),
				code.Instruction(code.OpConstant, 3),
				code.Instruction(code.OpIndex),
				code.Instruction(code.OpPop),
			},
		},
		{
			input:             "[1, 2, 3][1 + 1]",
			expectedConstants: []any{1, 2, 3, 1, 1},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpConstant, 1),
				code.Instruction(code.OpConstant, 2),
				code.Instruction(code.OpArray, 3),
				code.Instruction(code.OpConstant, 3),
				code.Instruction(code.OpConstant, 4),
				code.Instruction(code.OpAdd),
				code.Instruction(code.OpIndex),
				code.Instruction(code.OpPop),
			},
		},
		{
			input:             "{1: 12}[1]",
			expectedConstants: []any{1, 12, 1},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpConstant, 1),
				code.Instruction(code.OpMap, 1),
				code.Instruction(code.OpConstant, 2),
				code.Instruction(code.OpIndex),
				code.Instruction(code.OpPop),
			},
		},
		// TODO: this test fails intermittently due to random map access
		// so sometimes constants are out of order. not sure the best way to test that
		// {
		// 	input:             `{"foo": "bar", "baz": "bing"}["fo"+"o"]`,
		// 	expectedConstants: []any{"foo", "bar", "baz", "bing", "fo", "o"},
		// 	expectedInstructions: []code.Instructions{
		// 		code.Instruction(code.OpConstant, 0),
		// 		code.Instruction(code.OpConstant, 1),
		// 		code.Instruction(code.OpConstant, 2),
		// 		code.Instruction(code.OpConstant, 3),
		// 		code.Instruction(code.OpMap, 2),
		// 		code.Instruction(code.OpConstant, 4),
		// 		code.Instruction(code.OpConstant, 5),
		// 		code.Instruction(code.OpAdd),
		// 		code.Instruction(code.OpIndex),
		// 		code.Instruction(code.OpPop),
		// 	},
		// },
	}
	runCompilerTests(t, tests)
}

func TestMapLiterals(t *testing.T) {
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
		// TODO: this fails because of random order map access. not sure how to fix atm
		// {
		// 	input:             `{"test": 12, "thing": 44}`,
		// 	expectedConstants: []any{"test", 12, "thing", 44},
		// 	expectedInstructions: []code.Instructions{
		// 		code.Instruction(code.OpConstant, 0),
		// 		code.Instruction(code.OpConstant, 1),
		// 		code.Instruction(code.OpConstant, 2),
		// 		code.Instruction(code.OpConstant, 3),
		// 		code.Instruction(code.OpMap, 2),
		// 		code.Instruction(code.OpPop),
		// 	},
		// },
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

func TestReassignStatements(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
      one := 1;
      one = 2;
      one;
      `,
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpSetGlobal, 0),
				code.Instruction(code.OpConstant, 1),
				code.Instruction(code.OpSetGlobal, 0),
				code.Instruction(code.OpGetGlobal, 0),
				code.Instruction(code.OpPop),
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
      two := 2;
      two = one;
      two;
      `,
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Instruction(code.OpConstant, 0),
				code.Instruction(code.OpSetGlobal, 0),
				code.Instruction(code.OpConstant, 1),
				code.Instruction(code.OpSetGlobal, 1),
				code.Instruction(code.OpGetGlobal, 0),
				code.Instruction(code.OpSetGlobal, 1),
				code.Instruction(code.OpGetGlobal, 1),
				code.Instruction(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestForLoop(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
      a := 0;
      for i := 0; i < 10; i = i + 1; {
        a = a + 1;
      }
      a;`,
			expectedConstants: []any{0, 0, 1, 10, 1},
			expectedInstructions: []code.Instructions{
				// 0000 OpConstant 0
				code.Instruction(code.OpConstant, 0),
				// 0003 OpSetGlobal 0
				code.Instruction(code.OpSetGlobal, 0),
				// 0006 OpConstant 1
				code.Instruction(code.OpConstant, 1),
				// 0009 OpSetGlobal 1
				code.Instruction(code.OpSetGlobal, 1),
				// 0012 OpJump 25
				code.Instruction(code.OpJump, 25),
				// 0015 OpGetGlobal 1
				code.Instruction(code.OpGetGlobal, 1),
				// 0018 OpConstant 2
				code.Instruction(code.OpConstant, 2),
				// 0021 OpAdd
				code.Instruction(code.OpAdd),
				// 0022 OpSetGlobal 1
				code.Instruction(code.OpSetGlobal, 1),
				// 0025 OpConstant 3
				code.Instruction(code.OpConstant, 3),
				// 0028 OpGetGlobal 1
				code.Instruction(code.OpGetGlobal, 1),
				// 0031 OpGT
				code.Instruction(code.OpGT),
				// 0032 OpJumpNotTruthy 48
				code.Instruction(code.OpJumpNotTruthy, 48),
				// 0035 OpGetGlobal 0
				code.Instruction(code.OpGetGlobal, 0),
				// 0038 OpConstant 4
				code.Instruction(code.OpConstant, 4),
				// 0041 OpAdd
				code.Instruction(code.OpAdd),
				// 0042 OpSetGlobal 0
				code.Instruction(code.OpSetGlobal, 0),
				// 0045 OpJump 15
				code.Instruction(code.OpJump, 15),
				// 0048 OpGetGlobal 0
				code.Instruction(code.OpGetGlobal, 0),
				// 0051 OpPop
				code.Instruction(code.OpPop),
			},
		},
		{
			input: `
		    a := 0;
		    for i := 0; i < 10; i = i + 1; {
		      if i == 5 {
		        continue;
		      }
		      if i == 8 {
		        break;
		      }
		      a = a + 1;
		    }
		    a;`,
			expectedConstants: []any{0, 0, 1, 10, 5, 8, 1},
			expectedInstructions: []code.Instructions{
				// 0000 OpConstant 0
				code.Instruction(code.OpConstant, 0),
				// 0003 OpSetGlobal 0
				code.Instruction(code.OpSetGlobal, 0),
				// 0006 OpConstant 1
				code.Instruction(code.OpConstant, 1),
				// 0009 OpSetGlobal 1
				code.Instruction(code.OpSetGlobal, 1),
				// 0012 OpJump 25
				code.Instruction(code.OpJump, 25),
				// 0015 OpGetGlobal 1
				code.Instruction(code.OpGetGlobal, 1),
				// 0018 OpConstant 2
				code.Instruction(code.OpConstant, 2),
				// 0021 OpAdd
				code.Instruction(code.OpAdd),
				// 0022 OpSetGlobal 1
				code.Instruction(code.OpSetGlobal, 1),
				// 0025 OpConstant 3
				code.Instruction(code.OpConstant, 3),
				// 0028 OpGetGlobal 1
				code.Instruction(code.OpGetGlobal, 1),
				// 0031 OpGT
				code.Instruction(code.OpGT),
				// 0032 OpJumpNotTruthy 74
				code.Instruction(code.OpJumpNotTruthy, 74),
				// 0035 OpGetGlobal 1
				code.Instruction(code.OpGetGlobal, 1),
				// 0038 OpConstant 4
				code.Instruction(code.OpConstant, 4),
				// 0041 OpEQ
				code.Instruction(code.OpEQ),
				// 0042 OpJumpNotTruthy 48
				code.Instruction(code.OpJumpNotTruthy, 48),
				// 0045 OpJump 15
				code.Instruction(code.OpJump, 15),
				// 0048 OpGetGlobal 1
				code.Instruction(code.OpGetGlobal, 1),
				// 0051 OpConstant 5
				code.Instruction(code.OpConstant, 5),
				// 0054 OpEQ
				code.Instruction(code.OpEQ),
				// 0055 OpJumpNotTruthy 61
				code.Instruction(code.OpJumpNotTruthy, 61),
				// 0058 OpJump 74
				code.Instruction(code.OpJump, 74),
				// 0061 OpGetGlobal 0
				code.Instruction(code.OpGetGlobal, 0),
				// 0064 OpConstant 6
				code.Instruction(code.OpConstant, 6),
				// 0067 OpAdd
				code.Instruction(code.OpAdd),
				// 0068 OpSetGlobal 0
				code.Instruction(code.OpSetGlobal, 0),
				// 0071 OpJump 15
				code.Instruction(code.OpJump, 15),
				// 0074 OpGetGlobal 0
				code.Instruction(code.OpGetGlobal, 0),
				// 0077 OpPop
				code.Instruction(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestWhileLoop(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
      a := 0;
      while true {
        a = a + 1;
      }
      a;`,
			expectedConstants: []any{0, 1},
			expectedInstructions: []code.Instructions{
				// 0000
				code.Instruction(code.OpConstant, 0),
				// 0003
				code.Instruction(code.OpSetGlobal, 0),
				// 0006
				code.Instruction(code.OpTrue),
				// 0007
				code.Instruction(code.OpJumpNotTruthy, 23),
				// 0010
				code.Instruction(code.OpGetGlobal, 0),
				// 0013
				code.Instruction(code.OpConstant, 1),
				// 0016
				code.Instruction(code.OpAdd),
				// 0017
				code.Instruction(code.OpSetGlobal, 0),
				// 0020
				code.Instruction(code.OpJump, 6),
				// 0023
				code.Instruction(code.OpGetGlobal, 0),
				// 0025
				code.Instruction(code.OpPop),
			},
		},
		{
			input: `
      a := 0;
      while true {
        a = a + 1;
      }
      a;`,
			expectedConstants: []any{0, 1},
			expectedInstructions: []code.Instructions{
				// 0000
				code.Instruction(code.OpConstant, 0),
				// 0003
				code.Instruction(code.OpSetGlobal, 0),
				// 0006
				code.Instruction(code.OpTrue),
				// 0007
				code.Instruction(code.OpJumpNotTruthy, 23),
				// 0010
				code.Instruction(code.OpGetGlobal, 0),
				// 0013
				code.Instruction(code.OpConstant, 1),
				// 0016
				code.Instruction(code.OpAdd),
				// 0017
				code.Instruction(code.OpSetGlobal, 0),
				// 0020
				code.Instruction(code.OpJump, 6),
				// 0023
				code.Instruction(code.OpGetGlobal, 0),
				// 0026
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
				code.Instruction(code.OpJumpNotTruthy, 8),
				// 0004
				code.Instruction(code.OpConstant, 0),
				// 0007
				code.Instruction(code.OpPop),
				// 0008
				code.Instruction(code.OpConstant, 1),
				// 0010
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
				code.Instruction(code.OpJumpNotTruthy, 11),
				// 0004
				code.Instruction(code.OpConstant, 0),
				// 0007
				code.Instruction(code.OpPop),
				// 0008
				code.Instruction(code.OpJump, 15),
				// 0011
				code.Instruction(code.OpConstant, 1),
				// 0014
				code.Instruction(code.OpPop),
				// 0015
				code.Instruction(code.OpConstant, 2),
				// 0018
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
		return fmt.Errorf("wrong lengths:\ngot  %q\nwant %q", got, joined)
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
		case []code.Instructions:
			result, ok := got[i].(*object.CompiledFunction)
			if !ok {
				return fmt.Errorf("object is not a function: got %T (%v)", got, got)
			}
			if err := testInstructions(constant, result.Instructions); err != nil {
				return err
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
