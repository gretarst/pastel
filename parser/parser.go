package parser

import (
	"fmt"
	"pastel/ast"
	"pastel/lexer"
	"pastel/token"
	"strconv"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []*ParserError
}

// New creates a new Parser instance with the given lexer.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) HasErrors() bool {
	return len(p.errors) > 0
}

func (p *Parser) Errors() []*ParserError {
	return p.errors
}

// ParseExpression parses an expression in Pascal.
// Expressions include arithmetic operations like addition, subtraction, multiplication, and division.
func (p *Parser) ParseExpression() ast.Expr {
	return p.parseAddition()
}

// ParseProgram parses a complete Pascal program.
// A Pascal program starts with the 'program' keyword, followed by declarations and a main compound statement.
func (p *Parser) ParseProgram() *ast.Program {
	prog := &ast.Program{}

	if p.curToken.Type == token.PROGRAM {
		// Advance to the next token after 'program' keyword
		p.nextToken()
	} else {
		p.addError(
			"Expected 'program' keyword",
			fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
			"A Pascal program must start with the 'program' keyword.",
		)
		return nil
	}

	if p.curToken.Type != token.IDENT {
		p.addError(
			"Expected program name",
			fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
			"The 'program' keyword must be followed by an identifier.",
		)
		return nil
	}

	// Advance to the next token after the program name
	prog.Name = p.curToken.Literal
	p.nextToken()

	if p.curToken.Type != token.SEMICOLON {
		p.addError(
			"Expected semicolon",
			fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
			"Statements must end with a semicolon.",
		)
		return nil
	}

	// Advance to the next token after the semicolon
	p.nextToken()

	var decls []ast.Stmt
	for p.curToken.Type == token.VAR {
		decl := p.parseVarDecl()
		if decl != nil {
			decls = append(decls, decl)
		}
	}
	prog.Declarations = decls

	if p.curToken.Type != token.BEGIN {
		p.addError(
			"Expected 'begin' block",
			fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
			"A Pascal program must have a 'begin' block to define its main body.",
		)
		return nil
	}

	// Parse the compound statement starting with 'begin'
	stmt := p.parseCompound()
	compound, ok := stmt.(*ast.CompoundStmt)
	if !ok {
		p.addError(
			"Expected compound statement",
			"The main body of the program must be a compound statement.",
			"Ensure the program's main body starts with 'begin' and ends with 'end'.",
		)
		return nil
	}

	prog.Main = compound

	if p.curToken.Type != token.DOT {
		p.addError(
			"Expected '.' at the end of the program",
			fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
			"A Pascal program must end with a period ('.').",
		)
		return nil
	}

	return prog
}

// ParseStatement parses a single Pascal statement.
// Statements include assignments, compound statements, and print statements.
func (p *Parser) parseStatement() ast.Stmt {
	switch p.curToken.Type {
	case token.IDENT:
		// Look ahead to see if this is an assignment (IDENT := ...)
		if p.peekToken.Type == token.ASSIGN {
			return p.parseAssignment()
		}
		p.addError(
			fmt.Sprintf("Unexpected identifier '%s'", p.curToken.Literal),
			"This identifier is not part of an assignment or recognized statement.",
			"Make sure you're using ':=' for assignments or a known keyword like 'writeln'.",
		)
		p.nextToken()
		return nil

	case token.WRITELN:
		return p.parsePrint()

	case token.BEGIN:
		return p.parseCompound()

	default:
		return nil
	}
}

// ParseAssignment parses an assignment statement in Pascal.
// Assignment statements use the ':=' operator to assign values to variables.
func (p *Parser) parseAssignment() ast.Stmt {
	name := p.curToken.Literal // We are on IDENT

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	value := p.ParseExpression()

	if !p.curTokenIs(token.SEMICOLON) {
		p.addError(
			"Expected semicolon at the end of assignment",
			fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
			"Assignments must end with a semicolon.",
		)
		return nil
	}

	return &ast.AssignStmt{Name: name, Value: value}
}

// ParseCompound parses a compound statement in Pascal.
// Compound statements start with 'begin', contain multiple statements, and end with 'end'.
func (p *Parser) parseCompound() ast.Stmt {
	stmts := []ast.Stmt{}

	p.nextToken()

	for p.curToken.Type != token.END && p.curToken.Type != token.EOF && p.curToken.Type != token.DOT {
		stmt := p.parseStatement()
		if stmt != nil {
			stmts = append(stmts, stmt)
		}

		p.nextToken()
	}

	// Advance past 'end' token
	if p.curToken.Type == token.END {
		p.nextToken()
	}

	return &ast.CompoundStmt{Statements: stmts}
}

