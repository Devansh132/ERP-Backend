package main

import (
	"log"
	"school-erp-backend/config"
	"school-erp-backend/docs"
	"school-erp-backend/internal/handlers"
	"school-erp-backend/internal/middleware"
	"school-erp-backend/internal/models"
	"school-erp-backend/pkg/database"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           School ERP System API
// @version         1.0
// @description     Comprehensive School ERP System API with Admin, Teacher, and Student roles
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@schoolerp.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api
// @schemes   http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token. Example: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

func main() {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate database
	if err := database.DB.AutoMigrate(
		&models.User{},
		&models.Session{},
		&models.Student{},
		&models.Teacher{},
		&models.Class{},
		&models.Section{},
		&models.ClassSection{},
		&models.Attendance{},
		&models.Subject{},
		&models.Exam{},
		&models.Mark{},
		&models.Assignment{},
		&models.AssignmentSubmission{},
		&models.Timetable{},
		&models.Notice{},
		&models.CalendarEvent{},
		&models.LeaveRequest{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize handlers
	authHandler := handlers.NewAuthHandler()
	userHandler := handlers.NewUserHandler()
	studentHandler := handlers.NewStudentHandler()
	teacherHandler := handlers.NewTeacherHandler()
	classHandler := handlers.NewClassHandler()
	sectionHandler := handlers.NewSectionHandler()
	subjectHandler := handlers.NewSubjectHandler()
	attendanceHandler := handlers.NewAttendanceHandler()

	// Setup router
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := router.Group("/api")
	{
		// Authentication routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), authHandler.Register)
		}

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
		{
			// Users
			users := admin.Group("/users")
			{
				users.GET("", userHandler.GetUsers)
				users.GET("/:id", userHandler.GetUser)
				users.POST("", userHandler.CreateUser)
				users.PUT("/:id", userHandler.UpdateUser)
				users.DELETE("/:id", userHandler.DeleteUser)
			}

			// Students
			students := admin.Group("/students")
			{
				students.GET("", studentHandler.GetStudents)
				students.GET("/:id", studentHandler.GetStudent)
				students.POST("", studentHandler.CreateStudent)
				students.PUT("/:id", studentHandler.UpdateStudent)
				students.DELETE("/:id", studentHandler.DeleteStudent)
			}

			// Teachers
			teachers := admin.Group("/teachers")
			{
				teachers.GET("", teacherHandler.GetTeachers)
				teachers.GET("/:id", teacherHandler.GetTeacher)
				teachers.POST("", teacherHandler.CreateTeacher)
				teachers.PUT("/:id", teacherHandler.UpdateTeacher)
				teachers.DELETE("/:id", teacherHandler.DeleteTeacher)
			}

			// Classes
			classes := admin.Group("/classes")
			{
				classes.GET("", classHandler.GetClasses)
				classes.GET("/:id", classHandler.GetClass)
				classes.POST("", classHandler.CreateClass)
				classes.PUT("/:id", classHandler.UpdateClass)
				classes.DELETE("/:id", classHandler.DeleteClass)
			}

			// Sections
			sections := admin.Group("/sections")
			{
				sections.GET("", sectionHandler.GetSections)
				sections.GET("/:id", sectionHandler.GetSection)
				sections.POST("", sectionHandler.CreateSection)
				sections.PUT("/:id", sectionHandler.UpdateSection)
				sections.DELETE("/:id", sectionHandler.DeleteSection)
				sections.POST("/assign", sectionHandler.AssignSectionToClass)
			}

			// Subjects
			subjects := admin.Group("/subjects")
			{
				subjects.GET("", subjectHandler.GetSubjects)
				subjects.GET("/:id", subjectHandler.GetSubject)
				subjects.POST("", subjectHandler.CreateSubject)
				subjects.PUT("/:id", subjectHandler.UpdateSubject)
				subjects.DELETE("/:id", subjectHandler.DeleteSubject)
			}

			// Attendance
			attendance := admin.Group("/attendance")
			{
				attendance.POST("/mark", attendanceHandler.MarkAttendance)
				attendance.GET("/class/:class_id", attendanceHandler.GetAttendanceByClass)
				attendance.GET("/student/:student_id", attendanceHandler.GetAttendanceByStudent)
				attendance.GET("/statistics", attendanceHandler.GetAttendanceStatistics)
				attendance.GET("/reports", attendanceHandler.GetAttendanceReports)
				attendance.PUT("/:id", attendanceHandler.UpdateAttendance)
			}
		}

		// Teacher routes
		teacher := api.Group("/teacher")
		teacher.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("teacher"))
		{
			// Add teacher routes here
		}

		// Student routes
		student := api.Group("/student")
		student.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("student"))
		{
			// Add student routes here
		}
	}

	// Swagger documentation
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Host = "localhost:8080"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	port := ":" + config.AppConfig.ServerPort
	log.Printf("Server starting on port %s", port)
	log.Printf("Swagger documentation available at http://localhost%s/swagger/index.html", port)
	
	if err := router.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

