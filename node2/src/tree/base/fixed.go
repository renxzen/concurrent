package base

import (
	"bytes"
	"fmt"
)

type FixedAttributeGroup struct {
	parent     DataGrid
	attributes []Attribute
	size       int
	alloc      []byte
	maxRow     int
}

func (f *FixedAttributeGroup) String() string {
	return "FixedAttributeGroup"
}

func (f *FixedAttributeGroup) RowSizeInBytes() int {
	return len(f.attributes) * f.size
}

func (f *FixedAttributeGroup) Attributes() []Attribute {
	ret := make([]Attribute, len(f.attributes))

	for i, a := range f.attributes {
		ret[i] = a
	}
	return ret
}

func (f *FixedAttributeGroup) AddAttribute(a Attribute) error {
	f.attributes = append(f.attributes, a)
	return nil
}

func (f *FixedAttributeGroup) setStorage(a []byte) {
	f.alloc = a
}

func (f *FixedAttributeGroup) Storage() []byte {
	return f.alloc
}

func (f *FixedAttributeGroup) offset(col, row int) int {
	return row*f.RowSizeInBytes() + col*f.size
}

func (f *FixedAttributeGroup) set(col int, row int, val []byte) {

	if len(val) != f.size {
		panic(fmt.Sprintf("Tried to call set() with %d bytes, should be %d", len(val), f.size))
	}

	offset := f.offset(col, row)

	copied := copy(f.alloc[offset:], val)
	if copied != f.size {
		panic(fmt.Sprintf("set(%d) terminated by only copying %d bytes", copied, f.size))
	}

	row++
	if row > f.maxRow {
		f.maxRow = row
	}
}

func (f *FixedAttributeGroup) get(col int, row int) []byte {
	offset := f.offset(col, row)
	return f.alloc[offset : offset+f.size]
}

func (f *FixedAttributeGroup) appendToRowBuf(row int, buffer *bytes.Buffer) {
	for i, a := range f.attributes {
		postfix := " "
		if i == len(f.attributes)-1 {
			postfix = ""
		}
		buffer.WriteString(fmt.Sprintf("%s%s", a.GetStringFromSysVal(f.get(i, row)), postfix))
	}
}

func (f *FixedAttributeGroup) resize(add int) {
	newAlloc := make([]byte, len(f.alloc)+add)
	copy(newAlloc, f.alloc)
	f.alloc = newAlloc
}
