package compiler

import (
	"fmt"
	"strings"

	"github.com/jimmykodes/joker/ast"
	"github.com/jimmykodes/joker/code"
	"github.com/jimmykodes/joker/object"
)

type CompilationScope struct {
	instructions code.Instructions
	// ultInst is the last (ultimate) instruction emitted
	ultInst EmittedInstruction
	// penultInst is the second to last (penultimate) instruction emitted
	penultInst EmittedInstruction
	startPos   int
	setEndPos  []int
}

type Compiler struct {
	constants   []object.Object
	symbolTable *SymbolTable

	scopes   []CompilationScope
	scopeIdx int
}

func New() *Compiler {
	return &Compiler{
		symbolTable: NewSymbolTable(),
		scopes:      []CompilationScope{{}},
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
		if _, ok := node.Expression.(*ast.CommentLiteral); ok {
			return nil
		}
		if err := c.Compile(node.Expression); err != nil {
			return err
		}
		c.emit(code.OpPop)

	case *ast.Identifier:
		sym, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("could not resolve identifier: %s", node.Value)
		}
		switch sym.Scope {
		case GlobalScope:
			c.emit(code.OpGetGlobal, sym.Index)
		case LocalScope:
			c.emit(code.OpGetLocal, sym.Index)
		}

	case *ast.ReassignStatement:
		if err := c.Compile(node.Value); err != nil {
			return err
		}
		sym, ok := c.symbolTable.Resolve(node.Name.Value)
		if !ok {
			return fmt.Errorf("cannot resolve symbol %s", node.Name.Value)
		}
		switch sym.Scope {
		case GlobalScope:
			c.emit(code.OpSetGlobal, sym.Index)
		case LocalScope:
			c.emit(code.OpSetLocal, sym.Index)
		}

	case *ast.LetStatement:
		if err := c.Compile(node.Value); err != nil {
			return err
		}
		sym := c.symbolTable.Define(node.Name.Value)
		switch sym.Scope {
		case GlobalScope:
			c.emit(code.OpSetGlobal, sym.Index)
		case LocalScope:
			c.emit(code.OpSetLocal, sym.Index)
		}

	case *ast.DefineStatement:
		if err := c.Compile(node.Value); err != nil {
			return err
		}
		sym := c.symbolTable.Define(node.Name.Value)
		switch sym.Scope {
		case GlobalScope:
			c.emit(code.OpSetGlobal, sym.Index)
		case LocalScope:
			c.emit(code.OpSetLocal, sym.Index)
		}

	case *ast.FuncStatement:
		sym := c.symbolTable.Define(node.Name.Value)
		if err := c.Compile(node.Fn); err != nil {
			return err
		}
		switch sym.Scope {
		case GlobalScope:
			c.emit(code.OpSetGlobal, sym.Index)
		case LocalScope:
			c.emit(code.OpSetLocal, sym.Index)
		}

	// expressions
	case *ast.CallExpression:
		if err := c.Compile(node.Function); err != nil {
			return err
		}
		for _, arg := range node.Arguments {
			if err := c.Compile(arg); err != nil {
				return err
			}
		}
		c.emit(code.OpCall, len(node.Arguments))

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

	case *ast.IndexExpression:
		if err := c.Compile(node.Left); err != nil {
			return err
		}
		if err := c.Compile(node.Index); err != nil {
			return err
		}
		c.emit(code.OpIndex)

		// Conditionals
	case *ast.IfExpression:
		if err := c.Compile(node.Condition); err != nil {
			return err
		}

		jmpNTPos := c.emit(code.OpJumpNotTruthy, 0)

		if err := c.Compile(node.Consequence); err != nil {
			return err
		}

		c.removeLastInstruction(code.OpPop)
		jmpPos := c.emit(code.OpJump, 0)

		c.replaceOperand(jmpNTPos, len(c.scopes[c.scopeIdx].instructions))

		if node.Alternative == nil {
			c.emit(code.OpNull)
		} else {
			if err := c.Compile(node.Alternative); err != nil {
				return err
			}

			c.removeLastInstruction(code.OpPop)
		}

		c.replaceOperand(jmpPos, len(c.scopes[c.scopeIdx].instructions))

	case *ast.WhileExpression:
		oldStart := c.scopes[c.scopeIdx].startPos

		startPos := len(c.scopes[c.scopeIdx].instructions)
		c.scopes[c.scopeIdx].startPos = startPos

		if err := c.Compile(node.Condition); err != nil {
			return err
		}
		jntPos := c.emit(code.OpJumpNotTruthy, 0)

		if err := c.Compile(node.Body); err != nil {
			return err
		}

		c.emit(code.OpJump, startPos)
		endPos := len(c.scopes[c.scopeIdx].instructions)
		c.replaceOperand(jntPos, endPos)
		for _, setEndPos := range c.scopes[c.scopeIdx].setEndPos {
			c.replaceOperand(setEndPos, endPos)
		}
		c.emit(code.OpNull)
		c.scopes[c.scopeIdx].startPos = oldStart
		c.scopes[c.scopeIdx].setEndPos = nil

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
	case *ast.ArrayLiteral:
		for _, elem := range node.Elements {
			if err := c.Compile(elem); err != nil {
				return err
			}
		}
		c.emit(code.OpArray, len(node.Elements))
	case *ast.MapLiteral:
		for k, v := range node.Pairs {
			if err := c.Compile(k); err != nil {
				return err
			}
			if err := c.Compile(v); err != nil {
				return err
			}
		}
		c.emit(code.OpMap, len(node.Pairs))

	case *ast.FunctionLiteral:
		c.enterScope()

		for _, ident := range node.Parameters {
			c.symbolTable.Define(ident.Value)
		}

		if err := c.Compile(node.Body); err != nil {
			return err
		}

		if c.scopes[c.scopeIdx].ultInst.Opcode != code.OpReturn {
			c.emit(code.OpNull)
			c.emit(code.OpReturn)
		}
		numLocals := len(c.symbolTable.store)
		scope := c.leaveScope()

		cf := c.addConstant(&object.CompiledFunction{Instructions: scope.instructions, NumLocals: numLocals})
		c.emit(code.OpConstant, cf)

	case *ast.ReturnStatement:
		if c.scopeIdx == 0 {
			return fmt.Errorf("top level returns are not allowed")
		}
		if err := c.Compile(node.Value); err != nil {
			return err
		}
		c.emit(code.OpReturn)

	case *ast.BreakStatement:
		jmpPos := c.emit(code.OpJump, 0)
		c.scopes[c.scopeIdx].setEndPos = append(c.scopes[c.scopeIdx].setEndPos, jmpPos)

	case *ast.ContinueStatement:
		c.emit(code.OpJump, c.scopes[c.scopeIdx].startPos)

	default:
		return fmt.Errorf("unknown node: %T", node)

	}
	return nil
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.scopes[0].instructions,
		Constants:    c.constants,
	}
}

