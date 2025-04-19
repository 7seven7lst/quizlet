package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"quizlet/internal/models/quiz_attempt"
	"quizlet/internal/service"
)

type QuizAttemptHandler struct {
	quizAttemptService service.QuizAttemptService
}

func NewQuizAttemptHandler(quizAttemptService service.QuizAttemptService) *QuizAttemptHandler {
	return &QuizAttemptHandler{
		quizAttemptService: quizAttemptService,
	}
}

// getUserIDFromContext extracts and validates the user ID from the context
func (h *QuizAttemptHandler) getUserIDFromContext(c *gin.Context) (int64, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, errors.New("unauthorized")
	}
	
	uid, ok := userID.(int64)
	if !ok {
		return 0, errors.New("unauthorized")
	}
	
	return uid, nil
}

// ListQuizAttempts godoc
// @Summary List quiz attempts for a quiz suite
// @Description Get all quiz attempts for a specific quiz suite that belong to the authenticated user
// @Tags quiz-attempts
// @Accept json
// @Produce json
// @Param id path int true "Quiz Suite ID"
// @Security BearerAuth
// @Success 200 {array} quiz_attempt.QuizAttempt
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /quiz-suites/{id}/attempts [get]
func (h *QuizAttemptHandler) ListQuizAttempts(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	quizSuiteID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quiz suite id"})
		return
	}

	attempts, err := h.quizAttemptService.ListByQuizSuite(c.Request.Context(), quizSuiteID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, attempts)
}

// CreateQuizAttempt godoc
// @Summary Create a new quiz attempt
// @Description Create a new quiz attempt for a specific quiz suite
// @Tags quiz-attempts
// @Accept json
// @Produce json
// @Param id path int true "Quiz Suite ID"
// @Param request body quiz_attempt.CreateQuizAttemptRequest true "Quiz attempt creation request"
// @Security BearerAuth
// @Success 201 {object} quiz_attempt.QuizAttempt
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /quiz-suites/{id}/attempts [post]
func (h *QuizAttemptHandler) CreateQuizAttempt(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	quizSuiteID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quiz suite id"})
		return
	}

	var req quiz_attempt.CreateQuizAttemptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	attempt, err := h.quizAttemptService.Create(c.Request.Context(), quizSuiteID, userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, attempt)
}

// GetQuizAttempt godoc
// @Summary Get a specific quiz attempt
// @Description Get a specific quiz attempt by ID
// @Tags quiz-attempts
// @Accept json
// @Produce json
// @Param id path int true "Quiz Suite ID"
// @Param attemptId path int true "Quiz Attempt ID"
// @Security BearerAuth
// @Success 200 {object} quiz_attempt.QuizAttempt
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /quiz-suites/{id}/attempts/{attemptId} [get]
func (h *QuizAttemptHandler) GetQuizAttempt(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	attemptID, err := strconv.ParseInt(c.Param("attemptId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid attempt id"})
		return
	}

	attempt, err := h.quizAttemptService.Get(c.Request.Context(), attemptID, userID)
	if err != nil {
		if err == service.ErrQuizAttemptNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, attempt)
}

// UpdateQuizAttempt godoc
// @Summary Update a quiz attempt
// @Description Update a specific quiz attempt by ID
// @Tags quiz-attempts
// @Accept json
// @Produce json
// @Param id path int true "Quiz Suite ID"
// @Param attemptId path int true "Quiz Attempt ID"
// @Param request body quiz_attempt.UpdateQuizAttemptRequest true "Quiz attempt update request"
// @Security BearerAuth
// @Success 200 {object} quiz_attempt.QuizAttempt
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /quiz-suites/{id}/attempts/{attemptId} [put]
func (h *QuizAttemptHandler) UpdateQuizAttempt(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	attemptID, err := strconv.ParseInt(c.Param("attemptId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid attempt id"})
		return
	}

	var req quiz_attempt.UpdateQuizAttemptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	attempt, err := h.quizAttemptService.Update(c.Request.Context(), attemptID, userID, req)
	if err != nil {
		if err == service.ErrQuizAttemptNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, attempt)
}

// DeleteQuizAttempt godoc
// @Summary Delete a quiz attempt
// @Description Delete a specific quiz attempt by ID
// @Tags quiz-attempts
// @Accept json
// @Produce json
// @Param id path int true "Quiz Suite ID"
// @Param attemptId path int true "Quiz Attempt ID"
// @Security BearerAuth
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /quiz-suites/{id}/attempts/{attemptId} [delete]
func (h *QuizAttemptHandler) DeleteQuizAttempt(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	attemptID, err := strconv.ParseInt(c.Param("attemptId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid attempt id"})
		return
	}

	err = h.quizAttemptService.Delete(c.Request.Context(), attemptID, userID)
	if err != nil {
		if err == service.ErrQuizAttemptNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
} 