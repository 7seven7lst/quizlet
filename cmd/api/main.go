package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "quizlet/docs"
	"quizlet/internal/handlers"
	"quizlet/internal/repository"
	"quizlet/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"quizlet/internal/auth"
)

// @title           Quizlet API
// @version         1.0
// @description     API for managing quizzes, quiz suites, and users
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	quizRepo := repository.NewQuizRepository(db)
	quizSuiteRepo := repository.NewQuizSuiteRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	quizAttemptRepo := repository.NewQuizAttemptRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo, refreshTokenRepo)
	quizService := service.NewQuizService(quizRepo)
	quizSuiteService := service.NewQuizSuiteService(quizSuiteRepo, quizRepo)
	quizAttemptService := service.NewQuizAttemptService(quizAttemptRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	quizHandler := handlers.NewQuizHandler(quizService)
	quizSuiteHandler := handlers.NewQuizSuiteHandler(quizSuiteService)
	quizAttemptHandler := handlers.NewQuizAttemptHandler(quizAttemptService)

	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}))

	// Debug middleware to log all requests
	r.Use(func(c *gin.Context) {
		log.Printf("Request: %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
	})

	// Swagger setup
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API routes
	api := r.Group("/api")
	{
		// Public user routes (no auth required)
		api.POST("/users/login", userHandler.Login)
		api.POST("/users/refresh", userHandler.RefreshToken)
		api.POST("/users", userHandler.CreateUser)

		// Protected routes
		protected := api.Group("")
		protected.Use(auth.AuthMiddleware())
		{
			// Protected user routes
			protected.GET("/users/me", userHandler.GetCurrentUser)
			protected.GET("/users/:id", userHandler.GetUser)
			protected.PUT("/users/:id", userHandler.UpdateUser)
			protected.DELETE("/users/:id", userHandler.DeleteUser)
			protected.POST("/users/logout", userHandler.Logout)

			// Quiz routes
			protected.POST("/quizzes", quizHandler.CreateQuiz)
			protected.GET("/quizzes/:id", quizHandler.GetQuiz)
			protected.PUT("/quizzes/:id", quizHandler.UpdateQuiz)
			protected.DELETE("/quizzes/:id", quizHandler.DeleteQuiz)
			protected.POST("/quizzes/:id/selections", quizHandler.AddSelection)
			protected.DELETE("/quizzes/:id/selections/:selectionId", quizHandler.RemoveSelection)
			protected.GET("/quizzes/user", quizHandler.GetQuizzes)

			// Quiz Suite routes
			protected.POST("/quiz-suites", quizSuiteHandler.CreateQuizSuite)
			protected.GET("/quiz-suites", quizSuiteHandler.GetQuizSuites)
			protected.GET("/quiz-suites/:id", quizSuiteHandler.GetQuizSuite)
			protected.PUT("/quiz-suites/:id", quizSuiteHandler.UpdateQuizSuite)
			protected.DELETE("/quiz-suites/:id", quizSuiteHandler.DeleteQuizSuite)
			protected.POST("/quiz-suites/:id/quizzes/:quizId", quizSuiteHandler.AddQuizToSuite)
			protected.DELETE("/quiz-suites/:id/quizzes/:quizId", quizSuiteHandler.RemoveQuizFromSuite)

			// Quiz Attempt routes
			protected.GET("/quiz-suites/:id/attempts", quizAttemptHandler.ListQuizAttempts)
			protected.POST("/quiz-suites/:id/attempts", quizAttemptHandler.CreateQuizAttempt)
			protected.GET("/quiz-suites/:id/attempts/:attemptId", quizAttemptHandler.GetQuizAttempt)
			protected.PUT("/quiz-suites/:id/attempts/:attemptId", quizAttemptHandler.UpdateQuizAttempt)
			protected.DELETE("/quiz-suites/:id/attempts/:attemptId", quizAttemptHandler.DeleteQuizAttempt)
		}
	}

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
} 