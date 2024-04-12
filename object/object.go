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
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
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
func (f FunctionObj) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n")

	return out.String()
}

// ---------------------------------

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING }
func (s *String) Inspect() string  { return s.Value }
