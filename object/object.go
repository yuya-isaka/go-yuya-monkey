package object

import "fmt"

const (
	NULL   = "NULL"
	INT    = "INT"
	BOOL   = "BOOL"
	RETURN = "RETURN"
	ERROR  = "ERROR"
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
