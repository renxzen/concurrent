package response

import (
	"encoding/json"
	"net/http"
)

func SendResponse(w http.ResponseWriter, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Code)
	json.NewEncoder(w).Encode(response)
}

func SendData(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
