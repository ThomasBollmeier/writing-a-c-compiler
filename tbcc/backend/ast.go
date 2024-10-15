package backend

type AsmAstType int

const (
	AsmProgram AsmAstType = iota
	AsmFunctionDef
	AsmMov
	AsmUnary
	AsmBinary
	AsmIDiv
	AsmCdq
	AsmAllocStack
	AsmReturn
	AsmNeg
	AsmNot
	AsmAdd
	AsmSub
	AsmMul
	AsmImmediate
	AsmRegister
	AsmPseudoReg
	AsmStack
)

const (
	RegAX  string = "AX"
	RegDX         = "DX"
	RegR10        = "R10"
	RegR11        = "R11"
)

type AST interface {
	GetType() AsmAstType
	Accept(visitor AsmVisitor)
}

type AsmVisitor interface {
	VisitProgram(p *Program)
	VisitFunctionDef(f *FunctionDef)
	VisitMov(m *Mov)
	VisitUnary(u *Unary)
	VisitBinary(b *Binary)
	VisitIDiv(i *IDiv)
	VisitCdq(c *Cdq)
	VisitAllocStack(a *AllocStack)
	VisitReturn()
	VisitNeg(n *Neg)
	VisitNot(n *Not)
	VisitAdd(a *Add)
	VisitSub(s *Sub)
	VisitMul(m *Mul)
	VisitImmediate(i *Immediate)
	VisitRegister(r *Register)
	VisitPseudoReg(p *PseudoReg)
	VisitStack(s *Stack)
}

type Program struct {
	FuncDef FunctionDef
}

func NewProgram(funcDef FunctionDef) *Program {
	return &Program{funcDef}
}

func (p *Program) GetType() AsmAstType {
	return AsmProgram
}
func (p *Program) Accept(visitor AsmVisitor) {
	visitor.VisitProgram(p)
}

type FunctionDef struct {
	Name         string
	Instructions []Instruction
}

func NewFunctionDef(name string, instructions []Instruction) *FunctionDef {
	return &FunctionDef{name, instructions}
}

func (f *FunctionDef) GetType() AsmAstType {
	return AsmFunctionDef
}

func (f *FunctionDef) Accept(visitor AsmVisitor) {
	visitor.VisitFunctionDef(f)
}

type Instruction interface {
	AST
}

type Mov struct {
	Src Operand
	Dst Operand
}

func NewMov(src, dst Operand) *Mov {
	return &Mov{Src: src, Dst: dst}
}

func (m *Mov) GetType() AsmAstType {
	return AsmMov
}

func (m *Mov) Accept(visitor AsmVisitor) {
	visitor.VisitMov(m)
}

type Unary struct {
	Op      UnaryOp
	Operand Operand
}

func NewUnary(op UnaryOp, operand Operand) *Unary {
	return &Unary{op, operand}
}

func (u *Unary) GetType() AsmAstType {
	return AsmUnary
}

func (u *Unary) Accept(visitor AsmVisitor) {
	visitor.VisitUnary(u)
}

type Binary struct {
	Op       BinaryOp
	Operand1 Operand
	Operand2 Operand
}

func NewBinary(op BinaryOp, operand1 Operand, operand2 Operand) *Binary {
	return &Binary{op, operand1, operand2}
}

func (b *Binary) GetType() AsmAstType {
	return AsmBinary
}

func (b *Binary) Accept(visitor AsmVisitor) {
	visitor.VisitBinary(b)
}

type IDiv struct {
	Operand Operand
}

func NewIDiv(operand Operand) *IDiv {
	return &IDiv{operand}
}

func (i *IDiv) GetType() AsmAstType {
	return AsmIDiv
}

func (i *IDiv) Accept(visitor AsmVisitor) {
	visitor.VisitIDiv(i)
}

type Cdq struct{}

func NewCdq() *Cdq {
	return &Cdq{}
}

func (c *Cdq) GetType() AsmAstType {
	return AsmCdq
}

func (c *Cdq) Accept(visitor AsmVisitor) {
	visitor.VisitCdq(c)
}

type AllocStack struct {
	N int
}

func NewAllocStack(n int) *AllocStack {
	return &AllocStack{n}
}

func (a *AllocStack) GetType() AsmAstType {
	return AsmAllocStack
}

func (a *AllocStack) Accept(visitor AsmVisitor) {
	visitor.VisitAllocStack(a)
}

type Return struct{}

func NewReturn() *Return {
	return &Return{}
}

func (r *Return) GetType() AsmAstType {
	return AsmReturn
}

func (r *Return) Accept(visitor AsmVisitor) {
	visitor.VisitReturn()
}

type UnaryOp interface {
	AST
}

type Neg struct{}

func NewNeg() *Neg {
	return &Neg{}
}

func (n *Neg) GetType() AsmAstType {
	return AsmNeg
}

func (n *Neg) Accept(visitor AsmVisitor) {
	visitor.VisitNeg(n)
}

type Not struct{}

func NewNot() *Not {
	return &Not{}
}

func (n *Not) GetType() AsmAstType {
	return AsmNot
}

func (n *Not) Accept(visitor AsmVisitor) {
	visitor.VisitNot(n)
}

type BinaryOp interface {
	AST
}

type Add struct{}

func NewAdd() *Add {
	return &Add{}
}

func (a *Add) GetType() AsmAstType {
	return AsmAdd
}

func (a *Add) Accept(visitor AsmVisitor) {
	visitor.VisitAdd(a)
}

type Sub struct{}

func NewSub() *Sub {
	return &Sub{}
}

func (s *Sub) GetType() AsmAstType {
	return AsmSub
}

func (s *Sub) Accept(visitor AsmVisitor) {
	visitor.VisitSub(s)
}

type Mul struct{}

func NewMul() *Mul {
	return &Mul{}
}

func (m *Mul) GetType() AsmAstType {
	return AsmMul
}

func (m *Mul) Accept(visitor AsmVisitor) {
	visitor.VisitMul(m)
}

type Operand interface {
	AST
}

type Immediate struct {
	Value int
}

func NewImmediate(value int) *Immediate {
	return &Immediate{Value: value}
}

func (i *Immediate) GetType() AsmAstType {
	return AsmImmediate
}

func (i *Immediate) Accept(visitor AsmVisitor) {
	visitor.VisitImmediate(i)
}

type Register struct {
	Name string
}

func NewRegister(name string) *Register {
	return &Register{name}
}

func (r *Register) GetType() AsmAstType {
	return AsmRegister
}

func (r *Register) Accept(visitor AsmVisitor) {
	visitor.VisitRegister(r)
}

type PseudoReg struct {
	Ident string
}

func NewPseudoReg(ident string) *PseudoReg {
	return &PseudoReg{ident}
}

func (p *PseudoReg) GetType() AsmAstType {
	return AsmPseudoReg
}

func (p *PseudoReg) Accept(visitor AsmVisitor) {
	visitor.VisitPseudoReg(p)
}

type Stack struct {
	N int
}

func NewStack(n int) *Stack {
	return &Stack{n}
}

func (s *Stack) GetType() AsmAstType {
	return AsmStack
}

func (s *Stack) Accept(visitor AsmVisitor) {
	visitor.VisitStack(s)
}
