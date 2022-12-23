package main

import (
	"log"
	"net/http"

	"github.com/popooq/collectimg-ma/internal/server/router"
)

func main() {
	r := router.NewRouter()

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", r))
}
