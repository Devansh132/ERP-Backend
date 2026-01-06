package models

import (
	"time"
	"gorm.io/gorm"
)

type LeaveRequest struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	TeacherID   uint           `gorm:"not null" json:"teacher_id"`
	LeaveType   string         `gorm:"not null" json:"leave_type"` // sick, casual, personal, emergency
	StartDate   time.Time      `gorm:"not null" json:"start_date"`
	EndDate     time.Time      `gorm:"not null" json:"end_date"`
	Reason      string         `gorm:"type:text" json:"reason"`
	Status      string         `gorm:"default:pending" json:"status"` // pending, approved, rejected
	ApprovedBy  *uint          `json:"approved_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Teacher   Teacher `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	Approver  *User   `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`
}



