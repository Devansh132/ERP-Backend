package repository

import (
	"school-erp-backend/internal/models"
	"time"
	"gorm.io/gorm"
)

type AttendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) *AttendanceRepository {
	return &AttendanceRepository{db: db}
}

func (r *AttendanceRepository) Create(attendance *models.Attendance) error {
	return r.db.Create(attendance).Error
}

func (r *AttendanceRepository) CreateBatch(attendances []models.Attendance) error {
	return r.db.Create(&attendances).Error
}

func (r *AttendanceRepository) FindByID(id uint) (*models.Attendance, error) {
	var attendance models.Attendance
	err := r.db.Preload("Student").Preload("Student.User").Preload("Class").Preload("Section").
		First(&attendance, id).Error
	return &attendance, err
}

func (r *AttendanceRepository) FindByStudentID(studentID uint, startDate, endDate time.Time) ([]models.Attendance, error) {
	var attendances []models.Attendance
	query := r.db.Where("student_id = ?", studentID)
	
	if !startDate.IsZero() {
		query = query.Where("date >= ?", startDate)
	}
	if !endDate.IsZero() {
		query = query.Where("date <= ?", endDate)
	}
	
	err := query.Preload("Class").Preload("Section").
		Order("date DESC").Find(&attendances).Error
	return attendances, err
}

func (r *AttendanceRepository) FindByClassID(classID uint, date time.Time) ([]models.Attendance, error) {
	var attendances []models.Attendance
	query := r.db.Where("class_id = ?", classID)
	
	if !date.IsZero() {
		query = query.Where("DATE(date) = DATE(?)", date)
	}
	
	err := query.Preload("Student").Preload("Student.User").Preload("Section").
		Order("student_id ASC").Find(&attendances).Error
	return attendances, err
}

func (r *AttendanceRepository) FindByClassAndSection(classID, sectionID uint, date time.Time) ([]models.Attendance, error) {
	var attendances []models.Attendance
	query := r.db.Where("class_id = ? AND section_id = ?", classID, sectionID)
	
	if !date.IsZero() {
		query = query.Where("DATE(date) = DATE(?)", date)
	}
	
	err := query.Preload("Student").Preload("Student.User").
		Order("student_id ASC").Find(&attendances).Error
	return attendances, err
}

func (r *AttendanceRepository) FindByDateRange(startDate, endDate time.Time, classID, sectionID *uint) ([]models.Attendance, error) {
	var attendances []models.Attendance
	query := r.db.Where("date >= ? AND date <= ?", startDate, endDate)
	
	if classID != nil {
		query = query.Where("class_id = ?", *classID)
	}
	if sectionID != nil {
		query = query.Where("section_id = ?", *sectionID)
	}
	
	err := query.Preload("Student").Preload("Student.User").Preload("Class").Preload("Section").
		Order("date DESC, student_id ASC").Find(&attendances).Error
	return attendances, err
}

func (r *AttendanceRepository) FindExisting(studentID uint, date time.Time) (*models.Attendance, error) {
	var attendance models.Attendance
	err := r.db.Where("student_id = ? AND DATE(date) = DATE(?)", studentID, date).
		First(&attendance).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &attendance, err
}

func (r *AttendanceRepository) Update(attendance *models.Attendance) error {
	return r.db.Save(attendance).Error
}

func (r *AttendanceRepository) Delete(id uint) error {
	return r.db.Delete(&models.Attendance{}, id).Error
}

// GetStatistics returns attendance statistics for a student
func (r *AttendanceRepository) GetStatistics(studentID uint, startDate, endDate time.Time) (map[string]interface{}, error) {
	var stats struct {
		Total     int64
		Present   int64
		Absent    int64
		Late      int64
		Excused   int64
	}
	
	baseQuery := r.db.Model(&models.Attendance{}).Where("student_id = ?", studentID)
	
	if !startDate.IsZero() {
		baseQuery = baseQuery.Where("date >= ?", startDate)
	}
	if !endDate.IsZero() {
		baseQuery = baseQuery.Where("date <= ?", endDate)
	}
	
	baseQuery.Count(&stats.Total)
	
	// Count by status
	baseQuery.Where("status = ?", "present").Count(&stats.Present)
	baseQuery.Where("status = ?", "absent").Count(&stats.Absent)
	baseQuery.Where("status = ?", "late").Count(&stats.Late)
	baseQuery.Where("status = ?", "excused").Count(&stats.Excused)
	
	var percentage float64
	if stats.Total > 0 {
		percentage = float64(stats.Present) / float64(stats.Total) * 100
	}
	
	return map[string]interface{}{
		"total":      stats.Total,
		"present":    stats.Present,
		"absent":     stats.Absent,
		"late":       stats.Late,
		"excused":    stats.Excused,
		"percentage": percentage,
	}, nil
}

// GetClassStatistics returns attendance statistics for a class
func (r *AttendanceRepository) GetClassStatistics(classID, sectionID uint, startDate, endDate time.Time) (map[string]interface{}, error) {
	var stats struct {
		Total     int64
		Present   int64
		Absent    int64
		Late      int64
		Excused   int64
	}
	
	baseQuery := r.db.Model(&models.Attendance{}).Where("class_id = ?", classID)
	
	if sectionID > 0 {
		baseQuery = baseQuery.Where("section_id = ?", sectionID)
	}
	if !startDate.IsZero() {
		baseQuery = baseQuery.Where("date >= ?", startDate)
	}
	if !endDate.IsZero() {
		baseQuery = baseQuery.Where("date <= ?", endDate)
	}
	
	baseQuery.Count(&stats.Total)
	
	// Count by status
	baseQuery.Where("status = ?", "present").Count(&stats.Present)
	baseQuery.Where("status = ?", "absent").Count(&stats.Absent)
	baseQuery.Where("status = ?", "late").Count(&stats.Late)
	baseQuery.Where("status = ?", "excused").Count(&stats.Excused)
	
	var percentage float64
	if stats.Total > 0 {
		percentage = float64(stats.Present) / float64(stats.Total) * 100
	}
	
	return map[string]interface{}{
		"total":      stats.Total,
		"present":    stats.Present,
		"absent":     stats.Absent,
		"late":       stats.Late,
		"excused":    stats.Excused,
		"percentage": percentage,
	}, nil
}

