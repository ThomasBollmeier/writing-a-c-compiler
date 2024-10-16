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

type TacNode interface {
	GetType() TacType
}

type Program struct {
	Fun Function
}

func (p *Program) GetType() TacType {
	return TacProgram
}

type Function struct {
	Ident string
	Body  []Instruction
}

func (f *Function) GetType() TacType {
	return TacFunction
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

type Unary struct {
	Op  UnaryOp
	Src Value
	Dst Value
}

func (u *Unary) GetType() TacType {
	return TacUnary
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

type Copy struct {
	Src Value
	Dst Value
}

func (c *Copy) GetType() TacType {
	return TacCopy
}

type Jump struct {
	Target string
}

func (j *Jump) GetType() TacType {
	return TacJump
}

type JumpIfZero struct {
	Condition Value
	Target    string
}

func (j *JumpIfZero) GetType() TacType {
	return TacJumpIfZero
}

type JumpIfNotZero struct {
	Condition Value
	Target    string
}

func (j *JumpIfNotZero) GetType() TacType {
	return TacJumpIfNotZero
}

type Label struct {
	Name string
}

func (l *Label) GetType() TacType {
	return TacLabel
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

type Var struct {
	Ident string
}

func (v *Var) GetType() TacType {
	return TacVar
}

type UnaryOp interface {
	TacNode
}

type Complement struct{}

func (c *Complement) GetType() TacType {
	return TacComplement
}

type Negate struct{}

func (n *Negate) GetType() TacType {
	return TacNegate
}

type Not struct{}

func (n *Not) GetType() TacType {
	return TacNot
}

type BinaryOp interface {
	TacNode
}

type Add struct{}

func (a *Add) GetType() TacType {
	return TacAdd
}

type Sub struct{}

func (s *Sub) GetType() TacType {
	return TacSub
}

type Mul struct{}

func (m *Mul) GetType() TacType {
	return TacMul
}

type Div struct{}

func (d *Div) GetType() TacType {
	return TacDiv
}

type Remainder struct{}

func (r *Remainder) GetType() TacType {
	return TacRemainder
}

type BitAnd struct{}

func (b *BitAnd) GetType() TacType {
	return TacBitAnd
}

type BitOr struct{}

func (b *BitOr) GetType() TacType {
	return TacBitOr
}

type BitXor struct{}

func (b *BitXor) GetType() TacType {
	return TacBitXor
}

type BitShiftLeft struct{}

func (b *BitShiftLeft) GetType() TacType {
	return TacBitShiftLeft
}

type BitShiftRight struct{}

func (b *BitShiftRight) GetType() TacType {
	return TacBitShiftRight
}

type And struct{}

func (a *And) GetType() TacType {
	return TacAnd
}

type Or struct{}

func (p *Or) GetType() TacType {
	return TacOr
}

type Equal struct{}

func (e *Equal) GetType() TacType {
	return TacEq
}

type NotEqual struct{}

func (n *NotEqual) GetType() TacType {
	return TacNotEq
}

type Greater struct{}

func (g *Greater) GetType() TacType {
	return TacGt
}

type GreaterEq struct{}

func (g *GreaterEq) GetType() TacType {
	return TacGtEq
}

type Less struct{}

func (l *Less) GetType() TacType {
	return TacLt
}

type LessEq struct{}

func (l *LessEq) GetType() TacType {
	return TacLtEq
}
