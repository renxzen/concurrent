package base

import (
	"fmt"
)

type AttributeSpec struct {
	pond     int
	position int
	attr     Attribute
}

func (a *AttributeSpec) GetAttribute() Attribute {
	return a.attr
}

func (a *AttributeSpec) String() string {
	return fmt.Sprintf("AttributeSpec(Attribute: '%s', Pond: %d/%d)", a.attr, a.pond, a.position)
}
