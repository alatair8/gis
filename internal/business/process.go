// Package business содержит вспомогательные функции для моделирования
// государственных и муниципальных бизнес-процессов.
package business

import (
	"time"

	"zemlyaprosto/internal/model"
	"zemlyaprosto/internal/util"
)

// NewDefaultProcess создаёт типовой бизнес-процесс предоставления земельного участка.
//
// Процесс состоит из этапов, соответствующих ключевым шагам, которые выполняют
// уполномоченные органы. В реальной системе набор этапов мог бы настраиваться
// в зависимости от региона, статуса земельного участка и иных параметров.
func NewDefaultProcess(name string) model.BusinessProcess {
	now := time.Now()
	stages := []model.BusinessStage{
		{
			ID:          util.NewID(),
			Name:        "Приём и регистрация обращения",
			Description: "Уполномоченный орган регистрирует поступившее заявление",
			Status:      model.StagePending,
			UpdatedAt:   now,
		},
		{
			ID:          util.NewID(),
			Name:        "Рассмотрение и согласование",
			Description: "Специалисты анализируют документы и принимают решение",
			Status:      model.StagePending,
			UpdatedAt:   now,
		},
		{
			ID:          util.NewID(),
			Name:        "Подготовка итоговых документов",
			Description: "Формируются документы для постановки на кадастровый учёт и регистрации прав",
			Status:      model.StagePending,
			UpdatedAt:   now,
		},
	}

	return model.BusinessProcess{
		ID:        util.NewID(),
		Name:      name,
		Stages:    stages,
		CreatedAt: now,
	}
}

// AdvanceToNextStage переводит первый незавершённый этап в состояние "в работе".
//
// Функция возвращает обновлённый процесс и булево значение, показывающее был ли найден этап.
func AdvanceToNextStage(process model.BusinessProcess) (model.BusinessProcess, bool) {
	for i, stage := range process.Stages {
		if stage.Status == model.StagePending {
			stage.Status = model.StageInProgress
			stage.UpdatedAt = time.Now()
			process.Stages[i] = stage
			return process, true
		}
		if stage.Status == model.StageInProgress {
			// Этап уже выполняется, переход не требуется.
			return process, true
		}
	}
	return process, false
}

// CompleteStage помечает этап завершённым.
func CompleteStage(process model.BusinessProcess, stageID string, success bool) model.BusinessProcess {
	for i, stage := range process.Stages {
		if stage.ID == stageID {
			if success {
				stage.Status = model.StageCompleted
			} else {
				stage.Status = model.StageRejected
			}
			stage.UpdatedAt = time.Now()
			process.Stages[i] = stage
			break
		}
	}
	return process
}
