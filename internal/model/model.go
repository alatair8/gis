// Package model описывает ключевые доменные сущности сервиса «Земля просто».
//
// Модели вынесены в отдельный пакет, чтобы их можно было переиспользовать
// в различных уровнях приложения: хранилищах, сервисах, HTTP-обработчиках.
// Подробные комментарии к структурам помогают будущим разработчикам быстро
// погрузиться в предметную область.
package model

import "time"

// ContourSource описывает источник, из которого был создан контур участка.
//
// Значение "drawn" соответствует интерактивному рисованию пользователем на карте.
// Значение "coordinates" указывает на загрузку списка координат характерных точек.
// Значение "imported" обозначает импорт границы из внешней ГИС.
type ContourSource string

const (
	// ContourSourceDrawn — контур создан методом «рисования» на карте.
	ContourSourceDrawn ContourSource = "drawn"
	// ContourSourceCoordinates — контур создан путем ввода/загрузки координат.
	ContourSourceCoordinates ContourSource = "coordinates"
	// ContourSourceImported — контур загружен из внешней геоинформационной системы.
	ContourSourceImported ContourSource = "imported"
)

// Point описывает координату характерной точки границы земельного участка.
type Point struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Contour содержит описание границ земельного участка.
//
// В боевой системе идентификатор генерировался бы с помощью базы данных, здесь
// используется строковый идентификатор, чтобы упростить демонстрационную реализацию.
type Contour struct {
	ID          string        `json:"id"`
	Source      ContourSource `json:"source"`
	Description string        `json:"description"`
	Points      []Point       `json:"points"`
	CreatedAt   time.Time     `json:"created_at"`
}

// Attribute описывает пару ключ-значение в информационной карточке земельного участка.
//
// Использование структуры вместо простого map[string]string позволяет добавлять
// дополнительную метаинформацию (например, источник атрибута) без изменения
// сигнатур публичных методов.
type Attribute struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Source  string `json:"source"`
	Comment string `json:"comment,omitempty"`
}

// InformationCard хранит набор атрибутов для созданного контура земельного участка.
//
// Поля AutoAttributes и ManualAttributes разделены для того, чтобы можно было
// прозрачно отслеживать происхождение данных.
type InformationCard struct {
	ID               string      `json:"id"`
	ContourID        string      `json:"contour_id"`
	AutoAttributes   []Attribute `json:"auto_attributes"`
	ManualAttributes []Attribute `json:"manual_attributes"`
	CreatedAt        time.Time   `json:"created_at"`
}

// ParcelCategory определяет категорию готового земельного участка, доступного для выбора.
type ParcelCategory string

const (
	// ParcelCategoryConstruction описывает участки из перечня «Участки для стройки».
	ParcelCategoryConstruction ParcelCategory = "construction"
	// ParcelCategoryTourism описывает участки из перечня «Участки для туризма».
	ParcelCategoryTourism ParcelCategory = "tourism"
)

// ReadyParcel описывает готовый к предоставлению участок из справочников.
type ReadyParcel struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Category    ParcelCategory `json:"category"`
	Location    string         `json:"location"`
	Description string         `json:"description"`
	Contour     Contour        `json:"contour"`
	Available   bool           `json:"available"`
}

// Document описывает отдельный документ или материал, входящий в комплект.
type Document struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Source      string `json:"source"`
}

// DocumentPackage представляет комплект документов для подачи обращения.
type DocumentPackage struct {
	ID          string     `json:"id"`
	ParcelID    string     `json:"parcel_id"`
	ContourID   string     `json:"contour_id"`
	Documents   []Document `json:"documents"`
	CreatedAt   time.Time  `json:"created_at"`
	GeneratedBy string     `json:"generated_by"`
}

// BusinessStageStatus описывает состояние этапа бизнес-процесса.
type BusinessStageStatus string

const (
	// StagePending — этап ожидает выполнения.
	StagePending BusinessStageStatus = "pending"
	// StageInProgress — этап находится в работе.
	StageInProgress BusinessStageStatus = "in_progress"
	// StageCompleted — этап успешно завершен.
	StageCompleted BusinessStageStatus = "completed"
	// StageRejected — этап остановлен с отрицательным результатом.
	StageRejected BusinessStageStatus = "rejected"
)

// BusinessStage описывает отдельный этап бизнес-процесса предоставления услуги.
type BusinessStage struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Status      BusinessStageStatus `json:"status"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

// BusinessProcess агрегирует этапы государственной или муниципальной услуги.
type BusinessProcess struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Stages    []BusinessStage `json:"stages"`
	CreatedAt time.Time       `json:"created_at"`
}

// AssistantSuggestion описывает подсказку цифрового помощника.
type AssistantSuggestion struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Action      string `json:"action"`
}

// LayerFeature описывает объект слоя «Земля просто» для отображения на публичном портале.
type LayerFeature struct {
	ID         string            `json:"id"`
	Geometry   Contour           `json:"geometry"`
	Properties map[string]string `json:"properties"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

// Layer агрегирует набор объектов слоя «Земля просто».
type Layer struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Features []LayerFeature `json:"features"`
}
