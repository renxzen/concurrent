package info

import (
	"net/http"
	res "nodo1/src/response"
)

func CheckJWT(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")

		_, err := JWT.CheckToken(authorizationHeader)
		if err != nil {
			response := res.NewResponse(http.StatusMethodNotAllowed, "invalid token")
			res.SendResponse(w, response)
			return
		}

		next(w, r)
	})
}
