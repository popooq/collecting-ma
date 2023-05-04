// Пакет config нужен для инициализации конфига сервиса
package config

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

// Config определяет конфигурацию агента
type Config struct {
	Address        string        `env:"ADDRESS"`         // Address - адрес сервера
	ReportInterval time.Duration `env:"REPORT_INTERVAL"` // ReportInterval - интервал отправки метрик на сервер
	PollInterval   time.Duration `env:"POLL_INTERVAL"`   // PollInterval - интервал сбора метрик из рантайма
	Key            string        `env:"KEY"`             // Key - ключ шщифрования
	Rate           int           `env:"RATE_LIMIT"`      // Rate - количество единовременных потоков
	CryptoKey      string        `env:"CRYPTO_KEY"`      // CryptoKey - путь до файла с публичным ключом
}

// New создает новый конфиг
func New() *Config {
	var (
		cfg        Config
		pollTime   = time.Second * 2
		reportTime = time.Second * 10
	)

	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "address of the server")
	flag.StringVar(&cfg.Key, "k", "", "hashing key")
	flag.StringVar(&cfg.CryptoKey, "crypto-key", "", "public key file")
	flag.DurationVar(&cfg.PollInterval, "p", pollTime, "metric collection timer")
	flag.DurationVar(&cfg.ReportInterval, "r", reportTime, "metric send timer")
	flag.IntVar(&cfg.Rate, "l", 100, "worker rate")
	flag.Parse()

	if err := env.Parse(&cfg); err != nil {
		log.Printf("env parse failed :%s", err)
	}

	return &cfg
}
