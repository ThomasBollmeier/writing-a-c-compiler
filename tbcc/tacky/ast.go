package tacky

// Intermediate representation (IR): three address code (TAC)

type TacType int

const (
	TacProgram TacType = iota
	TacFunction
	TacReturn
	TacUnary
	TacBinary
	TacIntConstant
	TacVar
	TacComplement
	TacNegate
	TacAdd
	TacSub
	TacMul
	TacDiv
	TacRemainder
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

type BinaryOp interface{}

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
