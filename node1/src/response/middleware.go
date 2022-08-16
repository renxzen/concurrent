package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func checkMethod(next http.HandlerFunc, method string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			message := fmt.Sprintf("metodo %s no permitido", r.Method)
			res := NewResponse(http.StatusMethodNotAllowed, message)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(res)
			return
		}

		next(w, r)
	})
}

func GET(next http.HandlerFunc) http.HandlerFunc {
	return checkMethod(next, http.MethodGet)
}

func POST(next http.HandlerFunc) http.HandlerFunc {
	return checkMethod(next, http.MethodPost)
}

func PUT(next http.HandlerFunc) http.HandlerFunc {
	return checkMethod(next, http.MethodPut)
}

func DELETE(next http.HandlerFunc) http.HandlerFunc {
	return checkMethod(next, http.MethodDelete)
}
