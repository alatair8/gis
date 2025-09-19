// Package store содержит примитивную in-memory реализацию хранилища данных.
//
// Такое хранилище удобно использовать на ранних этапах проектирования, когда
// требуется продемонстрировать бизнес-логику без настройки полноценной БД.
package store

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"zemlyaprosto/internal/model"
	"zemlyaprosto/internal/util"
)

// MemoryStore — простое потокобезопасное хранилище данных в памяти процесса.
type MemoryStore struct {
	mu           sync.RWMutex
	contours     map[string]model.Contour
	cards        map[string]model.InformationCard
	readyParcels map[string]model.ReadyParcel
	docPackages  map[string]model.DocumentPackage
	processes    map[string]model.BusinessProcess
	layer        model.Layer
}

// NewMemoryStore инициализирует хранилище с небольшим набором демонстрационных данных.
func NewMemoryStore() *MemoryStore {
	layer := model.Layer{
		ID:       "layer-1",
		Name:     "Земля просто",
		Features: []model.LayerFeature{},
	}

	store := &MemoryStore{
		contours:     make(map[string]model.Contour),
		cards:        make(map[string]model.InformationCard),
		readyParcels: make(map[string]model.ReadyParcel),
		docPackages:  make(map[string]model.DocumentPackage),
		processes:    make(map[string]model.BusinessProcess),
		layer:        layer,
	}

	store.seedReadyParcels()
	return store
}

// seedReadyParcels заполняет хранилище примерами готовых участков для стройки и туризма.
func (m *MemoryStore) seedReadyParcels() {
	m.readyParcels["construction-1"] = model.ReadyParcel{
		ID:       "construction-1",
		Name:     "Промышленный парк «Северный»",
		Category: model.ParcelCategoryConstruction,
		Location: "Архангельская область, г. Архангельск",
		Description: "Участок в промышленной зоне с готовыми инженерными коммуникациями, " +
			"подходит для размещения производства.",
		Contour: model.Contour{
			ID:          "contour-construction-1",
			Source:      model.ContourSourceImported,
			CreatedAt:   time.Now(),
			Description: "Импортированная граница промышленного участка",
			Points:      []model.Point{{Latitude: 64.54, Longitude: 40.55}, {Latitude: 64.55, Longitude: 40.56}},
		},
		Available: true,
	}

	m.readyParcels["tourism-1"] = model.ReadyParcel{
		ID:          "tourism-1",
		Name:        "Туристический кластер «Бирюзовая Катунь»",
		Category:    model.ParcelCategoryTourism,
		Location:    "Республика Алтай",
		Description: "Живописный участок в горной местности, предназначенный для развития туризма.",
		Contour: model.Contour{
			ID:          "contour-tourism-1",
			Source:      model.ContourSourceImported,
			CreatedAt:   time.Now(),
			Description: "Контур сформирован органом власти и опубликован в справочнике",
			Points:      []model.Point{{Latitude: 51.99, Longitude: 85.85}, {Latitude: 52.0, Longitude: 85.86}},
		},
		Available: true,
	}
}

// SaveContour сохраняет контур участка и возвращает его копию с присвоенным идентификатором.
func (m *MemoryStore) SaveContour(contour model.Contour) model.Contour {
	m.mu.Lock()
	defer m.mu.Unlock()

	if contour.ID == "" {
		contour.ID = util.NewID()
	}
	contour.CreatedAt = time.Now()

	m.contours[contour.ID] = contour
	return contour
}

// GetContourByID возвращает контур по идентификатору.
func (m *MemoryStore) GetContourByID(id string) (model.Contour, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	contour, ok := m.contours[id]
	if !ok {
		return model.Contour{}, ErrNotFound
	}
	return contour, nil
}

// ListContours возвращает срез всех сохраненных контуров.
func (m *MemoryStore) ListContours() []model.Contour {
	m.mu.RLock()
	defer m.mu.RUnlock()

	contours := make([]model.Contour, 0, len(m.contours))
	for _, contour := range m.contours {
		contours = append(contours, contour)
	}
	return contours
}

// SaveInformationCard сохраняет информационную карточку.
func (m *MemoryStore) SaveInformationCard(card model.InformationCard) model.InformationCard {
	m.mu.Lock()
	defer m.mu.Unlock()

	if card.ID == "" {
		card.ID = util.NewID()
	}
	card.CreatedAt = time.Now()

	m.cards[card.ID] = card
	return card
}