func (c *Compiler) enterScope() {
	c.symbolTable = NewSymbolTable(OuterSymbolTable(c.symbolTable))
	c.scopes = append(c.scopes, CompilationScope{})
	c.scopeIdx++
}

func (c *Compiler) leaveScope() CompilationScope {
	c.symbolTable = c.symbolTable.outer
	scope := c.scopes[c.scopeIdx]
	c.scopes = c.scopes[:c.scopeIdx]
	c.scopeIdx--
	return scope
}

func (c *Compiler) removeLastInstruction(op code.Opcode) {
	scope := c.scopes[c.scopeIdx]
	if scope.ultInst.Opcode != op {
		return
	}
	scope.instructions = scope.instructions[:scope.ultInst.Position]
	scope.ultInst = scope.penultInst
	c.scopes[c.scopeIdx] = scope
}

func (c *Compiler) replaceInstruction(pos int, inst code.Instructions) {
	for i, n := range inst {
		c.scopes[c.scopeIdx].instructions[pos+i] = n
	}
}

func (c *Compiler) replaceOperand(pos int, operand int) {
	op := code.Opcode(c.scopes[c.scopeIdx].instructions[pos])
	newInst := code.Instruction(op, operand)
	c.replaceInstruction(pos, newInst)
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Instruction(op, operands...)
	pos := c.addInstruction(ins)
	c.setUltInst(op, pos)
	return pos
}

func (c *Compiler) setUltInst(op code.Opcode, pos int) {
	c.scopes[c.scopeIdx].penultInst = c.scopes[c.scopeIdx].ultInst
	c.scopes[c.scopeIdx].ultInst = EmittedInstruction{Opcode: op, Position: pos}
}

func (c *Compiler) addInstruction(ins code.Instructions) int {
	pos := len(c.scopes[c.scopeIdx].instructions)
	c.scopes[c.scopeIdx].instructions = append(c.scopes[c.scopeIdx].instructions, ins...)
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

func (b Bytecode) String() string {
	var sb strings.Builder
	sb.WriteString("Constants:\n")
	for _, constant := range b.Constants {
		fmt.Fprintf(&sb, "\t%+v\n", constant)
	}
	sb.WriteString("Instructions:\n")
	sb.WriteString(b.Instructions.String())

	return sb.String()
}

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}
