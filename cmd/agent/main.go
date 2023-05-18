package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/popooq/collectimg-ma/internal/agent/config"
	"github.com/popooq/collectimg-ma/internal/agent/metricsreader"
	"github.com/popooq/collectimg-ma/internal/agent/sender"
	"github.com/popooq/collectimg-ma/internal/utils/hasher"
)

func main() {
	cfg := config.New()
	hshr := hasher.Mew(cfg.Key)

	sndr := sender.New(hshr, cfg.Address, cfg.CryptoKey)
	reader := metricsreader.New(sndr, cfg.PollInterval, cfg.ReportInterval, cfg.Address, cfg.Rate)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	reader.Run(sigs)
}
