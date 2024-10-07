package frontend

type AstType int

const (
	AstProgram AstType = iota
	AstFunction
	AstReturn
	AstInteger
	AstIdentifier
)

type AST interface {
	GetType() AstType
	Accept(visitor AstVisitor)
}

type AstVisitor interface {
	VisitProgram(p *Program)
	VisitFunction(f *Function)
	VisitReturn(r *ReturnStmt)
	VisitInteger(i *IntegerLiteral)
	VisitIdentifier(id *Identifier)
}

type Program struct {
	Func Function
}

func (p *Program) GetType() AstType {
	return AstProgram
}
func (p *Program) Accept(visitor AstVisitor) {
	visitor.VisitProgram(p)
}

type Function struct {
	Name string
	Body Statement
}

func (f *Function) GetType() AstType {
	return AstFunction
}

func (f *Function) Accept(visitor AstVisitor) {
	visitor.VisitFunction(f)
}

type Statement interface {
	AST
}

type ReturnStmt struct {
	expression Expression
}

func (r *ReturnStmt) GetType() AstType {
	return AstReturn
}

func (r *ReturnStmt) Accept(visitor AstVisitor) {
	visitor.VisitReturn(r)
}

type Expression interface {
	AST
}

type IntegerLiteral struct {
	Value int
}

func (i *IntegerLiteral) GetType() AstType {
	return AstInteger
}

func (i *IntegerLiteral) Accept(visitor AstVisitor) {
	visitor.VisitInteger(i)
}

type Identifier struct {
	Value string
}

func (id *Identifier) GetType() AstType {
	return AstIdentifier
}

func (id *Identifier) Accept(visitor AstVisitor) {
	visitor.VisitIdentifier(id)
}
