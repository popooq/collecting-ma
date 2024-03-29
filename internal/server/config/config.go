package config

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

// Config определяет конфигурацию агента
type Config struct {
	Address       string        `env:"ADDRESS"`        // Address - адрес сервера
	StoreInterval time.Duration `env:"STORE_INTERVAL"` // StoreInterval - период сохранения метрик
	StoreFile     string        `env:"STORE_FILE"`     // Storefile - адрес файла для хранения метрик
	Restore       bool          `env:"RESTORE"`        // Restore - восстановление из бекапа bool
	Key           string        `env:"KEY"`            //Key - ключ шифрования
	DBAddress     string        `env:"DATABASE_DSN"`   // DBAddress - адрес базы данных
}

// New создает новый конфиг
func New() *Config {
	var (
		cfg       Config
		storeTime = time.Second * 5
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
