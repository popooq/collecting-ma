package config

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	Storefile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
	Key           string        `env:"KEY"`
}

func NewServerConfig() *Config {
	var cfg Config
	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "set server listening address")
	flag.StringVar(&cfg.Key, "k", "", "hashing key")
	flag.DurationVar(&cfg.StoreInterval, "i", time.Second*1, "metric backup timer")
	flag.StringVar(&cfg.Storefile, "f", "/tmp/devops-metrics-db.json", "directory for saving metrics")
	flag.BoolVar(&cfg.Restore, "r", true, "recovering from backup before start")
	flag.Parse()
	err := env.Parse(&cfg)
	if err != nil {
		log.Printf("env parse failed :%s", err)
	}
	return &cfg
}
