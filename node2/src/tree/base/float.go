package base

import (
	"fmt"
	"strconv"
)

type FloatAttribute struct {
	Name      string
	Precision int
}

func NewFloatAttribute(name string) *FloatAttribute {
	return &FloatAttribute{name, 2}
}

func (Attr *FloatAttribute) Compatible(other Attribute) bool {
	_, ok := other.(*FloatAttribute)
	return ok
}

func (Attr *FloatAttribute) Equals(other Attribute) bool {
	_, ok := other.(*FloatAttribute)
	if !ok {

		return false
	}
	if Attr.GetName() != other.GetName() {
		return false
	}
	return true
}

func (Attr *FloatAttribute) GetName() string {
	return Attr.Name
}

func (Attr *FloatAttribute) SetName(name string) {
	Attr.Name = name
}

func (Attr *FloatAttribute) GetType() int {
	return Float64Type
}

func (Attr *FloatAttribute) String() string {
	return fmt.Sprintf("FloatAttribute(%s)", Attr.Name)
}

func (Attr *FloatAttribute) CheckSysValFromString(rawVal string) ([]byte, error) {
	f, err := strconv.ParseFloat(rawVal, 64)
	if err != nil {
		return nil, err
	}

	ret := PackFloatToBytes(f)
	return ret, nil
}

func (Attr *FloatAttribute) GetSysValFromString(rawVal string) []byte {
	f, err := Attr.CheckSysValFromString(rawVal)
	if err != nil {
		panic(err)
	}
	return f
}

func (Attr *FloatAttribute) GetFloatFromSysVal(rawVal []byte) float64 {
	return UnpackBytesToFloat(rawVal)
}

func (Attr *FloatAttribute) GetStringFromSysVal(rawVal []byte) string {
	f := UnpackBytesToFloat(rawVal)
	formatString := fmt.Sprintf("%%.%df", Attr.Precision)
	return fmt.Sprintf(formatString, f)
}
