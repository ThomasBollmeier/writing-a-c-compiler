package frontend

type TypeId int

const (
	TypeInt TypeId = iota
	TypeFunc
)

type TypeInfo interface {
	GetTypeId() TypeId
	Equal(other TypeInfo) bool
}

type IntInfo struct{}

func (i *IntInfo) GetTypeId() TypeId {
	return TypeInt
}

func (i *IntInfo) Equal(TypeInfo) bool {
	return true
}

type FuncInfo struct {
	NumParams int
	IsDefined bool
}

func (f *FuncInfo) GetTypeId() TypeId {
	return TypeFunc
}

func (f *FuncInfo) Equal(other TypeInfo) bool {
	otherFunc := other.(*FuncInfo)
	return f.NumParams == otherFunc.NumParams
}
