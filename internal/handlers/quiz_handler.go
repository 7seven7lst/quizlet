package handlers

import (
	"net/http"
	"quizlet/internal/models"
	"quizlet/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QuizHandler struct {
	quizService service.QuizService
}

func NewQuizHandler(quizService service.QuizService) *QuizHandler {
	return &QuizHandler{
		quizService: quizService,
	}
}

// @Summary Create a new quiz
// @Description Create a new quiz with the provided information
// @Tags quizzes
// @Accept json
// @Produce json
// @Param quiz body models.Quiz true "Quiz information"
// @Success 201 {object} models.Quiz
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /quizzes [post]
func (h *QuizHandler) CreateQuiz(c *gin.Context) {
	var quiz models.Quiz
	if err := c.ShouldBindJSON(&quiz); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (assuming you have middleware that sets this)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	quiz.CreatedByID = userID.(uint)

	if err := h.quizService.CreateQuiz(&quiz); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, quiz)
}

// @Summary Get a quiz by ID
// @Description Get quiz information by quiz ID
// @Tags quizzes
// @Produce json
// @Param id path int true "Quiz ID"
// @Success 200 {object} models.Quiz
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /quizzes/{id} [get]
func (h *QuizHandler) GetQuiz(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quiz id"})
		return
	}

	quiz, err := h.quizService.GetQuizByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "quiz not found"})
		return
	}

	c.JSON(http.StatusOK, quiz)
}

// @Summary Get user's quizzes
// @Description Get all quizzes created by the authenticated user
// @Tags quizzes
// @Produce json
// @Success 200 {array} models.Quiz
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /quizzes/user [get]
func (h *QuizHandler) GetUserQuizzes(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	quizzes, err := h.quizService.GetQuizzesByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quizzes)
}

// @Summary Update a quiz
// @Description Update quiz information by quiz ID
// @Tags quizzes
// @Accept json
// @Produce json
// @Param id path int true "Quiz ID"
// @Param quiz body models.Quiz true "Quiz information"
// @Success 200 {object} models.Quiz
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /quizzes/{id} [put]
func (h *QuizHandler) UpdateQuiz(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quiz id"})
		return
	}

	var quiz models.Quiz
	if err := c.ShouldBindJSON(&quiz); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quiz.ID = uint(id)
	if err := h.quizService.UpdateQuiz(&quiz); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quiz)
}

// @Summary Delete a quiz
// @Description Delete a quiz by quiz ID
// @Tags quizzes
// @Produce json
// @Param id path int true "Quiz ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /quizzes/{id} [delete]
func (h *QuizHandler) DeleteQuiz(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quiz id"})
		return
	}

	if err := h.quizService.DeleteQuiz(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "quiz deleted successfully"})
}

// @Summary Add a selection to a quiz
// @Description Add a new selection to an existing quiz
// @Tags quizzes
// @Accept json
// @Produce json
// @Param id path int true "Quiz ID"
// @Param selection body models.QuizSelection true "Selection information"
// @Success 201 {object} models.QuizSelection
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /quizzes/{id}/selections [post]
func (h *QuizHandler) AddSelection(c *gin.Context) {
	quizID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quiz id"})
		return
	}

	var selection models.QuizSelection
	if err := c.ShouldBindJSON(&selection); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.quizService.AddSelection(uint(quizID), &selection); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, selection)
}

// @Summary Remove a selection from a quiz
// @Description Remove a selection from an existing quiz
// @Tags quizzes
// @Produce json
// @Param id path int true "Quiz ID"
// @Param selectionId path int true "Selection ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /quizzes/{id}/selections/{selectionId} [delete]
func (h *QuizHandler) RemoveSelection(c *gin.Context) {
	quizID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quiz id"})
		return
	}

	selectionID, err := strconv.ParseUint(c.Param("selectionId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid selection id"})
		return
	}

	if err := h.quizService.RemoveSelection(uint(quizID), uint(selectionID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "selection removed successfully"})
} 