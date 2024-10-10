package tacky

// Intermediate representation (IR): three address code (TAC)

type TacType int

const (
	TacProgram TacType = iota
	TacFunction
	TacReturn
	TacUnary
	TacIntConstant
	TacVar
	TacComplement
	TacNegate
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
