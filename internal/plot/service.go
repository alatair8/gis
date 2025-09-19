package plot

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// Contour описывает сохранённый контур земельного участка.
type Contour struct {
	ID         string            `json:"id"`
	OwnerID    string            `json:"ownerId"`
	Geometry   string            `json:"geometry"`
	Attributes map[string]string `json:"attributes"`
	CreatedAt  time.Time         `json:"createdAt"`
}

// ContourDraft содержит данные для создания контура.
type ContourDraft struct {
	OwnerID    string            `json:"ownerId"`
	Geometry   string            `json:"geometry"`
	Attributes map[string]string `json:"attributes"`
}

// Service определяет операции сервиса моделирования участков.
type Service interface {
	CreateContour(ctx context.Context, draft ContourDraft) (Contour, error)
	ListContours(ctx context.Context, ownerID string) ([]Contour, error)
}

// InMemoryService — временная реализация для прототипирования.
type InMemoryService struct {
	mu        sync.RWMutex
	sequence  int
	contours  map[string]Contour
	byOwnerID map[string][]string
}

// NewInMemoryService создаёт in-memory реализацию.
func NewInMemoryService() *InMemoryService {
	return &InMemoryService{
		contours:  make(map[string]Contour),
		byOwnerID: make(map[string][]string),
	}
}

// CreateContour сохраняет контур и возвращает его идентификатор.
func (s *InMemoryService) CreateContour(ctx context.Context, draft ContourDraft) (Contour, error) {
	if draft.Geometry == "" {
		return Contour{}, errors.New("geometry is required")
	}
	if err := ctx.Err(); err != nil {
		return Contour{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.sequence++
	contour := Contour{
		ID:         fmt.Sprintf("ctr-%d", s.sequence),
		OwnerID:    draft.OwnerID,
		Geometry:   draft.Geometry,
		Attributes: draft.Attributes,
		CreatedAt:  time.Now().UTC(),
	}
	s.contours[contour.ID] = contour
	if contour.OwnerID != "" {
		s.byOwnerID[contour.OwnerID] = append(s.byOwnerID[contour.OwnerID], contour.ID)
	}
	return contour, nil
}

// ListContours возвращает сохранённые контуры.
func (s *InMemoryService) ListContours(ctx context.Context, ownerID string) ([]Contour, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	if ownerID == "" {
		result := make([]Contour, 0, len(s.contours))
		for _, c := range s.contours {
			result = append(result, c)
		}
		return result, nil
	}

	ids := s.byOwnerID[ownerID]
	result := make([]Contour, 0, len(ids))
	for _, id := range ids {
		result = append(result, s.contours[id])
	}
	return result, nil
}

var _ Service = (*InMemoryService)(nil)
