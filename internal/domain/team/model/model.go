package model

import (
	"konsultn-api/internal/domain/user"
	"konsultn-api/internal/shared"
	"time"
)

type Team struct {
	shared.ULID `gorm:"embedded"`
	Name        string
	Slug        string
	OwnerID     string       // who created/administers the team
	Owner       user.User    `gorm:"foreignKey:OwnerID"`
	Members     []TeamMember `gorm:"foreignKey:TeamID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TeamMember struct {
	shared.ULID `gorm:"embedded"`
	TeamID      string
	UserID      string
	Role        string // e.g. "owner", "admin", "member", "viewer"
	JoinedAt    time.Time
	User        user.User `gorm:"foreignKey:UserID"`
}
