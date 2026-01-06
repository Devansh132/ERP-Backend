package models

import (
	"time"
	"gorm.io/gorm"
)

type Student struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	UserID          uint           `gorm:"not null;unique" json:"user_id"`
	AdmissionNumber string         `gorm:"unique;not null" json:"admission_number"`
	FirstName       string         `gorm:"not null" json:"first_name"`
	LastName        string         `gorm:"not null" json:"last_name"`
	DateOfBirth     time.Time      `json:"date_of_birth"`
	Gender          string         `json:"gender"`
	Address         string         `json:"address"`
	Phone           string         `json:"phone"`
	ParentName      string         `json:"parent_name"`
	ParentPhone     string         `json:"parent_phone"`
	ClassID         uint           `gorm:"not null" json:"class_id"`
	SectionID       uint           `gorm:"not null" json:"section_id"`
	Status          string         `gorm:"default:active" json:"status"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Class   Class   `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	Section Section `gorm:"foreignKey:SectionID" json:"section,omitempty"`
}