// GetInformationCardByContour возвращает карточку для конкретного контура, если она есть.
func (m *MemoryStore) GetInformationCardByContour(contourID string) (model.InformationCard, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, card := range m.cards {
		if card.ContourID == contourID {
			return card, true
		}
	}
	return model.InformationCard{}, false
}

// ListReadyParcels возвращает готовые участки по заданной категории.
func (m *MemoryStore) ListReadyParcels(category model.ParcelCategory) []model.ReadyParcel {
	m.mu.RLock()
	defer m.mu.RUnlock()

	parcels := make([]model.ReadyParcel, 0)
	for _, parcel := range m.readyParcels {
		if category == "" || parcel.Category == category {
			parcels = append(parcels, parcel)
		}
	}
	return parcels
}

// SaveDocumentPackage сохраняет комплект документов.
func (m *MemoryStore) SaveDocumentPackage(pkg model.DocumentPackage) model.DocumentPackage {
	m.mu.Lock()
	defer m.mu.Unlock()

	if pkg.ID == "" {
		pkg.ID = util.NewID()
	}
	pkg.CreatedAt = time.Now()

	m.docPackages[pkg.ID] = pkg
	return pkg
}

// ListDocumentPackages возвращает все сформированные комплекты документов.
func (m *MemoryStore) ListDocumentPackages() []model.DocumentPackage {
	m.mu.RLock()
	defer m.mu.RUnlock()

	packages := make([]model.DocumentPackage, 0, len(m.docPackages))
	for _, pkg := range m.docPackages {
		packages = append(packages, pkg)
	}
	return packages
}

// SaveBusinessProcess сохраняет бизнес-процесс.
func (m *MemoryStore) SaveBusinessProcess(process model.BusinessProcess) model.BusinessProcess {
	m.mu.Lock()
	defer m.mu.Unlock()

	if process.ID == "" {
		process.ID = util.NewID()
	}
	if process.CreatedAt.IsZero() {
		process.CreatedAt = time.Now()
	}

	m.processes[process.ID] = process
	return process
}

// GetBusinessProcessByID возвращает бизнес-процесс по идентификатору.
func (m *MemoryStore) GetBusinessProcessByID(id string) (model.BusinessProcess, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	process, ok := m.processes[id]
	if !ok {
		return model.BusinessProcess{}, ErrNotFound
	}
	return process, nil
}

// UpdateBusinessStage изменяет статус конкретного этапа бизнес-процесса.
func (m *MemoryStore) UpdateBusinessStage(processID, stageID string, status model.BusinessStageStatus) (model.BusinessProcess, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	process, ok := m.processes[processID]
	if !ok {
		return model.BusinessProcess{}, fmt.Errorf("process %s not found", processID)
	}

	updated := false
	for i, stage := range process.Stages {
		if stage.ID == stageID {
			stage.Status = status
			stage.UpdatedAt = time.Now()
			process.Stages[i] = stage
			updated = true
			break
		}
	}

	if !updated {
		return model.BusinessProcess{}, fmt.Errorf("stage %s not found", stageID)
	}

	m.processes[process.ID] = process
	return process, nil
}

// AddLayerFeature добавляет новый объект в слой «Земля просто».
func (m *MemoryStore) AddLayerFeature(feature model.LayerFeature) model.LayerFeature {
	m.mu.Lock()
	defer m.mu.Unlock()

	if feature.ID == "" {
		feature.ID = util.NewID()
	}
	feature.UpdatedAt = time.Now()

	m.layer.Features = append(m.layer.Features, feature)
	return feature
}

// GetLayer возвращает копию слоя «Земля просто».
func (m *MemoryStore) GetLayer() model.Layer {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.layer
}

// ErrNotFound используется в сервисах для единообразной обработки отсутствия данных.
var ErrNotFound = errors.New("not found")

// GetReadyParcelByID возвращает готовый участок по идентификатору.
func (m *MemoryStore) GetReadyParcelByID(id string) (model.ReadyParcel, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	parcel, ok := m.readyParcels[id]
	if !ok {
		return model.ReadyParcel{}, ErrNotFound
	}
	return parcel, nil
}
