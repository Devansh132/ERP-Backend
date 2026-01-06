package handlers

import (
	"net/http"
	"strconv"
	"school-erp-backend/internal/models"
	"school-erp-backend/internal/repository"
	"school-erp-backend/pkg/database"
	"github.com/gin-gonic/gin"
)

type SubjectHandler struct {
	subjectRepo *repository.SubjectRepository
}

func NewSubjectHandler() *SubjectHandler {
	return &SubjectHandler{
		subjectRepo: repository.NewSubjectRepository(database.DB),
	}
}

// GetSubjects godoc
// @Summary Get all subjects
// @Description Get list of all subjects
// @Tags Admin - Subjects
// @Accept json
// @Produce json
// @Success 200 {array} models.Subject
// @Failure 500 {object} ErrorResponse
// @Router /admin/subjects [get]
// @Security BearerAuth
func (h *SubjectHandler) GetSubjects(c *gin.Context) {
	subjects, err := h.subjectRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subjects)
}

// GetSubject godoc
// @Summary Get subject by ID
// @Description Get a specific subject by ID
// @Tags Admin - Subjects
// @Accept json
// @Produce json
// @Param id path int true "Subject ID"
// @Success 200 {object} models.Subject
// @Failure 404 {object} ErrorResponse
// @Router /admin/subjects/{id} [get]
// @Security BearerAuth
func (h *SubjectHandler) GetSubject(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	subject, err := h.subjectRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subject not found"})
		return
	}

	c.JSON(http.StatusOK, subject)
}

// CreateSubject godoc
// @Summary Create a new subject
// @Description Create a new subject record
// @Tags Admin - Subjects
// @Accept json
// @Produce json
// @Param subject body CreateSubjectRequest true "Subject data"
// @Success 201 {object} models.Subject
// @Failure 400 {object} ErrorResponse
// @Router /admin/subjects [post]
// @Security BearerAuth
func (h *SubjectHandler) CreateSubject(c *gin.Context) {
	var req CreateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subject := &models.Subject{
		Name:   req.Name,
		Code:   req.Code,
		Status: "active",
	}

	if err := h.subjectRepo.Create(subject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, subject)
}

// UpdateSubject godoc
// @Summary Update subject
// @Description Update an existing subject record
// @Tags Admin - Subjects
// @Accept json
// @Produce json
// @Param id path int true "Subject ID"
// @Param subject body UpdateSubjectRequest true "Subject data"
// @Success 200 {object} models.Subject
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/subjects/{id} [put]
// @Security BearerAuth
func (h *SubjectHandler) UpdateSubject(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	var req UpdateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subject, err := h.subjectRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subject not found"})
		return
	}

	// Update fields
	if req.Name != "" {
		subject.Name = req.Name
	}
	if req.Code != "" {
		subject.Code = req.Code
	}
	if req.Status != "" {
		subject.Status = req.Status
	}

	if err := h.subjectRepo.Update(subject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subject)
}

// DeleteSubject godoc
// @Summary Delete subject
// @Description Delete a subject record
// @Tags Admin - Subjects
// @Accept json
// @Produce json
// @Param id path int true "Subject ID"
// @Success 200 {object} SuccessResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/subjects/{id} [delete]
// @Security BearerAuth
func (h *SubjectHandler) DeleteSubject(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	if err := h.subjectRepo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subject not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subject deleted successfully"})
}

// Request Types
type CreateSubjectRequest struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code" binding:"required"`
}

type UpdateSubjectRequest struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	Status string `json:"status"`
}



