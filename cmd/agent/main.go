package main

import (
	"github.com/popooq/collectimg-ma/internal/agent/config"
	"github.com/popooq/collectimg-ma/internal/agent/metricsreader"
	"github.com/popooq/collectimg-ma/internal/agent/sender"
	"github.com/popooq/collectimg-ma/internal/utils/hasher"
)

func main() {
	cfg := config.New()
	hshr := hasher.Mew(cfg.Key)
	sndr := sender.New(hshr, cfg.Address)
	reader := metricsreader.New(sndr, cfg.PollInterval, cfg.ReportInterval, cfg.Address, cfg.Rate)
	reader.Run()
}
