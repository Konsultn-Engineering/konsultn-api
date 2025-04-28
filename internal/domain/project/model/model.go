package model

import (
	"konsultn-api/internal/domain/task"
	"konsultn-api/internal/shared"
)

type Project struct {
	shared.ULID `gorm:"embedded"`
	Name        string      `gorm:"size:255"`
	Tasks       []task.Task `gorm:"foreignKey:ProjectID"`
}
