package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/yuya-isaka/go-yuya-monkey/ast"
)

const (
	NULL     = "NULL"
	INT      = "INT"
	BOOL     = "BOOL"
	RETURN   = "RETURN"
	ERROR    = "ERROR"
	FUNCTION = "FUNCTION"
	STRING   = "STRING"
	BUILTIN  = "BUILTIN"
	ARRAY    = "ARRAY"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string // インタプリタで反復的に表示するやつ
}

// ---------------------------------

type IntObj struct {
	Value int64
}

func (i IntObj) Type() ObjectType { return INT }
func (i IntObj) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

// ---------------------------------

type BoolObj struct {
	Value bool
}

func (b BoolObj) Type() ObjectType { return BOOL }
func (b BoolObj) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

// ---------------------------------

type NullObj struct{}

func (n NullObj) Type() ObjectType { return NULL }
func (n NullObj) Inspect() string  { return "null" }

// ---------------------------------

type ReturnObj struct {
	Value Object
}

func (r ReturnObj) Type() ObjectType { return RETURN }
func (r ReturnObj) Inspect() string  { return r.Value.Inspect() }

// ---------------------------------

type ErrorObj struct {
	Value string
}

func (e ErrorObj) Type() ObjectType { return ERROR }
func (e ErrorObj) Inspect() string  { return "ERROR: " + e.Value }

// ---------------------------------

// 変数の内部表現はいらないんやなーー

// ---------------------------------

// Callの時に評価したいから、そのままノードを持っておかないといけない
type FunctionObj struct {
	Parameters []*ast.IdentNode
	Body       *ast.BlockNode
	Env        *Environment // クロージャだ
}

func (f FunctionObj) Type() ObjectType { return FUNCTION }

// 関数の表示だるう
func (f FunctionObj) Inspect() string {
	var out bytes.Buffer

	// 『var params []string』でも同じ
	// なぜなら、今回はappendで追加しているから
	// 他の使い方をするなら注意。
	params := make([]string, len(f.Parameters))
	for i, p := range f.Parameters {
		params[i] = p.String()
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n  ")
	out.WriteString(f.Body.String())
	out.WriteString("\n")
	out.WriteString("}")

	return out.String()
}

// ---------------------------------

type StringObj struct {
	Value string
}

func (s *StringObj) Type() ObjectType { return STRING }
func (s *StringObj) Inspect() string  { return s.Value }

// ---------------------------------

type BuiltinFunction func(args ...Object) Object

type BuiltinObj struct {
	Fn BuiltinFunction
}

func (b *BuiltinObj) Type() ObjectType { return BUILTIN }
func (b *BuiltinObj) Inspect() string  { return "builtin function" }

// ---------------------------------

type ArrayObj struct {
	Values []Object
}

func (a ArrayObj) Type() ObjectType { return ARRAY }
func (a ArrayObj) Inspect() string {
	var out bytes.Buffer

	values := make([]string, len(a.Values))
	for i, v := range a.Values {
		values[i] = v.Inspect()
	}

	out.WriteString("[")
	out.WriteString(strings.Join(values, ", "))
	out.WriteString("]")

	return out.String()
}
