package dto

import (
	"konsultn-api/internal/domain/project/model"
	"konsultn-api/internal/domain/task"
)

func FromModelProject(project *model.Project) ProjectDTO {
	tasks := make([]TaskDTO, len(project.Tasks))
	for i, task := range project.Tasks {
		tasks[i] = FromModelTask(&task)
	}

	return ProjectDTO{
		ID:    project.ID,
		Name:  project.Name,
		Tasks: tasks,
	}
}

// FromModelTask maps a Task GORM model to a TaskDTO
func FromModelTask(task *task.Task) TaskDTO {
	return TaskDTO{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		AssigneeID:  task.AssigneeID,
	}
}

func FromDTOProject(project ProjectDTO) model.Project {
	tasks := make([]task.Task, len(project.Tasks))
	for i, t := range project.Tasks {
		tasks[i] = FromDTOTask(t)
	}

	return model.Project{
		Name:  project.Name,
		Tasks: tasks,
	}
}

func FromDTOTask(dto TaskDTO) task.Task {
	return task.Task{
		Title:       dto.Title,
		Description: dto.Description,
		AssigneeID:  dto.AssigneeID, // assuming it's already the correct ULID type or will be parsed later
	}
}
