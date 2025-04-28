package task

import (
	"gorm.io/gorm"
	"konsultn-api/internal/domain/user"
	"konsultn-api/internal/shared"
	"time"
)

type Task struct {
	shared.ULID  `gorm:"embedded"`
	ProjectID    *string        `gorm:"type:varchar(26)" json:"project_id"`
	Title        string         `gorm:"type:varchar(255);not null" json:"title"`
	Description  string         `gorm:"type:varchar(2048); not null" json:"description"`
	Status       string         `gorm:"varchar(50)" json:"status"`
	Priority     string         `gorm:"type:varchar(20)" json:"priority"`
	DueDate      *time.Time     `gorm:"type:timestamp" json:"due_date"`
	AssigneeID   *string        `gorm:"type:varchar(26);" json:"assignee_id"`
	Assignee     user.User      `json:"assignee"`
	ParentTaskID *string        `gorm:"type:varchar(26)" json:"parent_task_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" swaggerignore:"true"`
}
