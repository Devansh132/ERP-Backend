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

type StudentHandler struct {
	studentRepo *repository.StudentRepository
}

func NewStudentHandler() *StudentHandler {
	return &StudentHandler{
		studentRepo: repository.NewStudentRepository(database.DB),
	}
}

// GetStudents godoc
// @Summary Get all students
// @Description Get list of all students with optional filters
// @Tags Admin - Students
// @Accept json
// @Produce json
// @Param class_id query int false "Filter by class ID"
// @Param section_id query int false "Filter by section ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {array} models.Student
// @Failure 500 {object} ErrorResponse
// @Router /admin/students [get]
// @Security BearerAuth
func (h *StudentHandler) GetStudents(c *gin.Context) {
	classID := c.Query("class_id")
	sectionID := c.Query("section_id")
	
	var students []models.Student
	query := database.DB.Preload("User").Preload("Class").Preload("Section")

	if classID != "" {
		query = query.Where("class_id = ?", classID)
	}
	if sectionID != "" {
		query = query.Where("section_id = ?", sectionID)
	}

	if err := query.Find(&students).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, students)
}

// GetStudent godoc
// @Summary Get student by ID
// @Description Get a specific student by ID
// @Tags Admin - Students
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Success 200 {object} models.Student
// @Failure 404 {object} ErrorResponse
// @Router /admin/students/{id} [get]
// @Security BearerAuth
func (h *StudentHandler) GetStudent(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	student, err := h.studentRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	c.JSON(http.StatusOK, student)
}

// CreateStudent godoc
// @Summary Create a new student
// @Description Create a new student record
// @Tags Admin - Students
// @Accept json
// @Produce json
// @Param student body CreateStudentRequest true "Student data"
// @Success 201 {object} models.Student
// @Failure 400 {object} ErrorResponse
// @Router /admin/students [post]
// @Security BearerAuth
func (h *StudentHandler) CreateStudent(c *gin.Context) {
	var req CreateStudentRequest
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

	student := &models.Student{
		UserID:          req.UserID,
		AdmissionNumber: req.AdmissionNumber,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		DateOfBirth:     dateOfBirth,
		Gender:          req.Gender,
		Address:         req.Address,
		Phone:           req.Phone,
		ParentName:      req.ParentName,
		ParentPhone:     req.ParentPhone,
		ClassID:         req.ClassID,
		SectionID:       req.SectionID,
		Status:          "active",
	}

	if err := h.studentRepo.Create(student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, student)
}

// UpdateStudent godoc
// @Summary Update student
// @Description Update an existing student record
// @Tags Admin - Students
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Param student body UpdateStudentRequest true "Student data"
// @Success 200 {object} models.Student
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/students/{id} [put]
// @Security BearerAuth
func (h *StudentHandler) UpdateStudent(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	var req UpdateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	student, err := h.studentRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	// Update fields
	if req.FirstName != "" {
		student.FirstName = req.FirstName
	}
	if req.LastName != "" {
		student.LastName = req.LastName
	}
	if req.DateOfBirth != "" {
		dateOfBirth, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
			return
		}
		student.DateOfBirth = dateOfBirth
	}
	if req.Gender != "" {
		student.Gender = req.Gender
	}
	if req.Address != "" {
		student.Address = req.Address
	}
	if req.Phone != "" {
		student.Phone = req.Phone
	}
	if req.ParentName != "" {
		student.ParentName = req.ParentName
	}
	if req.ParentPhone != "" {
		student.ParentPhone = req.ParentPhone
	}
	if req.ClassID != 0 {
		student.ClassID = req.ClassID
	}
	if req.SectionID != 0 {
		student.SectionID = req.SectionID
	}
	if req.Status != "" {
		student.Status = req.Status
	}

	if err := h.studentRepo.Update(student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, student)
}

// DeleteStudent godoc
// @Summary Delete student
// @Description Delete a student record
// @Tags Admin - Students
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Success 200 {object} SuccessResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/students/{id} [delete]
// @Security BearerAuth
func (h *StudentHandler) DeleteStudent(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	if err := h.studentRepo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student deleted successfully"})
}

// Request Types
type CreateStudentRequest struct {
	UserID          uint   `json:"user_id" binding:"required"`
	AdmissionNumber string `json:"admission_number" binding:"required"`
	FirstName       string `json:"first_name" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
	DateOfBirth     string `json:"date_of_birth" binding:"required"`
	Gender          string `json:"gender"`
	Address         string `json:"address"`
	Phone           string `json:"phone"`
	ParentName      string `json:"parent_name"`
	ParentPhone     string `json:"parent_phone"`
	ClassID         uint   `json:"class_id" binding:"required"`
	SectionID       uint   `json:"section_id" binding:"required"`
}

type UpdateStudentRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth"`
	Gender      string `json:"gender"`
	Address     string `json:"address"`
	Phone       string `json:"phone"`
	ParentName  string `json:"parent_name"`
	ParentPhone string `json:"parent_phone"`
	ClassID     uint   `json:"class_id"`
	SectionID   uint   `json:"section_id"`
	Status      string `json:"status"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

