package documents

import (
	"context"
	"fmt"
	"time"
)

// PackageRequest описывает данные, необходимые для формирования комплекта документов.
type PackageRequest struct {
	ContourID string            `json:"contourId"`
	Template  string            `json:"template"`
	Metadata  map[string]string `json:"metadata"`
}

// Package содержит результат формирования документов.
type Package struct {
	ID        string    `json:"id"`
	ContourID string    `json:"contourId"`
	Status    string    `json:"status"`
	Files     []string  `json:"files"`
	CreatedAt time.Time `json:"createdAt"`
}

// Service описывает операции генерации комплектов документов.
type Service interface {
	PreparePackage(ctx context.Context, req PackageRequest) (Package, error)
}

// StubService возвращает заранее подготовленные ответы, имитируя работу генератора документов.
type StubService struct{}

// NewStubService создаёт заглушку.
func NewStubService() *StubService { return &StubService{} }

// PreparePackage формирует фиктивный набор документов.
func (s *StubService) PreparePackage(ctx context.Context, req PackageRequest) (Package, error) {
	if req.ContourID == "" {
		return Package{}, fmt.Errorf("contourId is required")
	}
	if err := ctx.Err(); err != nil {
		return Package{}, err
	}
	createdAt := time.Now().UTC()
	pkg := Package{
		ID:        fmt.Sprintf("pkg-%d", createdAt.UnixNano()),
		ContourID: req.ContourID,
		Status:    "READY",
		Files: []string{
			fmt.Sprintf("%s/statement.pdf", req.ContourID),
			fmt.Sprintf("%s/plot-plan.pdf", req.ContourID),
		},
		CreatedAt: createdAt,
	}
	return pkg, nil
}

var _ Service = (*StubService)(nil)
