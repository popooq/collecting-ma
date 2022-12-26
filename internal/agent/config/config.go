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
}

func NewAgentConfig() *Config {
	var cfg Config
	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "address of the server")
	flag.DurationVar(&cfg.PollInterval, "p", time.Second*2, "metric collection timer")
	flag.DurationVar(&cfg.ReportInterval, "r", time.Second*10, "metric send timer")
	flag.Parse()
	err := env.Parse(&cfg)
	if err != nil {
		log.Printf("env parse failed :%s", err)
	}
	return &cfg
}
