package model

import (
	"gorm.io/gorm"
	"konsultn-api/internal/shared"
	"time"
)

type TeamInvitation struct {
	shared.ULID `gorm:"embedded"`
	FromUserID  string         `gorm:"not null"`
	ToUserID    string         `gorm:"uniqueIndex:idx_invite_team_user;not null"`
	TeamID      string         `gorm:"uniqueIndex:idx_invite_team_user;not null"`
	Message     *string        `gorm:"size:255"` // Optional message
	Status      string         `gorm:"not null"`
	Role        string         `gorm:"not null"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `swaggerignore:"true"`
	ExpiresAt   *time.Time     `gorm:"index"` // Expiry date for invitation, optional
}
