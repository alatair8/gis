package catalog

import (
	"context"
	"sync"
)

// Plot описывает участок из перечня готовых участков.
type Plot struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Purpose       string   `json:"purpose"`
	LifeSituation string   `json:"lifeSituation"`
	Geometry      string   `json:"geometry"`
	Tags          []string `json:"tags"`
}

// Service описывает операции работы с каталогом участков.
type Service interface {
	ListPlots(ctx context.Context, situation string) ([]Plot, error)
}

// InMemoryService хранит небольшой набор участков для демонстрации.
type InMemoryService struct {
	plots []Plot
	mu    sync.RWMutex
}

// NewInMemoryService создаёт предварительно заполненный каталог.
func NewInMemoryService() *InMemoryService {
	return &InMemoryService{
		plots: []Plot{
			{ID: "tour-1", Name: "Турбаза на озере", Purpose: "туризм", LifeSituation: "tourism", Geometry: "{...}", Tags: []string{"дорога", "электричество"}},
			{ID: "build-1", Name: "Микрорайон Северный", Purpose: "строительство", LifeSituation: "construction", Geometry: "{...}", Tags: []string{"газ", "школа"}},
		},
	}
}

// ListPlots возвращает участки по жизненной ситуации (или все, если не указана).
func (s *InMemoryService) ListPlots(ctx context.Context, situation string) ([]Plot, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()

	if situation == "" {
		result := make([]Plot, len(s.plots))
		copy(result, s.plots)
		return result, nil
	}

	var result []Plot
	for _, plot := range s.plots {
		if plot.LifeSituation == situation {
			result = append(result, plot)
		}
	}
	return result, nil
}

var _ Service = (*InMemoryService)(nil)
