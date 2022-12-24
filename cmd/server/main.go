package main

import (
	"log"
	"net/http"

	"github.com/popooq/collectimg-ma/internal/server/handlers"
	"github.com/popooq/collectimg-ma/internal/server/router"
	"github.com/popooq/collectimg-ma/internal/utils/backuper"
	"github.com/popooq/collectimg-ma/internal/utils/encoder"
	"github.com/popooq/collectimg-ma/internal/utils/env"
	"github.com/popooq/collectimg-ma/internal/utils/storage"
)

func main() {
	storage := storage.NewMetricStorage()
	encoder := encoder.NewEncoderMetricsStruct()
	handler := handlers.NewMetricStorage(storage, encoder)
	r := router.NewRouter(handler)
	env := env.ServerConfig()
	safe, err := backuper.NewSaver(storage, env)
	if err != nil {
		log.Printf("error during create new saver %s", err)
	}
	loader, err := backuper.NewLoader(storage, env, encoder)
	if err != nil {
		log.Printf("error during create new loader %s", err)
	}
	if env.Restore {
		loader.LoadFromFile()
	}
	go safe.Saver()
	log.Fatal(http.ListenAndServe(env.Address, r))
}
