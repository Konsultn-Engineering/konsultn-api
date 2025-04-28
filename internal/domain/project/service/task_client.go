package service

import (
	"konsultn-api/internal/domain/project/dto"
	"konsultn-api/internal/domain/project/model"
)

type TaskClient interface {
	CreateTaskForProject(project *model.Project, task dto.TaskDTO) error
}
