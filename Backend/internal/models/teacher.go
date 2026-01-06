package models

import (
	"time"
	"gorm.io/gorm"
)

type Teacher struct {
	ID                  uint           `gorm:"primaryKey" json:"id"`
	UserID              uint           `gorm:"not null;unique" json:"user_id"`
	EmployeeID          string         `gorm:"unique;not null" json:"employee_id"`
	FirstName           string         `gorm:"not null" json:"first_name"`
	LastName            string         `gorm:"not null" json:"last_name"`
	DateOfBirth         time.Time      `json:"date_of_birth"`
	Gender              string         `json:"gender"`
	Address             string         `json:"address"`
	Phone               string         `json:"phone"`
	Qualification      string         `json:"qualification"`
	Experience          int            `json:"experience"` // years
	SubjectSpecialization string       `json:"subject_specialization"`
	Status              string         `gorm:"default:active" json:"status"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}



