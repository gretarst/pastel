package interpreter

import (
	"bytes"
	"io"
	"os"
	"pastel/lexer"
	"pastel/parser"
	"strings"
	"testing"
)

func TestInterpreter_VariableDeclarationAndAssignment(t *testing.T) {
	input := `program test;
var x: integer;
begin
  x := 42;
end.`

	_, err := runProgram(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInterpreter_Addition(t *testing.T) {
	input := `program test;
var x: integer;
begin
  x := 5 + 3;
  writeln(x);
end.`

	output, err := runProgram(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "8\n"
	if output != expected {
		t.Fatalf("output wrong. expected=%q, got=%q", expected, output)
	}
}

func TestInterpreter_Subtraction(t *testing.T) {
	input := `program test;
var x: integer;
begin
  x := 10 - 3;
  writeln(x);
end.`

	output, err := runProgram(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "7\n"
	if output != expected {
		t.Fatalf("output wrong. expected=%q, got=%q", expected, output)
	}
}

func TestInterpreter_Multiplication(t *testing.T) {
	input := `program test;
var x: integer;
begin
  x := 4 * 5;
  writeln(x);
end.`

	output, err := runProgram(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "20\n"
	if output != expected {
		t.Fatalf("output wrong. expected=%q, got=%q", expected, output)
	}
}

func TestInterpreter_Division(t *testing.T) {
	input := `program test;
var x: integer;
begin
  x := 20 / 4;
  writeln(x);
end.`

	output, err := runProgram(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "5\n"
	if output != expected {
		t.Fatalf("output wrong. expected=%q, got=%q", expected, output)
	}
}

func TestInterpreter_OperatorPrecedence(t *testing.T) {
	input := `program test;
var x: integer;
begin
  x := 2 + 3 * 4;
  writeln(x);
end.`

	output, err := runProgram(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 2 + (3 * 4) = 2 + 12 = 14
	expected := "14\n"
	if output != expected {
		t.Fatalf("output wrong. expected=%q, got=%q", expected, output)
	}
}

func TestInterpreter_Parentheses(t *testing.T) {
	input := `program test;
var x: integer;
begin
  x := (2 + 3) * 4;
  writeln(x);
end.`

	output, err := runProgram(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// (2 + 3) * 4 = 5 * 4 = 20
	expected := "20\n"
	if output != expected {
		t.Fatalf("output wrong. expected=%q, got=%q", expected, output)
	}
}

func TestInterpreter_VariableInExpression(t *testing.T) {
	input := `program test;
var x: integer;
var y: integer;
begin
  x := 10;
  y := x + 5;
  writeln(y);
end.`

	output, err := runProgram(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "15\n"
	if output != expected {
		t.Fatalf("output wrong. expected=%q, got=%q", expected, output)
	}
}

func TestInterpreter_MultipleStatements(t *testing.T) {
	input := `program test;
var x: integer;
var y: integer;
begin
  x := 5;
  y := 10;
  x := x + y;
  writeln(x);
end.`

	output, err := runProgram(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "15\n"
	if output != expected {
		t.Fatalf("output wrong. expected=%q, got=%q", expected, output)
	}
}

func TestInterpreter_ComplexExpression(t *testing.T) {
	input := `program test;
var result: integer;
begin
  result := (10 + 5) * 2 - 6 / 3;
  writeln(result);
end.`

	output, err := runProgram(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// (10 + 5) * 2 - 6 / 3 = 15 * 2 - 2 = 30 - 2 = 28
	expected := "28\n"
	if output != expected {
		t.Fatalf("output wrong. expected=%q, got=%q", expected, output)
	}
}

func TestInterpreter_DivisionByZero(t *testing.T) {
	input := `program test;
var x: integer;
begin
  x := 10 / 0;
end.`

	_, err := runProgram(input)
	if err == nil {
		t.Fatalf("expected division by zero error, got none")
	}

	if !strings.Contains(err.Error(), "Division by zero") {
		t.Fatalf("expected 'Division by zero' error, got: %v", err)
	}
}

func TestInterpreter_UndeclaredVariable(t *testing.T) {
	input := `program test;
begin
  x := 10;
end.`

	_, err := runProgram(input)
	if err == nil {
		t.Fatalf("expected undeclared variable error, got none")
	}

	if !strings.Contains(err.Error(), "Undeclared variable") {
		t.Fatalf("expected 'Undeclared variable' error, got: %v", err)
	}
}

func TestInterpreter_UndefinedVariable(t *testing.T) {
	input := `program test;
var x: integer;
begin
  writeln(y);
end.`

	_, err := runProgram(input)
	if err == nil {
		t.Fatalf("expected undefined variable error, got none")
	}

	if !strings.Contains(err.Error(), "Undefined variable") {
		t.Fatalf("expected 'Undefined variable' error, got: %v", err)
	}
}

func TestInterpreter_DefaultVariableValue(t *testing.T) {
	input := `program test;
var x: integer;
begin
  writeln(x);
end.`

	output, err := runProgram(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Variables should default to 0
	expected := "0\n"
	if output != expected {
		t.Fatalf("output wrong. expected=%q, got=%q", expected, output)
	}
}

func TestInterpreter_MultipleWriteln(t *testing.T) {
	input := `program test;
var a: integer;
var b: integer;
begin
  a := 1;
  b := 2;
  writeln(a);
  writeln(b);
  writeln(a + b);
end.`

	output, err := runProgram(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "1\n2\n3\n"
	if output != expected {
		t.Fatalf("output wrong. expected=%q, got=%q", expected, output)
	}
}

// runProgram parses and executes a Pascal program, returning its output
func runProgram(input string) (string, error) {
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if p.HasErrors() {
		return "", p.Errors()[0]
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	interp := New()
	err := interp.Run(prog)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String(), err
}
