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
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	LT       = "<"
	GT       = ">"
	EQ       = "=="
	NOT_EQ   = "!="

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
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

type Token struct {
	Type    TokenType
	Literal string
}

// distingusih user-defined identifiers apart from language keyword
// assign one type to another -> maps DS
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"olaf":   LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// checks the keywords map to see whether
// the given identifier is in fact a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
