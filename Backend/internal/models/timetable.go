package models

import (
	"time"
	"gorm.io/gorm"
)

type Timetable struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ClassID     uint           `gorm:"not null" json:"class_id"`
	SectionID   uint           `gorm:"not null" json:"section_id"`
	Day         string         `gorm:"not null" json:"day"` // Monday, Tuesday, etc.
	PeriodNumber int           `gorm:"not null" json:"period_number"`
	SubjectID   uint           `gorm:"not null" json:"subject_id"`
	TeacherID   uint           `gorm:"not null" json:"teacher_id"`
	StartTime   string         `gorm:"not null" json:"start_time"` // HH:MM format
	EndTime     string         `gorm:"not null" json:"end_time"`   // HH:MM format
	RoomNumber  string         `json:"room_number"`
	AcademicYear string        `gorm:"not null" json:"academic_year"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Class   Class   `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	Section Section `gorm:"foreignKey:SectionID" json:"section,omitempty"`
	Subject Subject `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
	Teacher Teacher `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
}

