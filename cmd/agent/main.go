package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/popooq/collectimg-ma/internal/agent/config"
	"github.com/popooq/collectimg-ma/internal/agent/metricsreader"
	"github.com/popooq/collectimg-ma/internal/agent/protoreader"
	"github.com/popooq/collectimg-ma/internal/agent/sender"
	"github.com/popooq/collectimg-ma/internal/utils/hasher"
	pb "github.com/popooq/collectimg-ma/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.New()
	hshr := hasher.Mew(cfg.Key)

	sndr := sender.New(hshr, cfg.Address, cfg.CryptoKey)
	reader := metricsreader.New(sndr, cfg.PollInterval, cfg.ReportInterval, cfg.Address, cfg.Rate)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	reader.Run(sigs)

	conn, err := grpc.Dial(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := pb.NewMetricsClient(conn)

	protoreader.MetricRequest(client)
}
