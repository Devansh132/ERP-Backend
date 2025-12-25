package models

import (
	"time"
	"gorm.io/gorm"
)

type Exam struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"not null" json:"name"`
	ExamType     string         `gorm:"not null" json:"exam_type"` // e.g., "midterm", "final", "quiz"
	ClassID      uint           `gorm:"not null" json:"class_id"`
	SubjectID    uint           `gorm:"not null" json:"subject_id"`
	ExamDate     time.Time      `json:"exam_date"`
	TotalMarks   float64        `gorm:"not null" json:"total_marks"`
	PassingMarks float64        `json:"passing_marks"`
	AcademicYear string         `gorm:"not null" json:"academic_year"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Class   Class   `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	Subject Subject `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
}

type Mark struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	StudentID    uint           `gorm:"not null" json:"student_id"`
	SubjectID    uint           `gorm:"not null" json:"subject_id"`
	ExamID       uint           `gorm:"not null" json:"exam_id"`
	ExamType     string         `gorm:"not null" json:"exam_type"`
	MarksObtained float64       `gorm:"not null" json:"marks_obtained"`
	TotalMarks   float64        `gorm:"not null" json:"total_marks"`
	Percentage   float64        `json:"percentage"`
	Grade        string         `json:"grade"`
	AcademicYear string         `gorm:"not null" json:"academic_year"`
	CreatedBy    uint           `gorm:"not null" json:"created_by"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Student Student `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Subject Subject `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
	Exam    Exam    `gorm:"foreignKey:ExamID" json:"exam,omitempty"`
}

type Subject struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	Code      string         `gorm:"unique;not null" json:"code"`
	Status    string         `gorm:"default:active" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

