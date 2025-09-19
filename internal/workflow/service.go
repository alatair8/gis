package workflow

import (
	"context"
	"log"
	"time"
)

// Service описывает взаимодействие с оркестратором бизнес-процессов.
type Service interface {
	NotifyPackageReady(ctx context.Context, packageID string) error
}

// StubService логирует вызовы без реальной интеграции.
type StubService struct{}

// NewStubService создаёт заглушку сервиса.
func NewStubService() *StubService { return &StubService{} }

// NotifyPackageReady протоколирует событие формирования пакета документов.
func (s *StubService) NotifyPackageReady(ctx context.Context, packageID string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		log.Printf("workflow: package %s ready at %s", packageID, time.Now().UTC().Format(time.RFC3339))
		return nil
	}
}

var _ Service = (*StubService)(nil)
