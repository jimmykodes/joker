package vm

import (
	"errors"
	"fmt"
	"math"

	"github.com/jimmykodes/joker/builtins"
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
	fn := &object.CompiledFunction{Instructions: bytecode.Instructions}
	vm.pushFrame(NewFrame(&object.Closure{Fn: fn}, 0))
	return vm
}

func (vm *VM) Run() error {
	if err := vm.run(); err != nil {
		return fmt.Errorf("vm: %w", err)
	}
	return nil
}

var errStop = errors.New("program complete")

func (vm *VM) run() error {
	for {
		if err := vm.ExecuteInstruction(); err != nil {
			if errors.Is(err, errStop) {
				return nil
			}
			return err
		}
	}
}

func (vm *VM) Debug() error {
	for {
		r := make([]byte, 1)
		_, err := fmt.Scanln(&r)
		if err != nil {
			return err
		}
		fmt.Println("-----")
		switch r[0] {
		case 'n':
			if err := vm.ExecuteInstruction(); err != nil {
				if errors.Is(err, errStop) {
					return nil
				}
				return err
			}
		case 's':
			for i := 0; i < vm.sp; i++ {
				if obj := vm.stack[i]; obj != nil {
					fmt.Println(obj.Inspect())
				} else {
					fmt.Println(nil)
				}
			}
		case 'g':
			i := 0
			for {
				if vm.globals[i] == nil {
					break
				}
				fmt.Println(i, "-", vm.globals[i].Inspect())
				i++
			}
		case 'i':
			fmt.Println(vm.currentFrame().ip + 1)
			fmt.Println(vm.currentFrame().Instructions())
		case 'c':
			for i, constant := range vm.constants {
				fmt.Println(i, "-", constant.Inspect())
			}
		case 'h':
			fallthrough
		default:
			fmt.Println("")
		}
	}
}

func (vm *VM) ExecuteInstruction() error {
	var (
		ip  int
		ins code.Instructions
		op  code.Opcode
	)
	vm.currentFrame().ip++
	ip = vm.currentFrame().ip
	ins = vm.currentFrame().Instructions()
	if ip >= len(ins) {
		return errStop
	}
	op = code.Opcode(ins[ip])

	switch op {
	// stack manipulation
	case code.OpConstant:
		constIdx := code.ReadUint16(ins[ip+1:])
		vm.currentFrame().ip += 2
		if err := vm.push(vm.constants[constIdx]); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	case code.OpPop:
		vm.pop()

		// infix
	case code.OpAdd, code.OpSub, code.OpMult, code.OpDiv, code.OpMod, code.OpEQ, code.OpNEQ, code.OpGT, code.OpGTE:
		if err := vm.executeBinaryOperation(op); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		// prefix
	case code.OpBang, code.OpMinus:
		if err := vm.executePrefixOperator(op); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		// bools
	case code.OpTrue:
		if err := vm.push(object.True); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	case code.OpFalse:
		if err := vm.push(object.False); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	case code.OpNull:
		if err := vm.push(Null); err != nil {
			return fmt.Errorf("%s: %w", op, err)
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
			return fmt.Errorf("%s: %w", op, err)
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
			return fmt.Errorf("%s: %w", op, err)
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
			return fmt.Errorf("%s: %w", op, err)
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
			return fmt.Errorf("%s: %w", op, err)
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
			return fmt.Errorf("%s: %w", op, err)
		}

		// Function
	case code.OpCall:
		numElems := int(code.ReadUint8(ins[ip+1:]))
		vm.currentFrame().ip++

		obj := vm.stack[vm.sp-1-numElems]
		switch obj := obj.(type) {
		case *object.Closure:
			if numElems != obj.Fn.NumParams {
				return fmt.Errorf("invalid number of args, got %d - want %d", numElems, obj.Fn.NumParams)
			}

			fr := NewFrame(obj, vm.sp-numElems)
			vm.pushFrame(fr)
			vm.sp = fr.basePointer + obj.Fn.NumLocals
		case *object.Builtin:
			args := vm.stack[vm.sp-numElems : vm.sp]
			vm.sp = vm.sp - 1 - numElems
			res := obj.Fn(args...)
			if res == nil {
				res = Null
			}
			if err := vm.push(res); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}

		default:
			return fmt.Errorf("%s: invalid object on stack: %s is not callable", op, obj.Type())
		}

	case code.OpGetBuiltin:
		builtin := int(code.ReadUint8(ins[ip+1:]))
		vm.currentFrame().ip++
		obj, ok := builtins.Func(builtin)
		if !ok {
			return fmt.Errorf("invalid builtin: %d", builtin)
		}
		if err := vm.push(obj); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

	case code.OpClosure:
		constIdx := code.ReadUint16(ins[ip+1:])
		numFree := int(code.ReadUint8(ins[ip+3:]))
		vm.currentFrame().ip += 3

		obj := vm.constants[constIdx]
		fn, ok := obj.(*object.CompiledFunction)
		if !ok {
			return fmt.Errorf("%s: invalid object on stack: %s is not callable", op, obj.Type())
		}
		free := make([]object.Object, numFree)
		for i := range free {
			free[i] = vm.stack[vm.sp-numFree+i]
		}
		vm.sp = vm.sp - numFree

		if err := vm.push(&object.Closure{Fn: fn, Free: free}); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

	case code.OpGetFree:
		freeIdx := code.ReadUint8(ins[ip+1:])
		vm.currentFrame().ip++

		currentClosure := vm.currentFrame().cl
		if err := vm.push(currentClosure.Free[freeIdx]); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

	case code.OpSetFree:
		freeIdx := code.ReadUint8(ins[ip+1:])
		vm.currentFrame().ip++

		currentClosure := vm.currentFrame().cl
		val := vm.pop()
		currentClosure.Free[freeIdx] = val

		// when setting a free variable, the original variable
		// has to be set in the outer frame, too. otherwise
		// it won't be changed. so in the compiler, we emit "Set<Scope>"
		// codes until we get to something _not_ FreeScope-d so here, we need
		// to, for each instruction that is SetFree, look back one frame and
		// and set the index, continuously until we hit a local scope var
		// TODO: might be nice to make this recursive rather than iterative?
		df := 1
		for {
			curIns := code.Opcode(ins[vm.currentFrame().ip])
			if curIns == code.OpSetLocal {
				idx := code.ReadUint8(ins[ip+1:])
				vm.currentFrame().ip++
				fr := vm.frames[vm.framesIdx-df]
				vm.stack[fr.basePointer+int(idx)] = val
				break
			} else if curIns == code.OpSetFree {
				idx := code.ReadUint8(ins[ip+1:])
				vm.currentFrame().ip++
				fr := vm.frames[vm.framesIdx-1-df]
				fr.cl.Free[idx] = val
				df++
			} else {
				break
			}
		}

	case code.OpReturn:
		val := vm.pop()
		fr := vm.popFrame()
		vm.sp = fr.basePointer - 1
		if err := vm.push(val); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

	default:
		return fmt.Errorf("invalid op: %q", op)
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
