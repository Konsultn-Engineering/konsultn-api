package shared

import (
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type ULID struct {
	ID string `gorm:"primaryKey" json:"id"`
}

func (m *ULID) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		t := time.Now()
		entropy := rand.New(rand.NewSource(t.UnixNano()))
		m.ID = ulid.MustNew(ulid.Timestamp(t), entropy).String()
	}
	return
}
