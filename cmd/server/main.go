package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"

	"github.com/popooq/collectimg-ma/internal/server/config"
	"github.com/popooq/collectimg-ma/internal/server/handlers"
	"github.com/popooq/collectimg-ma/internal/server/handlersproto"
	"github.com/popooq/collectimg-ma/internal/storage"
	"github.com/popooq/collectimg-ma/internal/utils/dbsaver"
	"github.com/popooq/collectimg-ma/internal/utils/filesaver"
	"github.com/popooq/collectimg-ma/internal/utils/hasher"
	pb "github.com/popooq/collectimg-ma/proto"
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

	handler := handlers.New(Storage, hasher, config.Restore, config.TrustedSubnet, config.CryptoKey)
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

	server := &http.Server{
		Addr:    config.Address,
		Handler: router,
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	idleConnsClosed := make(chan struct{})

	go func() {
		<-sig

		if err := server.Shutdown(context); err != nil {

			log.Printf("\nHTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if !config.GRPC {
		log.Println("starting server over http/1")
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}

	log.Println("starting server over http/2")
	err := protoserver(Storage)
	if err != nil {
		log.Fatal(err)
	}

	<-idleConnsClosed
}

func protoserver(storage *storage.MetricsStorage) error {
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()

	metricServer := handlersproto.NewMetricServer(storage)
	pb.RegisterMetricsServer(s, &metricServer)

	fmt.Println("Сервер gRPC начал работу")

	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
		return err
	}
	return err
}
