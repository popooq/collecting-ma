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
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
	Key           string        `env:"KEY"`
	DBAddress     string        `env:"DATABASE_DSN"`
}

func New() *Config {
	var (
		cfg       Config
		storeTime = time.Second * 300
	)

	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "set server listening address")
	flag.StringVar(&cfg.Key, "k", "", "hashing key")
	flag.DurationVar(&cfg.StoreInterval, "i", storeTime, "metric backup timer")
	flag.StringVar(&cfg.StoreFile, "f", "/tmp/devops-metrics-db.json", "directory for saving metrics")
	flag.BoolVar(&cfg.Restore, "r", true, "recovering from backup before start")
	flag.StringVar(&cfg.DBAddress, "d", "", "set the DB address")
	flag.Parse()

	if err := env.Parse(&cfg); err != nil {
		log.Printf("env parse failed :%s", err)
	}

	return &cfg
}
