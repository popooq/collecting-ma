package main

import (
	"context"
	"log"
	"net/http"

	"github.com/popooq/collectimg-ma/internal/server/config"
	"github.com/popooq/collectimg-ma/internal/server/handlers"
	"github.com/popooq/collectimg-ma/internal/server/router"
	"github.com/popooq/collectimg-ma/internal/storage"
	"github.com/popooq/collectimg-ma/internal/utils/backuper"
	"github.com/popooq/collectimg-ma/internal/utils/encoder"
	"github.com/popooq/collectimg-ma/internal/utils/hasher"
	"github.com/popooq/collectimg-ma/internal/utils/pgdb"
)

func main() {
	context := context.Background()
	storage := storage.New()
	encoder := encoder.New()
	config := config.New()
	hasher := hasher.Mew(config.Key)
	database := pgdb.New(context, config, storage)
	handler := handlers.New(storage, hasher, database)
	router := router.New(handler)
	saver, err := backuper.NewSaver(storage, config, encoder, database)
	if err != nil {
		log.Printf("error during create new saver %s", err)
	}

	if config.Restore {
		loader, err := backuper.NewLoader(storage, config, encoder)
		if err != nil {
			log.Printf("error during create new loader %s", err)
		}

		err = loader.LoadFromFile()
		if err != nil {
			log.Printf("error during load from file %s", err)
		}
	}

	go saver.GoFile()

	log.Fatal(http.ListenAndServe(config.Address, handlers.GzipHandler(router)))
}
