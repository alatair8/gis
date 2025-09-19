// Package assistant реализует простейший цифровой помощник сервиса.
//
// В полноценной системе помощник мог бы использовать модели машинного обучения
// или сложные правила. Для демонстрационных целей достаточно нескольких
// эвристик, которые подсказывают пользователю следующие шаги.
package assistant

import "zemlyaprosto/internal/model"

// Request описывает вводные данные, которые пользователь предоставляет цифровому помощнику.
type Request struct {
	// Goal — желаемый результат пользователя (создать контур, выбрать участок и т. д.).
	Goal string `json:"goal"`
	// PreferredCategory — категория готовых участков, если пользователь хочет выбрать готовый вариант.
	PreferredCategory model.ParcelCategory `json:"preferred_category"`
	// HasContour показывает, существует ли уже созданный пользователем контур.
	HasContour bool `json:"has_contour"`
}

// DigitalAssistant генерирует подсказки на основании ввода пользователя и доступных данных.
type DigitalAssistant struct{}

// NewDigitalAssistant создаёт нового помощника.
func NewDigitalAssistant() *DigitalAssistant {
	return &DigitalAssistant{}
}

// Suggest формирует набор подсказок.
func (d *DigitalAssistant) Suggest(req Request, parcels []model.ReadyParcel) []model.AssistantSuggestion {
	suggestions := make([]model.AssistantSuggestion, 0)

	switch req.Goal {
	case "create_contour":
		suggestions = append(suggestions, model.AssistantSuggestion{
			Title:       "Нарисуйте границы на карте",
			Description: "Используйте инструмент рисования, чтобы очертить примерные границы участка.",
			Action:      "draw_contour",
		})
		suggestions = append(suggestions, model.AssistantSuggestion{
			Title:       "Загрузите координаты",
			Description: "Если у вас есть координаты в формате CSV или JSON, загрузите их для точного построения контура.",
			Action:      "upload_coordinates",
		})
	case "choose_parcel":
		suggestions = append(suggestions, model.AssistantSuggestion{
			Title:       "Просмотрите готовые участки",
			Description: "Сервис содержит перечни участков для стройки и туризма. Выберите подходящий вариант.",
			Action:      "list_ready_parcels",
		})
		for _, parcel := range parcels {
			if req.PreferredCategory == "" || parcel.Category == req.PreferredCategory {
				suggestions = append(suggestions, model.AssistantSuggestion{
					Title:       parcel.Name,
					Description: parcel.Description,
					Action:      "select_parcel:" + parcel.ID,
				})
			}
		}
	case "prepare_documents":
		if req.HasContour {
			suggestions = append(suggestions, model.AssistantSuggestion{
				Title:       "Сформируйте комплект документов",
				Description: "Используйте созданный контур, чтобы автоматически собрать необходимый пакет.",
				Action:      "generate_documents",
			})
		} else {
			suggestions = append(suggestions, model.AssistantSuggestion{
				Title:       "Выберите основу для документов",
				Description: "Создайте контур или выберите готовый участок, чтобы приступить к подготовке документов.",
				Action:      "choose_source",
			})
		}
	default:
		suggestions = append(suggestions, model.AssistantSuggestion{
			Title:       "Изучите руководство",
			Description: "Не уверены, с чего начать? Ознакомьтесь с описанием сервиса и примерами жизненных ситуаций.",
			Action:      "open_help",
		})
	}

	return suggestions
}
