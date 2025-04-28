package taskservice

import (
	"context"
	"konsultn-api/internal/domain/task"
)

type TaskService interface {
	CreateTask(req task.Task) (*task.Task, error)
	GetTaskByID(id string) (*task.Task, error)
}

type TaskUserService interface {
	ValidateAssignee(ctx context.Context, userID string) error
}
