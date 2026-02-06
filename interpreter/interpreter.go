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
			i.env.Set(v.Name, defaultValue(v.Type))
		}
	}

	if err := i.evalStmt(prog.Main); err != nil {
		return err
	}

	return nil
}

func defaultValue(typeName string) Value {
	switch typeName {
	case "integer":
		return &IntegerValue{Val: 0}
	case "real":
		return &RealValue{Val: 0.0}
	case "boolean":
		return &BooleanValue{Val: false}
	case "char":
		return &CharValue{Val: ' '}
	case "string":
		return &StringValue{Val: ""}
	default:
		return &IntegerValue{Val: 0}
	}
}

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
		fmt.Println(val.String())

	default:
		return &PascalError{
			Msg:    "Unknown statement type",
			Detail: fmt.Sprintf("Encountered an unsupported statement: %T", stmt),
			Hint:   "Ensure all statements are valid Pascal constructs.",
		}
	}

	return nil
}

func (i *Interpreter) evalExpr(expr ast.Expr) (Value, error) {
	switch e := expr.(type) {
	case *ast.IntegerLiteral:
		return &IntegerValue{Val: e.Value}, nil

	case *ast.RealLiteral:
		return &RealValue{Val: e.Value}, nil

	case *ast.BooleanLiteral:
		return &BooleanValue{Val: e.Value}, nil

	case *ast.CharLiteral:
		return &CharValue{Val: e.Value}, nil

	case *ast.StringLiteral:
		return &StringValue{Val: e.Value}, nil

	case *ast.BinaryExpr:
		left, err := i.evalExpr(e.Left)
		if err != nil {
			return nil, err
		}

		right, err := i.evalExpr(e.Right)
		if err != nil {
			return nil, err
		}

		return i.evalBinaryOp(e.Operator, left, right)

	case *ast.Identifier:
		val, ok := i.env.Get(e.Value)
		if !ok {
			return nil, &PascalError{
				Msg:    fmt.Sprintf("Undefined variable '%s'", e.Value),
				Detail: "This variable is being used but was never declared or assigned a value.",
				Hint:   fmt.Sprintf("Declare the variable using `var %s: integer;` and assign it a value before use.", e.Value),
			}
		}
		return val, nil

	default:
		return nil, &PascalError{
			Msg:    "Unknown expression type",
			Detail: fmt.Sprintf("Encountered an unsupported expression: %T", expr),
			Hint:   "Ensure all expressions are valid Pascal constructs.",
		}
	}
}

func (i *Interpreter) evalBinaryOp(op token.Token, left, right Value) (Value, error) {
	switch op.Type {
	case token.PLUS:
		return i.evalPlus(left, right)
	case token.MINUS:
		return i.evalMinus(left, right)
	case token.STAR:
		return i.evalStar(left, right)
	case token.SLASH:
		return i.evalSlash(left, right)
	default:
		return nil, &PascalError{
			Msg:    "Unknown operator",
			Detail: fmt.Sprintf("Operator '%s' is not supported.", op.Literal),
			Hint:   "Use valid operators such as +, -, *, or /.",
		}
	}
}

func (i *Interpreter) evalPlus(left, right Value) (Value, error) {
	switch l := left.(type) {
	case *IntegerValue:
		switch r := right.(type) {
		case *IntegerValue:
			return &IntegerValue{Val: l.Val + r.Val}, nil
		case *RealValue:
			return &RealValue{Val: float64(l.Val) + r.Val}, nil
		}
	case *RealValue:
		switch r := right.(type) {
		case *IntegerValue:
			return &RealValue{Val: l.Val + float64(r.Val)}, nil
		case *RealValue:
			return &RealValue{Val: l.Val + r.Val}, nil
		}
	case *StringValue:
		switch r := right.(type) {
		case *StringValue:
			return &StringValue{Val: l.Val + r.Val}, nil
		case *CharValue:
			return &StringValue{Val: l.Val + string(r.Val)}, nil
		}
	case *CharValue:
		switch r := right.(type) {
		case *StringValue:
			return &StringValue{Val: string(l.Val) + r.Val}, nil
		case *CharValue:
			return &StringValue{Val: string(l.Val) + string(r.Val)}, nil
		}
	}
	return nil, &PascalError{
		Msg:    "Type mismatch in addition",
		Detail: fmt.Sprintf("Cannot add %s and %s.", left.Type(), right.Type()),
		Hint:   "Ensure both operands are numeric types or both are strings/chars.",
	}
}

func (i *Interpreter) evalMinus(left, right Value) (Value, error) {
	switch l := left.(type) {
	case *IntegerValue:
		switch r := right.(type) {
		case *IntegerValue:
			return &IntegerValue{Val: l.Val - r.Val}, nil
		case *RealValue:
			return &RealValue{Val: float64(l.Val) - r.Val}, nil
		}
	case *RealValue:
		switch r := right.(type) {
		case *IntegerValue:
			return &RealValue{Val: l.Val - float64(r.Val)}, nil
		case *RealValue:
			return &RealValue{Val: l.Val - r.Val}, nil
		}
	}
	return nil, &PascalError{
		Msg:    "Type mismatch in subtraction",
		Detail: fmt.Sprintf("Cannot subtract %s from %s.", right.Type(), left.Type()),
		Hint:   "Ensure both operands are numeric types.",
	}
}

func (i *Interpreter) evalStar(left, right Value) (Value, error) {
	switch l := left.(type) {
	case *IntegerValue:
		switch r := right.(type) {
		case *IntegerValue:
			return &IntegerValue{Val: l.Val * r.Val}, nil
		case *RealValue:
			return &RealValue{Val: float64(l.Val) * r.Val}, nil
		}
	case *RealValue:
		switch r := right.(type) {
		case *IntegerValue:
			return &RealValue{Val: l.Val * float64(r.Val)}, nil
		case *RealValue:
			return &RealValue{Val: l.Val * r.Val}, nil
		}
	}
	return nil, &PascalError{
		Msg:    "Type mismatch in multiplication",
		Detail: fmt.Sprintf("Cannot multiply %s and %s.", left.Type(), right.Type()),
		Hint:   "Ensure both operands are numeric types.",
	}
}

func (i *Interpreter) evalSlash(left, right Value) (Value, error) {
	switch l := left.(type) {
	case *IntegerValue:
		switch r := right.(type) {
		case *IntegerValue:
			if r.Val == 0 {
				return nil, divisionByZeroError()
			}
			return &IntegerValue{Val: l.Val / r.Val}, nil
		case *RealValue:
			if r.Val == 0 {
				return nil, divisionByZeroError()
			}
			return &RealValue{Val: float64(l.Val) / r.Val}, nil
		}
	case *RealValue:
		switch r := right.(type) {
		case *IntegerValue:
			if r.Val == 0 {
				return nil, divisionByZeroError()
			}
			return &RealValue{Val: l.Val / float64(r.Val)}, nil
		case *RealValue:
			if r.Val == 0 {
				return nil, divisionByZeroError()
			}
			return &RealValue{Val: l.Val / r.Val}, nil
		}
	}
	return nil, &PascalError{
		Msg:    "Type mismatch in division",
		Detail: fmt.Sprintf("Cannot divide %s by %s.", left.Type(), right.Type()),
		Hint:   "Ensure both operands are numeric types.",
	}
}

func divisionByZeroError() *PascalError {
	return &PascalError{
		Msg:    "Division by zero",
		Detail: "An attempt was made to divide by zero.",
		Hint:   "Ensure the divisor is not zero before performing division.",
	}
}
