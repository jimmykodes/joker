package compiler

import (
	"fmt"

	"github.com/jimmykodes/joker/ast"
	"github.com/jimmykodes/joker/code"
	"github.com/jimmykodes/joker/object"
)

type Compiler struct {
	instructions code.Instructions
	constants    []object.Object
	symbolTable  *SymbolTable

	// ultInst is the last (ultimate) instruction emitted
	ultInst EmittedInstruction
	// prenultInst is the second to last (penultimate) instruction emitted
	prenultInst EmittedInstruction
}

func New() *Compiler {
	return &Compiler{
		symbolTable: NewSymbolTable(),
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			if err := c.Compile(s); err != nil {
				return err
			}
		}

	case *ast.BlockStatement:
		for _, s := range node.Statements {
			if err := c.Compile(s); err != nil {
				return err
			}
		}

	case *ast.ExpressionStatement:
		if err := c.Compile(node.Expression); err != nil {
			return err
		}
		c.emit(code.OpPop)

	case *ast.Identifier:
		sym, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("could not resolve identifier: %s", node.Value)
		}
		c.emit(code.OpGetGlobal, sym.Index)

	case *ast.LetStatement:
		if err := c.Compile(node.Value); err != nil {
			return err
		}
		sym := c.symbolTable.Define(node.Name.Value)
		c.emit(code.OpSetGlobal, sym.Index)

	case *ast.DefineStatement:
		if err := c.Compile(node.Value); err != nil {
			return err
		}
		sym := c.symbolTable.Define(node.Name.Value)
		c.emit(code.OpSetGlobal, sym.Index)

		// expressions
	case *ast.InfixExpression:
		if node.Operator == "<" || node.Operator == "<=" {
			if err := c.Compile(node.Right); err != nil {
				return err
			}
			if err := c.Compile(node.Left); err != nil {
				return err
			}
			switch node.Operator {
			case "<":
				c.emit(code.OpGT)
			case "<=":
				c.emit(code.OpGTE)
			}
			return nil
		}
		if err := c.Compile(node.Left); err != nil {
			return err
		}
		if err := c.Compile(node.Right); err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMult)
		case "/":
			c.emit(code.OpDiv)
		case "%":
			c.emit(code.OpMod)
		case "==":
			c.emit(code.OpEQ)
		case "!=":
			c.emit(code.OpNEQ)
		case ">":
			c.emit(code.OpGT)
		case ">=":
			c.emit(code.OpGTE)

		default:
			return fmt.Errorf("unknown operator: %s", node.Operator)
		}

	case *ast.PrefixExpression:
		if err := c.Compile(node.Right); err != nil {
			return err
		}
		switch node.Operator {
		case "!":
			c.emit(code.OpBang)
		case "-":
			c.emit(code.OpMinus)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

		// Conditionals
	case *ast.IfExpression:
		if err := c.Compile(node.Condition); err != nil {
			return err
		}

		jmpNTPos := c.emit(code.OpJumpNotTruthy, 0)

		if err := c.Compile(node.Consequence); err != nil {
			return err
		}

		if c.ultInst.Opcode == code.OpPop {
			c.removeLastInstruction()
		}
		jmpPos := c.emit(code.OpJump, 0)

		c.replaceOperand(jmpNTPos, len(c.instructions))

		if node.Alternative == nil {
			c.emit(code.OpNull)
		} else {
			if err := c.Compile(node.Alternative); err != nil {
				return err
			}

			if c.ultInst.Opcode == code.OpPop {
				c.removeLastInstruction()
			}
		}

		c.replaceOperand(jmpPos, len(c.instructions))

		// Literals
	case *ast.IntegerLiteral:
		obj := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(obj))
	case *ast.FloatLiteral:
		obj := &object.Float{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(obj))
	case *ast.StringLiteral:
		obj := &object.String{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(obj))
	case *ast.BooleanLiteral:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}

	}
	return nil
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

func (c *Compiler) removeLastInstruction() {
	c.instructions = c.instructions[:c.ultInst.Position]
	c.ultInst = c.prenultInst
}

func (c *Compiler) replaceInstruction(pos int, inst code.Instructions) {
	for i, n := range inst {
		c.instructions[pos+i] = n
	}
}

func (c *Compiler) replaceOperand(pos int, operand int) {
	op := code.Opcode(c.instructions[pos])
	newInst := code.Make(op, operand)
	c.replaceInstruction(pos, newInst)
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)
	c.setUltInst(op, pos)
	return pos
}

func (c *Compiler) setUltInst(op code.Opcode, pos int) {
	c.prenultInst = c.ultInst
	c.ultInst = EmittedInstruction{Opcode: op, Position: pos}
}

func (c *Compiler) addInstruction(ins code.Instructions) int {
	pos := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return pos
}

func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}
