package main

import (
	"fmt"
	"log"
	"os"

	"quizlet/internal/handlers"
	"quizlet/internal/repository"
	"quizlet/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

	// Initialize services
	userService := service.NewUserService(userRepo)
	quizService := service.NewQuizService(quizRepo)
	quizSuiteService := service.NewQuizSuiteService(quizSuiteRepo, quizRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	quizHandler := handlers.NewQuizHandler(quizService)
	quizSuiteHandler := handlers.NewQuizSuiteHandler(quizSuiteService)

	// Initialize router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// User routes
	users := r.Group("/api/users")
	{
		users.POST("/", userHandler.CreateUser)
		users.GET("/:id", userHandler.GetUser)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
		users.POST("/login", userHandler.Login)
	}

	// Quiz routes
	quizzes := r.Group("/api/quizzes")
	{
		quizzes.POST("/", quizHandler.CreateQuiz)
		quizzes.GET("/:id", quizHandler.GetQuiz)
		quizzes.PUT("/:id", quizHandler.UpdateQuiz)
		quizzes.DELETE("/:id", quizHandler.DeleteQuiz)
		quizzes.GET("/user", quizHandler.GetUserQuizzes)
		quizzes.POST("/:id/selections", quizHandler.AddSelection)
		quizzes.DELETE("/:id/selections/:selectionId", quizHandler.RemoveSelection)
	}

	// Quiz Suite routes
	quizSuites := r.Group("/api/quiz-suites")
	{
		quizSuites.POST("/", quizSuiteHandler.CreateQuizSuite)
		quizSuites.GET("/:id", quizSuiteHandler.GetQuizSuite)
		quizSuites.PUT("/:id", quizSuiteHandler.UpdateQuizSuite)
		quizSuites.DELETE("/:id", quizSuiteHandler.DeleteQuizSuite)
		quizSuites.GET("/user", quizSuiteHandler.GetUserQuizSuites)
		quizSuites.POST("/:id/quizzes/:quizId", quizSuiteHandler.AddQuizToSuite)
		quizSuites.DELETE("/:id/quizzes/:quizId", quizSuiteHandler.RemoveQuizFromSuite)
	}

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
} 