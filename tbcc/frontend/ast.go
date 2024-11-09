package frontend

type AstType int

const (
	AstProgram AstType = iota
	AstFunction
	AstVarDecl
	AstReturn
	AstExprStmt
	AstIfStmt
	AstBlockStmt
	AstGotoStmt
	AstLabelStmt
	AstDoWhileStmt
	AstWhileStmt
	AstForStmt
	AstBreakStmt
	AstContinueStmt
	AstSwitchStmt
	AstCaseStmt
	AstNullStmt
	AstInteger
	AstVariable
	AstUnary
	AstPostfixIncDec
	AstBinary
	AstConditional
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
	VisitIfStmt(i *IfStmt)
	VisitBlockStmt(b *BlockStmt)
	VisitGotoStmt(g *GotoStmt)
	VisitLabelStmt(l *LabelStmt)
	VisitDoWhileStmt(d *DoWhileStmt)
	VisitWhileStmt(w *WhileStmt)
	VisitForStmt(f *ForStmt)
	VisitBreakStmt(b *BreakStmt)
	VisitContinueStmt(c *ContinueStmt)
	VisitSwitchStmt(s *SwitchStmt)
	VisitCaseStmt(c *CaseStmt)
	VisitNullStmt()
	VisitInteger(i *IntegerLiteral)
	VisitVariable(v *Variable)
	VisitUnary(u *UnaryExpression)
	VisitPostfixIncDec(p *PostfixIncDec)
	VisitBinary(b *BinaryExpression)
	VisitConditional(c *Conditional)
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
	Body BlockStmt
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

type IfStmt struct {
	Condition  Expression
	Consequent Statement
	Alternate  Statement
}

func (i *IfStmt) GetType() AstType {
	return AstIfStmt
}

func (i *IfStmt) Accept(visitor AstVisitor) {
	visitor.VisitIfStmt(i)
}

type BlockStmt struct {
	Items []BodyItem
}

func (b *BlockStmt) GetType() AstType {
	return AstBlockStmt
}

func (b *BlockStmt) Accept(visitor AstVisitor) {
	visitor.VisitBlockStmt(b)
}

type GotoStmt struct {
	Target string
}

func (g *GotoStmt) GetType() AstType {
	return AstGotoStmt
}

func (g *GotoStmt) Accept(visitor AstVisitor) {
	visitor.VisitGotoStmt(g)
}

type LabelStmt struct {
	Name string
}

func (l *LabelStmt) GetType() AstType {
	return AstLabelStmt
}

func (l *LabelStmt) Accept(visitor AstVisitor) {
	visitor.VisitLabelStmt(l)
}

type DoWhileStmt struct {
	Condition Expression
	Body      Statement
	Label     string
}

func (d *DoWhileStmt) GetType() AstType {
	return AstDoWhileStmt
}

func (d *DoWhileStmt) Accept(visitor AstVisitor) {
	visitor.VisitDoWhileStmt(d)
}

type WhileStmt struct {
	Condition Expression
	Body      Statement
	Label     string
}

func (w *WhileStmt) GetType() AstType {
	return AstWhileStmt
}

func (w *WhileStmt) Accept(visitor AstVisitor) {
	visitor.VisitWhileStmt(w)
}

type ForStmt struct {
	InitStmt  BodyItem
	Condition Expression
	Post      Expression
	Body      Statement
	Label     string
}

func (f *ForStmt) GetType() AstType {
	return AstForStmt
}

func (f *ForStmt) Accept(visitor AstVisitor) {
	visitor.VisitForStmt(f)
}

type BreakStmt struct {
	Label string
}

func (b *BreakStmt) GetType() AstType {
	return AstBreakStmt
}

func (b *BreakStmt) Accept(visitor AstVisitor) {
	visitor.VisitBreakStmt(b)
}

type ContinueStmt struct {
	Label string
}

func (c *ContinueStmt) GetType() AstType {
	return AstContinueStmt
}

func (c *ContinueStmt) Accept(visitor AstVisitor) {
	visitor.VisitContinueStmt(c)
}

type SwitchStmt struct {
	Expr           Expression
	Body           Statement
	Label          string
	FirstCaseLabel string
}

func (s *SwitchStmt) GetType() AstType {
	return AstSwitchStmt
}

func (s *SwitchStmt) Accept(visitor AstVisitor) {
	visitor.VisitSwitchStmt(s)
}

type CaseStmt struct {
	Value         Expression
	Label         string
	PrevCaseLabel string
	NextCaseLabel string
}

func (c *CaseStmt) GetType() AstType {
	return AstCaseStmt
}

func (c *CaseStmt) Accept(visitor AstVisitor) {
	visitor.VisitCaseStmt(c)
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

type Conditional struct {
	Condition  Expression
	Consequent Expression
	Alternate  Expression
}

func (c *Conditional) GetType() AstType {
	return AstConditional
}

func (c *Conditional) Accept(visitor AstVisitor) {
	visitor.VisitConditional(c)
}
