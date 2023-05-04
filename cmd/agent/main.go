package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/popooq/collectimg-ma/internal/agent/config"
	"github.com/popooq/collectimg-ma/internal/agent/metricsreader"
	"github.com/popooq/collectimg-ma/internal/agent/sender"
	"github.com/popooq/collectimg-ma/internal/utils/encryptor"
	"github.com/popooq/collectimg-ma/internal/utils/hasher"
)

func main() {
	cfg := config.New()
	hshr := hasher.Mew(cfg.Key)
	enc, err := encryptor.New(cfg.CryptoKey, "public")
	if err != nil {
		log.Fatalf("%s", err)
	}
	sndr := sender.New(hshr, cfg.Address, enc)
	reader := metricsreader.New(sndr, cfg.PollInterval, cfg.ReportInterval, cfg.Address, cfg.Rate)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	reader.Run(sigs)
}
