package handlers

import (
	"net/http"
	"strconv"
	"school-erp-backend/internal/models"
	"school-erp-backend/internal/repository"
	"school-erp-backend/pkg/database"
	"github.com/gin-gonic/gin"
)

type ClassHandler struct {
	classRepo *repository.ClassRepository
}

func NewClassHandler() *ClassHandler {
	return &ClassHandler{
		classRepo: repository.NewClassRepository(database.DB),
	}
}

// GetClasses godoc
// @Summary Get all classes
// @Description Get list of all classes with their sections
// @Tags Admin - Classes
// @Accept json
// @Produce json
// @Success 200 {array} models.Class
// @Failure 500 {object} ErrorResponse
// @Router /admin/classes [get]
// @Security BearerAuth
func (h *ClassHandler) GetClasses(c *gin.Context) {
	classes, err := h.classRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, classes)
}

// GetClass godoc
// @Summary Get class by ID
// @Description Get a specific class by ID with its sections
// @Tags Admin - Classes
// @Accept json
// @Produce json
// @Param id path int true "Class ID"
// @Success 200 {object} models.Class
// @Failure 404 {object} ErrorResponse
// @Router /admin/classes/{id} [get]
// @Security BearerAuth
func (h *ClassHandler) GetClass(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	class, err := h.classRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
		return
	}

	c.JSON(http.StatusOK, class)
}

// CreateClass godoc
// @Summary Create a new class
// @Description Create a new class record
// @Tags Admin - Classes
// @Accept json
// @Produce json
// @Param class body CreateClassRequest true "Class data"
// @Success 201 {object} models.Class
// @Failure 400 {object} ErrorResponse
// @Router /admin/classes [post]
// @Security BearerAuth
func (h *ClassHandler) CreateClass(c *gin.Context) {
	var req CreateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	class := &models.Class{
		Name:     req.Name,
		Level:    req.Level,
		Capacity: req.Capacity,
		Status:   "active",
	}

	if err := h.classRepo.Create(class); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, class)
}

// UpdateClass godoc
// @Summary Update class
// @Description Update an existing class record
// @Tags Admin - Classes
// @Accept json
// @Produce json
// @Param id path int true "Class ID"
// @Param class body UpdateClassRequest true "Class data"
// @Success 200 {object} models.Class
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/classes/{id} [put]
// @Security BearerAuth
func (h *ClassHandler) UpdateClass(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	var req UpdateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	class, err := h.classRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
		return
	}

	// Update fields
	if req.Name != "" {
		class.Name = req.Name
	}
	if req.Level != 0 {
		class.Level = req.Level
	}
	if req.Capacity != 0 {
		class.Capacity = req.Capacity
	}
	if req.Status != "" {
		class.Status = req.Status
	}

	if err := h.classRepo.Update(class); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, class)
}

// DeleteClass godoc
// @Summary Delete class
// @Description Delete a class record
// @Tags Admin - Classes
// @Accept json
// @Produce json
// @Param id path int true "Class ID"
// @Success 200 {object} SuccessResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/classes/{id} [delete]
// @Security BearerAuth
func (h *ClassHandler) DeleteClass(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	if err := h.classRepo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Class deleted successfully"})
}

// Request Types
type CreateClassRequest struct {
	Name     string `json:"name" binding:"required"`
	Level    int    `json:"level" binding:"required"`
	Capacity int    `json:"capacity"`
}

type UpdateClassRequest struct {
	Name     string `json:"name"`
	Level    int    `json:"level"`
	Capacity int    `json:"capacity"`
	Status   string `json:"status"`
}



