package main

import (
	"log"
	"net/http"

	handlers "github.com/popooq/collectimg-ma/internal/server/handlers"
	storage "github.com/popooq/collectimg-ma/internal/server/storage"
)

func main() {
	memS := storage.NewMemStorage()
	handler := handlers.NewmetricStorage(memS)
	mux := http.NewServeMux()
	mux.Handle("/update/", handler)

	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}
