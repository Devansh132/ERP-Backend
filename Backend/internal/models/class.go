package models

import (
	"time"
	"gorm.io/gorm"
)

type Class struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"not null" json:"name"`      // e.g., "1st", "2nd", "3rd"
	Level     int            `gorm:"not null" json:"level"`     // Numeric level for sorting
	Capacity  int            `json:"capacity"`
	Status    string         `gorm:"default:active" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Sections []Section `gorm:"foreignKey:ClassID" json:"sections,omitempty"`
}

type Section struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ClassID   uint           `gorm:"not null" json:"class_id"`
	Name      string         `gorm:"not null" json:"name"` // e.g., "A", "B", "C"
	Capacity  int            `json:"capacity"`
	Status    string         `gorm:"default:active" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Class Class `gorm:"foreignKey:ClassID" json:"class,omitempty"`
}

type ClassSection struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ClassID     uint      `gorm:"not null" json:"class_id"`
	SectionID   uint      `gorm:"not null" json:"section_id"`
	AcademicYear string   `gorm:"not null" json:"academic_year"`
	Status      string    `gorm:"default:active" json:"status"`
	CreatedAt   time.Time `json:"created_at"`

	Class   Class   `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	Section Section `gorm:"foreignKey:SectionID" json:"section,omitempty"`
}



