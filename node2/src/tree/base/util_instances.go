package base

import (
	"fmt"
	"math/rand"
)

func GeneratePredictionVector(from FixedDataGrid) UpdatableDataGrid {
	classAttrs := from.AllClassAttributes()
	_, rowCount := from.Size()
	ret := NewDenseInstances()
	for _, a := range classAttrs {
		ret.AddAttribute(a)
		ret.AddClassAttribute(a)
	}
	ret.Extend(rowCount)
	return ret
}

func GetClass(from DataGrid, row int) string {

	classAttrs := from.AllClassAttributes()
	if len(classAttrs) > 1 {
		panic("More than one class defined")
	} else if len(classAttrs) == 0 {
		panic("No class defined!")
	}
	classAttr := classAttrs[0]

	classAttrSpec, err := from.GetAttribute(classAttr)
	if err != nil {
		panic(fmt.Errorf("can't resolve class attribute %s", err))
	}

	classVal := from.Get(classAttrSpec, row)
	if classVal == nil {
		panic("Class values shouldn't be missing")
	}

	return classAttr.GetStringFromSysVal(classVal)
}

func SetClass(at UpdatableDataGrid, row int, class string) {

	classAttrs := at.AllClassAttributes()
	if len(classAttrs) > 1 {
		panic("More than one class defined")
	} else if len(classAttrs) == 0 {
		panic("No class Attributes are defined")
	}

	classAttr := classAttrs[0]

	classAttrSpec, err := at.GetAttribute(classAttr)
	if err != nil {
		panic(fmt.Errorf("can't resolve class attribute %s", err))
	}

	classBytes := classAttr.GetSysValFromString(class)
	at.Set(classAttrSpec, row, classBytes)
}

func GetAttributeByName(inst FixedDataGrid, name string) Attribute {
	for _, a := range inst.AllAttributes() {
		if a.GetName() == name {
			return a
		}
	}
	return nil
}

func GetClassDistribution(inst FixedDataGrid) map[string]int {
	ret := make(map[string]int)
	_, rows := inst.Size()
	for i := 0; i < rows; i++ {
		cls := GetClass(inst, i)
		ret[cls]++
	}
	return ret
}

func GetClassDistributionAfterThreshold(inst FixedDataGrid, at Attribute, val float64) map[string]map[string]int {
	ret := make(map[string]map[string]int)

	attrSpec, err := inst.GetAttribute(at)
	if err != nil {
		panic(fmt.Sprintf("Invalid attribute %s (%s)", at, err))
	}

	if _, ok := at.(*FloatAttribute); !ok {
		panic("Must be numeric!")
	}

	_, rows := inst.Size()

	for i := 0; i < rows; i++ {
		splitVal := UnpackBytesToFloat(inst.Get(attrSpec, i)) > val
		splitVar := "0"
		if splitVal {
			splitVar = "1"
		}
		classVar := GetClass(inst, i)
		if _, ok := ret[splitVar]; !ok {
			ret[splitVar] = make(map[string]int)
			i--
			continue
		}
		ret[splitVar][classVar]++
	}

	return ret
}

func GetClassDistributionAfterSplit(inst FixedDataGrid, at Attribute) map[string]map[string]int {

	ret := make(map[string]map[string]int)

	attrSpec, err := inst.GetAttribute(at)
	if err != nil {
		panic(fmt.Sprintf("Invalid attribute %s (%s)", at, err))
	}

	_, rows := inst.Size()

	for i := 0; i < rows; i++ {
		splitVar := at.GetStringFromSysVal(inst.Get(attrSpec, i))
		classVar := GetClass(inst, i)
		if _, ok := ret[splitVar]; !ok {
			ret[splitVar] = make(map[string]int)
			i--
			continue
		}
		ret[splitVar][classVar]++
	}

	return ret
}

func DecomposeOnNumericAttributeThreshold(inst FixedDataGrid, at Attribute, val float64) map[string]FixedDataGrid {

	if _, ok := at.(*FloatAttribute); !ok {
		panic("Invalid argument")
	}

	attrSpec, err := inst.GetAttribute(at)
	if err != nil {
		panic(fmt.Sprintf("Invalid Attribute index %s", at))
	}

	newAttrs := make([]Attribute, 0)
	for _, a := range inst.AllAttributes() {
		if a.Equals(at) {
			continue
		}
		newAttrs = append(newAttrs, a)
	}

	ret := make(map[string]FixedDataGrid)

	rowMaps := make(map[string][]int)

	fullAttrSpec := ResolveAttributes(inst, newAttrs)
	fullAttrSpec = append(fullAttrSpec, attrSpec)

	inst.MapOverRows(fullAttrSpec, func(row [][]byte, rowNo int) (bool, error) {

		targetBytes := row[len(row)-1]
		targetVal := UnpackBytesToFloat(targetBytes)
		val := targetVal > val
		targetSet := "0"
		if val {
			targetSet = "1"
		}
		rowMap := rowMaps[targetSet]
		rowMaps[targetSet] = append(rowMap, rowNo)
		return true, nil
	})

	for a := range rowMaps {
		ret[a] = NewInstancesViewFromVisible(inst, rowMaps[a], newAttrs)
	}

	return ret
}

