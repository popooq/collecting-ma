package config

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	Key            string        `env:"KEY"`
	Rate           int           `env:"RATE_LIMIT"`
}

func New() *Config {
	var (
		cfg        Config
		pollTime   = time.Second * 2
		reportTime = time.Second * 10
	)

	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "address of the server")
	flag.StringVar(&cfg.Key, "k", "", "hashing key")
	flag.DurationVar(&cfg.PollInterval, "p", pollTime, "metric collection timer")
	flag.DurationVar(&cfg.ReportInterval, "r", reportTime, "metric send timer")
	flag.IntVar(&cfg.Rate, "l", 5, "worker rate")
	flag.Parse()

	if err := env.Parse(&cfg); err != nil {
		log.Printf("env parse failed :%s", err)
	}

	return &cfg
}
