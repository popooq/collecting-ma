package env

import (
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

type ConfigAgent struct {
	Address        string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"2s"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"10s"`
}

type ConfigServer struct {
	Address       string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"100s"`
	Storefile     string        `env:"STORE_FILE,notEmpty" envDefault:"/tmp/devops-metrics-db.json" `
	Restore       bool          `env:"RESTORE" envDefault:"true"`
}

func AgentConfig() *ConfigAgent {
	var cfg ConfigAgent
	err := env.Parse(&cfg)
	if err != nil {
		log.Printf("env parse failed :%s", err)
	}
	return &cfg
}

func ServerConfig() *ConfigServer {
	var cfg ConfigServer
	err := env.Parse(&cfg)
	if err != nil {
		log.Printf("env parse failed :%s", err)
	}
	return &cfg
}
