package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/yourusername/quizlet/docs" // This will be replaced with your actual module name
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
	r := gin.Default()

	// Swagger documentation route
	r.GET("/api-docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := r.Group("/api")
	{
		// User routes
		users := api.Group("/users")
		{
			users.POST("/register", registerUser)
			users.POST("/login", loginUser)
			users.GET("/profile", authMiddleware(), getUserProfile)
			users.PUT("/profile", authMiddleware(), updateUserProfile)
		}

		// Quiz routes
		quizzes := api.Group("/quizzes")
		{
			quizzes.POST("/", authMiddleware(), createQuiz)
			quizzes.GET("/", authMiddleware(), getQuizzes)
			quizzes.GET("/:id", authMiddleware(), getQuiz)
			quizzes.PUT("/:id", authMiddleware(), updateQuiz)
			quizzes.DELETE("/:id", authMiddleware(), deleteQuiz)
		}

		// Quiz Suite routes
		quizSuites := api.Group("/quiz-suites")
		{
			quizSuites.POST("/", authMiddleware(), createQuizSuite)
			quizSuites.GET("/", authMiddleware(), getQuizSuites)
			quizSuites.GET("/:id", authMiddleware(), getQuizSuite)
			quizSuites.PUT("/:id", authMiddleware(), updateQuizSuite)
			quizSuites.DELETE("/:id", authMiddleware(), deleteQuizSuite)
		}
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
} 