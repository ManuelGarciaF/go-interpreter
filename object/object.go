package object

import (
	"fmt"
	"strings"

	"github.com/ManuelGarciaF/go-interpreter/ast"
)

type ObjectType int

const (
	INTEGER_OBJ ObjectType = iota
	STRING_OBJ
	ARRAY_OBJ
	BOOLEAN_OBJ
	NULL_OBJ
	RETURN_VALUE_OBJ
	FUNCTION_OBJ
	BUILTIN_OBJ
	ERROR_OBJ
)

// For pretty printing the enum values
var objectTypeStrings = map[ObjectType]string{
	INTEGER_OBJ:      "INTEGER",
	STRING_OBJ:       "STRING",
	ARRAY_OBJ:        "ARRAY",
	BOOLEAN_OBJ:      "BOOLEAN",
	NULL_OBJ:         "NULL",
	RETURN_VALUE_OBJ: "RETURN_VALUE",
	FUNCTION_OBJ:     "FUNCTION",
	BUILTIN_OBJ:      "BUILTIN",
	ERROR_OBJ:        "ERROR",
}

func (o ObjectType) String() string {
	return objectTypeStrings[o]
}

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (*Integer) Type() ObjectType  { return INTEGER_OBJ }
func (i *Integer) Inspect() string { return fmt.Sprint(i.Value) }

type String struct {
	Value string
}

func (*String) Type() ObjectType  { return STRING_OBJ }
func (s *String) Inspect() string { return fmt.Sprint("\"" + s.Value + "\"") }

type Array struct {
	Elements []Object
}

func (*Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var sb strings.Builder

	elements := make([]string, 0, len(a.Elements))
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	sb.WriteByte('[')
	sb.WriteString(strings.Join(elements, ", "))
	sb.WriteByte(']')

	return sb.String()
}

type Boolean struct {
	Value bool
}

func (*Boolean) Type() ObjectType  { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

type Null struct{}

func (*Null) Type() ObjectType  { return NULL_OBJ }
func (b *Null) Inspect() string { return "null" }

// Wraps an object
type ReturnValue struct {
	Value Object
}

func (*ReturnValue) Type() ObjectType   { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (*Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var sb strings.Builder

	params := make([]string, 0, len(f.Parameters))
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	sb.WriteString("fn(")
	sb.WriteString(strings.Join(params, ", "))
	sb.WriteString(") {\n")
	sb.WriteString(f.Body.String())
	sb.WriteString("\n}")

	return sb.String()
}

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (*Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (*Builtin) Inspect() string  { return "builtin function" }

type Error struct {
	Message string
}

func (*Error) Type() ObjectType  { return ERROR_OBJ }
func (e *Error) Inspect() string { return "ERROR: " + e.Message }
