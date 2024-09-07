package ast

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
