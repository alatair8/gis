// Package app инициализирует зависимости и запускает HTTP-сервер сервиса.
package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"zemlyaprosto/internal/assistant"
	httpapi "zemlyaprosto/internal/http"
	"zemlyaprosto/internal/layer"
	"zemlyaprosto/internal/service"
	"zemlyaprosto/internal/store"
)

// Application агрегирует все компоненты и управляет жизненным циклом сервиса.
type Application struct {
	server *http.Server
}

// New создаёт приложение с инициализированными зависимостями.
func New(addr string) *Application {
	store := store.NewMemoryStore()
	assistant := assistant.NewDigitalAssistant()
	layerManager := layer.NewManager()
	svc := service.New(store, assistant, layerManager)

	mux := http.NewServeMux()
	handler := httpapi.New(svc)
	handler.Register(mux)

	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Application{server: srv}
}

// Run запускает HTTP-сервер и блокирует выполнение до получения ошибки.
func (a *Application) Run() error {
	log.Printf("HTTP сервер запущен на %s", a.server.Addr)
	return a.server.ListenAndServe()
}

// Shutdown корректно останавливает сервер.
func (a *Application) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
