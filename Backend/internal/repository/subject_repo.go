package repository

import (
	"school-erp-backend/internal/models"
	"gorm.io/gorm"
)

type SubjectRepository struct {
	db *gorm.DB
}

func NewSubjectRepository(db *gorm.DB) *SubjectRepository {
	return &SubjectRepository{db: db}
}

func (r *SubjectRepository) Create(subject *models.Subject) error {
	return r.db.Create(subject).Error
}

func (r *SubjectRepository) FindByID(id uint) (*models.Subject, error) {
	var subject models.Subject
	err := r.db.First(&subject, id).Error
	return &subject, err
}

func (r *SubjectRepository) FindByCode(code string) (*models.Subject, error) {
	var subject models.Subject
	err := r.db.Where("code = ?", code).First(&subject).Error
	return &subject, err
}

func (r *SubjectRepository) Update(subject *models.Subject) error {
	return r.db.Save(subject).Error
}

func (r *SubjectRepository) Delete(id uint) error {
	return r.db.Delete(&models.Subject{}, id).Error
}

func (r *SubjectRepository) FindAll() ([]models.Subject, error) {
	var subjects []models.Subject
	err := r.db.Order("name ASC").Find(&subjects).Error
	return subjects, err
}



