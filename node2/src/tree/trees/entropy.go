package trees

import (
	"math"
	"sort"

	"nodo2/src/tree/base"
)

type InformationGainRuleGenerator struct {
}

func (r *InformationGainRuleGenerator) GenerateSplitRule(f base.FixedDataGrid) *DecisionTreeRule {
	attrs := f.AllAttributes()
	classAttrs := f.AllClassAttributes()
	candidates := base.AttributeDifferenceReferences(attrs, classAttrs)

	return r.GetSplitRuleFromSelection(candidates, f)
}

func (r *InformationGainRuleGenerator) GetSplitRuleFromSelection(consideredAttributes []base.Attribute, f base.FixedDataGrid) *DecisionTreeRule {
	var selectedAttribute base.Attribute

	if len(consideredAttributes) == 0 {
		panic("More Attributes should be considered")
	}

	maxGain := math.Inf(-1)
	selectedVal := math.Inf(1)

	classDist := base.GetClassDistribution(f)
	baseEntropy := getBaseEntropy(classDist)

	for _, s := range consideredAttributes {
		var informationGain float64
		var splitVal float64
		if fAttr, ok := s.(*base.FloatAttribute); ok {
			var attributeEntropy float64
			attributeEntropy, splitVal = getNumericAttributeEntropy(f, fAttr)
			informationGain = baseEntropy - attributeEntropy
		} else {
			proposedClassDist := base.GetClassDistributionAfterSplit(f, s)
			localEntropy := getSplitEntropy(proposedClassDist)
			informationGain = baseEntropy - localEntropy
		}

		if informationGain > maxGain {
			maxGain = informationGain
			selectedAttribute = s
			selectedVal = splitVal
		}
	}

	return &DecisionTreeRule{selectedAttribute, selectedVal}
}

type numericSplitRef struct {
	val   float64
	class string
}

type splitVec []numericSplitRef

func (a splitVec) Len() int           { return len(a) }
func (a splitVec) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a splitVec) Less(i, j int) bool { return a[i].val < a[j].val }

func getNumericAttributeEntropy(f base.FixedDataGrid, attr *base.FloatAttribute) (float64, float64) {
	attrSpec, err := f.GetAttribute(attr)
	if err != nil {
		panic(err)
	}

	_, rows := f.Size()
	refs := make([]numericSplitRef, rows)
	f.MapOverRows([]base.AttributeSpec{attrSpec}, func(val [][]byte, row int) (bool, error) {
		cls := base.GetClass(f, row)
		v := base.UnpackBytesToFloat(val[0])
		refs[row] = numericSplitRef{v, cls}
		return true, nil
	})

	sort.Sort(splitVec(refs))

	generateCandidateSplitDistribution := func(val float64) map[string]map[string]int {
		presplit := make(map[string]int)
		postplit := make(map[string]int)
		for _, i := range refs {
			if i.val < val {
				presplit[i.class]++
			} else {
				postplit[i.class]++
			}
		}
		ret := make(map[string]map[string]int)
		ret["0"] = presplit
		ret["1"] = postplit
		return ret
	}

	minSplitEntropy := math.Inf(1)
	minSplitVal := math.Inf(1)

	for i := 0; i < len(refs)-1; i++ {
		val := refs[i].val + refs[i+1].val
		val /= 2
		splitDist := generateCandidateSplitDistribution(val)
		splitEntropy := getSplitEntropy(splitDist)
		if splitEntropy < minSplitEntropy {
			minSplitEntropy = splitEntropy
			minSplitVal = val
		}
	}

	return minSplitEntropy, minSplitVal
}

func getSplitEntropy(s map[string]map[string]int) float64 {
	ret := 0.0
	count := 0
	for a := range s {
		for c := range s[a] {
			count += s[a][c]
		}
	}
	for a := range s {
		total := 0.0
		for c := range s[a] {
			total += float64(s[a][c])
		}
		for c := range s[a] {
			ret -= float64(s[a][c]) / float64(count) * math.Log(float64(s[a][c])/float64(count)) / math.Log(2)
		}
		ret += total / float64(count) * math.Log(total/float64(count)) / math.Log(2)
	}
	return ret
}

func getBaseEntropy(s map[string]int) float64 {
	ret := 0.0
	count := 0
	for k := range s {
		count += s[k]
	}
	for k := range s {
		ret -= float64(s[k]) / float64(count) * math.Log(float64(s[k])/float64(count)) / math.Log(2)
	}
	return ret
}
