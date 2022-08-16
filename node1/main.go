package main

import (
	"log"
	"net/http"
	"time"

	"nodo1/src/info"
)

func main() {
	handler := http.NewServeMux()
	server := &http.Server{
		Addr:           ":8081",
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	info.EnableInfoController(handler)

	log.Println("Http server started on port 8081...")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Server failed: ", err.Error())
	}
}
