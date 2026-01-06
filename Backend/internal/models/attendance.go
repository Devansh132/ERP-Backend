package models

import (
	"time"
	"gorm.io/gorm"
)

type Attendance struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	StudentID uint          `gorm:"not null" json:"student_id"`
	ClassID   uint          `gorm:"not null" json:"class_id"`
	SectionID uint          `gorm:"not null" json:"section_id"`
	Date      time.Time     `gorm:"not null" json:"date"`
	Status    string        `gorm:"not null;check:status IN ('present','absent','late','excused')" json:"status"`
	MarkedBy  uint          `gorm:"not null" json:"marked_by"` // User ID of teacher/admin
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Student Student `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Class   Class   `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	Section Section `gorm:"foreignKey:SectionID" json:"section,omitempty"`
}



