package trees

import (
	"bytes"
	"fmt"
	"sort"

	"nodo2/src/tree/base"
	"nodo2/src/tree/evaluation"
)

type NodeType int

const (
	LeafNode NodeType = 1
	RuleNode NodeType = 2
)

type RuleGenerator interface {
	GenerateSplitRule(base.FixedDataGrid) *DecisionTreeRule
}

type DecisionTreeRule struct {
	SplitAttr base.Attribute
	SplitVal  float64
}

func (d *DecisionTreeRule) String() string {
	if _, ok := d.SplitAttr.(*base.FloatAttribute); ok {
		return fmt.Sprintf("DecisionTreeRule(%s <= %f)", d.SplitAttr.GetName(), d.SplitVal)
	}
	return fmt.Sprintf("DecisionTreeRule(%s)", d.SplitAttr.GetName())
}

type DecisionTreeNode struct {
	Type      NodeType
	Children  map[string]*DecisionTreeNode
	ClassDist map[string]int
	Class     string
	ClassAttr base.Attribute
	SplitRule *DecisionTreeRule
}

func getClassAttr(from base.FixedDataGrid) base.Attribute {
	allClassAttrs := from.AllClassAttributes()
	return allClassAttrs[0]
}

func InferID3Tree(from base.FixedDataGrid, with RuleGenerator) *DecisionTreeNode {

	classes := base.GetClassDistribution(from)

	if len(classes) == 1 {
		maxClass := ""
		for i := range classes {
			maxClass = i
		}
		ret := &DecisionTreeNode{
			LeafNode,
			nil,
			classes,
			maxClass,
			getClassAttr(from),
			&DecisionTreeRule{nil, 0.0},
		}
		return ret
	}

	maxVal := 0
	maxClass := ""
	for i := range classes {
		if classes[i] > maxVal {
			maxClass = i
			maxVal = classes[i]
		}
	}

	cols, _ := from.Size()
	if cols == 2 {
		ret := &DecisionTreeNode{
			LeafNode,
			nil,
			classes,
			maxClass,
			getClassAttr(from),
			&DecisionTreeRule{nil, 0.0},
		}
		return ret
	}

	ret := &DecisionTreeNode{
		RuleNode,
		nil,
		classes,
		maxClass,
		getClassAttr(from),
		nil,
	}

	splitRule := with.GenerateSplitRule(from)
	if splitRule == nil {

		return ret
	}

	var splitInstances map[string]base.FixedDataGrid
	if _, ok := splitRule.SplitAttr.(*base.FloatAttribute); ok {
		splitInstances = base.DecomposeOnNumericAttributeThreshold(from,
			splitRule.SplitAttr, splitRule.SplitVal)
	} else {
		splitInstances = base.DecomposeOnAttributeValues(from, splitRule.SplitAttr)
	}

	ret.Children = make(map[string]*DecisionTreeNode)
	for k := range splitInstances {
		newInstances := splitInstances[k]
		ret.Children[k] = InferID3Tree(newInstances, with)
	}
	ret.SplitRule = splitRule
	return ret
}

func (d *DecisionTreeNode) getNestedString(level int) string {
	buf := bytes.NewBuffer(nil)
	tmp := bytes.NewBuffer(nil)
	for i := 0; i < level; i++ {
		tmp.WriteString("\t")
	}
	buf.WriteString(tmp.String())
	if d.Children == nil {
		buf.WriteString(fmt.Sprintf("Leaf(%s)", d.Class))
	} else {
		var keys []string
		buf.WriteString(fmt.Sprintf("Rule(%s)", d.SplitRule))
		for k := range d.Children {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			buf.WriteString("\n")
			buf.WriteString(tmp.String())
			buf.WriteString("\t")
			buf.WriteString(k)
			buf.WriteString("\n")
			buf.WriteString(d.Children[k].getNestedString(level + 1))
		}
	}
	return buf.String()
}

func (d *DecisionTreeNode) String() string {
	return d.getNestedString(0)
}

