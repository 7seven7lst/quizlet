package handlers

import (
	"errors"
	"net/http"
	"quizlet/internal/models/quiz_suite"
	"quizlet/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type QuizSuiteHandler struct {
	quizSuiteService service.QuizSuiteService
}

func NewQuizSuiteHandler(quizSuiteService service.QuizSuiteService) *QuizSuiteHandler {
	return &QuizSuiteHandler{
		quizSuiteService: quizSuiteService,
	}
}

// getUserIDFromContext extracts and validates the user ID from the context
func (h *QuizSuiteHandler) getUserIDFromContext(c *gin.Context) (uint, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, errors.New("unauthorized")
	}
	
	uid, ok := userID.(uint)
	if !ok {
		return 0, errors.New("unauthorized")
	}
	
	return uid, nil
}

// validateQuizSuiteAccess checks if the user has access to the quiz suite
func (h *QuizSuiteHandler) validateQuizSuiteAccess(suiteID, userID uint) error {
	suite, err := h.quizSuiteService.GetQuizSuite(suiteID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("quiz suite not found")
		}
		return err
	}
	
	if suite.CreatedByID != userID {
		return errors.New("unauthorized: you don't have permission to access this quiz suite")
	}
	
	return nil
}

// @Summary Create a new quiz suite
// @Description Create a new quiz suite with the provided details
// @Tags quiz-suites
// @Accept json
// @Produce json
// @Param quiz_suite body quiz_suite.CreateQuizSuiteRequest true "Quiz Suite object"
// @Success 201 {object} quiz_suite.QuizSuite
// @Failure 400 {object} map[string]string "Bad Request - Title is required"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Security BearerAuth
// @Router /quiz-suites [post]
func (h *QuizSuiteHandler) CreateQuizSuite(c *gin.Context) {
	var req quiz_suite.CreateQuizSuiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	qs := &quiz_suite.QuizSuite{
		Title:       req.Title,
		Description: req.Description,
		CreatedByID: userID,
	}

	if qs.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quiz suite title is required"})
		return
	}

	if err := h.quizSuiteService.CreateQuizSuite(qs); err != nil {
		if err == gorm.ErrInvalidDB {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gorm: invalid db"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, qs)
}

func (h *QuizSuiteHandler) GetQuizSuites(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	quizSuites, err := h.quizSuiteService.GetUserQuizSuites(userID)
	if err != nil {
		if err == gorm.ErrInvalidDB {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gorm: invalid db"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if quizSuites == nil {
		quizSuites = []*quiz_suite.QuizSuite{}
	}

	c.JSON(http.StatusOK, gin.H{"quiz_suites": quizSuites})
}

// @Summary Get a quiz suite by ID
// @Description Retrieve a specific quiz suite by its ID
// @Tags quiz-suites
// @Produce json
// @Param id path int true "Quiz Suite ID"
// @Success 200 {object} quiz_suite.QuizSuite
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /quiz-suites/{id} [get]
func (h *QuizSuiteHandler) GetQuizSuite(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quiz suite id"})
		return
	}

	quizSuite, err := h.quizSuiteService.GetQuizSuite(uint(id))
	if err != nil {
		if err == gorm.ErrInvalidDB {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gorm: invalid db"})
			return
		}
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "quiz suite not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quizSuite)
}

func (h *QuizSuiteHandler) GetUserQuizSuites(c *gin.Context) {
	// Get user ID from context (assuming you have middleware that sets this)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	quizSuites, err := h.quizSuiteService.GetUserQuizSuites(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"quiz_suites": quizSuites})
}

// @Summary Update a quiz suite
// @Description Update an existing quiz suite with the provided details
// @Tags quiz-suites
// @Accept json
// @Produce json
// @Param id path int true "Quiz Suite ID"
// @Param quiz_suite body quiz_suite.UpdateQuizSuiteRequest true "Quiz Suite update object"
// @Success 200 {object} quiz_suite.QuizSuite
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /quiz-suites/{id} [put]
func (h *QuizSuiteHandler) UpdateQuizSuite(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quiz suite id"})
		return
	}

	var req quiz_suite.UpdateQuizSuiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// First check if the quiz suite exists
	existingSuite, err := h.quizSuiteService.GetQuizSuite(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "quiz suite not found"})
			return
		}
		if err == gorm.ErrInvalidDB {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gorm: invalid db"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update the existing suite with new values
	existingSuite.Title = req.Title
	existingSuite.Description = req.Description

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	existingSuite.CreatedByID = userID

	err = h.quizSuiteService.UpdateQuizSuite(existingSuite)
	if err != nil {
		if err == gorm.ErrInvalidDB {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gorm: invalid db"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existingSuite)
}

// @Summary Delete a quiz suite
// @Description Delete a quiz suite by its ID
// @Tags quiz-suites
// @Produce json
// @Param id path int true "Quiz Suite ID"
// @Success 200 {object} SuccessResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /quiz-suites/{id} [delete]
func (h *QuizSuiteHandler) DeleteQuizSuite(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quiz suite id"})
		return
	}

	err = h.quizSuiteService.DeleteQuizSuite(uint(id))
	if err != nil {
		if err == gorm.ErrInvalidDB {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gorm: invalid db"})
			return
		}
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "quiz suite not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "quiz suite deleted successfully"})
}

// @Summary Add a quiz to a quiz suite
// @Description Add an existing quiz to a quiz suite
// @Tags quiz-suites
// @Produce json
// @Param id path int true "Quiz Suite ID"
// @Param quizId path int true "Quiz ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /quiz-suites/{id}/quizzes/{quizId} [post]
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

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	qs, err := h.quizSuiteService.GetQuizSuite(uint(suiteID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "quiz suite not found"})
			return
		}
		if err == gorm.ErrInvalidDB {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gorm: invalid db"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if qs.CreatedByID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: you don't have permission to access this quiz suite"})
		return
	}

	if err := h.quizSuiteService.AddQuizToSuite(uint(suiteID), uint(quizID)); err != nil {
		if err == gorm.ErrInvalidDB {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gorm: invalid db"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "quiz added to suite successfully"})
}

// @Summary Remove a quiz from a quiz suite
// @Description Remove a quiz from an existing quiz suite
// @Tags quiz-suites
// @Produce json
// @Param id path int true "Quiz Suite ID"
// @Param quizId path int true "Quiz ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /quiz-suites/{id}/quizzes/{quizId} [delete]
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