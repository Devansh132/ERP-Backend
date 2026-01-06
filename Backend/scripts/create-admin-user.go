package main

import (
	"fmt"
	"log"
	"os"
	"school-erp-backend/config"
	"school-erp-backend/internal/models"
	"school-erp-backend/pkg/database"
	"school-erp-backend/pkg/utils"
)

func main() {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Get email and password from command line or use defaults
	email := "admin@school.com"
	password := "admin123"

	if len(os.Args) > 1 {
		email = os.Args[1]
	}
	if len(os.Args) > 2 {
		password = os.Args[2]
	}

	// Check if user already exists
	var existingUser models.User
	result := database.DB.Where("email = ?", email).First(&existingUser)
	if result.Error == nil {
		fmt.Printf("User with email '%s' already exists!\n", email)
		fmt.Println("To update password, delete the user first or use a different email.")
		return
	}

	// Hash password
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	// Create user
	user := models.User{
		Email:        email,
		PasswordHash: passwordHash,
		Role:         "admin",
		Status:       "active",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		log.Fatal("Failed to create user:", err)
	}

	fmt.Println("✅ Admin user created successfully!")
	fmt.Printf("Email: %s\n", email)
	fmt.Printf("Password: %s\n", password)
	fmt.Printf("Role: admin\n")
	fmt.Println("\nYou can now login with these credentials in Swagger.")
}



