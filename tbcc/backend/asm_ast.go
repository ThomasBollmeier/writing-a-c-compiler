package backend

type AsmAstType int

const (
	AsmProgram AsmAstType = iota
	AsmFunctionDef
	AsmMov
	AsmReturn
	AsmImmediate
	AsmRegister
)

type AsmAST interface {
	GetType() AsmAstType
	Accept(visitor AsmVisitor)
}

type AsmVisitor interface {
	VisitProgram(p *Program)
	VisitFunctionDef(f *FunctionDef)
	VisitMov(m *Mov)
	VisitReturn()
	VisitImmediate(i *Immediate)
	VisitRegister()
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
	AsmAST
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

type Operand interface {
	AsmAST
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

type Register struct{}

func NewRegister() *Register {
	return &Register{}
}

func (r *Register) GetType() AsmAstType {
	return AsmRegister
}

func (r *Register) Accept(visitor AsmVisitor) {
	visitor.VisitRegister()
}
