package vm

import (
	"fmt"
	"math"

	"github.com/jimmykodes/joker/code"
	"github.com/jimmykodes/joker/compiler"
	"github.com/jimmykodes/joker/object"
)

const (
	GlobalSize     = math.MaxUint16
	StackSize      = 2048
	FrameStackSize = 1024
)

var Null = &object.Null{}

type VM struct {
	constants []object.Object
	globals   [GlobalSize]object.Object

	stack [StackSize]object.Object
	sp    int

	frames    [FrameStackSize]*Frame
	framesIdx int
}

func New(bytecode *compiler.Bytecode) *VM {
	vm := &VM{constants: bytecode.Constants}

	vm.pushFrame(NewFrame(bytecode.Instructions, 0))

	return vm
}

func (vm *VM) Run() error {
	var (
		ip  int
		ins code.Instructions
		op  code.Opcode
	)
	for {
		vm.currentFrame().ip++
		ip = vm.currentFrame().ip
		ins = vm.currentFrame().instructions
		if ip >= len(ins) {
			break
		}
		op = code.Opcode(ins[ip])
		switch op {
		// stack manipulation
		case code.OpConstant:
			constIdx := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2
			if err := vm.push(vm.constants[constIdx]); err != nil {
				return err
			}
		case code.OpPop:
			vm.pop()

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
			pos := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip = pos - 1
		case code.OpJumpNotTruthy:
			condition := vm.pop()
			if condition == object.False || condition == Null {
				pos := int(code.ReadUint16(ins[ip+1:]))
				vm.currentFrame().ip = pos - 1
			} else {
				vm.currentFrame().ip += 2
			}

			// variables
		case code.OpSetGlobal:
			idx := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2
			vm.globals[idx] = vm.pop()
		case code.OpGetGlobal:
			idx := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2
			if err := vm.push(vm.globals[idx]); err != nil {
				return err
			}

		case code.OpSetLocal:
			idx := code.ReadUint8(ins[ip+1:])
			fr := vm.currentFrame()
			fr.ip++
			vm.stack[fr.basePointer+int(idx)] = vm.pop()

		case code.OpGetLocal:
			idx := code.ReadUint8(ins[ip+1:])
			fr := vm.currentFrame()
			fr.ip++

			if err := vm.push(vm.stack[fr.basePointer+int(idx)]); err != nil {
				return err
			}

			// Composites
		case code.OpArray:
			numElems := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2

			elems := make([]object.Object, 0, numElems)
			for i := vm.sp - numElems; i < vm.sp; i++ {
				elems = append(elems, vm.stack[i])
			}

			vm.sp -= numElems

			if err := vm.push(&object.Array{Elements: elems}); err != nil {
				return err
			}
		case code.OpMap:
			numElems := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2

			pairs := make(map[object.HashKey]object.HashPair, numElems)
			for i := 0; i < numElems; i++ {
				offset := vm.sp - (i * 2)

				val := vm.stack[offset-1]
				key := vm.stack[offset-2]
				hashKey, ok := key.(object.Hashable)
				if !ok {
					return fmt.Errorf("invalid object on stack, %s is not hashable and cannot be used as a map key", key.Type())
				}
				pairs[hashKey.HashKey()] = object.HashPair{Key: key, Value: val}
			}
			vm.sp -= numElems * 2

			if err := vm.push(&object.Map{Pairs: pairs}); err != nil {
				return err
			}

			// Access
		case code.OpIndex:
			idx := vm.pop()
			obj := vm.pop()

			res, ok := obj.(object.Indexer)
			if !ok {
				return fmt.Errorf("invalid object on stack: %s is not indexable", obj.Type())
			}
			if err := vm.push(res.Idx(idx)); err != nil {
				return err
			}

			// Function
		case code.OpCall:
			numElems := int(code.ReadUint8(ins[ip+1:]))
			vm.currentFrame().ip++

			obj := vm.stack[vm.sp-1-numElems]
			res, ok := obj.(*object.CompiledFunction)
			if !ok {
				return fmt.Errorf("invalid object on stack: %s is not callable", obj.Type())
			}

			if numElems != res.NumParams {
				return fmt.Errorf("invalid number of args, got %d - want %d", numElems, res.NumParams)
			}

			fr := NewFrame(res.Instructions, vm.sp-numElems)
			vm.pushFrame(fr)
			vm.sp = fr.basePointer + res.NumLocals

		case code.OpReturn:
			val := vm.pop()
			fr := vm.popFrame()
			vm.sp = fr.basePointer - 1
			if err := vm.push(val); err != nil {
				return err
			}

		default:
			return fmt.Errorf("invalid op: %q", op)
		}
	}
	return nil
}

func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.framesIdx-1]
}

func (vm *VM) pushFrame(f *Frame) {
	vm.frames[vm.framesIdx] = f
	vm.framesIdx++
}

func (vm *VM) popFrame() *Frame {
	f := vm.currentFrame()
	vm.frames[vm.framesIdx-1] = nil
	vm.framesIdx--
	return f
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
