package ast

import "pastel/token"

// IntegerLiteral represents an integer literal value.
type IntegerLiteral struct {
	Value int
}

func (*IntegerLiteral) node()     {}
func (*IntegerLiteral) exprNode() {}

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
