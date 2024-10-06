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
	visitProgram(p *Program)
	visitFunction(f *Function)
	visitReturn(r *ReturnStmt)
	visitInteger(i *IntegerLiteral)
	visitIdentifier(id *Identifier)
}

type Program struct {
	function Function
}

func (p *Program) GetType() AstType {
	return AstProgram
}
func (p *Program) Accept(visitor AstVisitor) {
	visitor.visitProgram(p)
}

type Function struct {
	name string
	body Statement
}

func (f *Function) GetType() AstType {
	return AstFunction
}

func (f *Function) Accept(visitor AstVisitor) {
	visitor.visitFunction(f)
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
	visitor.visitReturn(r)
}

type Expression interface {
	AST
}

type IntegerLiteral struct {
	value int
}

func (i *IntegerLiteral) GetType() AstType {
	return AstInteger
}

func (i *IntegerLiteral) Accept(visitor AstVisitor) {
	visitor.visitInteger(i)
}

type Identifier struct {
	value string
}

func (id *Identifier) GetType() AstType {
	return AstIdentifier
}

func (id *Identifier) Accept(visitor AstVisitor) {
	visitor.visitIdentifier(id)
}
