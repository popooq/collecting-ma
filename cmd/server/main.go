package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/popooq/collectimg-ma/internal/server/config"
	"github.com/popooq/collectimg-ma/internal/server/handlers"
	"github.com/popooq/collectimg-ma/internal/storage"
	"github.com/popooq/collectimg-ma/internal/utils/dbsaver"
	"github.com/popooq/collectimg-ma/internal/utils/filesaver"
	"github.com/popooq/collectimg-ma/internal/utils/hasher"
)

func main() {
	var Storage *storage.MetricsStorage
	context := context.Background()
	config := config.New()
	hasher := hasher.Mew(config.Key)
	if config.DBAddress != "" {
		dbsaver, err := dbsaver.New(context, config)
		if err != nil {
			log.Println(err)
		}
		Storage = storage.New(dbsaver, *config)
		if dbsaver != nil {
			dbsaver.CreateTable()
		}
	} else if config.StoreFile != "" {
		saver, err := filesaver.New(config)
		if err != nil {
			log.Println(err)
		}
		Storage = storage.New(saver, *config)
	}

	handler := handlers.New(Storage, hasher)
	router := chi.NewRouter()
	router.Mount("/", handler.Route())

	log.Fatal(http.ListenAndServe(config.Address, handlers.GzipHandler(router)))
}
