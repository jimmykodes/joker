package vm

import (
	"fmt"

	"github.com/jimmykodes/joker/code"
	"github.com/jimmykodes/joker/compiler"
	"github.com/jimmykodes/joker/object"
)

const StackSize = 2048

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
		case code.OpAdd, code.OpSub, code.OpMult, code.OpDiv, code.OpMod:
			r := vm.pop()
			l := vm.pop()
			if err := vm.executeBinaryOperation(op, l, r); err != nil {
				return err
			}
		case code.OpTrue:
			if err := vm.push(object.True); err != nil {
				return err
			}
		case code.OpFalse:
			if err := vm.push(object.False); err != nil {
				return err
			}

		case code.OpPop:
			vm.pop()
		}
	}
	return nil
}

func (vm *VM) executeBinaryOperation(op code.Opcode, l, r object.Object) error {
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
