package main

import (
	"log"
	"net/http"

	"github.com/popooq/collectimg-ma/internal/server/router"
	"github.com/popooq/collectimg-ma/internal/utils/env"
)

func main() {
	r := router.NewRouter()
	cfg := env.ServerConfig()
	log.Fatal(http.ListenAndServe(cfg.Address, r))
}
