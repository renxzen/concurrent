package base

import (
	"fmt"
)

func NonClassFloatAttributes(d DataGrid) []Attribute {
	classAttrs := d.AllClassAttributes()
	allAttrs := d.AllAttributes()
	ret := make([]Attribute, 0)
	for _, a := range allAttrs {
		matched := false
		if _, ok := a.(*FloatAttribute); !ok {
			continue
		}
		for _, b := range classAttrs {
			if a.Equals(b) {
				matched = true
				break
			}
		}
		if !matched {
			ret = append(ret, a)
		}
	}
	return ret
}

func NonClassAttributes(d DataGrid) []Attribute {
	classAttrs := d.AllClassAttributes()
	allAttrs := d.AllAttributes()
	return AttributeDifferenceReferences(allAttrs, classAttrs)
}

func ResolveAttributes(d DataGrid, attrs []Attribute) []AttributeSpec {
	ret := make([]AttributeSpec, len(attrs))
	for i, a := range attrs {
		spec, err := d.GetAttribute(a)
		if err != nil {
			panic(fmt.Errorf("error resolving attribute %s: %s", a, err))
		}
		ret[i] = spec
	}
	return ret
}

func ResolveAllAttributes(d DataGrid) []AttributeSpec {
	return ResolveAttributes(d, d.AllAttributes())
}

func buildAttrSet(a []Attribute) map[Attribute]bool {
	ret := make(map[Attribute]bool)
	for _, a := range a {
		ret[a] = true
	}
	return ret
}

func AttributeIntersect(a1, a2 []Attribute) []Attribute {
	ret := make([]Attribute, 0)
	for _, a := range a1 {
		matched := false
		for _, b := range a2 {
			if a.Equals(b) {
				matched = true
				break
			}
		}
		if matched {
			ret = append(ret, a)
		}
	}
	return ret
}

func AttributeIntersectReferences(a1, a2 []Attribute) []Attribute {
	a1b := buildAttrSet(a1)
	a2b := buildAttrSet(a2)
	ret := make([]Attribute, 0)
	for a := range a1b {
		if _, ok := a2b[a]; ok {
			ret = append(ret, a)
		}
	}
	return ret
}

func AttributeDifference(a1, a2 []Attribute) []Attribute {
	ret := make([]Attribute, 0)
	for _, a := range a1 {
		matched := false
		for _, b := range a2 {
			if a.Equals(b) {
				matched = true
				break
			}
		}
		if !matched {
			ret = append(ret, a)
		}
	}
	return ret
}

func AttributeDifferenceReferences(a1, a2 []Attribute) []Attribute {
	a1b := buildAttrSet(a1)
	a2b := buildAttrSet(a2)
	ret := make([]Attribute, 0)
	for a := range a1b {
		if _, ok := a2b[a]; !ok {
			ret = append(ret, a)
		}
	}
	return ret
}
