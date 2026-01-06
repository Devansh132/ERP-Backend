//go:build script
// +build script

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

	// Get email and password from command line
	email := "admin@school.com"
	password := "password123"

	if len(os.Args) > 1 {
		email = os.Args[1]
	}
	if len(os.Args) > 2 {
		password = os.Args[2]
	}

	// Find user
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		log.Fatalf("User with email '%s' not found: %v", email, err)
	}

	// Hash new password
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	// Update password
	user.PasswordHash = passwordHash
	if err := database.DB.Save(&user).Error; err != nil {
		log.Fatal("Failed to update password:", err)
	}

	fmt.Println("✅ Password updated successfully!")
	fmt.Printf("Email: %s\n", email)
	fmt.Printf("New Password: %s\n", password)
	fmt.Println("\nYou can now login with these credentials.")
}
