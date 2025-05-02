package model

import (
	"konsultn-api/internal/shared"
	"time"
)

type Team struct {
	shared.ULID `gorm:"embedded"`
	Name        string
	Slug        string
	OwnerID     string       // who created/administers the team
	Owner       *UserView    `gorm:"-"`
	Members     []TeamMember `gorm:"foreignKey:TeamID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
