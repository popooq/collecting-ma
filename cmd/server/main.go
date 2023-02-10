package main

import (
	"context"
	"log"
	"net/http"

	"github.com/popooq/collectimg-ma/internal/server/config"
	"github.com/popooq/collectimg-ma/internal/server/handlers"
	"github.com/popooq/collectimg-ma/internal/server/router"
	"github.com/popooq/collectimg-ma/internal/storage"
	"github.com/popooq/collectimg-ma/internal/utils/dbsaver"
	"github.com/popooq/collectimg-ma/internal/utils/filesaver"
	"github.com/popooq/collectimg-ma/internal/utils/hasher"
)

func main() {
	context := context.Background()
	config := config.New()
	hasher := hasher.Mew(config.Key)
	dbsaver, err := dbsaver.NewSaver(context, config)
	if err != nil {
		log.Printf("error during create new dbsaver %s", err)
	}
	saver, err := filesaver.New(config)
	if err != nil {
		log.Printf("error during create new saver %s", err)
	}
	var Storage *storage.MetricsStorage
	if config.DBAddress != "" {
		Storage = storage.New(dbsaver, *config)
	} else {
		Storage = storage.New(saver, *config)
	}

	handler := handlers.New(Storage, hasher)
	router := router.New(handler)
	if config.Restore {
		err = Storage.Load()
		if err != nil {
			log.Printf("error during load from file %s", err)
		}
	}

	if dbsaver != nil {
		dbsaver.CreateTable()
	}

	log.Fatal(http.ListenAndServe(config.Address, handlers.GzipHandler(router)))
}
