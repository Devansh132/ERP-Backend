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

type AttendanceHandler struct {
	attendanceRepo *repository.AttendanceRepository
	studentRepo    *repository.StudentRepository
}

func NewAttendanceHandler() *AttendanceHandler {
	return &AttendanceHandler{
		attendanceRepo: repository.NewAttendanceRepository(database.DB),
		studentRepo:    repository.NewStudentRepository(database.DB),
	}
}

// MarkAttendance godoc
// @Summary Mark attendance for a class/section
// @Description Mark attendance for multiple students in a class/section on a specific date
// @Tags Admin - Attendance
// @Accept json
// @Produce json
// @Param attendance body MarkAttendanceRequest true "Attendance data"
// @Success 201 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /admin/attendance/mark [post]
// @Security BearerAuth
func (h *AttendanceHandler) MarkAttendance(c *gin.Context) {
	var req MarkAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Get user ID from token (marked_by)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	markedBy := userID.(uint)

	// Get students for the class and section
	students, err := h.studentRepo.FindByClassAndSection(req.ClassID, req.SectionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to fetch students"})
		return
	}

	if len(students) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No students found for the selected class and section"})
		return
	}

	// Create attendance records
	var attendances []models.Attendance
	for _, student := range students {
		// Check if attendance already exists
		existing, _ := h.attendanceRepo.FindExisting(student.ID, date)
		if existing != nil {
			// Update existing record
			existing.Status = req.Attendance[student.ID]
			if err := h.attendanceRepo.Update(existing); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update attendance"})
				return
			}
			continue
		}

		// Create new record
		status := req.Attendance[student.ID]
		if status == "" {
			status = "absent" // Default to absent if not specified
		}

		attendance := models.Attendance{
			StudentID: student.ID,
			ClassID:   req.ClassID,
			SectionID: req.SectionID,
			Date:      date,
			Status:    status,
			MarkedBy:  markedBy,
		}
		attendances = append(attendances, attendance)
	}

	if len(attendances) > 0 {
		if err := h.attendanceRepo.CreateBatch(attendances); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create attendance records"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Attendance marked successfully"})
}

// GetAttendanceByClass godoc
// @Summary Get attendance by class
// @Description Get attendance records for a specific class and optional date
// @Tags Admin - Attendance
// @Accept json
// @Produce json
// @Param class_id path int true "Class ID"
// @Param section_id query int false "Section ID"
// @Param date query string false "Date (YYYY-MM-DD)"
// @Success 200 {array} models.Attendance
// @Failure 400 {object} ErrorResponse
// @Router /admin/attendance/class/{class_id} [get]
// @Security BearerAuth
func (h *AttendanceHandler) GetAttendanceByClass(c *gin.Context) {
	classID, _ := strconv.ParseUint(c.Param("class_id"), 10, 32)
	sectionIDStr := c.Query("section_id")
	dateStr := c.Query("date")

	var sectionID uint
	if sectionIDStr != "" {
		sectionIDUint, _ := strconv.ParseUint(sectionIDStr, 10, 32)
		sectionID = uint(sectionIDUint)
	}

	var date time.Time
	if dateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err == nil {
			date = parsedDate
		}
	}

	var attendances []models.Attendance
	var err error

	if sectionID > 0 && !date.IsZero() {
		attendances, err = h.attendanceRepo.FindByClassAndSection(uint(classID), sectionID, date)
	} else if sectionID > 0 {
		attendances, err = h.attendanceRepo.FindByClassAndSection(uint(classID), sectionID, time.Time{})
	} else if !date.IsZero() {
		attendances, err = h.attendanceRepo.FindByClassID(uint(classID), date)
	} else {
		attendances, err = h.attendanceRepo.FindByClassID(uint(classID), time.Time{})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, attendances)
}

// GetAttendanceByStudent godoc
// @Summary Get attendance by student
// @Description Get attendance records for a specific student
// @Tags Admin - Attendance
// @Accept json
// @Produce json
// @Param student_id path int true "Student ID"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {array} models.Attendance
// @Failure 400 {object} ErrorResponse
// @Router /admin/attendance/student/{student_id} [get]
// @Security BearerAuth
func (h *AttendanceHandler) GetAttendanceByStudent(c *gin.Context) {
	studentID, _ := strconv.ParseUint(c.Param("student_id"), 10, 32)
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	if startDateStr != "" {
		parsed, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			startDate = parsed
		}
	}
	if endDateStr != "" {
		parsed, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			endDate = parsed
		}
	}

	attendances, err := h.attendanceRepo.FindByStudentID(uint(studentID), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, attendances)
}

