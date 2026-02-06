package parser

import (
	"pastel/ast"
	"pastel/lexer"
	"testing"
)

func TestParseProgram_MinimalProgram(t *testing.T) {
	input := `program test;
begin
end.`

	l := lexer.New(input)
	p := New(l)
	prog := p.ParseProgram()

	checkParserErrors(t, p)

	if prog == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if prog.Name != "test" {
		t.Fatalf("program name wrong. expected=%q, got=%q", "test", prog.Name)
	}

	if len(prog.Declarations) != 0 {
		t.Fatalf("expected 0 declarations, got %d", len(prog.Declarations))
	}

	if prog.Main == nil {
		t.Fatalf("expected main compound statement, got nil")
	}
}

func TestParseProgram_WithVariable(t *testing.T) {
	input := `program test;
var x: integer;
begin
end.`

	l := lexer.New(input)
	p := New(l)
	prog := p.ParseProgram()

	checkParserErrors(t, p)

	if prog == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(prog.Declarations) != 1 {
		t.Fatalf("expected 1 declaration, got %d", len(prog.Declarations))
	}

	varDecl, ok := prog.Declarations[0].(*ast.VarDecl)
	if !ok {
		t.Fatalf("expected *ast.VarDecl, got %T", prog.Declarations[0])
	}

	if varDecl.Name != "x" {
		t.Fatalf("variable name wrong. expected=%q, got=%q", "x", varDecl.Name)
	}

	if varDecl.Type != "integer" {
		t.Fatalf("variable type wrong. expected=%q, got=%q", "integer", varDecl.Type)
	}
}

func TestParseProgram_WithAssignment(t *testing.T) {
	input := `program test;
var x: integer;
begin
  x := 42;
end.`

	l := lexer.New(input)
	p := New(l)
	prog := p.ParseProgram()

	checkParserErrors(t, p)

	if prog == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(prog.Main.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Main.Statements))
	}

	assignStmt, ok := prog.Main.Statements[0].(*ast.AssignStmt)
	if !ok {
		t.Fatalf("expected *ast.AssignStmt, got %T", prog.Main.Statements[0])
	}

	if assignStmt.Name != "x" {
		t.Fatalf("assignment name wrong. expected=%q, got=%q", "x", assignStmt.Name)
	}

	intLit, ok := assignStmt.Value.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expected *ast.IntegerLiteral, got %T", assignStmt.Value)
	}

	if intLit.Value != 42 {
		t.Fatalf("integer value wrong. expected=%d, got=%d", 42, intLit.Value)
	}
}

func TestParseProgram_WithWriteln(t *testing.T) {
	input := `program test;
var x: integer;
begin
  x := 10;
  writeln(x);
end.`

	l := lexer.New(input)
	p := New(l)
	prog := p.ParseProgram()

	checkParserErrors(t, p)

	if prog == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(prog.Main.Statements) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(prog.Main.Statements))
	}

	printStmt, ok := prog.Main.Statements[1].(*ast.PrintStmt)
	if !ok {
		t.Fatalf("expected *ast.PrintStmt, got %T", prog.Main.Statements[1])
	}

	ident, ok := printStmt.Argument.(*ast.Identifier)
	if !ok {
		t.Fatalf("expected *ast.Identifier, got %T", printStmt.Argument)
	}

	if ident.Value != "x" {
		t.Fatalf("identifier value wrong. expected=%q, got=%q", "x", ident.Value)
	}
}

func TestParseExpression_IntegerLiteral(t *testing.T) {
	input := `program test;
var x: integer;
begin
  x := 123;
end.`

	l := lexer.New(input)
	p := New(l)
	prog := p.ParseProgram()

	checkParserErrors(t, p)

	assignStmt := prog.Main.Statements[0].(*ast.AssignStmt)
	intLit, ok := assignStmt.Value.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expected *ast.IntegerLiteral, got %T", assignStmt.Value)
	}

	if intLit.Value != 123 {
		t.Fatalf("integer value wrong. expected=%d, got=%d", 123, intLit.Value)
	}
}

func TestParseExpression_BinaryExpressionAddition(t *testing.T) {
	input := `program test;
var x: integer;
begin
  x := 5 + 3;
end.`

	l := lexer.New(input)
	p := New(l)
	prog := p.ParseProgram()

	checkParserErrors(t, p)

	assignStmt := prog.Main.Statements[0].(*ast.AssignStmt)
	binExpr, ok := assignStmt.Value.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", assignStmt.Value)
	}

	leftInt, ok := binExpr.Left.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expected left to be *ast.IntegerLiteral, got %T", binExpr.Left)
	}
	if leftInt.Value != 5 {
		t.Fatalf("left value wrong. expected=%d, got=%d", 5, leftInt.Value)
	}

	if binExpr.Operator.Literal != "+" {
		t.Fatalf("operator wrong. expected=%q, got=%q", "+", binExpr.Operator.Literal)
	}

	rightInt, ok := binExpr.Right.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expected right to be *ast.IntegerLiteral, got %T", binExpr.Right)
	}
	if rightInt.Value != 3 {
		t.Fatalf("right value wrong. expected=%d, got=%d", 3, rightInt.Value)
	}
}

