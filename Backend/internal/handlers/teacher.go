package handlers

import (
	"net/http"
	"strconv"
	"time"
	"school-erp-backend/internal/models"
	"school-erp-backend/internal/repository"
	"school-erp-backend/pkg/database"
	"github.com/gin-gonic/gin"
)

type TeacherHandler struct {
	teacherRepo *repository.TeacherRepository
}

func NewTeacherHandler() *TeacherHandler {
	return &TeacherHandler{
		teacherRepo: repository.NewTeacherRepository(database.DB),
	}
}

// GetTeachers godoc
// @Summary Get all teachers
// @Description Get list of all teachers with optional filters
// @Tags Admin - Teachers
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {array} models.Teacher
// @Failure 500 {object} ErrorResponse
// @Router /admin/teachers [get]
// @Security BearerAuth
func (h *TeacherHandler) GetTeachers(c *gin.Context) {
	var teachers []models.Teacher
	query := database.DB.Preload("User")

	if err := query.Find(&teachers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teachers)
}

// GetTeacher godoc
// @Summary Get teacher by ID
// @Description Get a specific teacher by ID
// @Tags Admin - Teachers
// @Accept json
// @Produce json
// @Param id path int true "Teacher ID"
// @Success 200 {object} models.Teacher
// @Failure 404 {object} ErrorResponse
// @Router /admin/teachers/{id} [get]
// @Security BearerAuth
func (h *TeacherHandler) GetTeacher(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	teacher, err := h.teacherRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher not found"})
		return
	}

	c.JSON(http.StatusOK, teacher)
}

// CreateTeacher godoc
// @Summary Create a new teacher
// @Description Create a new teacher record
// @Tags Admin - Teachers
// @Accept json
// @Produce json
// @Param teacher body CreateTeacherRequest true "Teacher data"
// @Success 201 {object} models.Teacher
// @Failure 400 {object} ErrorResponse
// @Router /admin/teachers [post]
// @Security BearerAuth
func (h *TeacherHandler) CreateTeacher(c *gin.Context) {
	var req CreateTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse date of birth
	dateOfBirth, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	teacher := &models.Teacher{
		UserID:              req.UserID,
		EmployeeID:          req.EmployeeID,
		FirstName:           req.FirstName,
		LastName:            req.LastName,
		DateOfBirth:         dateOfBirth,
		Gender:              req.Gender,
		Address:             req.Address,
		Phone:               req.Phone,
		Qualification:      req.Qualification,
		Experience:          req.Experience,
		SubjectSpecialization: req.SubjectSpecialization,
		Status:              "active",
	}

	if err := h.teacherRepo.Create(teacher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, teacher)
}

// UpdateTeacher godoc
// @Summary Update teacher
// @Description Update an existing teacher record
// @Tags Admin - Teachers
// @Accept json
// @Produce json
// @Param id path int true "Teacher ID"
// @Param teacher body UpdateTeacherRequest true "Teacher data"
// @Success 200 {object} models.Teacher
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/teachers/{id} [put]
// @Security BearerAuth
func (h *TeacherHandler) UpdateTeacher(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	var req UpdateTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teacher, err := h.teacherRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher not found"})
		return
	}

	// Update fields
	if req.FirstName != "" {
		teacher.FirstName = req.FirstName
	}
	if req.LastName != "" {
		teacher.LastName = req.LastName
	}
	if req.DateOfBirth != "" {
		dateOfBirth, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
			return
		}
		teacher.DateOfBirth = dateOfBirth
	}
	if req.Gender != "" {
		teacher.Gender = req.Gender
	}
	if req.Address != "" {
		teacher.Address = req.Address
	}
	if req.Phone != "" {
		teacher.Phone = req.Phone
	}
	if req.Qualification != "" {
		teacher.Qualification = req.Qualification
	}
	if req.Experience != 0 {
		teacher.Experience = req.Experience
	}
	if req.SubjectSpecialization != "" {
		teacher.SubjectSpecialization = req.SubjectSpecialization
	}
	if req.Status != "" {
		teacher.Status = req.Status
	}

	if err := h.teacherRepo.Update(teacher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teacher)
}

// DeleteTeacher godoc
// @Summary Delete teacher
// @Description Delete a teacher record
// @Tags Admin - Teachers
// @Accept json
// @Produce json
// @Param id path int true "Teacher ID"
// @Success 200 {object} SuccessResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/teachers/{id} [delete]
// @Security BearerAuth
func (h *TeacherHandler) DeleteTeacher(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	if err := h.teacherRepo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Teacher deleted successfully"})
}

// Request Types
type CreateTeacherRequest struct {
	UserID              uint   `json:"user_id" binding:"required"`
	EmployeeID          string `json:"employee_id" binding:"required"`
	FirstName           string `json:"first_name" binding:"required"`
	LastName            string `json:"last_name" binding:"required"`
	DateOfBirth         string `json:"date_of_birth" binding:"required"`
	Gender              string `json:"gender"`
	Address             string `json:"address"`
	Phone               string `json:"phone"`
	Qualification       string `json:"qualification"`
	Experience          int    `json:"experience"`
	SubjectSpecialization string `json:"subject_specialization"`
}

type UpdateTeacherRequest struct {
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name"`
	DateOfBirth         string `json:"date_of_birth"`
	Gender              string `json:"gender"`
	Address             string `json:"address"`
	Phone               string `json:"phone"`
	Qualification       string `json:"qualification"`
	Experience          int    `json:"experience"`
	SubjectSpecialization string `json:"subject_specialization"`
	Status              string `json:"status"`
}

