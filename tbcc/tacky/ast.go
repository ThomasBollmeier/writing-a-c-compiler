package tacky

// Intermediate representation (IR): three address code (TAC)

type TacType int

const (
	TacProgram TacType = iota
	TacFunction
	TacReturn
	TacUnary
	TacBinary
	TacCopy
	TacJump
	TacJumpIfZero
	TacJumpIfNotZero
	TacLabel
	TacFunCall
	TacIntConstant
	TacVar
	TacComplement
	TacNegate
	TacAdd
	TacSub
	TacMul
	TacDiv
	TacRemainder
	TacBitAnd
	TacBitOr
	TacBitXor
	TacBitShiftLeft
	TacBitShiftRight
	TacNot
	TacAnd
	TacOr
	TacEq
	TacNotEq
	TacGt
	TacGtEq
	TacLt
	TacLtEq
)

type TacVisitor interface {
	visitProgram(p *Program)
	visitFunction(f *Function)
	visitReturn(r *Return)
	visitUnary(u *Unary)
	visitBinary(b *Binary)
	visitCopy(c *Copy)
	visitJump(j *Jump)
	visitJumpIfZero(j *JumpIfZero)
	visitJumpIfNotZero(j *JumpIfNotZero)
	visitLabel(l *Label)
	visitFunctionCall(f *FunctionCall)
	visitIntConstant(i *IntConstant)
	visitVar(v *Var)
	visitComplement()
	visitNegate()
	visitNot()
	visitAdd()
	visitSub()
	visitMul()
	visitDiv()
	visitRemainder()
	visitBitAnd()
	visitBitOr()
	visitBitXor()
	visitBitShiftLeft()
	visitBitShiftRight()
	visitAnd()
	visitOr()
	visitEqual()
	visitNotEqual()
	visitGreater()
	visitGreaterEq()
	visitLess()
	visitLessEq()
}

type TacNode interface {
	GetType() TacType
	Accept(visitor TacVisitor)
}

type Program struct {
	Funs []Function
}

func (p *Program) GetType() TacType {
	return TacProgram
}

func (p *Program) Accept(visitor TacVisitor) {
	visitor.visitProgram(p)
}

type Function struct {
	Ident      string
	Parameters []string
	Body       []Instruction
}

func (f *Function) GetType() TacType {
	return TacFunction
}

func (f *Function) Accept(visitor TacVisitor) {
	visitor.visitFunction(f)
}

type Instruction interface {
	TacNode
}

type Return struct {
	Val Value
}

func (r *Return) GetType() TacType {
	return TacReturn
}

func (r *Return) Accept(visitor TacVisitor) {
	visitor.visitReturn(r)
}

type Unary struct {
	Op  UnaryOp
	Src Value
	Dst Value
}

func (u *Unary) GetType() TacType {
	return TacUnary
}

func (u *Unary) Accept(visitor TacVisitor) {
	visitor.visitUnary(u)
}

type Binary struct {
	Op   BinaryOp
	Src1 Value
	Src2 Value
	Dst  Value
}

func (b *Binary) GetType() TacType {
	return TacBinary
}

func (b *Binary) Accept(visitor TacVisitor) {
	visitor.visitBinary(b)
}

type Copy struct {
	Src Value
	Dst Value
}

func (c *Copy) GetType() TacType {
	return TacCopy
}

func (c *Copy) Accept(visitor TacVisitor) {
	visitor.visitCopy(c)
}

type Jump struct {
	Target string
}

func (j *Jump) GetType() TacType {
	return TacJump
}

func (j *Jump) Accept(visitor TacVisitor) {
	visitor.visitJump(j)
}

type JumpIfZero struct {
	Condition Value
	Target    string
}

func (j *JumpIfZero) GetType() TacType {
	return TacJumpIfZero
}

func (j *JumpIfZero) Accept(visitor TacVisitor) {
	visitor.visitJumpIfZero(j)
}

type JumpIfNotZero struct {
	Condition Value
	Target    string
}

func (j *JumpIfNotZero) GetType() TacType {
	return TacJumpIfNotZero
}

func (j *JumpIfNotZero) Accept(visitor TacVisitor) {
	visitor.visitJumpIfNotZero(j)
}

type Label struct {
	Name string
}

func (l *Label) GetType() TacType {
	return TacLabel
}

func (l *Label) Accept(visitor TacVisitor) {
	visitor.visitLabel(l)
}

type FunctionCall struct {
	Name string
	Args []Value
	Dst  Value
}

func (f *FunctionCall) GetType() TacType {
	return TacFunCall
}

func (f *FunctionCall) Accept(visitor TacVisitor) {
	visitor.visitFunctionCall(f)
}

