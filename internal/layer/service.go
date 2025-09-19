package layer

import "context"

// Service управляет подготовкой слоя «Земля просто».
type Service interface {
	Publish(ctx context.Context) error
	Close() error
}

// StubService имитирует обновление слоя.
type StubService struct{}

// NewStubService создаёт заглушку.
func NewStubService() *StubService { return &StubService{} }

// Publish ничего не делает в заглушке.
func (s *StubService) Publish(ctx context.Context) error {
	return ctx.Err()
}

// Close реализует интерфейс Service.
func (s *StubService) Close() error { return nil }

var _ Service = (*StubService)(nil)
