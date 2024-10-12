package backend

type AsmAstType int

const (
	AsmProgram AsmAstType = iota
	AsmFunctionDef
	AsmMov
	AsmUnary
	AsmAllocStack
	AsmReturn
	AsmNeg
	AsmNot
	AsmImmediate
	AsmRegister
	AsmPseudoReg
	AsmStack
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
	VisitAllocStack(a *AllocStack)
	VisitReturn()
	VisitNeg(n *Neg)
	VisitNot(n *Not)
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
