// Package service реализует бизнес-логику сервиса «Земля просто».
package service

import (
	"errors"
	"fmt"
	"strings"

	"zemlyaprosto/internal/assistant"
	"zemlyaprosto/internal/business"
	"zemlyaprosto/internal/layer"
	"zemlyaprosto/internal/model"
	"zemlyaprosto/internal/store"
	"zemlyaprosto/internal/util"
)

// Service объединяет работу хранилища, цифрового помощника и других компонентов.
type Service struct {
	store        *store.MemoryStore
	assistant    *assistant.DigitalAssistant
	layerManager *layer.Manager
}

// New создаёт новый экземпляр бизнес-сервиса.
func New(store *store.MemoryStore, assistant *assistant.DigitalAssistant, layerManager *layer.Manager) *Service {
	return &Service{
		store:        store,
		assistant:    assistant,
		layerManager: layerManager,
	}
}

// CreateContourFromDrawing регистрирует контур, нарисованный пользователем на карте.
func (s *Service) CreateContourFromDrawing(description string, points []model.Point) (model.Contour, error) {
	if len(points) < 3 {
		return model.Contour{}, errors.New("для построения контура необходимо минимум 3 точки")
	}
	contour := model.Contour{
		Description: description,
		Points:      points,
		Source:      model.ContourSourceDrawn,
	}
	return s.store.SaveContour(contour), nil
}

// CreateContourFromCoordinates создаёт контур на основе списка координат.
func (s *Service) CreateContourFromCoordinates(description string, points []model.Point) (model.Contour, error) {
	if len(points) == 0 {
		return model.Contour{}, errors.New("список координат не может быть пустым")
	}
	contour := model.Contour{
		Description: description,
		Points:      points,
		Source:      model.ContourSourceCoordinates,
	}
	return s.store.SaveContour(contour), nil
}

// ImportContour загружает контур из внешней системы.
func (s *Service) ImportContour(description string, points []model.Point) (model.Contour, error) {
	if len(points) == 0 {
		return model.Contour{}, errors.New("импортированный контур не содержит точек")
	}
	contour := model.Contour{
		Description: description,
		Points:      points,
		Source:      model.ContourSourceImported,
	}
	return s.store.SaveContour(contour), nil
}

// ListContours возвращает все созданные контуры.
func (s *Service) ListContours() []model.Contour {
	return s.store.ListContours()
}

// CreateInformationCard формирует информационную карточку для контура.
func (s *Service) CreateInformationCard(contourID string, autoAttrs, manualAttrs []model.Attribute) (model.InformationCard, error) {
	if contourID == "" {
		return model.InformationCard{}, errors.New("не указан идентификатор контура")
	}
	if _, err := s.store.GetContourByID(contourID); err != nil {
		return model.InformationCard{}, fmt.Errorf("контур не найден: %w", err)
	}

	card := model.InformationCard{
		ContourID:        contourID,
		AutoAttributes:   autoAttrs,
		ManualAttributes: manualAttrs,
	}

	return s.store.SaveInformationCard(card), nil
}

// ListReadyParcels возвращает перечень готовых участков по категории.
func (s *Service) ListReadyParcels(category model.ParcelCategory) []model.ReadyParcel {
	return s.store.ListReadyParcels(category)
}

