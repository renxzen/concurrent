package base

import (
	"bytes"
)

type AttributeGroup interface {
	appendToRowBuf(row int, buffer *bytes.Buffer)
	AddAttribute(Attribute) error
	Attributes() []Attribute
	get(int, int) []byte
	set(int, int, []byte)
	setStorage([]byte)
	RowSizeInBytes() int
	resize(int)
	Storage() []byte
	String() string
}
