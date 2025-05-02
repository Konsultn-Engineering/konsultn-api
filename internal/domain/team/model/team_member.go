package model

import (
	"konsultn-api/internal/shared"
	"time"
)

type TeamMember struct {
	shared.ULID `gorm:"embedded"`
	TeamID      string `gorm:"uniqueIndex:idx_team_member;not null"`
	UserID      string `gorm:"uniqueIndex:idx_team_member;not null"`
	Role        string
	JoinedAt    time.Time `gorm:"autoCreateTime"`
	User        UserView  `gorm:"-"`
}
