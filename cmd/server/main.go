package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"zemlyaprosto/internal/app"
)

// main запускает HTTP-сервер сервиса «Земля просто».
func main() {
	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	application := app.New(addr)

	// Запускаем сервер в отдельной горутине, чтобы можно было корректно завершить работу.
	go func() {
		if err := application.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("сервер завершился с ошибкой: %v", err)
		}
	}()

	log.Printf("Сервис «Земля просто» доступен по адресу %s", addr)

	// Ожидаем сигнал завершения (Ctrl+C или SIGTERM).
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Получен сигнал завершения, останавливаем сервер...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := application.Shutdown(ctx); err != nil {
		log.Fatalf("ошибка при завершении сервера: %v", err)
	}
	log.Println("Сервис остановлен корректно")
}
