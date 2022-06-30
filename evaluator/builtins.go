package evaluator

import (
	"fmt"
	"strconv"

	"github.com/jimmykodes/joker/object"
)

func nArgs(n int, args []object.Object) object.Object {
	if len(args) != n {
		return newError("invalid number of args. got %d - want %d", len(args), n)
	}
	return nil
}

var builtins = map[string]*object.Builtin{
	"int": {
		Fn: func(_ *object.Environment, args ...object.Object) object.Object {
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
	"float": {
		Fn: func(_ *object.Environment, args ...object.Object) object.Object {
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
	"string": {
		Fn: func(_ *object.Environment, args ...object.Object) object.Object {
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
	"len": {
		Fn: func(_ *object.Environment, args ...object.Object) object.Object {
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
	"del": {
		Fn: func(_ *object.Environment, args ...object.Object) object.Object {
			if err := nArgs(2, args); err != nil {
				return err
			}
			m, ok := args[0].(*object.Map)
			if !ok {
				return newError("invalid type for del. got %s, want %s", args[0].Type(), object.MapType)
			}
			k, ok := args[1].(object.Hashable)
			if !ok {
				return newError("invalid key type")
			}
			delete(m.Pairs, k.HashKey())
			return Null
		},
	},
	"print": {
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			out := make([]any, len(args))
			for i, arg := range args {
				out[i] = arg.Inspect()
			}
			fmt.Fprintln(env.Out(), out...)
			return Null
		},
	},
	"append": {
		Fn: func(_ *object.Environment, args ...object.Object) object.Object {
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
	"slice": {
		Fn: func(_ *object.Environment, args ...object.Object) object.Object {
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
}
