package ast

import "pastel/token"

// IntegerLiteral represents an integer literal value.
type IntegerLiteral struct {
	Value int
}

func (*IntegerLiteral) node()     {}
func (*IntegerLiteral) exprNode() {}

// RealLiteral represents a floating-point literal value.
type RealLiteral struct {
	Value float64
}

func (*RealLiteral) node()     {}
func (*RealLiteral) exprNode() {}

// BooleanLiteral represents a boolean literal value.
type BooleanLiteral struct {
	Value bool
}

func (*BooleanLiteral) node()     {}
func (*BooleanLiteral) exprNode() {}

// CharLiteral represents a single character literal.
type CharLiteral struct {
	Value rune
}

func (*CharLiteral) node()     {}
func (*CharLiteral) exprNode() {}

// StringLiteral represents a string literal value.
type StringLiteral struct {
	Value string
}

func (*StringLiteral) node()     {}
func (*StringLiteral) exprNode() {}

// BinaryExpr represents a binary expression (e.g., a + b).
type BinaryExpr struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (*BinaryExpr) node()     {}
func (*BinaryExpr) exprNode() {}

// Identifier represents a variable reference.
type Identifier struct {
	Value string
}

func (*Identifier) node()     {}
func (*Identifier) exprNode() {}
