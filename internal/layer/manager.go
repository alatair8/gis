// Package layer содержит функции для сопровождения слоя «Земля просто».
package layer

import (
	"time"

	"zemlyaprosto/internal/model"
	"zemlyaprosto/internal/util"
)

// Manager агрегирует функциональность по управлению слоем.
type Manager struct{}

// NewManager создаёт менеджер слоя.
func NewManager() *Manager {
	return &Manager{}
}

// BuildFeature формирует объект слоя по данным контура и дополнительным атрибутам.
func (m *Manager) BuildFeature(contour model.Contour, attributes map[string]string) model.LayerFeature {
	feature := model.LayerFeature{
		ID:         util.NewID(),
		Geometry:   contour,
		Properties: make(map[string]string),
		UpdatedAt:  time.Now(),
	}

	for key, value := range attributes {
		feature.Properties[key] = value
	}

	return feature
}
