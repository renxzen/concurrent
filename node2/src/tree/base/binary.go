package base

import (
	"fmt"
	"strconv"
)

type BinaryAttribute struct {
	Name string
}

func NewBinaryAttribute(name string) *BinaryAttribute {
	return &BinaryAttribute{
		name,
	}
}

func (b *BinaryAttribute) GetName() string {
	return b.Name
}

func (b *BinaryAttribute) SetName(name string) {
	b.Name = name
}

func (b *BinaryAttribute) GetType() int {
	return BinaryType
}

func (b *BinaryAttribute) GetSysValFromString(userVal string) []byte {
	f, err := strconv.ParseFloat(userVal, 64)
	if err != nil {
		panic(err)
	}
	ret := make([]byte, 1)
	if f > 0 {
		ret[0] = 1
	}
	return ret
}

func (b *BinaryAttribute) GetStringFromSysVal(val []byte) string {
	if val[0] > 0 {
		return "1"
	}
	return "0"
}

func (b *BinaryAttribute) Equals(other Attribute) bool {
	if a, ok := other.(*BinaryAttribute); !ok {
		return false
	} else {
		return a.Name == b.Name
	}
}

func (b *BinaryAttribute) Compatible(other Attribute) bool {
	if _, ok := other.(*BinaryAttribute); !ok {
		return false
	} else {
		return true
	}
}

func (b *BinaryAttribute) String() string {
	return fmt.Sprintf("BinaryAttribute(%s)", b.Name)
}
