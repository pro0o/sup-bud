package ast

import (
	"main/token"
)

type Node interface {
	TokenLiteral() string
}

// root node of AST
type Program struct {
	// slice of statement nodes i.e.
	// every valid program is a series of statement
	Statements []Statement
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// name of the variable
type Identifier struct {
	Token token.Token // token.IDENT in question
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// statement nodes
func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

// parsing a let statement
type LetStatement struct {
	Token token.Token // the token.LET in question.
	Name  *Identifier
	Value Expression // evaluated value
}

type ReturnStatement struct {
	Token       token.Token // return token type
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