func computeAccuracy(predictions base.FixedDataGrid, from base.FixedDataGrid) float64 {
	cf, _ := evaluation.GetConfusionMatrix(from, predictions)
	return evaluation.GetAccuracy(cf)
}

func (d *DecisionTreeNode) Prune(using base.FixedDataGrid) {

	if d.Children == nil {
		return
	}
	if d.SplitRule == nil {
		return
	}

	sub := base.DecomposeOnAttributeValues(using, d.SplitRule.SplitAttr)
	for k := range d.Children {
		if sub[k] == nil {
			continue
		}
		subH, subV := sub[k].Size()
		if subH == 0 || subV == 0 {
			continue
		}
		d.Children[k].Prune(sub[k])
	}

	predictions, _ := d.Predict(using)
	baselineAccuracy := computeAccuracy(predictions, using)

	tmpChildren := d.Children
	d.Children = nil

	predictions, _ = d.Predict(using)
	newAccuracy := computeAccuracy(predictions, using)

	if newAccuracy < baselineAccuracy {
		d.Children = tmpChildren
	}
}

func (d *DecisionTreeNode) Predict(what base.FixedDataGrid) (base.FixedDataGrid, error) {
	predictions := base.GeneratePredictionVector(what)
	classAttr := getClassAttr(predictions)
	classAttrSpec, err := predictions.GetAttribute(classAttr)
	if err != nil {
		panic(err)
	}
	predAttrs := base.AttributeDifferenceReferences(what.AllAttributes(), predictions.AllClassAttributes())
	predAttrSpecs := base.ResolveAttributes(what, predAttrs)
	what.MapOverRows(predAttrSpecs, func(row [][]byte, rowNo int) (bool, error) {
		cur := d
		for {
			if cur.Children == nil {
				predictions.Set(classAttrSpec, rowNo, classAttr.GetSysValFromString(cur.Class))
				break
			} else {
				splitVal := cur.SplitRule.SplitVal
				at := cur.SplitRule.SplitAttr
				ats, err := what.GetAttribute(at)
				if err != nil {
					panic(err)
				}

				var classVar string
				if _, ok := ats.GetAttribute().(*base.FloatAttribute); ok {

					classVal := base.UnpackBytesToFloat(what.Get(ats, rowNo))
					if classVal > splitVal {
						classVar = "1"
					} else {
						classVar = "0"
					}
				} else {
					classVar = ats.GetAttribute().GetStringFromSysVal(what.Get(ats, rowNo))
				}
				if next, ok := cur.Children[classVar]; ok {
					cur = next
				} else {

					var bestChild string
					for c := range cur.Children {
						bestChild = c
						if c > classVar {
							break
						}
					}
					cur = cur.Children[bestChild]
				}
			}
		}
		return true, nil
	})
	return predictions, nil
}

type ID3DecisionTree struct {
	base.BaseClassifier
	Root       *DecisionTreeNode
	PruneSplit float64
	Rule       RuleGenerator
}

func NewID3DecisionTree(prune float64) *ID3DecisionTree {
	return &ID3DecisionTree{
		base.BaseClassifier{},
		nil,
		prune,
		new(InformationGainRuleGenerator),
	}
}

func NewID3DecisionTreeFromRule(prune float64, rule RuleGenerator) *ID3DecisionTree {
	return &ID3DecisionTree{
		base.BaseClassifier{},
		nil,
		prune,
		rule,
	}
}

func (t *ID3DecisionTree) Fit(on base.FixedDataGrid) error {
	if t.PruneSplit > 0.001 {
		trainData, testData := base.InstancesTrainTestSplit(on, t.PruneSplit)
		t.Root = InferID3Tree(trainData, t.Rule)
		t.Root.Prune(testData)
	} else {
		t.Root = InferID3Tree(on, t.Rule)
	}
	return nil
}

func (t *ID3DecisionTree) Predict(what base.FixedDataGrid) (base.FixedDataGrid, error) {
	return t.Root.Predict(what)
}

func (t *ID3DecisionTree) String() string {
	return fmt.Sprintf("ID3DecisionTree(%s\n)", t.Root)
}
