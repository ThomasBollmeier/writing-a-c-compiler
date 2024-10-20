package frontend

type AstType int

const (
	AstProgram AstType = iota
	AstFunction
	AstReturn
	AstInteger
	AstIdentifier
	AstUnary
	AstBinary
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
	VisitUnary(u *UnaryExpression)
	VisitBinary(b *BinaryExpression)
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
	Expression Expression
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

type UnaryExpression struct {
	Operator string
	Right    Expression
}

func (u *UnaryExpression) GetType() AstType {
	return AstUnary
}

func (u *UnaryExpression) Accept(visitor AstVisitor) {
	visitor.VisitUnary(u)
}

type BinaryExpression struct {
	Operator string
	Left     Expression
	Right    Expression
}

func (b *BinaryExpression) GetType() AstType {
	return AstBinary
}

func (b *BinaryExpression) Accept(visitor AstVisitor) {
	visitor.VisitBinary(b)
}
