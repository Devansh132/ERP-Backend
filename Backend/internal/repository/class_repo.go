package repository

import (
	"school-erp-backend/internal/models"
	"gorm.io/gorm"
)

type ClassRepository struct {
	db *gorm.DB
}

func NewClassRepository(db *gorm.DB) *ClassRepository {
	return &ClassRepository{db: db}
}

func (r *ClassRepository) Create(class *models.Class) error {
	return r.db.Create(class).Error
}

func (r *ClassRepository) FindByID(id uint) (*models.Class, error) {
	var class models.Class
	err := r.db.Preload("Sections").First(&class, id).Error
	return &class, err
}

func (r *ClassRepository) FindByName(name string) (*models.Class, error) {
	var class models.Class
	err := r.db.Where("name = ?", name).First(&class).Error
	return &class, err
}

func (r *ClassRepository) Update(class *models.Class) error {
	return r.db.Save(class).Error
}

func (r *ClassRepository) Delete(id uint) error {
	return r.db.Delete(&models.Class{}, id).Error
}

func (r *ClassRepository) FindAll() ([]models.Class, error) {
	var classes []models.Class
	err := r.db.Preload("Sections").Order("level ASC").Find(&classes).Error
	return classes, err
}

