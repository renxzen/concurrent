package base

import (
	"fmt"
)

type CategoricalAttribute struct {
	Name   string
	values []string
}

func NewCategoricalAttribute() *CategoricalAttribute {
	return &CategoricalAttribute{
		"",
		make([]string, 0),
	}
}

func (Attr *CategoricalAttribute) GetValues() []string {
	return Attr.values
}

func (Attr *CategoricalAttribute) GetName() string {
	return Attr.Name
}

func (Attr *CategoricalAttribute) SetName(name string) {
	Attr.Name = name
}

func (Attr *CategoricalAttribute) GetType() int {
	return CategoricalType
}

func (Attr *CategoricalAttribute) GetSysVal(userVal string) []byte {
	for idx, val := range Attr.values {
		if val == userVal {
			return PackU64ToBytes(uint64(idx))
		}
	}
	return nil
}

func (Attr *CategoricalAttribute) GetUsrVal(sysVal []byte) string {
	idx := UnpackBytesToU64(sysVal)
	return Attr.values[idx]
}

func (Attr *CategoricalAttribute) GetSysValFromString(rawVal string) []byte {

	catIndex := -1
	for i, s := range Attr.values {
		if s == rawVal {
			catIndex = i
			break
		}
	}
	if catIndex == -1 {
		Attr.values = append(Attr.values, rawVal)
		catIndex = len(Attr.values) - 1
	}

	ret := PackU64ToBytes(uint64(catIndex))
	return ret
}

func (Attr *CategoricalAttribute) String() string {
	return fmt.Sprintf("CategoricalAttribute(\"%s\", %s)", Attr.Name, Attr.values)
}

func (Attr *CategoricalAttribute) GetStringFromSysVal(rawVal []byte) string {
	convVal := int(UnpackBytesToU64(rawVal))
	if convVal >= len(Attr.values) {
		panic(fmt.Sprintf("Out of range: %d in %d (%s)", convVal, len(Attr.values), Attr))
	}
	return Attr.values[convVal]
}

func (Attr *CategoricalAttribute) Equals(other Attribute) bool {
	attribute, ok := other.(*CategoricalAttribute)
	if !ok {

		return false
	}
	if Attr.GetName() != attribute.GetName() {
		return false
	}

	if len(attribute.values) != len(Attr.values) {
		return false
	}

	for i, a := range Attr.values {
		if a != attribute.values[i] {
			return false
		}
	}

	return true
}

func (Attr *CategoricalAttribute) Compatible(other Attribute) bool {
	attribute, ok := other.(*CategoricalAttribute)
	if !ok {
		return false
	}

	if len(attribute.values) != len(Attr.values) {
		return false
	}

	for i, a := range Attr.values {
		if a != attribute.values[i] {
			return false
		}
	}

	return true
}
