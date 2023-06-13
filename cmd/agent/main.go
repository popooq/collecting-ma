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
	config := config.New()
	hshr := hasher.New(config.Key)

	sndr := sender.New(hshr, config.Address, config.CryptoKey)
	reader := metricsreader.New(sndr, config.PollInterval, config.ReportInterval, config.Address, config.Rate)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	if !config.GRPC {
		reader.Run(sigs)
	}
	conn, err := grpc.Dial(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := pb.NewMetricsClient(conn)

	protoreader.MetricRequest(client, sigs, config.PollInterval, config.ReportInterval, config.Rate, *hshr)
}
