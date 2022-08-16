package trees

import (
	"fmt"
	"log"

	"nodo2/src/tree/base"
	"nodo2/src/tree/model"
)

func GetSinglePrediction(request model.Info) string {
	age := fmt.Sprintf("%v", request.Age)
	gender := fmt.Sprintf("%v", request.Gender)
	firstVaccine := fmt.Sprintf("%v", request.FirstVaccine)
	secondVaccine := fmt.Sprintf("%v", request.SecondVaccine)

	newData := []string{age, gender, firstVaccine, secondVaccine, "0"}

	instances := base.NewStructuralCopy(TrainData)
	instances.Extend(1)

	n, _ := instances.Size()
	for i := 0; i < n; i++ {
		spec := base.ResolveAttributes(instances, instances.AllAttributes())[i]
		instances.Set(spec, 0, spec.GetAttribute().GetSysValFromString(newData[i]))
	}

	prediction, err := Tree.Predict(instances)
	if err != nil {
		log.Fatal(err)
	}

	// _, n := prediction.Size()
	// for i := 0; i < n; i++ {
	// 	spec := base.ResolveAttributes(prediction, prediction.AllAttributes())[i]
	// 	fmt.Println("Value:", spec.GetAttribute().GetStringFromSysVal(prediction.Get(spec, 0)))
	// }

	spec := base.ResolveAttributes(prediction, prediction.AllAttributes())[0]
	return spec.GetAttribute().GetStringFromSysVal(prediction.Get(spec, 0))
}
