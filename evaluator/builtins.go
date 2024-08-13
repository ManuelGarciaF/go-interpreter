package evaluator

import "github.com/ManuelGarciaF/go-interpreter/object"

var builtins = map[string]*object.Builtin{
	"len": {Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1",
				len(args))
		}

		switch arg := args[0].(type) {
		case *object.String:
			return nativeToIntegerObject(len(arg.Value))
		default:
			return newError("argument to `len` not supported, got %s", arg.Type())
		}
	}},
}
