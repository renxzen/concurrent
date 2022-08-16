package preprocessing

import (
	"fmt"
	"log"

	"net/http"
	"os"
	"strconv"

	"encoding/csv"
)

func Preprocess() {
	// resp, err := http.Get("https://raw.githubusercontent.com/renxzen/concurrent/main/backend/dataset/TB_HOSP_VAC_FALLECIDOS.csv")
	// file, err := os.Open("../../../../backend/dataset/TB_HOSP_VAC_FALLECIDOS.csv")
	// defer file.Close()
	// reader := csv.NewReader(file)

	resp, err := http.Get("https://raw.githubusercontent.com/renxzen/concurrent/main/backend/dataset/TB_HOSP_VAC_FALLECIDOS.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println("File downloaded. Preprocessing...")

	reader := csv.NewReader(resp.Body)

	csvBody, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	csvBody = csvBody[1:]

	indexes := []int{6, 7, 18, 20, 23}
	headers := []interface{}{"edad", "sexo", "fdosis1", "fdosis2", "fallecido"}
	data := make([][]interface{}, len(csvBody))

	dict := make(map[string]interface{})
	dict["F"] = 0
	dict["M"] = 1
	dict[""] = 0
	dict["PFIZER"] = 1
	dict["SINOPHARM"] = 2
	dict["ASTRAZENECA"] = 3
	// dict["0"] = "Fallecido"
	// dict["1"] = "Vivo"

	for i, row := range csvBody {
		array := make([]interface{}, len(indexes))

		for j, index := range indexes {
			if index != 6 && index != 23 {
				array[j] = dict[row[index]]
				continue
			}

			array[j], _ = strconv.Atoi(row[index])
		}

		data[i] = array
	}

	data = append([][]interface{}{headers}, data...)

	outputFile, err := os.Create("datasets/encoded.csv")
	if err != nil {
		log.Fatal("Unable to read input file", err)
	}

	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)

	defer writer.Flush()

	for _, value := range data {
		row := make([]string, len(value))
		for i, col := range value {
			row[i] = fmt.Sprintf("%v", col)
		}

		err := writer.Write(row)
		if err != nil {
			log.Fatal("Cannot write to file", err)
		}
	}

	fmt.Println("File saved.")
}
