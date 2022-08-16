package base

import (
	"bytes"
	"fmt"
)

type InstancesView struct {
	src        FixedDataGrid
	attrs      []AttributeSpec
	rows       map[int]int
	classAttrs map[Attribute]bool
	maskRows   bool
}

func (v *InstancesView) addClassAttrsFromSrc(src FixedDataGrid) {
	for _, a := range src.AllClassAttributes() {
		matched := true
		if v.attrs != nil {
			matched = false
			for _, b := range v.attrs {
				if b.attr.Equals(a) {
					matched = true
				}
			}
		}
		if matched {
			v.classAttrs[a] = true
		}
	}
}

func (v *InstancesView) resolveRow(origRow int) int {
	if v.rows != nil {
		if newRow, ok := v.rows[origRow]; !ok {
			if v.maskRows {
				return -1
			}
		} else {
			return newRow
		}
	}

	return origRow

}

func NewInstancesViewFromRows(src FixedDataGrid, rows map[int]int) *InstancesView {
	ret := &InstancesView{
		src,
		nil,
		rows,
		make(map[Attribute]bool),
		false,
	}

	ret.addClassAttrsFromSrc(src)
	return ret
}

func NewInstancesViewFromVisible(src FixedDataGrid, rows []int, attrs []Attribute) *InstancesView {
	ret := &InstancesView{
		src,
		ResolveAttributes(src, attrs),
		make(map[int]int),
		make(map[Attribute]bool),
		true,
	}

	for i, a := range rows {
		ret.rows[i] = a
	}

	ret.addClassAttrsFromSrc(src)
	return ret
}

func NewInstancesViewFromAttrs(src FixedDataGrid, attrs []Attribute) *InstancesView {
	ret := &InstancesView{
		src,
		ResolveAttributes(src, attrs),
		nil,
		make(map[Attribute]bool),
		false,
	}

	ret.addClassAttrsFromSrc(src)
	return ret
}

func (v *InstancesView) GetAttribute(a Attribute) (AttributeSpec, error) {
	if a == nil {
		return AttributeSpec{}, fmt.Errorf("Attribute can't be nil")
	}

	if v.attrs == nil {
		return v.src.GetAttribute(a)
	}

	for _, r := range v.attrs {

		if r.GetAttribute().Equals(a) {
			return r, nil
		}
	}
	return AttributeSpec{}, fmt.Errorf("requested attribute has been filtered")
}

func (v *InstancesView) AllAttributes() []Attribute {

	if v.attrs == nil {
		return v.src.AllAttributes()
	}

	ret := make([]Attribute, len(v.attrs))

	for i, a := range v.attrs {
		ret[i] = a.GetAttribute()
	}

	return ret
}

func (v *InstancesView) AddClassAttribute(a Attribute) error {

	matched := false
	for _, r := range v.AllAttributes() {
		if r.Equals(a) {
			matched = true
		}
	}
	if !matched {
		return fmt.Errorf("Attribute has been filtered")
	}

	v.classAttrs[a] = true
	return nil
}

func (v *InstancesView) RemoveClassAttribute(a Attribute) error {
	v.classAttrs[a] = false
	return nil
}

func (v *InstancesView) AllClassAttributes() []Attribute {
	ret := make([]Attribute, 0)
	for a := range v.classAttrs {
		if v.classAttrs[a] {
			ret = append(ret, a)
		}
	}
	return ret
}

func (v *InstancesView) Get(as AttributeSpec, row int) []byte {

	row = v.resolveRow(row)
	if row == -1 {
		panic("Out of range")
	}
	return v.src.Get(as, row)
}

func (v *InstancesView) MapOverRows(as []AttributeSpec, rowFunc func([][]byte, int) (bool, error)) error {
	if v.maskRows {
		rowBuf := make([][]byte, len(as))
		for r := range v.rows {
			row := v.rows[r]
			for i, a := range as {
				rowBuf[i] = v.src.Get(a, row)
			}
			ok, err := rowFunc(rowBuf, r)
			if err != nil {
				return err
			}
			if !ok {
				break
			}
		}
		return nil
	} else {
		return v.src.MapOverRows(as, rowFunc)
	}
}

func (v *InstancesView) Size() (int, int) {

	hSize, vSize := v.src.Size()

	if v.attrs != nil {
		hSize = len(v.attrs)
	}

	if v.rows != nil {
		if v.maskRows {
			vSize = len(v.rows)
		} else if len(v.rows) > vSize {
			vSize = len(v.rows)
		}
	}
	return hSize, vSize
}

func (v *InstancesView) String() string {
	var buffer bytes.Buffer
	maxRows := 30

	as := ResolveAllAttributes(v)

	cols, rows := v.Size()
	buffer.WriteString("InstancesView with ")
	buffer.WriteString(fmt.Sprintf("%d row(s) ", rows))
	buffer.WriteString(fmt.Sprintf("%d attribute(s)\n", cols))
	if v.attrs != nil {
		buffer.WriteString("With defined Attribute view\n")
	}
	if v.rows != nil {
		buffer.WriteString("With defined Row view\n")
	}
	if v.maskRows {
		buffer.WriteString("Row masking on.\n")
	}
	buffer.WriteString("Attributes:\n")

	for _, a := range as {
		prefix := "\t"
		if v.classAttrs[a.attr] {
			prefix = "*\t"
		}
		buffer.WriteString(fmt.Sprintf("%s%s\n", prefix, a.attr))
	}

	if rows < maxRows {
		maxRows = rows
	}
	buffer.WriteString("Data:")
	for i := 0; i < maxRows; i++ {
		buffer.WriteString("\t")
		for _, a := range as {
			val := v.Get(a, i)
			buffer.WriteString(fmt.Sprintf("%s ", a.attr.GetStringFromSysVal(val)))
		}
		buffer.WriteString("\n")
	}

	missingRows := rows - maxRows
	if missingRows != 0 {
		buffer.WriteString(fmt.Sprintf("\t...\n%d row(s) undisplayed", missingRows))
	} else {
		buffer.WriteString("All rows displayed")
	}

	return buffer.String()
}

func (v *InstancesView) RowString(row int) string {
	var buffer bytes.Buffer
	as := ResolveAllAttributes(v)
	first := true
	for _, a := range as {
		val := v.Get(a, row)
		prefix := " "
		if first {
			prefix = ""
			first = false
		}
		buffer.WriteString(fmt.Sprintf("%s%s", prefix, a.attr.GetStringFromSysVal(val)))
	}
	return buffer.String()
}