// GetAttendanceStatistics godoc
// @Summary Get attendance statistics
// @Description Get attendance statistics for a student or class
// @Tags Admin - Attendance
// @Accept json
// @Produce json
// @Param student_id query int false "Student ID"
// @Param class_id query int false "Class ID"
// @Param section_id query int false "Section ID"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Router /admin/attendance/statistics [get]
// @Security BearerAuth
func (h *AttendanceHandler) GetAttendanceStatistics(c *gin.Context) {
	studentIDStr := c.Query("student_id")
	classIDStr := c.Query("class_id")
	sectionIDStr := c.Query("section_id")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	if startDateStr != "" {
		parsed, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			startDate = parsed
		}
	}
	if endDateStr != "" {
		parsed, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			endDate = parsed
		}
	}

	var stats map[string]interface{}
	var err error

	if studentIDStr != "" {
		studentID, _ := strconv.ParseUint(studentIDStr, 10, 32)
		stats, err = h.attendanceRepo.GetStatistics(uint(studentID), startDate, endDate)
	} else if classIDStr != "" {
		classID, _ := strconv.ParseUint(classIDStr, 10, 32)
		var sectionID uint
		if sectionIDStr != "" {
			sectionIDUint, _ := strconv.ParseUint(sectionIDStr, 10, 32)
			sectionID = uint(sectionIDUint)
		}
		stats, err = h.attendanceRepo.GetClassStatistics(uint(classID), sectionID, startDate, endDate)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Either student_id or class_id is required"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// UpdateAttendance godoc
// @Summary Update attendance record
// @Description Update an existing attendance record
// @Tags Admin - Attendance
// @Accept json
// @Produce json
// @Param id path int true "Attendance ID"
// @Param attendance body UpdateAttendanceRequest true "Attendance data"
// @Success 200 {object} models.Attendance
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/attendance/{id} [put]
// @Security BearerAuth
func (h *AttendanceHandler) UpdateAttendance(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req UpdateAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	attendance, err := h.attendanceRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attendance record not found"})
		return
	}

	if req.Status != "" {
		attendance.Status = req.Status
	}

	if err := h.attendanceRepo.Update(attendance); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, attendance)
}

// GetAttendanceReports godoc
// @Summary Get attendance reports
// @Description Get attendance reports with filters
// @Tags Admin - Attendance
// @Accept json
// @Produce json
// @Param class_id query int false "Class ID"
// @Param section_id query int false "Section ID"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {array} models.Attendance
// @Failure 400 {object} ErrorResponse
// @Router /admin/attendance/reports [get]
// @Security BearerAuth
func (h *AttendanceHandler) GetAttendanceReports(c *gin.Context) {
	classIDStr := c.Query("class_id")
	sectionIDStr := c.Query("section_id")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	if startDateStr != "" {
		parsed, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			startDate = parsed
		}
	} else {
		// Default to start of current month
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	}

	if endDateStr != "" {
		parsed, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			endDate = parsed
		}
	} else {
		// Default to end of current month
		now := time.Now()
		endDate = time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 0, now.Location())
	}

	var classID, sectionID *uint
	if classIDStr != "" {
		id, _ := strconv.ParseUint(classIDStr, 10, 32)
		classIDVal := uint(id)
		classID = &classIDVal
	}
	if sectionIDStr != "" {
		id, _ := strconv.ParseUint(sectionIDStr, 10, 32)
		sectionIDVal := uint(id)
		sectionID = &sectionIDVal
	}

	attendances, err := h.attendanceRepo.FindByDateRange(startDate, endDate, classID, sectionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, attendances)
}

// Request Types
type MarkAttendanceRequest struct {
	ClassID    uint              `json:"class_id" binding:"required"`
	SectionID  uint              `json:"section_id" binding:"required"`
	Date       string            `json:"date" binding:"required"`
	Attendance map[uint]string   `json:"attendance" binding:"required"` // student_id -> status
}

type UpdateAttendanceRequest struct {
	Status string `json:"status" binding:"required,oneof=present absent late excused"`
}

