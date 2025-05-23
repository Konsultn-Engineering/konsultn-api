package service

import (
	"gorm.io/gorm"
	"konsultn-api/internal/domain/project/dto"
	"konsultn-api/internal/domain/project/model"
)

type TaskLocal struct {
	db *gorm.DB
}

func NewTaskLocal(db *gorm.DB) *TaskLocal {
	return &TaskLocal{db: db}
}

func (c *TaskLocal) CreateTaskForProject(project *model.Project, t dto.TaskDTO) error {
	task := dto.FromDTOTask(t)
	task.ProjectID = &project.ID // Set foreign key

	// Save the task first
	if err := c.db.Create(&task).Error; err != nil {
		return err
	}

	// Associate task to project (this is optional if you set FK correctly, but still good practice)
	return c.db.Model(project).Association("Tasks").Append(&task)
}
