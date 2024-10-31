package frontend

type AstType int

const (
	AstProgram AstType = iota
	AstFunction
	AstVarDecl
	AstReturn
	AstExprStmt
	AstNullStmt
	AstInteger
	AstVariable
	AstUnary
	AstPostfixIncDec
	AstBinary
)

type AST interface {
	GetType() AstType
	Accept(visitor AstVisitor)
}

type AstVisitor interface {
	VisitProgram(p *Program)
	VisitFunction(f *Function)
	VisitVarDecl(v *VarDecl)
	VisitReturn(r *ReturnStmt)
	VisitExprStmt(e *ExpressionStmt)
	VisitNullStmt()
	VisitInteger(i *IntegerLiteral)
	VisitVariable(v *Variable)
	VisitUnary(u *UnaryExpression)
	VisitPostfixIncDec(p *PostfixIncDec)
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
	Body []BodyItem
}

func (f *Function) GetType() AstType {
	return AstFunction
}

func (f *Function) Accept(visitor AstVisitor) {
	visitor.VisitFunction(f)
}

type BodyItem interface {
	AST
}

type VarDecl struct {
	Name      string
	InitValue Expression
}

func (v *VarDecl) GetType() AstType {
	return AstVarDecl
}

func (v *VarDecl) Accept(visitor AstVisitor) {
	visitor.VisitVarDecl(v)
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

type ExpressionStmt struct {
	Expression Expression
}

func (e *ExpressionStmt) GetType() AstType {
	return AstExprStmt
}

func (e *ExpressionStmt) Accept(visitor AstVisitor) {
	visitor.VisitExprStmt(e)
}

type NullStmt struct{}

func (n *NullStmt) GetType() AstType {
	return AstNullStmt
}

func (n *NullStmt) Accept(visitor AstVisitor) {
	visitor.VisitNullStmt()
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

type Variable struct {
	Name string
}

func (v *Variable) GetType() AstType {
	return AstVariable
}

func (v *Variable) Accept(visitor AstVisitor) {
	visitor.VisitVariable(v)
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

type PostfixIncDec struct {
	Operator string
	Operand  Variable
}

func (p *PostfixIncDec) GetType() AstType {
	return AstPostfixIncDec
}

func (p *PostfixIncDec) Accept(visitor AstVisitor) {
	visitor.VisitPostfixIncDec(p)
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
