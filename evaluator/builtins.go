package evaluator

import (
	"fmt"

	"github.com/jimmykodes/jk/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("invalid number of args, got %d - want 1", len(args))
			}
			l, ok := args[0].(object.Lenner)
			if !ok {
				return newError("len() not supported on %s", args[0].Type())
			}
			return l.Len()
		},
	},
	"print": {
		Fn: func(args ...object.Object) object.Object {
			out := make([]any, len(args))
			for i, arg := range args {
				out[i] = arg.Inspect()
			}
			fmt.Println(out...)
			return Null
		},
	},
	"append": {
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
	"slice": {
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
}
