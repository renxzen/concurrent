package base

const (
	CategoricalType = iota
	Float64Type
	BinaryType
)

type Attribute interface {
	GetType() int
	GetName() string
	SetName(string)
	String() string
	GetSysValFromString(string) []byte
	GetStringFromSysVal([]byte) string
	Equals(Attribute) bool
	Compatible(Attribute) bool
}
