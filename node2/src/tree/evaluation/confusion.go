package evaluation

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"nodo2/src/tree/base"
)

type ConfusionMatrix map[string]map[string]int

func GetConfusionMatrix(ref base.FixedDataGrid, gen base.FixedDataGrid) (map[string]map[string]int, error) {
	_, refRows := ref.Size()
	_, genRows := gen.Size()

	if refRows != genRows {
		return nil, fmt.Errorf("row count mismatch: ref has %d rows, gen has %d rows", refRows, genRows)
	}

	ret := make(map[string]map[string]int)

	for i := 0; i < int(refRows); i++ {
		referenceClass := base.GetClass(ref, i)
		predictedClass := base.GetClass(gen, i)
		if _, ok := ret[referenceClass]; ok {
			ret[referenceClass][predictedClass] += 1
		} else {
			ret[referenceClass] = make(map[string]int)
			ret[referenceClass][predictedClass] = 1
		}
	}
	return ret, nil
}

func GetTruePositives(class string, c ConfusionMatrix) float64 {
	return float64(c[class][class])
}

func GetFalsePositives(class string, c ConfusionMatrix) float64 {
	ret := 0.0
	for k := range c {
		if k == class {
			continue
		}
		ret += float64(c[k][class])
	}
	return ret
}

func GetFalseNegatives(class string, c ConfusionMatrix) float64 {
	ret := 0.0
	for k := range c[class] {
		if k == class {
			continue
		}
		ret += float64(c[class][k])
	}
	return ret
}

func GetTrueNegatives(class string, c ConfusionMatrix) float64 {
	ret := 0.0
	for k := range c {
		if k == class {
			continue
		}
		for l := range c[k] {
			if l == class {
				continue
			}
			ret += float64(c[k][l])
		}
	}
	return ret
}

func GetPrecision(class string, c ConfusionMatrix) float64 {

	truePositives := GetTruePositives(class, c)
	falsePositives := GetFalsePositives(class, c)
	return truePositives / (truePositives + falsePositives)
}

func GetRecall(class string, c ConfusionMatrix) float64 {

	truePositives := GetTruePositives(class, c)
	falseNegatives := GetFalseNegatives(class, c)
	return truePositives / (truePositives + falseNegatives)
}

func GetF1Score(class string, c ConfusionMatrix) float64 {
	precision := GetPrecision(class, c)
	recall := GetRecall(class, c)
	return 2 * (precision * recall) / (precision + recall)
}

func GetAccuracy(c ConfusionMatrix) float64 {
	correct := 0
	total := 0
	for i := range c {
		for j := range c[i] {
			if i == j {
				correct += c[i][j]
			}
			total += c[i][j]
		}
	}
	return float64(correct) / float64(total)
}

func GetMicroPrecision(c ConfusionMatrix) float64 {
	truePositives := 0.0
	falsePositives := 0.0
	for k := range c {
		truePositives += GetTruePositives(k, c)
		falsePositives += GetFalsePositives(k, c)
	}
	return truePositives / (truePositives + falsePositives)
}

func GetMacroPrecision(c ConfusionMatrix) float64 {
	precisionVals := 0.0
	for k := range c {
		precisionVals += GetPrecision(k, c)
	}
	return precisionVals / float64(len(c))
}

func GetMicroRecall(c ConfusionMatrix) float64 {
	truePositives := 0.0
	falseNegatives := 0.0
	for k := range c {
		truePositives += GetTruePositives(k, c)
		falseNegatives += GetFalseNegatives(k, c)
	}
	return truePositives / (truePositives + falseNegatives)
}

func GetMacroRecall(c ConfusionMatrix) float64 {
	recallVals := 0.0
	for k := range c {
		recallVals += GetRecall(k, c)
	}
	return recallVals / float64(len(c))
}

func GetSummary(c ConfusionMatrix) string {
	var buffer bytes.Buffer
	w := new(tabwriter.Writer)
	w.Init(&buffer, 0, 8, 0, '\t', 0)

	fmt.Fprintln(w, "Reference Class\tTrue Positives\tFalse Positives\tTrue Negatives\tPrecision\tRecall\tF1 Score")
	fmt.Fprintln(w, "---------------\t--------------\t---------------\t--------------\t---------\t------\t--------")
	for k := range c {
		tp := GetTruePositives(k, c)
		fp := GetFalsePositives(k, c)
		tn := GetTrueNegatives(k, c)
		prec := GetPrecision(k, c)
		rec := GetRecall(k, c)
		f1 := GetF1Score(k, c)

		fmt.Fprintf(w, "%s\t%.0f\t%.0f\t%.0f\t%.4f\t%.4f\t%.4f\n", k, tp, fp, tn, prec, rec, f1)
	}
	w.Flush()
	buffer.WriteString(fmt.Sprintf("Overall accuracy: %.4f\n", GetAccuracy(c)))

	return buffer.String()
}
