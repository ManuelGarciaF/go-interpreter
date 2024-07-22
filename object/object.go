package object

import (
	"fmt"
	"strings"

	"github.com/ManuelGarciaF/go-interpreter/ast"
)

type ObjectType int

const (
	INTEGER_OBJ ObjectType = iota
	BOOLEAN_OBJ
	NULL_OBJ
	RETURN_VALUE_OBJ
	FUNCTION_OBJ
	ERROR_OBJ
)

// For pretty printing the enum values
var objectTypeStrings = map[ObjectType]string{
	INTEGER_OBJ:      "INTEGER",
	BOOLEAN_OBJ:      "BOOLEAN",
	NULL_OBJ:         "NULL",
	RETURN_VALUE_OBJ: "RETURN_VALUE",
	FUNCTION_OBJ:     "FUNCTION",
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

func (*Function) Type() ObjectType   { return FUNCTION_OBJ }
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

type Error struct {
	Message string
}

func (*Error) Type() ObjectType  { return ERROR_OBJ }
func (e *Error) Inspect() string { return "ERROR: " + e.Message }
