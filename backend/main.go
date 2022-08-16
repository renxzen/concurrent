package main

import (
	"backend/src/user"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rs/cors"
)

func main() {
	handler := http.NewServeMux()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:           ":" + port,
		Handler:        cors.Default().Handler(handler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	user.EnableUserController(handler)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Server failed: ", err.Error())
	}
}
