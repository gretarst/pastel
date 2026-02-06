package ast

// AssignStmt represents an assignment statement (name := value).
type AssignStmt struct {
	Name  string
	Value Expr
}

func (*AssignStmt) node()     {}
func (*AssignStmt) stmtNode() {}

// PrintStmt represents a writeln statement.
type PrintStmt struct {
	Argument Expr
}

func (*PrintStmt) node()     {}
func (*PrintStmt) stmtNode() {}

// CompoundStmt represents a begin...end block.
type CompoundStmt struct {
	Statements []Stmt
}

func (*CompoundStmt) node()     {}
func (*CompoundStmt) stmtNode() {}

// VarDecl represents a variable declaration.
type VarDecl struct {
	Name string
	Type string
}

func (*VarDecl) node()     {}
func (*VarDecl) stmtNode() {}

// Program represents a complete Pascal program.
type Program struct {
	Name         string
	Declarations []Stmt
	Main         *CompoundStmt
}

func (*Program) node()     {}
func (*Program) stmtNode() {}
