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
		case code.OpAdd:
			r := vm.pop()
			l := vm.pop()

			left, ok := l.(object.Adder)
			if !ok {
				return fmt.Errorf("invalid object on stack, %s does not implement add", l.Type())
			}
			if err := vm.push(left.Add(r)); err != nil {
				return err
			}
		case code.OpPop:
			vm.pop()
		}
	}
	return nil
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
