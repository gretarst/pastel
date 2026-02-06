package ast

// Node is the base interface for all AST nodes.
type Node interface {
	node()
}

// Expr is the interface for all expression nodes.
type Expr interface {
	Node
	exprNode()
}

// Stmt is the interface for all statement nodes.
type Stmt interface {
	Node
	stmtNode()
}
