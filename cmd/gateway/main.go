package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"zemlya-prosto/internal/app"
)

func main() {
	cfg := app.LoadConfig()

	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}

	mux := http.NewServeMux()
	application.RegisterRoutes(mux)

	server := &http.Server{
		Addr:              cfg.HTTPListenAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("HTTP server is listening on %s", cfg.HTTPListenAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}

	application.Close()
}
