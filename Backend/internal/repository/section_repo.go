package repository

import (
	"school-erp-backend/internal/models"
	"gorm.io/gorm"
)

type SectionRepository struct {
	db *gorm.DB
}

func NewSectionRepository(db *gorm.DB) *SectionRepository {
	return &SectionRepository{db: db}
}

func (r *SectionRepository) Create(section *models.Section) error {
	return r.db.Create(section).Error
}

func (r *SectionRepository) FindByID(id uint) (*models.Section, error) {
	var section models.Section
	err := r.db.Preload("Class").First(&section, id).Error
	return &section, err
}

func (r *SectionRepository) FindByClassID(classID uint) ([]models.Section, error) {
	var sections []models.Section
	err := r.db.Where("class_id = ?", classID).Preload("Class").Find(&sections).Error
	return sections, err
}

func (r *SectionRepository) FindByClassAndName(classID uint, name string) (*models.Section, error) {
	var section models.Section
	err := r.db.Where("class_id = ? AND name = ?", classID, name).First(&section).Error
	return &section, err
}

func (r *SectionRepository) Update(section *models.Section) error {
	return r.db.Save(section).Error
}

func (r *SectionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Section{}, id).Error
}

func (r *SectionRepository) FindAll() ([]models.Section, error) {
	var sections []models.Section
	err := r.db.Preload("Class").Find(&sections).Error
	return sections, err
}

