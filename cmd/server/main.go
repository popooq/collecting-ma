package main

import (
	"log"
	"net/http"

	"github.com/popooq/collectimg-ma/internal/server/config"
	"github.com/popooq/collectimg-ma/internal/server/handlers"
	"github.com/popooq/collectimg-ma/internal/server/router"
	"github.com/popooq/collectimg-ma/internal/storage"
	"github.com/popooq/collectimg-ma/internal/utils/backuper"
	"github.com/popooq/collectimg-ma/internal/utils/encoder"
)

func main() {
	storage := storage.NewMetricStorage()
	encoder := encoder.NewEncoderMetricsStruct()
	handler := handlers.NewMetricStorage(storage, encoder)
	r := router.NewRouter(handler)
	cfg := config.NewServerConfig()
	safe, err := backuper.NewSaver(storage, cfg, encoder)
	if err != nil {
		log.Printf("error during create new saver %s", err)
	}
	if cfg.Restore {
		loader, err := backuper.NewLoader(storage, cfg, encoder)
		if err != nil {
			log.Printf("error during create new loader %s", err)
		}
		err = loader.LoadFromFile()
		if err != nil {
			log.Printf("error during load from file %s", err)
		}
	}
	go safe.Saver()
	log.Fatal(http.ListenAndServe(cfg.Address, handlers.GzipHandler(r)))
}
