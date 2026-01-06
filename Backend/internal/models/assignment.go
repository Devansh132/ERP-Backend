package models

import (
	"time"
	"gorm.io/gorm"
)

type Assignment struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `json:"description"`
	SubjectID   uint           `gorm:"not null" json:"subject_id"`
	ClassID     uint           `gorm:"not null" json:"class_id"`
	TeacherID   uint           `gorm:"not null" json:"teacher_id"`
	DueDate     time.Time      `gorm:"not null" json:"due_date"`
	TotalMarks  float64        `json:"total_marks"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Subject Subject `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
	Class   Class   `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	Teacher Teacher `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
}

type AssignmentSubmission struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	AssignmentID  uint           `gorm:"not null" json:"assignment_id"`
	StudentID     uint           `gorm:"not null" json:"student_id"`
	SubmissionDate time.Time     `gorm:"not null" json:"submission_date"`
	FilePath      string         `json:"file_path"`
	MarksObtained float64        `json:"marks_obtained"`
	Feedback      string         `json:"feedback"`
	Status        string         `gorm:"default:pending" json:"status"` // pending, graded, late
	SubmittedAt   time.Time      `json:"submitted_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Assignment Assignment `gorm:"foreignKey:AssignmentID" json:"assignment,omitempty"`
	Student    Student    `gorm:"foreignKey:StudentID" json:"student,omitempty"`
}



