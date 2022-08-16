package base

type Classifier interface {
	Predict(FixedDataGrid) (FixedDataGrid, error)
	Fit(FixedDataGrid) error
	String() string
}

type BaseClassifier struct {
	TrainingData *DataGrid
}
