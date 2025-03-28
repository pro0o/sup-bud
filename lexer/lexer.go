package lexer

import (
	"fmt"
	"strings"

	"github.com/pro0o/sup-bud/token"
)

type Lexer struct {
	input         string
	position      int      // current char
	readPosition  int      // next char
	ch            byte     // current char in consideration
	line          int      // current line number
	column        int      // current column number
	linePositions []int    // tracks starting position of each line
	errors        []string // collect lexical errors
}

func New(input string) *Lexer {
	l := &Lexer{
		input:         input,
		line:          1,
		column:        0,
		linePositions: []int{0}, // First line starts at position 0
		errors:        []string{},
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	// Check for newline to track line numbers
	if l.ch == '\n' {
		l.line++
		l.column = 0
		l.linePositions = append(l.linePositions, l.readPosition+1)
	}

	// check EOF
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII NUL char
	} else {
		l.ch = l.input[l.readPosition]
		l.column++
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// Get line number for a given position
func (l *Lexer) GetLineNumber(position int) int {
	if position < 0 {
		return 1 // Default to first line for invalid positions
	}

	// Find the line that contains this position
	lineNum := 1
	for i := 1; i < len(l.linePositions); i++ {
		if position < l.linePositions[i] {
			break
		}
		lineNum++
	}
	return lineNum
}

func newToken(tokenType token.TokenType, ch byte, position int) token.Token {
	return token.Token{
		Type:     tokenType,
		Literal:  string(ch),
		Position: position,
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.eatWhitespace()

	switch l.ch {
	case '=':
		// check for EQ
		if l.peekChar() == '=' {
			ch := l.ch
			pos := l.position
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal, Position: pos}
		} else {
			tok = newToken(token.ASSIGN, l.ch, l.position)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			pos := l.position
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal, Position: pos}
		} else {
			tok = newToken(token.BANG, l.ch, l.position)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, l.position)
	case '(':
		tok = newToken(token.LPAREN, l.ch, l.position)
	case ')':
		tok = newToken(token.RPAREN, l.ch, l.position)
	case ',':
		tok = newToken(token.COMMA, l.ch, l.position)
	case '+':
		tok = newToken(token.PLUS, l.ch, l.position)
	case '-':
		tok = newToken(token.MINUS, l.ch, l.position)
	case '/':
		tok = newToken(token.SLASH, l.ch, l.position)
	case '*':
		tok = newToken(token.ASTERISK, l.ch, l.position)
	case '<':
		tok = newToken(token.LT, l.ch, l.position)
	case '>':
		tok = newToken(token.GT, l.ch, l.position)
	case '{':
		tok = newToken(token.LBRACE, l.ch, l.position)
	case '}':
		tok = newToken(token.RBRACE, l.ch, l.position)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Position = l.position
	default:
		if isLetter(l.ch) {
			pos := l.position
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.Position = pos
			return tok
		} else if isDigit(l.ch) {
			pos := l.position
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			tok.Position = pos
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch, l.position)
			l.addError(fmt.Sprintf("Line %d, Column %d: illegal character '%c' found",
				l.line, l.column, l.ch))
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) eatWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// advances lexer from curr position to next any valid letter
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// TODO: safe place to playaround with any symbols in keyword
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '~'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// similar to readChar except it doesnt increment l.position
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

// Error handling methods
func (l *Lexer) addError(msg string) {
	l.errors = append(l.errors, msg)
}

func (l *Lexer) Errors() []string {
	return l.errors
}

func (l *Lexer) HasErrors() bool {
	return len(l.errors) > 0
}

// Format all lexer errors into a single string
func (l *Lexer) FormatErrors() string {
	if len(l.errors) == 0 {
		return ""
	}

	var errorBuilder strings.Builder
	errorBuilder.WriteString("Lexer errors:\n")

	for _, err := range l.errors {
		errorBuilder.WriteString("  - ")
		errorBuilder.WriteString(err)
		errorBuilder.WriteString("\n")
	}

	return errorBuilder.String()
}
