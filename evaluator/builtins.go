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
		case *object.Array:
			return nativeToIntegerObject(len(arg.Elements))
		default:
			return newError("argument to `len` not supported, got %s", arg.Type())
		}
	}},
	"first": {Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1",
				len(args))
		}
		arr, ok := args[0].(*object.Array)
		if !ok {
			return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
		}

		if len(arr.Elements) == 0 {
			return NULL
		}

		return arr.Elements[0]

	}},
	"last": {Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1",
				len(args))
		}
		arr, ok := args[0].(*object.Array)
		if !ok {
			return newError("argument to `last` must be ARRAY, got %s", args[0].Type())
		}

		if len(arr.Elements) == 0 {
			return NULL
		}

		return arr.Elements[len(arr.Elements)-1]
	}},
	"tail": {Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1",
				len(args))
		}
		arr, ok := args[0].(*object.Array)
		if !ok {
			return newError("argument to `tail` must be ARRAY, got %s", args[0].Type())
		}

		length := len(arr.Elements)
		if length == 0 {
			return NULL
		}

		newElements := make([]object.Object, length-1)
		copy(newElements, arr.Elements[1:length])
		return &object.Array{Elements: newElements}
	}},
	"push": {Fn: func(args ...object.Object) object.Object {
		if len(args) != 2 {
			return newError("wrong number of arguments. got=%d, want=2",
				len(args))
		}
		arr, ok := args[0].(*object.Array)
		if !ok {
			return newError("first argument to `push` must be ARRAY, got %s", args[0].Type())
		}

		length := len(arr.Elements)
		newElements := make([]object.Object, length, length+1)
		copy(newElements, arr.Elements)
		newElements = append(newElements, args[1])

		return &object.Array{Elements: newElements}
	}},
}
