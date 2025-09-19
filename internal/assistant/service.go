package assistant

import "context"

// Recommendation описывает совет цифрового помощника.
type Recommendation struct {
	Message string `json:"message"`
}

// Service описывает интерфейс цифрового помощника.
type Service interface {
	NextStep(ctx context.Context, scenario string, state map[string]string) (Recommendation, error)
}

// ScenarioAssistant — минимальная реализация помощника на основе сценариев.
type ScenarioAssistant struct{}

// NewScenarioAssistant возвращает заглушку помощника.
func NewScenarioAssistant() *ScenarioAssistant { return &ScenarioAssistant{} }

// NextStep возвращает рекомендацию, основанную на активной жизненной ситуации.
func (s *ScenarioAssistant) NextStep(ctx context.Context, scenario string, state map[string]string) (Recommendation, error) {
	switch scenario {
	case "construction":
		return Recommendation{Message: "Проверьте наличие инженерных сетей и ограничений строительства."}, nil
	case "tourism":
		return Recommendation{Message: "Убедитесь в соблюдении природоохранных требований и доступности дорог."}, nil
	default:
		return Recommendation{Message: "Заполните обязательные атрибуты контура и сформируйте пакет документов."}, nil
	}
}

var _ Service = (*ScenarioAssistant)(nil)