func TestParseExpression_OperatorPrecedence(t *testing.T) {
	input := `program test;
var x: integer;
begin
  x := 2 + 3 * 4;
end.`

	l := lexer.New(input)
	p := New(l)
	prog := p.ParseProgram()

	checkParserErrors(t, p)

	assignStmt := prog.Main.Statements[0].(*ast.AssignStmt)
	binExpr, ok := assignStmt.Value.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", assignStmt.Value)
	}

	// Should be parsed as: 2 + (3 * 4)
	// So the top-level operator is +
	if binExpr.Operator.Literal != "+" {
		t.Fatalf("top-level operator wrong. expected=%q, got=%q", "+", binExpr.Operator.Literal)
	}

	// Left should be integer 2
	leftInt, ok := binExpr.Left.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expected left to be *ast.IntegerLiteral, got %T", binExpr.Left)
	}
	if leftInt.Value != 2 {
		t.Fatalf("left value wrong. expected=%d, got=%d", 2, leftInt.Value)
	}

	// Right should be binary expr (3 * 4)
	rightBin, ok := binExpr.Right.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected right to be *ast.BinaryExpr, got %T", binExpr.Right)
	}
	if rightBin.Operator.Literal != "*" {
		t.Fatalf("right operator wrong. expected=%q, got=%q", "*", rightBin.Operator.Literal)
	}
}

func TestParseExpression_Parentheses(t *testing.T) {
	input := `program test;
var x: integer;
begin
  x := (2 + 3) * 4;
end.`

	l := lexer.New(input)
	p := New(l)
	prog := p.ParseProgram()

	checkParserErrors(t, p)

	assignStmt := prog.Main.Statements[0].(*ast.AssignStmt)
	binExpr, ok := assignStmt.Value.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", assignStmt.Value)
	}

	// Should be parsed as: (2 + 3) * 4
	// So the top-level operator is *
	if binExpr.Operator.Literal != "*" {
		t.Fatalf("top-level operator wrong. expected=%q, got=%q", "*", binExpr.Operator.Literal)
	}

	// Left should be binary expr (2 + 3)
	leftBin, ok := binExpr.Left.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected left to be *ast.BinaryExpr, got %T", binExpr.Left)
	}
	if leftBin.Operator.Literal != "+" {
		t.Fatalf("left operator wrong. expected=%q, got=%q", "+", leftBin.Operator.Literal)
	}

	// Right should be integer 4
	rightInt, ok := binExpr.Right.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expected right to be *ast.IntegerLiteral, got %T", binExpr.Right)
	}
	if rightInt.Value != 4 {
		t.Fatalf("right value wrong. expected=%d, got=%d", 4, rightInt.Value)
	}
}

func TestParseExpression_IdentifierInExpression(t *testing.T) {
	input := `program test;
var x: integer;
var y: integer;
begin
  x := 10;
  y := x + 5;
end.`

	l := lexer.New(input)
	p := New(l)
	prog := p.ParseProgram()

	checkParserErrors(t, p)

	assignStmt := prog.Main.Statements[1].(*ast.AssignStmt)
	binExpr, ok := assignStmt.Value.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", assignStmt.Value)
	}

	leftIdent, ok := binExpr.Left.(*ast.Identifier)
	if !ok {
		t.Fatalf("expected left to be *ast.Identifier, got %T", binExpr.Left)
	}
	if leftIdent.Value != "x" {
		t.Fatalf("left identifier wrong. expected=%q, got=%q", "x", leftIdent.Value)
	}
}

func TestParserErrors_MissingProgramKeyword(t *testing.T) {
	input := `test;
begin
end.`

	l := lexer.New(input)
	p := New(l)
	p.ParseProgram()

	if !p.HasErrors() {
		t.Fatalf("expected parser errors, got none")
	}
}

func TestParserErrors_MissingSemicolon(t *testing.T) {
	input := `program test
begin
end.`

	l := lexer.New(input)
	p := New(l)
	p.ParseProgram()

	if !p.HasErrors() {
		t.Fatalf("expected parser errors, got none")
	}
}

func TestParserErrors_MissingDot(t *testing.T) {
	input := `program test;
begin
end`

	l := lexer.New(input)
	p := New(l)
	p.ParseProgram()

	if !p.HasErrors() {
		t.Fatalf("expected parser errors, got none")
	}
}

func TestParserErrors_ErrorLineAndColumn(t *testing.T) {
	input := `program test
begin
end.`

	l := lexer.New(input)
	p := New(l)
	p.ParseProgram()

	if !p.HasErrors() {
		t.Fatalf("expected parser errors, got none")
	}

	err := p.Errors()[0]
	if err.Line == 0 {
		t.Fatalf("expected error line to be set, got 0")
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	if !p.HasErrors() {
		return
	}

	t.Errorf("parser has %d errors", len(p.Errors()))
	for _, err := range p.Errors() {
		t.Errorf("parser error: %s", err.Error())
	}
	t.FailNow()
}
