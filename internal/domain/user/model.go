package user

import (
	"gorm.io/gorm"
	"konsultn-api/internal/shared"
	"time"
)

type User struct {
	shared.ULID          `gorm:"embedded"`
	UID                  string `gorm:"size:255;unique"`
	FirstName            string `gorm:"size:255"`
	LastName             string `gorm:"size:255"`
	Email                string `gorm:"unique;size:255;not null"`
	Password             string `json:"password"` // raw password from user
	PasswordHash         string `gorm:"size:255" json:"-"`
	PhoneNumber          string `gorm:"size:20"`
	ProfilePictureURL    string `gorm:"size:255"`
	SocialProvider       string `gorm:"size:50"`
	SocialID             string `gorm:"size:255"`
	SocialEmail          string `gorm:"size:255"`
	SocialProfilePicture string `gorm:"size:255"`
	AccessToken          string `gorm:"size:255"`
	RefreshToken         string `gorm:"size:255"`
	Status               string `gorm:"size:255;default:'PENDING VERIFICATION'"`
	LastLogin            *time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            gorm.DeletedAt `json:"deleted_at" swaggerignore:"true"`
	ResetToken           string         `gorm:"size:255"`
	ResetTokenExpiry     *time.Time
	TwoFactorEnabled     bool   `gorm:"default:false"`
	TwoFactorSecret      string `gorm:"size:255"`
}
