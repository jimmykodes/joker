package builtins

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jimmykodes/joker/object"
)

func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func nArgs(n int, args []object.Object) *object.Error {
	if len(args) != n {
		return newError("invalid number of args: got %d - want %d", len(args), n)
	}
	return nil
}

type builtin int

//go:generate stringer -type builtin -linecomment
const (
	start  builtin = iota
	Int            // int
	Float          // float
	String         // string
	Len            // len
	Pop            // pop
	Print          // print
	Append         // append
	Set            // set
	Slice          // slice
	Argv           // argv
	end
)

var lookups map[string]builtin

func init() {
	lookups = make(map[string]builtin, end)
	for i := start + 1; i < end; i++ {
		lookups[i.String()] = i
	}
}

func Func(i int) (*object.Builtin, bool) {
	b := builtin(i)
	if start >= b || b >= end {
		return nil, false
	}
	return builtins[b], true
}

func Lookup(name string) (int, bool) {
	val, ok := lookups[name]
	return int(val), ok
}

func LookupFunc(name string) (*object.Builtin, bool) {
	val, ok := lookups[name]
	return builtins[val], ok
}

var builtins = [...]*object.Builtin{
	Int: {
		Name: Int.String(),
		Fn: func(args ...object.Object) object.Object {
			if errOb := nArgs(1, args); errOb != nil {
				return errOb
			}
			switch a := args[0].(type) {
			case *object.Integer:
				return a
			case *object.Float:
				return &object.Integer{Value: int64(a.Value)}
			case *object.String:
				i, err := strconv.ParseInt(a.Value, 10, 64)
				if err != nil {
					if i, err := strconv.ParseFloat(a.Value, 64); err == nil {
						// if we got an error trying to parse it as an integer, attempt it as a float
						// and then cast to an int.
						return &object.Integer{Value: int64(i)}
					}
					return newError("invalid input")
				}
				return &object.Integer{Value: i}
			default:
				return object.ErrUnsupportedType
			}
		},
	},
	Float: {
		Name: Float.String(),
		Fn: func(args ...object.Object) object.Object {
			if errOb := nArgs(1, args); errOb != nil {
				return errOb
			}
			switch a := args[0].(type) {
			case *object.Integer:
				return &object.Float{Value: float64(a.Value)}
			case *object.Float:
				return a
			case *object.String:
				i, err := strconv.ParseFloat(a.Value, 64)
				if err != nil {
					return newError("invalid input")
				}
				return &object.Float{Value: i}
			default:
				return object.ErrUnsupportedType
			}
		},
	},
	String: {
		Name: String.String(),
		Fn: func(args ...object.Object) object.Object {
			if errOb := nArgs(1, args); errOb != nil {
				return errOb
			}
			switch a := args[0].(type) {
			case *object.Integer:
				return &object.String{Value: strconv.FormatInt(a.Value, 10)}
			case *object.Float:
				return &object.String{Value: fmt.Sprintf("%v", a.Value)}
			case *object.String:
				return a
			default:
				return object.ErrUnsupportedType
			}
		},
	},
	Len: {
		Name: Len.String(),
		Fn: func(args ...object.Object) object.Object {
			if errOb := nArgs(1, args); errOb != nil {
				return errOb
			}
			l, ok := args[0].(object.Lenner)
			if !ok {
				return newError("len() not supported on %s", args[0].Type())
			}
			return l.Len()
		},
	},
	Pop: {
		Name: Pop.String(),
		// TODO: create a "popable" interface and have this compare the object to the interface
		// implement the interface on both maps and slices
		Fn: func(args ...object.Object) object.Object {
			if err := nArgs(2, args); err != nil {
				return err
			}
			m, ok := args[0].(*object.Map)
			if !ok {
				return newError("invalid type for pop. got %s, want %s", args[0].Type(), object.MapType)
			}
			k, ok := args[1].(object.Hashable)
			if !ok {
				return newError("invalid key type")
			}
			obj, ok := m.Pairs[k.HashKey()]
			if !ok {
				return nil
			}
			delete(m.Pairs, k.HashKey())
			return obj.Value
		},
	},
	Print: {
		Name: Print.String(),
		Fn: func(args ...object.Object) object.Object {
			out := make([]any, len(args))
			for i, arg := range args {
				out[i] = arg.Inspect()
			}
			fmt.Fprintln(os.Stdout, out...)
			return nil
		},
	},
	Append: {
		Name: Append.String(),
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 2 {
				return newError("invalid number of args, got %d, want 2+", len(args))
			}
			source, ok := args[0].(*object.Array)
			if !ok {
				return newError("first argument of append must be an %s", object.ArrayType)
			}
			return &object.Array{Elements: append(source.Elements, args[1:]...)}
		},
	},
	Set: {
		Name: Set.String(),
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 3 {
				return newError("invalid number of args, got %d, want 3", len(args))
			}
			obj := args[0]

			settable, ok := obj.(object.Settable)
			if !ok {
				return newError("invalid object: %T is not Settable", obj)
			}
			return settable.Set(args[1], args[2])
		},
	},
	Slice: {
		Name: Slice.String(),
		Fn: func(args ...object.Object) object.Object {
			var (
				source object.Object
				start  int64
				end    int64
			)
			switch len(args) {
			case 0, 1:
				return newError("invalid number of args, got %d, want 2+", len(args))
			case 2:
				source = args[0]
				if args[1].Type() != object.IntegerType {
					return newError("cannot slice using type %s, must be %s", args[1].Type(), object.IntegerType)
				}
				end = args[1].(*object.Integer).Value
			case 3:
				source = args[0]
				if args[1].Type() != object.IntegerType || args[2].Type() != object.IntegerType {
					return newError("cannot slice using type %s, must be %s", args[1].Type(), object.IntegerType)
				}
				start = args[1].(*object.Integer).Value
				end = args[2].(*object.Integer).Value
			}
			if start < 0 {
				return newError("starting point of slice cannot be negative")
			}

			switch src := source.(type) {
			case *object.Array:
				if start > int64(len(src.Elements)) || end > int64(len(src.Elements)) {
					return newError("index out of range [%d] with length %d", end, len(src.Elements))
				}
				return &object.Array{Elements: src.Elements[start:end]}
			case *object.String:
				if start > int64(len(src.Value)) || end > int64(len(src.Value)) {
					return newError("index out of range [%d] with length %d", end, len(src.Value))
				}
				return &object.String{Value: src.Value[start:end]}
			default:
				return newError("invalid source for slice, must be %s or %s", object.ArrayType, object.StringType)
			}
		},
	},
	Argv: {
		Name: Argv.String(),
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 0 {
				return newError("invalid number of args, got %d, want 0", len(args))
			}
			argv := os.Args
			elements := make([]object.Object, len(os.Args))
			for i, arg := range argv {
				elements[i] = &object.String{Value: arg}
			}
			return &object.Array{Elements: elements}
		},
	},
}
