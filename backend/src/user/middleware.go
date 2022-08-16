package user

import (
	res "backend/src/response"
	"net/http"
)

func CheckJWT(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")

		username, err := JWT.CheckToken(authorizationHeader)
		if err != nil {
			response := res.NewResponse(http.StatusMethodNotAllowed, "invalid token")
			res.SendResponse(w, response)
			return
		}

		urlUsername := r.URL.Query().Get("email")
		if username != urlUsername {
			response := res.NewResponse(http.StatusMethodNotAllowed, "invalid authorization")
			res.SendResponse(w, response)
			return
		}

		next(w, r)
	})
}
