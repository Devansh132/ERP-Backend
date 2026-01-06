package repository

import (
	"school-erp-backend/internal/models"
	"gorm.io/gorm"
)

type StudentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) Create(student *models.Student) error {
	return r.db.Create(student).Error
}

func (r *StudentRepository) FindByID(id uint) (*models.Student, error) {
	var student models.Student
	err := r.db.Preload("User").Preload("Class").Preload("Section").First(&student, id).Error
	return &student, err
}

func (r *StudentRepository) FindByAdmissionNumber(admissionNumber string) (*models.Student, error) {
	var student models.Student
	err := r.db.Where("admission_number = ?", admissionNumber).First(&student).Error
	return &student, err
}

func (r *StudentRepository) Update(student *models.Student) error {
	return r.db.Save(student).Error
}

func (r *StudentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Student{}, id).Error
}

func (r *StudentRepository) FindByClassAndSection(classID, sectionID uint) ([]models.Student, error) {
	var students []models.Student
	err := r.db.Where("class_id = ? AND section_id = ?", classID, sectionID).Find(&students).Error
	return students, err
}



