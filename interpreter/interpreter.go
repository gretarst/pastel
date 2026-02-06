package interpreter

import (
	"fmt"
	"pastel/ast"
	"pastel/token"
)

// Interpreter holds the state for program execution.
type Interpreter struct {
	env *Environment
}

// New creates a new Interpreter instance with a fresh environment.
func New() *Interpreter {
	return &Interpreter{env: NewEnvironment()}
}

// Run executes a Pascal program.
func (i *Interpreter) Run(prog *ast.Program) error {
	for _, decl := range prog.Declarations {
		if v, ok := decl.(*ast.VarDecl); ok {
			i.env.Set(v.Name, 0) // default value is 0
		}
	}

	if err := i.evalStmt(prog.Main); err != nil {
		return err
	}

	return nil
}

// evalStmt evaluates a single statement.
func (i *Interpreter) evalStmt(stmt ast.Stmt) error {
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		if !i.env.Exists(s.Name) {
			return &PascalError{
				Msg:    fmt.Sprintf("Undeclared variable '%s'", s.Name),
				Detail: "This variable is being used but was never declared with a type.",
				Hint:   fmt.Sprintf("Try adding `var %s: integer;` at the top of your program.", s.Name),
			}
		}

		val, err := i.evalExpr(s.Value)
		if err != nil {
			return err
		}
		i.env.Set(s.Name, val)

	case *ast.CompoundStmt:
		for _, stmt := range s.Statements {
			if err := i.evalStmt(stmt); err != nil {
				return err
			}
		}

	case *ast.PrintStmt:
		val, err := i.evalExpr(s.Argument)
		if err != nil {
			return err
		}
		fmt.Println(val)

	default:
		return &PascalError{
			Msg:    "Unknown statement type",
			Detail: fmt.Sprintf("Encountered an unsupported statement: %T", stmt),
			Hint:   "Ensure all statements are valid Pascal constructs.",
		}
	}

	return nil
}

// evalExpr evaluates an expression and returns its integer value.
func (i *Interpreter) evalExpr(expr ast.Expr) (int, error) {
	switch e := expr.(type) {
	case *ast.IntegerLiteral:
		return e.Value, nil

	case *ast.BinaryExpr:
		left, err := i.evalExpr(e.Left)
		if err != nil {
			return 0, err
		}

		right, err := i.evalExpr(e.Right)
		if err != nil {
			return 0, err
		}

		switch e.Operator.Type {
		case token.PLUS:
			return left + right, nil
		case token.MINUS:
			return left - right, nil
		case token.STAR:
			return left * right, nil
		case token.SLASH:
			if right == 0 {
				return 0, &PascalError{
					Msg:    "Division by zero",
					Detail: "An attempt was made to divide by zero.",
					Hint:   "Ensure the divisor is not zero before performing division.",
				}
			}
			return left / right, nil
		default:
			return 0, &PascalError{
				Msg:    "Unknown operator",
				Detail: fmt.Sprintf("Operator '%s' is not supported.", e.Operator.Literal),
				Hint:   "Use valid operators such as +, -, *, or /.",
			}
		}

	case *ast.Identifier:
		val, ok := i.env.Get(e.Value)
		if !ok {
			return 0, &PascalError{
				Msg:    fmt.Sprintf("Undefined variable '%s'", e.Value),
				Detail: "This variable is being used but was never declared or assigned a value.",
				Hint:   fmt.Sprintf("Declare the variable using `var %s: integer;` and assign it a value before use.", e.Value),
			}
		}
		return val, nil

	default:
		return 0, &PascalError{
			Msg:    "Unknown expression type",
			Detail: fmt.Sprintf("Encountered an unsupported expression: %T", expr),
			Hint:   "Ensure all expressions are valid Pascal constructs.",
		}
	}
}
