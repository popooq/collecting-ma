package main

import (
	"log"
	"net/http"

	"github.com/popooq/collectimg-ma/internal/server/router"
)

func main() {
	r := router.NewRouter()

	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
