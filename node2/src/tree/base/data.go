package base

type SortDirection int

const (
	Descending SortDirection = 1
	Ascending  SortDirection = 2
)

type DataGrid interface {
	GetAttribute(Attribute) (AttributeSpec, error)
	AllAttributes() []Attribute
	AddClassAttribute(Attribute) error
	RemoveClassAttribute(Attribute) error
	AllClassAttributes() []Attribute
	Get(AttributeSpec, int) []byte
	MapOverRows([]AttributeSpec, func([][]byte, int) (bool, error)) error
}

type FixedDataGrid interface {
	DataGrid
	RowString(int) string
	Size() (int, int)
}

type UpdatableDataGrid interface {
	FixedDataGrid
	Set(AttributeSpec, int, []byte)
	AddAttribute(Attribute) AttributeSpec
	Extend(int) error
}
