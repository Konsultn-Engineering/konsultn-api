package service

import (
	"gorm.io/gorm"
	"konsultn-api/internal/domain/project/dto"
	"konsultn-api/internal/domain/project/model"
	"konsultn-api/internal/domain/project/repository"
	"konsultn-api/internal/shared/mapper"
)

type ProjectService struct {
	Repo       *repository.Repository[model.Project]
	taskClient TaskClient // interface
}

func NewProjectService(db *gorm.DB) *ProjectService {
	return &ProjectService{
		Repo:       repository.NewRepository(db),
		taskClient: NewTaskLocal(db),
	}
}

func (s *ProjectService) CreateProjectTask(projectID string, t dto.TaskDTO) (*dto.ProjectDTO, error) {

	project, err := s.Repo.FindById(projectID)
	if err != nil {
		return nil, err
	}

	err = s.taskClient.CreateTaskForProject(project, t)
	if err != nil {
		return nil, err
	}
	s.Repo.Preload(&project, []string{"Tasks", "Tasks.Assignee"}, "id", projectID)
	projectDTO, err := mapper.Convert[dto.ProjectDTO](project)

	return &projectDTO, nil
}
