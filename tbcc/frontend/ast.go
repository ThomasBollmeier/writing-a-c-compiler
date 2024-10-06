package frontend

type AstType int

type AST interface {
	GetType() AstType
	Accept(visitor AstVisitor)
}

type AstVisitor interface{}
