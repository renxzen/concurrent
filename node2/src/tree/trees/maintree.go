package trees

import (
	"fmt"
	"log"
	"math/rand"

	"nodo2/src/tree/base"
	"nodo2/src/tree/evaluation"
	"nodo2/src/tree/preprocessing"

	"time"
)

var Tree *ID3DecisionTree
var TrainData base.FixedDataGrid

func init() {
	Tree, TrainData = GetTree()
}

func GetTree() (*ID3DecisionTree, base.FixedDataGrid) {
	rand.Seed(time.Now().UnixNano())

	preprocessing.Preprocess()

	dataset, err := base.ParseCSVToInstances("datasets/encoded.csv", true)
	if err != nil {
		log.Fatal(err)
	}

	trainData, testData := base.InstancesTrainTestSplit(dataset, 0.80)
	tree := NewID3DecisionTree(0.8)

	err = tree.Fit(trainData)
	if err != nil {
		log.Fatal(err)
	}

	predictions, err := tree.Predict(testData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Decision Tree Performance (information gain)")

	cf, err := evaluation.GetConfusionMatrix(testData, predictions)
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to get confusion matrix: %s", err.Error()))
	}

	fmt.Println(evaluation.GetSummary(cf))

	return tree, trainData
}
