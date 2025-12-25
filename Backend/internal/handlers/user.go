package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"school-erp-backend/internal/models"
	"school-erp-backend/internal/repository"
	"school-erp-backend/pkg/database"
	"school-erp-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userRepo *repository.UserRepository
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userRepo: repository.NewUserRepository(database.DB),
	}
}

// GetUsers godoc
// @Summary Get all users
// @Description Get list of all users with optional filters
// @Tags Admin - Users
// @Accept json
// @Produce json
// @Param role query string false "Filter by role (admin, teacher, student)"
// @Param status query string false "Filter by status (active, inactive)"
// @Success 200 {array} UserResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/users [get]
// @Security BearerAuth
func (h *UserHandler) GetUsers(c *gin.Context) {
	role := c.Query("role")
	status := c.Query("status")

	var users []models.User
	query := database.DB

	if role != "" {
		query = query.Where("role = ?", role)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	var response []UserResponse
	for _, user := range users {
		response = append(response, UserResponse{
			ID:     user.ID,
			Email:  user.Email,
			Role:   user.Role,
			Status: user.Status,
		})
	}

	c.JSON(http.StatusOK, response)
}

// GetUser godoc
// @Summary Get user by ID
// @Description Get a single user by ID
// @Tags Admin - Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/users/{id} [get]
// @Security BearerAuth
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:     user.ID,
		Email:  user.Email,
		Role:   user.Role,
		Status: user.Status,
	})
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user account
// @Tags Admin - Users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User data"
// @Success 201 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Router /admin/users [post]
// @Security BearerAuth
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email already exists (case-insensitive, excluding soft-deleted)
	email := strings.ToLower(strings.TrimSpace(req.Email))
	existingUser, err := h.userRepo.FindByEmail(email)
	if err == nil && existingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Validate and normalize role
	role := strings.ToLower(strings.TrimSpace(req.Role))
	if role != "admin" && role != "teacher" && role != "student" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role. Must be admin, teacher, or student"})
		return
	}

	user := &models.User{
		Email:        email, // Use normalized email
		PasswordHash: passwordHash,
		Role:         role,
		Status:       req.Status,
	}

	if user.Status == "" {
		user.Status = "active"
	}

	if err := h.userRepo.Create(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, UserResponse{
		ID:     user.ID,
		Email:  user.Email,
		Role:   user.Role,
		Status: user.Status,
	})
}

// UpdateUser godoc
// @Summary Update user
// @Description Update an existing user
// @Tags Admin - Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body UpdateUserRequest true "User data"
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/users/{id} [put]
// @Security BearerAuth
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update email if provided and check for duplicates (case-insensitive)
	if req.Email != "" {
		email := strings.ToLower(strings.TrimSpace(req.Email))
		if email != strings.ToLower(user.Email) {
			existingUser, err := h.userRepo.FindByEmail(email)
			if err == nil && existingUser != nil && existingUser.ID != user.ID {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
				return
			}
			user.Email = email
		}
	}

	// Update password if provided
	if req.Password != "" {
		passwordHash, err := utils.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.PasswordHash = passwordHash
	}

	// Update role if provided
	if req.Role != "" {
		role := strings.ToLower(strings.TrimSpace(req.Role))
		if role != "admin" && role != "teacher" && role != "student" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role. Must be admin, teacher, or student"})
			return
		}
		user.Role = role
	}

	// Update status if provided
	if req.Status != "" {
		user.Status = req.Status
	}

	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:     user.ID,
		Email:  user.Email,
		Role:   user.Role,
		Status: user.Status,
	})
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user by ID
// @Tags Admin - Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/users/{id} [delete]
// @Security BearerAuth
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if user exists
	_, err = h.userRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := h.userRepo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// Request Types
type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required"`
	Status   string `json:"status"`
}

type UpdateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

