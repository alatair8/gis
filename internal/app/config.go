package app

import (
	"os"
	"time"
)

// Config описывает параметры запуска HTTP-приложения и интеграций.
type Config struct {
	HTTPListenAddr string
	ShutdownGrace  time.Duration
}

// LoadConfig формирует конфигурацию приложения из переменных окружения.
func LoadConfig() Config {
	listenAddr := os.Getenv("HTTP_LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":8080"
	}

	return Config{
		HTTPListenAddr: listenAddr,
		ShutdownGrace:  10 * time.Second,
	}
}
