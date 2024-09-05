package token

type TokenType string

// token tyeps
const (
	ILLEGAL = "ILLEGAL" // any token out of scope/ not defined by me.
	EOF     = "EOF"

	// identifiers
	IDENT = "IDENT"
	INT   = "INT"

	// operators
	ASSIGN = "="
	PLUS   = "+"

	// delimiters: separation between data streams
	COMMA     = ","
	SEMICOLON = ";"

	// parenthesis
	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
)

type Token struct {
	Type    TokenType
	Literal string
}
