package object

import "fmt"

const (
	NULL = "NULL"
	INT  = "INTEGER"
	BOOL = "BOOLEAN"
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