// GenerateDocumentPackage собирает комплект документов для обращения.
func (s *Service) GenerateDocumentPackage(contourID, parcelID string) (model.DocumentPackage, error) {
	if contourID == "" && parcelID == "" {
		return model.DocumentPackage{}, errors.New("необходимо указать контур или готовый участок")
	}

	documents := make([]model.Document, 0)
	generator := make([]string, 0)

	if contourID != "" {
		contour, err := s.store.GetContourByID(contourID)
		if err != nil {
			return model.DocumentPackage{}, fmt.Errorf("контур не найден: %w", err)
		}
		documents = append(documents, model.Document{
			ID:          util.NewID(),
			Name:        "Схема расположения земельного участка",
			Description: "Схема автоматически сформирована на основании созданного контура",
			Source:      "generated_from_contour",
		})
		documents = append(documents, model.Document{
			ID:          util.NewID(),
			Name:        "Координаты характерных точек",
			Description: "Ведомость координат для подачи в органы кадастрового учёта",
			Source:      "generated_from_contour",
		})
		generator = append(generator, "contour:"+contour.ID)
	}

	if parcelID != "" {
		parcel, err := s.store.GetReadyParcelByID(parcelID)
		if err != nil {
			return model.DocumentPackage{}, fmt.Errorf("готовый участок не найден: %w", err)
		}
		documents = append(documents, model.Document{
			ID:          util.NewID(),
			Name:        "Выписка из перечня готовых участков",
			Description: "Документ подтверждает параметры участка из перечня",
			Source:      "ready_parcel_registry",
		})
		generator = append(generator, "ready_parcel:"+parcel.ID)
	}

	// Общие документы, необходимые для подачи обращений.
	documents = append(documents, model.Document{
		ID:          util.NewID(),
		Name:        "Заявление",
		Description: "Черновик заявления на предоставление земельного участка",
		Source:      "template",
	})
	documents = append(documents, model.Document{
		ID:          util.NewID(),
		Name:        "Согласие на обработку персональных данных",
		Description: "Обязательный документ для подачи обращения",
		Source:      "template",
	})

	pkg := model.DocumentPackage{
		ParcelID:    parcelID,
		ContourID:   contourID,
		Documents:   documents,
		GeneratedBy: strings.Join(generator, ";"),
	}

	return s.store.SaveDocumentPackage(pkg), nil
}

// GetDocumentPackages возвращает сформированные ранее комплекты документов.
func (s *Service) GetDocumentPackages() []model.DocumentPackage {
	return s.store.ListDocumentPackages()
}

// GetAssistantSuggestions возвращает подсказки цифрового помощника.
func (s *Service) GetAssistantSuggestions(req assistant.Request) []model.AssistantSuggestion {
	parcels := s.store.ListReadyParcels(req.PreferredCategory)
	return s.assistant.Suggest(req, parcels)
}

// CreateBusinessProcess создаёт новый бизнес-процесс и сохраняет его в хранилище.
func (s *Service) CreateBusinessProcess(name string) model.BusinessProcess {
	process := business.NewDefaultProcess(name)
	return s.store.SaveBusinessProcess(process)
}

// AdvanceBusinessProcess переводит следующий этап процесса в работу.
func (s *Service) AdvanceBusinessProcess(processID string) (model.BusinessProcess, error) {
	process, err := s.store.GetBusinessProcessByID(processID)
	if err != nil {
		return model.BusinessProcess{}, fmt.Errorf("процесс не найден: %w", err)
	}

	updated, ok := business.AdvanceToNextStage(process)
	if !ok {
		return model.BusinessProcess{}, errors.New("в процессе отсутствуют этапы для запуска")
	}
	updated = s.store.SaveBusinessProcess(updated)
	return updated, nil
}

// CompleteBusinessStage завершает конкретный этап процесса.
func (s *Service) CompleteBusinessStage(processID, stageID string, success bool) (model.BusinessProcess, error) {
	process, err := s.store.GetBusinessProcessByID(processID)
	if err != nil {
		return model.BusinessProcess{}, fmt.Errorf("процесс не найден: %w", err)
	}

	updated := business.CompleteStage(process, stageID, success)
	updated = s.store.SaveBusinessProcess(updated)
	return updated, nil
}

// PublishContourToLayer добавляет контур в слой «Земля просто».
func (s *Service) PublishContourToLayer(contourID string, attributes map[string]string) (model.LayerFeature, error) {
	contour, err := s.store.GetContourByID(contourID)
	if err != nil {
		return model.LayerFeature{}, fmt.Errorf("контур не найден: %w", err)
	}
	feature := s.layerManager.BuildFeature(contour, attributes)
	return s.store.AddLayerFeature(feature), nil
}

// GetLayer возвращает текущее состояние слоя «Земля просто».
func (s *Service) GetLayer() model.Layer {
	return s.store.GetLayer()
}
