package handlers

import (
	"net/http"
	"strconv"
	"school-erp-backend/internal/models"
	"school-erp-backend/internal/repository"
	"school-erp-backend/pkg/database"
	"github.com/gin-gonic/gin"
)

type SectionHandler struct {
	sectionRepo *repository.SectionRepository
}

func NewSectionHandler() *SectionHandler {
	return &SectionHandler{
		sectionRepo: repository.NewSectionRepository(database.DB),
	}
}

// GetSections godoc
// @Summary Get all sections
// @Description Get list of all sections, optionally filtered by class_id
// @Tags Admin - Sections
// @Accept json
// @Produce json
// @Param class_id query int false "Filter by class ID"
// @Success 200 {array} models.Section
// @Failure 500 {object} ErrorResponse
// @Router /admin/sections [get]
// @Security BearerAuth
func (h *SectionHandler) GetSections(c *gin.Context) {
	classID := c.Query("class_id")
	
	var sections []models.Section
	var err error
	
	if classID != "" {
		classIDUint, _ := strconv.ParseUint(classID, 10, 32)
		sections, err = h.sectionRepo.FindByClassID(uint(classIDUint))
	} else {
		sections, err = h.sectionRepo.FindAll()
	}
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sections)
}

// GetSection godoc
// @Summary Get section by ID
// @Description Get a specific section by ID
// @Tags Admin - Sections
// @Accept json
// @Produce json
// @Param id path int true "Section ID"
// @Success 200 {object} models.Section
// @Failure 404 {object} ErrorResponse
// @Router /admin/sections/{id} [get]
// @Security BearerAuth
func (h *SectionHandler) GetSection(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	section, err := h.sectionRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Section not found"})
		return
	}

	c.JSON(http.StatusOK, section)
}

// CreateSection godoc
// @Summary Create a new section
// @Description Create a new section record
// @Tags Admin - Sections
// @Accept json
// @Produce json
// @Param section body CreateSectionRequest true "Section data"
// @Success 201 {object} models.Section
// @Failure 400 {object} ErrorResponse
// @Router /admin/sections [post]
// @Security BearerAuth
func (h *SectionHandler) CreateSection(c *gin.Context) {
	var req CreateSectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	section := &models.Section{
		ClassID:  req.ClassID,
		Name:     req.Name,
		Capacity: req.Capacity,
		Status:   "active",
	}

	if err := h.sectionRepo.Create(section); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, section)
}

// UpdateSection godoc
// @Summary Update section
// @Description Update an existing section record
// @Tags Admin - Sections
// @Accept json
// @Produce json
// @Param id path int true "Section ID"
// @Param section body UpdateSectionRequest true "Section data"
// @Success 200 {object} models.Section
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/sections/{id} [put]
// @Security BearerAuth
func (h *SectionHandler) UpdateSection(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	var req UpdateSectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	section, err := h.sectionRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Section not found"})
		return
	}

	// Update fields
	if req.Name != "" {
		section.Name = req.Name
	}
	if req.ClassID != 0 {
		section.ClassID = req.ClassID
	}
	if req.Capacity != 0 {
		section.Capacity = req.Capacity
	}
	if req.Status != "" {
		section.Status = req.Status
	}

	if err := h.sectionRepo.Update(section); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, section)
}

// DeleteSection godoc
// @Summary Delete section
// @Description Delete a section record
// @Tags Admin - Sections
// @Accept json
// @Produce json
// @Param id path int true "Section ID"
// @Success 200 {object} SuccessResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/sections/{id} [delete]
// @Security BearerAuth
func (h *SectionHandler) DeleteSection(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	if err := h.sectionRepo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Section not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Section deleted successfully"})
}

// AssignSectionToClass godoc
// @Summary Assign section to class
// @Description Create a ClassSection relationship for academic year
// @Tags Admin - Sections
// @Accept json
// @Produce json
// @Param assignment body AssignSectionRequest true "Assignment data"
// @Success 201 {object} models.ClassSection
// @Failure 400 {object} ErrorResponse
// @Router /admin/sections/assign [post]
// @Security BearerAuth
func (h *SectionHandler) AssignSectionToClass(c *gin.Context) {
	var req AssignSectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	classSection := &models.ClassSection{
		ClassID:      req.ClassID,
		SectionID:    req.SectionID,
		AcademicYear: req.AcademicYear,
		Status:       "active",
	}

	if err := database.DB.Create(classSection).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Preload relationships
	database.DB.Preload("Class").Preload("Section").First(classSection, classSection.ID)

	c.JSON(http.StatusCreated, classSection)
}

// Request Types
type CreateSectionRequest struct {
	ClassID  uint   `json:"class_id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Capacity int    `json:"capacity"`
}

type UpdateSectionRequest struct {
	ClassID  uint   `json:"class_id"`
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
	Status   string `json:"status"`
}

type AssignSectionRequest struct {
	ClassID      uint   `json:"class_id" binding:"required"`
	SectionID    uint   `json:"section_id" binding:"required"`
	AcademicYear string `json:"academic_year" binding:"required"`
}

