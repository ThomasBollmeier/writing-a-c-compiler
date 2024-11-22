package backend

type AsmAstType int

const (
	AsmProgram AsmAstType = iota
	AsmFunctionDef
	AsmMov
	AsmUnary
	AsmBinary
	AsmCmp
	AsmIDiv
	AsmCdq
	AsmJmp
	AsmJmpCC
	AsmSetCC
	AsmLabel
	AsmAllocStack
	AsmDeAllocStack
	AsmPush
	AsmCall
	AsmReturn
	AsmNeg
	AsmNot
	AsmAdd
	AsmSub
	AsmMul
	AsmBitAnd
	AsmBitOr
	AsmBitXor
	AsmBitShiftLeft
	AsmBitShiftRight
	AsmImmediate
	AsmRegister
	AsmPseudoReg
	AsmStack
)

type ConditionCode uint

const (
	CcEq ConditionCode = iota
	CcNotEq
	CcGt
	CcGtEq
	CcLt
	CcLtEq
)

const (
	RegAX  string = "AX"
	RegCX  string = "CX"
	RegDX         = "DX"
	RegDI         = "DI"
	RegSI         = "SI"
	RegR8         = "R8"
	RegR9         = "R9"
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
	VisitCmp(c *Cmp)
	VisitIDiv(i *IDiv)
	VisitCdq(c *Cdq)
	VisitJump(j *Jump)
	VisitJumpCC(j *JumpCC)
	VisitSetCC(s *SetCC)
	VisitLabel(l *Label)
	VisitAllocStack(a *AllocStack)
	VisitDeAllocStack(d *DeAllocStack)
	VisitPush(p *Push)
	VisitCall(c *Call)
	VisitReturn()
	VisitNeg(n *Neg)
	VisitNot(n *Not)
	VisitAdd(a *Add)
	VisitSub(s *Sub)
	VisitMul(m *Mul)
	VisitBitOp(bo BinaryOp)
	VisitImmediate(i *Immediate)
	VisitRegister(r *Register)
	VisitPseudoReg(p *PseudoReg)
	VisitStack(s *Stack)
}

type Program struct {
	FuncDefs []FunctionDef
}

func NewProgram(funcDefs []FunctionDef) *Program {
	return &Program{funcDefs}
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

type Cmp struct {
	Left  Operand
	Right Operand
}

func NewCmp(left Operand, right Operand) *Cmp {
	return &Cmp{left, right}
}

func (c *Cmp) GetType() AsmAstType {
	return AsmCmp
}

func (c *Cmp) Accept(visitor AsmVisitor) {
	visitor.VisitCmp(c)
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

type Jump struct {
	Identifier string
}

func NewJump(identifier string) *Jump {
	return &Jump{identifier}
}

func (j *Jump) GetType() AsmAstType {
	return AsmJmp
}

func (j *Jump) Accept(visitor AsmVisitor) {
	visitor.VisitJump(j)
}

type JumpCC struct {
	CondCode   ConditionCode
	Identifier string
}

func NewJumpCC(condCode ConditionCode, identifier string) *JumpCC {
	return &JumpCC{condCode, identifier}
}

func (j *JumpCC) GetType() AsmAstType {
	return AsmJmpCC
}

func (j *JumpCC) Accept(visitor AsmVisitor) {
	visitor.VisitJumpCC(j)
}

type SetCC struct {
	CondCode ConditionCode
	Op       Operand
}

func NewSetCC(condCode ConditionCode, op Operand) *SetCC {
	return &SetCC{condCode, op}
}

func (s *SetCC) GetType() AsmAstType {
	return AsmSetCC
}

func (s *SetCC) Accept(visitor AsmVisitor) {
	visitor.VisitSetCC(s)
}

type Label struct {
	Identifier string
}

func NewLabel(identifier string) *Label {
	return &Label{identifier}
}

func (l *Label) GetType() AsmAstType {
	return AsmLabel
}

func (l *Label) Accept(visitor AsmVisitor) {
	visitor.VisitLabel(l)
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

type DeAllocStack struct {
	N int
}

func NewDeAllocStack(n int) *DeAllocStack {
	return &DeAllocStack{n}
}

func (d *DeAllocStack) GetType() AsmAstType {
	return AsmDeAllocStack
}

func (d *DeAllocStack) Accept(visitor AsmVisitor) {
	visitor.VisitDeAllocStack(d)
}

type Push struct {
	Op Operand
}

func NewPush(op Operand) *Push {
	return &Push{op}
}

func (p *Push) GetType() AsmAstType {
	return AsmPush
}

func (p *Push) Accept(visitor AsmVisitor) {
	visitor.VisitPush(p)
}

type Call struct {
	Identifier string
}

func NewCall(identifier string) *Call {
	return &Call{identifier}
}

func (c *Call) GetType() AsmAstType {
	return AsmCall
}

func (c *Call) Accept(visitor AsmVisitor) {
	visitor.VisitCall(c)
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

type BitAnd struct{}

func NewBitAnd() *BitAnd {
	return &BitAnd{}
}

func (b *BitAnd) GetType() AsmAstType {
	return AsmBitAnd
}

func (b *BitAnd) Accept(visitor AsmVisitor) {
	visitor.VisitBitOp(b)
}

type BitOr struct{}

func NewBitOr() *BitOr {
	return &BitOr{}
}

func (b *BitOr) GetType() AsmAstType {
	return AsmBitOr
}

func (b *BitOr) Accept(visitor AsmVisitor) {
	visitor.VisitBitOp(b)
}

type BitXor struct{}

func NewBitXor() *BitXor {
	return &BitXor{}
}

func (b *BitXor) GetType() AsmAstType {
	return AsmBitXor
}

func (b *BitXor) Accept(visitor AsmVisitor) {
	visitor.VisitBitOp(b)
}

type BitShiftLeft struct{}

func NewBitShiftLeft() *BitShiftLeft {
	return &BitShiftLeft{}
}

func (b *BitShiftLeft) GetType() AsmAstType {
	return AsmBitShiftLeft
}

func (b *BitShiftLeft) Accept(visitor AsmVisitor) {
	visitor.VisitBitOp(b)
}

type BitShiftRight struct{}

func NewBitShiftRight() *BitShiftRight {
	return &BitShiftRight{}
}

func (b *BitShiftRight) GetType() AsmAstType {
	return AsmBitShiftRight
}

func (b *BitShiftRight) Accept(visitor AsmVisitor) {
	visitor.VisitBitOp(b)
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
