package models

import (
	"time"
	"gorm.io/gorm"
)

type Notice struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Title         string         `gorm:"not null" json:"title"`
	Content       string         `gorm:"type:text" json:"content"`
	Category      string         `json:"category"`
	Priority      string         `gorm:"default:normal" json:"priority"` // low, normal, high, urgent
	VisibilityType string        `gorm:"not null" json:"visibility_type"` // all, class, role
	TargetAudience string        `json:"target_audience"` // JSON string for class IDs or roles
	PublishedAt   time.Time      `json:"published_at"`
	CreatedBy     uint           `gorm:"not null" json:"created_by"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Creator User `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}



