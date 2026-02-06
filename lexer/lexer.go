package lexer

import (
	"pastel/token"
	"strings"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	line         int
	column       int
}

func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1

	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

func (l *Lexer) Ch() byte {
	return l.ch
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) newTokenWithPos(tokenType token.TokenType, ch byte, line, column int) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: line, Column: column}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return strings.ToLower(l.input[start:l.position])
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readNumber() (string, bool) {
	start := l.position
	isReal := false
	for isDigit(l.ch) {
		l.readChar()
	}
	if l.ch == '.' && isDigit(l.peekChar()) {
		isReal = true
		l.readChar()
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	return l.input[start:l.position], isReal
}

func (l *Lexer) readString() string {
	l.readChar()
	start := l.position
	for l.ch != '\'' && l.ch != 0 {
		l.readChar()
	}
	str := l.input[start:l.position]
	l.readChar()
	return str
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	// Store position before reading token
	line := l.line
	col := l.column

	// FIRST: letters (identifiers and keywords)
	if isLetter(l.ch) {
		literal := l.readIdentifier()
		tokType := token.LookupIdent(literal)
		return token.Token{Type: tokType, Literal: literal, Line: line, Column: col}
	}

	// SECOND: numbers
	if isDigit(l.ch) {
		literal, isReal := l.readNumber()
		if isReal {
			return token.Token{Type: token.REAL_LIT, Literal: literal, Line: line, Column: col}
		}
		return token.Token{Type: token.INT, Literal: literal, Line: line, Column: col}
	}

	// THIRD: string/char literals
	if l.ch == '\'' {
		literal := l.readString()
		if len(literal) == 1 {
			return token.Token{Type: token.CHAR_LIT, Literal: literal, Line: line, Column: col}
		}
		return token.Token{Type: token.STRING_LIT, Literal: literal, Line: line, Column: col}
	}

	var tok token.Token
	switch l.ch {
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.ASSIGN, Literal: literal, Line: line, Column: col}
		} else {
			tok = l.newTokenWithPos(token.COLON, l.ch, line, col)
		}
	case ';':
		tok = l.newTokenWithPos(token.SEMICOLON, l.ch, line, col)
	case ',':
		tok = l.newTokenWithPos(token.COMMA, l.ch, line, col)
	case '+':
		tok = l.newTokenWithPos(token.PLUS, l.ch, line, col)
	case '-':
		tok = l.newTokenWithPos(token.MINUS, l.ch, line, col)
	case '*':
		tok = l.newTokenWithPos(token.STAR, l.ch, line, col)
	case '/':
		tok = l.newTokenWithPos(token.SLASH, l.ch, line, col)
	case '(':
		tok = l.newTokenWithPos(token.LPAREN, l.ch, line, col)
	case ')':
		tok = l.newTokenWithPos(token.RPAREN, l.ch, line, col)
	case '.':
		tok = l.newTokenWithPos(token.DOT, l.ch, line, col)
	case 0:
		tok = token.Token{Type: token.EOF, Literal: "", Line: line, Column: col}
	default:
		tok = l.newTokenWithPos(token.ILLEGAL, l.ch, line, col)
	}
	l.readChar()
	return tok
}
