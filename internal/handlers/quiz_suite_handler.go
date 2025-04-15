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

func (h *QuizSuiteHandler) CreateQuizSuite(c *gin.Context) {
	var qs quiz_suite.QuizSuite
	if err := c.ShouldBindJSON(&qs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	qs.CreatedByID = userID

	if qs.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quiz suite title is required"})
		return
	}

	if err := h.quizSuiteService.CreateQuizSuite(&qs); err != nil {
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
// @Param quiz_suite body quiz_suite.QuizSuite true "Quiz Suite object"
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

	var updateQS quiz_suite.QuizSuite
	if err := c.ShouldBindJSON(&updateQS); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateQS.ID = uint(id)
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	updateQS.CreatedByID = userID

	err = h.quizSuiteService.UpdateQuizSuite(&updateQS)
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

	c.JSON(http.StatusOK, updateQS)
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