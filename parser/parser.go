package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pro0o/sup-bud/ast"
	"github.com/pro0o/sup-bud/lexer"
	"github.com/pro0o/sup-bud/token"
)

type Parser struct {
	l *lexer.Lexer

	curToken   token.Token
	peekToken  token.Token
	errors     []string
	lineErrors map[int][]string // Track errors by line number

	prefixParseFns map[token.TokenType]prefixParseFn // nuds <- null denotations
	infixParseFns  map[token.TokenType]infixParseFn  // leds <- left denotations
	debugMode      bool                              // New field for enabling debug output
	depth          int
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

func (p *Parser) registerPrefix(tokenType token.TokenType, bud prefixParseFn) {
	p.prefixParseFns[tokenType] = bud
}

func (p *Parser) registerInfix(tokenType token.TokenType, bud infixParseFn) {
	p.infixParseFns[tokenType] = bud
}

// operator precedence
// _ ident being 0 to call being 7
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // func(X)
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

// Add new method to enable/disable debug mode
func (p *Parser) SetDebugMode(enabled bool) {
	p.debugMode = enabled
}

// Add helper for indentation
func (p *Parser) indent() string {
	return strings.Repeat("  ", p.depth)
}

func New(lex *lexer.Lexer) *Parser {
	p := &Parser{
		l:          lex,
		errors:     []string{},
		lineErrors: make(map[int][]string),
	}

	// Check if lexer has errors before proceeding
	if lex.HasErrors() {
		p.errors = append(p.errors, lex.Errors()...)
	}

	// fill in cur and peek token
	p.nextToken()
	p.nextToken()

	// nuds <- value is 0
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	// leds <- traversal
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)

	return p
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		// set current token to peek and peek -> next
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

// Enhanced error reporting with token line information
func (p *Parser) peekError(t token.TokenType) {
	line := p.l.GetLineNumber(p.peekToken.Position)
	msg := fmt.Sprintf("Line %d: expected next token to be %s, got %s instead",
		line, t, p.peekToken.Type)

	p.errors = append(p.errors, msg)

	// Also track errors by line for contextual error reporting
	p.lineErrors[line] = append(p.lineErrors[line],
		fmt.Sprintf("expected %s, got %s", t, p.peekToken.Type))
}

func (p *Parser) ParseProgram() *ast.Program {
	if p.debugMode {
		fmt.Println("Starting program parse")
	}

	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// If we already have lexer errors, don't attempt to parse
	if len(p.errors) > 0 {
		return program
	}

	for p.curToken.Type != token.EOF {
		if p.debugMode {
			fmt.Printf("%sToken: %s (%s)\n", p.indent(), p.curToken.Type, p.curToken.Literal)
		}

		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
			if p.debugMode {
				fmt.Printf("%sStatement AST: %s\n", p.indent(), stmt.String())
			}
		}
		p.nextToken()
	}

	if p.debugMode {
		fmt.Println("\nFull Program AST:")
		fmt.Println(program.String())
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	if p.curTokenIs(token.SEMICOLON) {
		p.addError("return statement requires an expression")
		return stmt
	}

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if stmt.Expression == nil {
		return nil
	}

	// optional semicolons
	// makes REPL easier later on
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	// type assertion that if the curr token type prefix func exist in parser map.
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	// if the func exists, call it.
	leftExp := prefix()

	if leftExp == nil {
		return nil
	}

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		// type assertion that if the curr token type prefix func exist in parser map.
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		// calls the func for the curr infix exp.
		leftExp = infix(leftExp)

		if leftExp == nil {
			return nil
		}
	}

	return leftExp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	// p.curToken is either of type token.BANG or token.MINUS
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	if expression.Right == nil {
		line := p.l.GetLineNumber(p.curToken.Position)
		p.addError(fmt.Sprintf("Line %d: invalid expression after %s operator",
			line, expression.Operator))
		return nil
	}

	return expression
}

// string -> int64
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)

	if err != nil {
		line := p.l.GetLineNumber(p.curToken.Position)
		msg := fmt.Sprintf("Line %d: could not parse '%s' as integer",
			line, p.curToken.Literal)
		p.addError(msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	line := p.l.GetLineNumber(p.curToken.Position)
	msg := fmt.Sprintf("Line %d: no prefix parse function for %s found", line, t)
	p.addError(msg)
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()

	expression.Right = p.parseExpression(precedence)

	if expression.Right == nil {
		line := p.l.GetLineNumber(p.curToken.Position)
		p.addError(fmt.Sprintf("Line %d: invalid right expression in operator '%s'",
			line, expression.Operator))
		return nil
	}

	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if exp == nil {
		return nil
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	if p.curTokenIs(token.EOF) {
		line := p.l.GetLineNumber(p.curToken.Position)
		p.addError(fmt.Sprintf("Line %d: unclosed block statement, expected '}'", line))
	}

	return block
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if expression.Condition == nil {
		return nil
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if lit.Parameters == nil {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments()

	if exp.Arguments == nil {
		return nil
	}

	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	expr := p.parseExpression(LOWEST)

	if expr == nil {
		return nil
	}

	args = append(args, expr)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		expr := p.parseExpression(LOWEST)
		if expr == nil {
			return nil
		}

		args = append(args, expr)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}

// Helper method to add error messages with current context
func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, msg)
}

// Format all errors with contextual information
func (p *Parser) FormatErrors() string {
	if len(p.errors) == 0 {
		return ""
	}

	var errorBuilder strings.Builder
	errorBuilder.WriteString("Parser errors:\n")

	for _, err := range p.errors {
		errorBuilder.WriteString("  - ")
		errorBuilder.WriteString(err)
		errorBuilder.WriteString("\n")
	}

	return errorBuilder.String()
}
