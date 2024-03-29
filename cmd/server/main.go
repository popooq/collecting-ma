package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/popooq/collectimg-ma/internal/server/config"
	"github.com/popooq/collectimg-ma/internal/server/handlers"
	"github.com/popooq/collectimg-ma/internal/storage"
	"github.com/popooq/collectimg-ma/internal/utils/dbsaver"
	"github.com/popooq/collectimg-ma/internal/utils/filesaver"
	"github.com/popooq/collectimg-ma/internal/utils/hasher"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	var Storage *storage.MetricsStorage
	context := context.Background()
	config := config.New()
	hasher := hasher.Mew(config.Key)
	if config.DBAddress != "" {
		dbsaver, err := dbsaver.New(context, config.DBAddress)
		if err != nil {
			log.Println(err)
		}
		Storage = storage.New(dbsaver)
		dbsaver.Migrate()
	} else if config.StoreFile != "" {
		saver, err := filesaver.New(config.StoreFile)
		if err != nil {
			log.Println(err)
		}
		Storage = storage.New(saver)
	}

	handler := handlers.New(Storage, hasher, config.Restore)
	router := chi.NewRouter()
	router.Mount("/", handler.Route())
	router.Mount("/debug", middleware.Profiler())

	if buildVersion != "" {
		fmt.Println("Build version: ", buildVersion)
	}
	fmt.Println("Build version: N/A")

	if buildDate != "" {
		fmt.Println("Build date: ", buildDate)
	}
	fmt.Println("Build date: N/A")

	if buildCommit != "" {
		fmt.Println("Build commit: ", buildCommit)
	}
	fmt.Println("Build commit: N/A")

	log.Fatal(http.ListenAndServe(config.Address, handlers.GzipHandler(router)))
}
