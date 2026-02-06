package lexer

import (
	"pastel/token"
	"testing"
)

func TestNextToken_SingleCharacterTokens(t *testing.T) {
	input := `+-*/();:,.`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.STAR, "*"},
		{token.SLASH, "/"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.COLON, ":"},
		{token.COMMA, ","},
		{token.DOT, "."},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_MultiCharacterTokens(t *testing.T) {
	input := `:=`

	l := New(input)
	tok := l.NextToken()

	if tok.Type != token.ASSIGN {
		t.Fatalf("tokentype wrong. expected=%q, got=%q", token.ASSIGN, tok.Type)
	}

	if tok.Literal != ":=" {
		t.Fatalf("literal wrong. expected=%q, got=%q", ":=", tok.Literal)
	}
}

func TestNextToken_Keywords(t *testing.T) {
	input := `program var begin end writeln integer`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.PROGRAM, "program"},
		{token.VAR, "var"},
		{token.BEGIN, "begin"},
		{token.END, "end"},
		{token.WRITELN, "writeln"},
		{token.INTEGER, "integer"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_CaseInsensitiveKeywords(t *testing.T) {
	input := `PROGRAM Program pROGRAM`

	l := New(input)

	for i := 0; i < 3; i++ {
		tok := l.NextToken()
		if tok.Type != token.PROGRAM {
			t.Fatalf("expected PROGRAM, got %q", tok.Type)
		}
		if tok.Literal != "program" {
			t.Fatalf("expected lowercase 'program', got %q", tok.Literal)
		}
	}
}

func TestNextToken_Identifiers(t *testing.T) {
	input := `myVar x y1 counter`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENT, "myvar"},
		{token.IDENT, "x"},
		{token.IDENT, "y1"},
		{token.IDENT, "counter"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_Integers(t *testing.T) {
	input := `123 0 42 999`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INT, "123"},
		{token.INT, "0"},
		{token.INT, "42"},
		{token.INT, "999"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_WhitespaceHandling(t *testing.T) {
	input := "  \t\n  x  \n\n  y  "

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENT, "x"},
		{token.IDENT, "y"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_LineAndColumnTracking(t *testing.T) {
	input := `program test;
var x: integer;
begin
  x := 42;
end.`

	tests := []struct {
		expectedType   token.TokenType
		expectedLine   int
		expectedColumn int
	}{
		{token.PROGRAM, 1, 1},
		{token.IDENT, 1, 9}, // "test"
		{token.SEMICOLON, 1, 13},
		{token.VAR, 2, 1},
		{token.IDENT, 2, 5}, // "x"
		{token.COLON, 2, 6},
		{token.INTEGER, 2, 8},
		{token.SEMICOLON, 2, 15},
		{token.BEGIN, 3, 1},
		{token.IDENT, 4, 3},  // "x"
		{token.ASSIGN, 4, 5}, // ":="
		{token.INT, 4, 8},    // "42"
		{token.SEMICOLON, 4, 10},
		{token.END, 5, 1},
		{token.DOT, 5, 4},
		{token.EOF, 5, 5},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Line != tt.expectedLine {
			t.Fatalf("tests[%d] (%s) - line wrong. expected=%d, got=%d",
				i, tok.Type, tt.expectedLine, tok.Line)
		}

		if tok.Column != tt.expectedColumn {
			t.Fatalf("tests[%d] (%s) - column wrong. expected=%d, got=%d",
				i, tok.Type, tt.expectedColumn, tok.Column)
		}
	}
}

func TestNextToken_CompleteProgram(t *testing.T) {
	input := `program hello;
var x: integer;
begin
  x := 10 + 5 * 2;
  writeln(x);
end.`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.PROGRAM, "program"},
		{token.IDENT, "hello"},
		{token.SEMICOLON, ";"},
		{token.VAR, "var"},
		{token.IDENT, "x"},
		{token.COLON, ":"},
		{token.INTEGER, "integer"},
		{token.SEMICOLON, ";"},
		{token.BEGIN, "begin"},
		{token.IDENT, "x"},
		{token.ASSIGN, ":="},
		{token.INT, "10"},
		{token.PLUS, "+"},
		{token.INT, "5"},
		{token.STAR, "*"},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},
		{token.WRITELN, "writeln"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.END, "end"},
		{token.DOT, "."},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_IllegalCharacter(t *testing.T) {
	input := `@`

	l := New(input)
	tok := l.NextToken()

	if tok.Type != token.ILLEGAL {
		t.Fatalf("expected ILLEGAL token, got %q", tok.Type)
	}
}

func TestNextToken_RealLiterals(t *testing.T) {
	input := `3.14 0.5 123.456 42.0`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.REAL_LIT, "3.14"},
		{token.REAL_LIT, "0.5"},
		{token.REAL_LIT, "123.456"},
		{token.REAL_LIT, "42.0"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_StringAndCharLiterals(t *testing.T) {
	input := `'a' 'hello' 'world' 'x'`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.CHAR_LIT, "a"},
		{token.STRING_LIT, "hello"},
		{token.STRING_LIT, "world"},
		{token.CHAR_LIT, "x"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_BooleanKeywords(t *testing.T) {
	input := `true false TRUE FALSE True`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.TRUE, "true"},
		{token.FALSE, "false"},
		{token.TRUE, "true"},
		{token.FALSE, "false"},
		{token.TRUE, "true"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_TypeKeywords(t *testing.T) {
	input := `integer real boolean char string`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INTEGER, "integer"},
		{token.REAL, "real"},
		{token.BOOLEAN, "boolean"},
		{token.CHAR, "char"},
		{token.STRING, "string"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_IntegerNotReal(t *testing.T) {
	input := `42.method`

	l := New(input)

	tok := l.NextToken()
	if tok.Type != token.INT {
		t.Fatalf("expected INT, got %q", tok.Type)
	}
	if tok.Literal != "42" {
		t.Fatalf("expected '42', got %q", tok.Literal)
	}

	tok = l.NextToken()
	if tok.Type != token.DOT {
		t.Fatalf("expected DOT, got %q", tok.Type)
	}

	tok = l.NextToken()
	if tok.Type != token.IDENT {
		t.Fatalf("expected IDENT, got %q", tok.Type)
	}
}