func DecomposeOnAttributeValues(inst FixedDataGrid, at Attribute) map[string]FixedDataGrid {

	attrSpec, err := inst.GetAttribute(at)
	if err != nil {
		panic(fmt.Sprintf("Invalid Attribute index %s", at))
	}

	newAttrs := make([]Attribute, 0)
	for _, a := range inst.AllAttributes() {
		if a.Equals(at) {
			continue
		}
		newAttrs = append(newAttrs, a)
	}

	ret := make(map[string]FixedDataGrid)

	rowMaps := make(map[string][]int)

	fullAttrSpec := ResolveAttributes(inst, newAttrs)
	fullAttrSpec = append(fullAttrSpec, attrSpec)

	inst.MapOverRows(fullAttrSpec, func(row [][]byte, rowNo int) (bool, error) {

		targetBytes := row[len(row)-1]
		targetAttr := fullAttrSpec[len(fullAttrSpec)-1].attr
		targetSet := targetAttr.GetStringFromSysVal(targetBytes)
		if _, ok := rowMaps[targetSet]; !ok {
			rowMaps[targetSet] = make([]int, 0)
		}
		rowMap := rowMaps[targetSet]
		rowMaps[targetSet] = append(rowMap, rowNo)
		return true, nil
	})

	for a := range rowMaps {
		ret[a] = NewInstancesViewFromVisible(inst, rowMaps[a], newAttrs)
	}

	return ret
}

func InstancesTrainTestSplit(src FixedDataGrid, prop float64) (FixedDataGrid, FixedDataGrid) {
	trainingRows := make([]int, 0)
	testingRows := make([]int, 0)
	src = Shuffle(src)

	_, rows := src.Size()
	for i := 0; i < rows; i++ {
		trainOrTest := rand.Intn(101)
		if trainOrTest > int(100*prop) {
			trainingRows = append(trainingRows, i)
		} else {
			testingRows = append(testingRows, i)
		}
	}

	allAttrs := src.AllAttributes()

	return NewInstancesViewFromVisible(src, trainingRows, allAttrs), NewInstancesViewFromVisible(src, testingRows, allAttrs)

}

func LazyShuffle(from FixedDataGrid) FixedDataGrid {
	_, rows := from.Size()
	rowMap := make(map[int]int)
	for i := 0; i < rows; i++ {
		j := rand.Intn(i + 1)
		rowMap[i] = j
		rowMap[j] = i
	}
	return NewInstancesViewFromRows(from, rowMap)
}

func Shuffle(from FixedDataGrid) FixedDataGrid {
	_, rows := from.Size()
	if inst, ok := from.(*DenseInstances); ok {
		for i := 0; i < rows; i++ {
			j := rand.Intn(i + 1)
			inst.swapRows(i, j)
		}
		return inst
	} else {
		return LazyShuffle(from)
	}
}

func SampleWithReplacement(from FixedDataGrid, size int) FixedDataGrid {
	rowMap := make(map[int]int)
	_, rows := from.Size()
	for i := 0; i < size; i++ {
		srcRow := rand.Intn(rows)
		rowMap[i] = srcRow
	}
	return NewInstancesViewFromRows(from, rowMap)
}

func CheckCompatible(s1 FixedDataGrid, s2 FixedDataGrid) []Attribute {
	s1Attrs := s1.AllAttributes()
	s2Attrs := s2.AllAttributes()
	interAttrs := AttributeIntersect(s1Attrs, s2Attrs)
	if len(interAttrs) != len(s1Attrs) {
		return nil
	} else if len(interAttrs) != len(s2Attrs) {
		return nil
	}
	return interAttrs
}

func CheckStrictlyCompatible(s1 FixedDataGrid, s2 FixedDataGrid) bool {

	d1, ok1 := s1.(*DenseInstances)
	d2, ok2 := s2.(*DenseInstances)
	if !ok1 || !ok2 {
		return false
	}

	d1ags := d1.AllAttributeGroups()
	d2ags := d2.AllAttributeGroups()

	for a := range d1ags {
		_, ok := d2ags[a]
		if !ok {
			return false
		}
	}

	for a := range d2ags {
		_, ok := d1ags[a]
		if !ok {
			return false
		}
	}

	for a := range d1ags {
		ag1 := d1ags[a]
		ag2 := d2ags[a]
		a1 := ag1.Attributes()
		a2 := ag2.Attributes()
		for i := range a1 {
			at1 := a1[i]
			at2 := a2[i]
			if !at1.Equals(at2) {
				return false
			}
		}
	}

	return true
}
