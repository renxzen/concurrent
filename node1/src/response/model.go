package response

import (
	"net/http"
	"strings"
	"time"
)

type Response struct {
	Timestamp time.Time `json:"timestamp"`
	Code      int       `json:"code"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
}

func NewResponse(code int, message string) Response {
	return Response{
		Timestamp: time.Now(),
		Code:      code,
		Status:    http.StatusText(code),
		Message:   strings.Title(message),
	}
}
