package taskservice

import (
	"konsultn-api/internal/domain/task"
)

type service struct {
	repo *task.Repository[task.Task]
}

func NewService(repo *task.Repository[task.Task]) TaskService {
	return &service{repo: repo}
}

func (s *service) CreateTask(taskDto task.Task) (task.Task, error) {
	createdTask, err := s.repo.Save(taskDto)

	if err != nil {
		return task.Task{}, err
	}

	return createdTask, nil
}

func (s *service) GetTaskByID(id string) (task.Task, error) {
	createdTask, err := s.repo.FindById(id)
	return createdTask, err
}
