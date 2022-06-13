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
			var v int64
			switch a := args[0].(type) {
			case *object.String:
				v = int64(len(a.Value))
			default:
				return newError("cannot call len() on %s", a.Type())
			}
			return &object.Integer{Value: v}
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
}