// ParsePrint parses a print statement in Pascal.
// Print statements use the 'writeln' keyword to output values.
func (p *Parser) parsePrint() ast.Stmt {
	if !p.curTokenIs(token.WRITELN) {
		p.addError(
			"Expected 'writeln' keyword",
			fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
			"Use 'writeln' to print values.",
		)
		return nil
	}

	// Advance to the next token after 'writeln'
	p.nextToken()

	if !p.curTokenIs(token.LPAREN) {
		p.addError(
			"Expected '(' after 'writeln'",
			fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
			"The 'writeln' keyword must be followed by parentheses containing the argument.",
		)
		return nil
	}

	// Advance to the next token after '('
	p.nextToken()

	arg := p.ParseExpression()

	if !p.curTokenIs(token.RPAREN) {
		p.addError(
			"Expected ')' after writeln argument",
			fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
			"Ensure the argument to 'writeln' is enclosed in parentheses.",
		)
		return nil
	}

	// Advance to the next token after ')'
	p.nextToken()

	if !p.curTokenIs(token.SEMICOLON) {
		p.addError(
			"Expected ';' after writeln",
			fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
			"Statements must end with a semicolon.",
		)
		return nil
	}

	return &ast.PrintStmt{Argument: arg}
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) addError(msg, detail, hint string) {
	p.errors = append(p.errors, &ParserError{
		Msg:    msg,
		Detail: detail,
		Hint:   hint,
		Line:   p.curToken.Line,
		Column: p.curToken.Column,
	})
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	p.errors = append(p.errors, &ParserError{
		Msg:    fmt.Sprintf("Expected next token to be %s", t),
		Detail: fmt.Sprintf("Got %q (%s) instead.", p.peekToken.Literal, p.peekToken.Type),
		Hint:   "Check the syntax of your program.",
		Line:   p.peekToken.Line,
		Column: p.peekToken.Column,
	})
	return false
}

func (p *Parser) parseAddition() ast.Expr {
	left := p.parseMultiplication()

	for p.curTokenIs(token.PLUS) || p.curTokenIs(token.MINUS) {
		op := p.curToken
		p.nextToken()
		right := p.parseMultiplication()
		left = &ast.BinaryExpr{Left: left, Operator: op, Right: right}
	}

	return left
}

func (p *Parser) parseMultiplication() ast.Expr {
	left := p.parsePrimary()

	for p.curTokenIs(token.STAR) || p.curTokenIs(token.SLASH) {
		op := p.curToken
		p.nextToken()
		right := p.parsePrimary()
		left = &ast.BinaryExpr{Left: left, Operator: op, Right: right}
	}

	return left
}

func (p *Parser) parsePrimary() ast.Expr {
	switch p.curToken.Type {
	case token.LPAREN:
		p.nextToken() // Advance from '(' to first token inside

		expr := p.ParseExpression()

		if !p.curTokenIs(token.RPAREN) {
			p.addError(
				"Expected closing parenthesis",
				fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
				"Ensure all opening parentheses have matching closing parentheses.",
			)
			return nil
		}

		p.nextToken() // Consume ')'
		return expr

	case token.INT:
		val, _ := strconv.Atoi(p.curToken.Literal)
		lit := &ast.IntegerLiteral{Value: val}
		p.nextToken()
		return lit

	case token.REAL_LIT:
		val, _ := strconv.ParseFloat(p.curToken.Literal, 64)
		lit := &ast.RealLiteral{Value: val}
		p.nextToken()
		return lit

	case token.TRUE:
		lit := &ast.BooleanLiteral{Value: true}
		p.nextToken()
		return lit

	case token.FALSE:
		lit := &ast.BooleanLiteral{Value: false}
		p.nextToken()
		return lit

	case token.CHAR_LIT:
		lit := &ast.CharLiteral{Value: rune(p.curToken.Literal[0])}
		p.nextToken()
		return lit

	case token.STRING_LIT:
		lit := &ast.StringLiteral{Value: p.curToken.Literal}
		p.nextToken()
		return lit

	case token.IDENT:
		ident := &ast.Identifier{Value: p.curToken.Literal}
		p.nextToken()
		return ident

	default:
		p.addError(
			"Unexpected token in primary expression",
			fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
			"Check the syntax of your expression.",
		)
		return nil
	}
}

func (p *Parser) parseVarDecl() ast.Stmt {
	// Advance to the next token after 'var'
	p.nextToken()

	if p.curToken.Type != token.IDENT {
		p.addError(
			"Expected variable name after 'var'",
			fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
			"Variable declarations must start with a valid identifier.",
		)
		return nil
	}

	name := p.curToken.Literal

	// Advance to the next token after the variable name
	p.nextToken()

	if p.curToken.Type != token.COLON {
		p.addError(
			"Expected ':' after variable name",
			fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
			"Variable declarations must specify a type after the colon.",
		)
		return nil
	}

	// Advance to the next token after ':'
	p.nextToken()

	var varType string
	switch p.curToken.Type {
	case token.INTEGER, token.REAL, token.BOOLEAN, token.CHAR, token.STRING:
		varType = p.curToken.Literal
	default:
		p.addError(
			"Expected type for variable",
			fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
			"Supported types are: integer, real, boolean, char, string.",
		)
		return nil
	}

	// Advance to the next token after the type
	p.nextToken()

	if p.curToken.Type != token.SEMICOLON {
		p.addError(
			"Expected ';' after variable declaration",
			fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
			"Variable declarations must end with a semicolon.",
		)
		return nil
	}

	// Advance to the next token after the semicolon
	p.nextToken()

	return &ast.VarDecl{Name: name, Type: varType}
}
