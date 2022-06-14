package evaluator

import (
	"errors"
	"fmt"

	"github.com/jimmykodes/jk/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("invalid number of args, got %d - want 1", len(args))
			}
			l, err := args[0].Len()
			if errors.Is(err, object.ErrUnsupportedOperation) {
				return newError("len() not supported on %s", args[0].Type())
			} else if err != nil {
				return newError(err.Error())
			}
			return l
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
}
