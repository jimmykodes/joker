package vm

import (
	"fmt"

	"github.com/jimmykodes/joker/code"
	"github.com/jimmykodes/joker/compiler"
	"github.com/jimmykodes/joker/object"
)

const StackSize = 2048

var Null = &object.Null{}

type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack [StackSize]object.Object
	sp    int
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		constants:    bytecode.Constants,
		instructions: bytecode.Instructions,
	}
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])
		switch op {
		case code.OpConstant:
			constIdx := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			if err := vm.push(vm.constants[constIdx]); err != nil {
				return err
			}

			// infix
		case code.OpAdd, code.OpSub, code.OpMult, code.OpDiv, code.OpMod, code.OpEQ, code.OpNEQ, code.OpGT, code.OpGTE:
			if err := vm.executeBinaryOperation(op); err != nil {
				return err
			}

			// prefix
		case code.OpBang, code.OpMinus:
			if err := vm.executePrefixOperator(op); err != nil {
				return err
			}

			// bools
		case code.OpTrue:
			if err := vm.push(object.True); err != nil {
				return err
			}
		case code.OpFalse:
			if err := vm.push(object.False); err != nil {
				return err
			}
		case code.OpNull:
			if err := vm.push(Null); err != nil {
				return err
			}

			// jumps
		case code.OpJump:
			pos := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip = pos - 1
		case code.OpJumpNotTruthy:
			condition := vm.pop()
			if condition == object.False {
				pos := int(code.ReadUint16(vm.instructions[ip+1:]))
				ip = pos - 1
			} else {
				ip += 2
			}

		case code.OpPop:
			vm.pop()

		default:
			return fmt.Errorf("invalid op: %q", op)
		}
	}
	return nil
}

func (vm *VM) executeBinaryOperation(op code.Opcode) error {
	r, l := vm.pop(), vm.pop()
	var res object.Object
	switch op {
	case code.OpAdd:
		left, ok := l.(object.Adder)
		if !ok {
			return fmt.Errorf("invalid object on stack, %s does not implement add", l.Type())
		}
		res = left.Add(r)
	case code.OpSub:
		left, ok := l.(object.Subber)
		if !ok {
			return fmt.Errorf("invalid object on stack, %s does not implement sub", l.Type())
		}
		res = left.Sub(r)
	case code.OpMult:
		left, ok := l.(object.MultDiver)
		if !ok {
			return fmt.Errorf("invalid object on stack, %s does not implement multiplication", l.Type())
		}
		res = left.Mult(r)
	case code.OpDiv:
		left, ok := l.(object.MultDiver)
		if !ok {
			return fmt.Errorf("invalid object on stack, %s does not implement division", l.Type())
		}
		res = left.Div(r)
	case code.OpMod:
		left, ok := l.(object.Modder)
		if !ok {
			return fmt.Errorf("invalid object on stack, %s does not implement modular division", l.Type())
		}
		res = left.Mod(r)
	case code.OpEQ:
		left, ok := l.(object.Equal)
		if !ok {
			return fmt.Errorf("invalid object on stack, %s does not implement equality", l.Type())
		}
		res = left.EQ(r)
	case code.OpNEQ:
		left, ok := l.(object.Equal)
		if !ok {
			return fmt.Errorf("invalid object on stack, %s does not implement inequality", l.Type())
		}
		res = left.NEQ(r)
	case code.OpGT:
		left, ok := l.(object.Inequality)
		if !ok {
			return fmt.Errorf("invalid object on stack, %s does not implement comparison", l.Type())
		}
		res = left.GT(r)
	case code.OpGTE:
		left, ok := l.(object.Inequality)
		if !ok {
			return fmt.Errorf("invalid object on stack, %s does not implement comparison", l.Type())
		}
		res = left.GTE(r)
	default:
		return fmt.Errorf("invalid op: %q", op)

	}
	return vm.push(res)
}

func (vm *VM) executePrefixOperator(op code.Opcode) error {
	r := vm.pop()
	var res object.Object
	switch op {
	case code.OpMinus:
		right, ok := r.(object.Negater)
		if !ok {
			return fmt.Errorf("invalid object on stack, %s does not implement negation", r.Type())
		}
		res = right.Negative()
	case code.OpBang:
		right, ok := r.(object.Booler)
		if !ok {
			return fmt.Errorf("invalid object on stack, %s does not implement ! inversion", r.Type())
		}
		res = right.Bool().Invert()
	default:
		return fmt.Errorf("invalid op: %q", op)
	}
	return vm.push(res)
}

func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}

	return vm.stack[vm.sp-1]
}

func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp]
}

func (vm *VM) push(obj object.Object) error {
	if vm.sp > StackSize {
		return fmt.Errorf("stack overflow")
	}
	vm.stack[vm.sp] = obj
	vm.sp++
	return nil
}

func (vm *VM) pop() object.Object {
	o := vm.stack[vm.sp-1]
	vm.sp--
	return o
}
