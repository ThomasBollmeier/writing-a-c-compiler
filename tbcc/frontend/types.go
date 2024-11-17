package frontend

type TypeId int

const (
	TypeInt TypeId = iota
	TypeFunc
)

type typeInfo interface {
	getTypeId() TypeId
	equal(other typeInfo) bool
}

type intInfo struct{}

func (i *intInfo) getTypeId() TypeId {
	return TypeInt
}

func (i *intInfo) equal(typeInfo) bool {
	return true
}

type funcInfo struct {
	numParams int
	isDefined bool
}

func (f *funcInfo) getTypeId() TypeId {
	return TypeFunc
}

func (f *funcInfo) equal(other typeInfo) bool {
	otherFunc := other.(*funcInfo)
	return f.numParams == otherFunc.numParams
}
