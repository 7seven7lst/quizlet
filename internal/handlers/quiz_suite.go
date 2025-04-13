package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"quizlet/internal/models"
	"gorm.io/gorm"
)

// @Summary      Create a new quiz suite
// @Description  Create a new quiz suite with the provided details
// @Tags         quiz-suites
// @Accept       json
// @Produce      json
// @Param        quiz_suite  body      models.QuizSuite  true  "Quiz Suite object"
// @Success      201        {object}  models.QuizSuite
// @Failure      400        {object}  ErrorResponse
// @Failure      401        {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /quiz-suites [post]
func CreateQuizSuite(c *gin.Context) {
	var quizSuite models.QuizSuite
	if err := c.ShouldBindJSON(&quizSuite); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}
	quizSuite.CreatedByID = userID.(uint)

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Create(&quizSuite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, quizSuite)
}

// @Summary      Get all quiz suites
// @Description  Retrieve all quiz suites for the authenticated user
// @Tags         quiz-suites
// @Produce      json
// @Success      200  {array}   models.QuizSuite
// @Failure      401  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /quiz-suites [get]
func GetQuizSuites(c *gin.Context) {
	var quizSuites []models.QuizSuite
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("created_by_id = ?", userID).Find(&quizSuites).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, quizSuites)
}

// @Summary      Get a quiz suite by ID
// @Description  Retrieve a specific quiz suite by its ID
// @Tags         quiz-suites
// @Produce      json
// @Param        id   path      integer  true  "Quiz Suite ID"
// @Success      200  {object}  models.QuizSuite
// @Failure      401  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /quiz-suites/{id} [get]
func GetQuizSuite(c *gin.Context) {
	var quizSuite models.QuizSuite
	id := c.Param("id")

	db := c.MustGet("db").(*gorm.DB)
	if err := db.First(&quizSuite, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "quiz suite not found"})
		return
	}

	c.JSON(http.StatusOK, quizSuite)
}

// @Summary      Update a quiz suite
// @Description  Update an existing quiz suite with the provided details
// @Tags         quiz-suites
// @Accept       json
// @Produce      json
// @Param        id          path      integer        true  "Quiz Suite ID"
// @Param        quiz_suite  body      models.QuizSuite  true  "Quiz Suite object"
// @Success      200        {object}  models.QuizSuite
// @Failure      400        {object}  ErrorResponse
// @Failure      401        {object}  ErrorResponse
// @Failure      404        {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /quiz-suites/{id} [put]
func UpdateQuizSuite(c *gin.Context) {
	var quizSuite models.QuizSuite
	id := c.Param("id")

	db := c.MustGet("db").(*gorm.DB)
	if err := db.First(&quizSuite, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "quiz suite not found"})
		return
	}

	if err := c.ShouldBindJSON(&quizSuite); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := db.Save(&quizSuite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, quizSuite)
}

// @Summary      Delete a quiz suite
// @Description  Delete an existing quiz suite by its ID
// @Tags         quiz-suites
// @Produce      json
// @Param        id   path      integer  true  "Quiz Suite ID"
// @Success      200  {object}  SuccessResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /quiz-suites/{id} [delete]
func DeleteQuizSuite(c *gin.Context) {
	var quizSuite models.QuizSuite
	id := c.Param("id")

	db := c.MustGet("db").(*gorm.DB)
	if err := db.First(&quizSuite, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "quiz suite not found"})
		return
	}

	if err := db.Delete(&quizSuite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "quiz suite deleted successfully"})
} 