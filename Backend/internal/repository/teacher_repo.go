package repository

import (
	"school-erp-backend/internal/models"
	"gorm.io/gorm"
)

type TeacherRepository struct {
	db *gorm.DB
}

func NewTeacherRepository(db *gorm.DB) *TeacherRepository {
	return &TeacherRepository{db: db}
}

func (r *TeacherRepository) Create(teacher *models.Teacher) error {
	return r.db.Create(teacher).Error
}

func (r *TeacherRepository) FindByID(id uint) (*models.Teacher, error) {
	var teacher models.Teacher
	err := r.db.Preload("User").First(&teacher, id).Error
	return &teacher, err
}

func (r *TeacherRepository) FindByEmployeeID(employeeID string) (*models.Teacher, error) {
	var teacher models.Teacher
	err := r.db.Where("employee_id = ?", employeeID).First(&teacher).Error
	return &teacher, err
}

func (r *TeacherRepository) FindByUserID(userID uint) (*models.Teacher, error) {
	var teacher models.Teacher
	err := r.db.Where("user_id = ?", userID).Preload("User").First(&teacher).Error
	return &teacher, err
}

func (r *TeacherRepository) Update(teacher *models.Teacher) error {
	return r.db.Save(teacher).Error
}

func (r *TeacherRepository) Delete(id uint) error {
	return r.db.Delete(&models.Teacher{}, id).Error
}

func (r *TeacherRepository) FindAll() ([]models.Teacher, error) {
	var teachers []models.Teacher
	err := r.db.Preload("User").Find(&teachers).Error
	return teachers, err
}

