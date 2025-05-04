package model

import (
	"gorm.io/gorm"
	"konsultn-api/internal/shared"
	"time"
)

type Team struct {
	// Identifiers
	shared.ULID `gorm:"embedded"`
	Name        string       `gorm:"type:varchar(255);not null"`
	Slug        string       `gorm:"type:varchar(255);uniqueIndex;not null"`
	Description string       `gorm:"type:varchar(255)"`
	OwnerID     string       `gorm:"type:varchar(255);not null;index"`
	Owner       *UserView    `gorm:"-"` // ignored by GORM; populated manually
	Members     []TeamMember `gorm:"foreignKey:TeamID;constraint:OnDelete:CASCADE"`
	UpdatedBy   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `swaggerignore:"true"`
}
