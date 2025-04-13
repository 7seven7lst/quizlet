package handlers

import (
	"net/http"
	"quizlet/internal/models"
	"quizlet/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QuizSuiteHandler struct {
	quizSuiteService service.QuizSuiteService
}

func NewQuizSuiteHandler(quizSuiteService service.QuizSuiteService) *QuizSuiteHandler {
	return &QuizSuiteHandler{
		quizSuiteService: quizSuiteService,
	}
}

func (h *QuizSuiteHandler) CreateQuizSuite(c *gin.Context) {
	var quizSuite models.QuizSuite
	if err := c.ShouldBindJSON(&quizSuite); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (assuming you have middleware that sets this)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	quizSuite.CreatedByID = userID.(uint)

	if err := h.quizSuiteService.CreateQuizSuite(&quizSuite); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, quizSuite)
}

func (h *QuizSuiteHandler) GetQuizSuite(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quiz suite id"})
		return
	}

	quizSuite, err := h.quizSuiteService.GetQuizSuiteByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "quiz suite not found"})
		return
	}

	c.JSON(http.StatusOK, quizSuite)
}

func (h *QuizSuiteHandler) GetUserQuizSuites(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	quizSuites, err := h.quizSuiteService.GetQuizSuitesByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quizSuites)
}

func (h *QuizSuiteHandler) UpdateQuizSuite(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quiz suite id"})
		return
	}

	var quizSuite models.QuizSuite
	if err := c.ShouldBindJSON(&quizSuite); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quizSuite.ID = uint(id)
	if err := h.quizSuiteService.UpdateQuizSuite(&quizSuite); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quizSuite)
}

func (h *QuizSuiteHandler) DeleteQuizSuite(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quiz suite id"})
		return
	}

	if err := h.quizSuiteService.DeleteQuizSuite(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "quiz suite deleted successfully"})
}

func (h *QuizSuiteHandler) AddQuizToSuite(c *gin.Context) {
	suiteID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quiz suite id"})
		return
	}

	quizID, err := strconv.ParseUint(c.Param("quizId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quiz id"})
		return
	}

	if err := h.quizSuiteService.AddQuizToSuite(uint(suiteID), uint(quizID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "quiz added to suite successfully"})
}

func (h *QuizSuiteHandler) RemoveQuizFromSuite(c *gin.Context) {
	suiteID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quiz suite id"})
		return
	}

	quizID, err := strconv.ParseUint(c.Param("quizId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quiz id"})
		return
	}

	if err := h.quizSuiteService.RemoveQuizFromSuite(uint(suiteID), uint(quizID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "quiz removed from suite successfully"})
} 