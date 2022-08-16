package base

import (
	"bytes"
	"fmt"
)

type BinaryAttributeGroup struct {
	parent     DataGrid
	attributes []Attribute
	size       int
	alloc      []byte
	maxRow     int
}

func (b *BinaryAttributeGroup) String() string {
	return "BinaryAttributeGroup"
}

func (b *BinaryAttributeGroup) RowSizeInBytes() int {
	return (len(b.attributes) + 7) / 8
}

func (b *BinaryAttributeGroup) Attributes() []Attribute {
	ret := make([]Attribute, len(b.attributes))
	for i, a := range b.attributes {
		ret[i] = a
	}
	return ret
}

func (b *BinaryAttributeGroup) AddAttribute(a Attribute) error {
	b.attributes = append(b.attributes, a)
	return nil
}

func (b *BinaryAttributeGroup) Storage() []byte {
	return b.alloc
}

func (b *BinaryAttributeGroup) setStorage(a []byte) {
	b.alloc = a
}

func (b *BinaryAttributeGroup) getByteOffset(col, row int) int {
	return row*b.RowSizeInBytes() + col/8
}

func (b *BinaryAttributeGroup) set(col, row int, val []byte) {

	offset := b.getByteOffset(col, row)

	if val[0] > 0 {
		b.alloc[offset] |= (1 << (uint(col) % 8))
	} else {

		b.alloc[offset] &= ^(1 << (uint(col) % 8))
	}

	row++
	if row > b.maxRow {
		b.maxRow = row
	}
}

func (b *BinaryAttributeGroup) get(col, row int) []byte {
	offset := b.getByteOffset(col, row)
	if b.alloc[offset]&(1<<(uint(col%8))) > 0 {
		return []byte{1}
	} else {
		return []byte{0}
	}
}

func (b *BinaryAttributeGroup) appendToRowBuf(row int, buffer *bytes.Buffer) {
	for i, a := range b.attributes {
		postfix := " "
		if i == len(b.attributes)-1 {
			postfix = ""
		}
		buffer.WriteString(fmt.Sprintf("%s%s",
			a.GetStringFromSysVal(b.get(i, row)), postfix))
	}
}

func (b *BinaryAttributeGroup) resize(add int) {
	newAlloc := make([]byte, len(b.alloc)+add)
	copy(newAlloc, b.alloc)
	b.alloc = newAlloc
}
