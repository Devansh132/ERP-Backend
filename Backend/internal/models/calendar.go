package models

import (
	"time"
	"gorm.io/gorm"
)

type CalendarEvent struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	StartDate   time.Time      `gorm:"not null" json:"start_date"`
	EndDate     time.Time      `json:"end_date"`
	EventType   string         `gorm:"not null" json:"event_type"` // holiday, exam, event, meeting
	Visibility  string         `gorm:"default:all" json:"visibility"` // all, admin, teacher, student
	CreatedBy   uint           `gorm:"not null" json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Creator User `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}