type Value interface {
	TacNode
}

type IntConstant struct {
	Val int
}

func (i *IntConstant) GetType() TacType {
	return TacIntConstant
}

func (i *IntConstant) Accept(visitor TacVisitor) {
	visitor.visitIntConstant(i)
}

type Var struct {
	Ident string
}

func (v *Var) GetType() TacType {
	return TacVar
}

func (v *Var) Accept(visitor TacVisitor) {
	visitor.visitVar(v)
}

type UnaryOp interface {
	TacNode
}

type Complement struct{}

func (c *Complement) GetType() TacType {
	return TacComplement
}

func (c *Complement) Accept(visitor TacVisitor) {
	visitor.visitComplement()
}

type Negate struct{}

func (n *Negate) GetType() TacType {
	return TacNegate
}

func (n *Negate) Accept(visitor TacVisitor) {
	visitor.visitNegate()
}

type Not struct{}

func (n *Not) GetType() TacType {
	return TacNot
}

func (n *Not) Accept(visitor TacVisitor) {
	visitor.visitNot()
}

type BinaryOp interface {
	TacNode
}

type Add struct{}

func (a *Add) GetType() TacType {
	return TacAdd
}

func (a *Add) Accept(visitor TacVisitor) {
	visitor.visitAdd()
}

type Sub struct{}

func (s *Sub) GetType() TacType {
	return TacSub
}

func (s *Sub) Accept(visitor TacVisitor) {
	visitor.visitSub()
}

type Mul struct{}

func (m *Mul) GetType() TacType {
	return TacMul
}

func (m *Mul) Accept(visitor TacVisitor) {
	visitor.visitMul()
}

type Div struct{}

func (d *Div) GetType() TacType {
	return TacDiv
}

func (d *Div) Accept(visitor TacVisitor) {
	visitor.visitDiv()
}

type Remainder struct{}

func (r *Remainder) GetType() TacType {
	return TacRemainder
}

func (r *Remainder) Accept(visitor TacVisitor) {
	visitor.visitRemainder()
}

type BitAnd struct{}

func (b *BitAnd) GetType() TacType {
	return TacBitAnd
}

func (b *BitAnd) Accept(visitor TacVisitor) {
	visitor.visitBitAnd()
}

type BitOr struct{}

func (b *BitOr) GetType() TacType {
	return TacBitOr
}

func (b *BitOr) Accept(visitor TacVisitor) {
	visitor.visitBitOr()
}

type BitXor struct{}

func (b *BitXor) GetType() TacType {
	return TacBitXor
}

func (b *BitXor) Accept(visitor TacVisitor) {
	visitor.visitBitXor()
}

type BitShiftLeft struct{}

func (b *BitShiftLeft) GetType() TacType {
	return TacBitShiftLeft
}

func (b *BitShiftLeft) Accept(visitor TacVisitor) {
	visitor.visitBitShiftLeft()
}

type BitShiftRight struct{}

func (b *BitShiftRight) GetType() TacType {
	return TacBitShiftRight
}

func (b *BitShiftRight) Accept(visitor TacVisitor) {
	visitor.visitBitShiftRight()
}

type And struct{}

func (a *And) GetType() TacType {
	return TacAnd
}

func (a *And) Accept(visitor TacVisitor) {
	visitor.visitAnd()
}

type Or struct{}

func (p *Or) GetType() TacType {
	return TacOr
}

func (p *Or) Accept(visitor TacVisitor) {
	visitor.visitOr()
}

type Equal struct{}

func (e *Equal) GetType() TacType {
	return TacEq
}

func (e *Equal) Accept(visitor TacVisitor) {
	visitor.visitEqual()
}

type NotEqual struct{}

func (n *NotEqual) GetType() TacType {
	return TacNotEq
}

func (n *NotEqual) Accept(visitor TacVisitor) {
	visitor.visitNotEqual()
}

type Greater struct{}

func (g *Greater) GetType() TacType {
	return TacGt
}

func (g *Greater) Accept(visitor TacVisitor) {
	visitor.visitGreater()
}

type GreaterEq struct{}

func (g *GreaterEq) GetType() TacType {
	return TacGtEq
}

func (g *GreaterEq) Accept(visitor TacVisitor) {
	visitor.visitGreaterEq()
}

type Less struct{}

func (l *Less) GetType() TacType {
	return TacLt
}

func (l *Less) Accept(visitor TacVisitor) {
	visitor.visitLess()
}

type LessEq struct{}

func (l *LessEq) GetType() TacType {
	return TacLtEq
}

func (l *LessEq) Accept(visitor TacVisitor) {
	visitor.visitLessEq()
}
