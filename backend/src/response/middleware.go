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

		// CORS(&w)
		next(w, r)
	})
}

func CORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Content-Type, Accept, Jwt-Token, Authorization, Origin, Accept, X-Requested-With, Access-Control-Request-Method, Access-Control-Request-Headers")
	(*w).Header().Set("Access-Control-Expose-Headers", "Origin, Content-Type, Accept, Jwt-Token, Authorization, Access-Control-Allow-Origin, Access-Control-Allow-Origin, Access-Control-Allow-Credentials")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
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
