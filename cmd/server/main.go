package main

import (
	"log"
	"net/http"

	"github.com/popooq/collectimg-ma/internal/server/config"
	"github.com/popooq/collectimg-ma/internal/server/handlers"
	"github.com/popooq/collectimg-ma/internal/server/router"
	"github.com/popooq/collectimg-ma/internal/utils/backuper"
	"github.com/popooq/collectimg-ma/internal/utils/encoder"
	"github.com/popooq/collectimg-ma/internal/utils/storage"
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
		loader.LoadFromFile()
	}
	go safe.Saver()
	log.Fatal(http.ListenAndServe(cfg.Address, handlers.GzipHandler(r)))
}
