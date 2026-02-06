package interpreter

import "fmt"

// ValueType represents the type of a Value.
type ValueType string

const (
	IntegerType ValueType = "integer"
	RealType    ValueType = "real"
	BooleanType ValueType = "boolean"
	CharType    ValueType = "char"
	StringType  ValueType = "string"
)

// Value represents a runtime value in the interpreter.
type Value interface {
	Type() ValueType
	String() string
}

// IntegerValue holds an integer value.
type IntegerValue struct{ Val int }

func (v *IntegerValue) Type() ValueType { return IntegerType }
func (v *IntegerValue) String() string  { return fmt.Sprintf("%d", v.Val) }

// RealValue holds a floating-point value.
type RealValue struct{ Val float64 }

func (v *RealValue) Type() ValueType { return RealType }
func (v *RealValue) String() string  { return fmt.Sprintf("%g", v.Val) }

// BooleanValue holds a boolean value.
type BooleanValue struct{ Val bool }

func (v *BooleanValue) Type() ValueType { return BooleanType }
func (v *BooleanValue) String() string {
	if v.Val {
		return "true"
	}
	return "false"
}

// CharValue holds a single character.
type CharValue struct{ Val rune }

func (v *CharValue) Type() ValueType { return CharType }
func (v *CharValue) String() string  { return string(v.Val) }

// StringValue holds a string value.
type StringValue struct{ Val string }

func (v *StringValue) Type() ValueType { return StringType }
func (v *StringValue) String() string  { return v.Val }
