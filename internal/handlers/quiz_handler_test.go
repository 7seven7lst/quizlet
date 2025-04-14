package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"quizlet/internal/models"
	"quizlet/tests/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateQuiz(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new controller for mocking
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock quiz service
	mockQuizService := mocks.NewMockQuizService(ctrl)

	// Create quiz handler with mock service
	handler := NewQuizHandler(mockQuizService)

	// Test cases
	tests := []struct {
		name           string
		inputQuiz      models.Quiz
		setupAuth      func(*gin.Context)
		mockService    func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "Successful quiz creation",
			inputQuiz: models.Quiz{
				Question:      "What is the capital of France?",
				QuizType:      models.QuizTypeSingleChoice,
				CorrectAnswer: "Paris",
			},
			setupAuth: func(c *gin.Context) {
				c.Set("userID", uint(1))
			},
			mockService: func() {
				mockQuizService.EXPECT().
					CreateQuiz(gomock.Any()).
					Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"question":       "What is the capital of France?",
				"quiz_type":      "single_choice",
				"correct_answer": "Paris",
				"created_by_id":  float64(1),
			},
		},
		{
			name: "Unauthorized - no user ID",
			inputQuiz: models.Quiz{
				Question:      "What is the capital of France?",
				QuizType:      models.QuizTypeSingleChoice,
				CorrectAnswer: "Paris",
			},
			setupAuth: func(c *gin.Context) {
				// Don't set userID
			},
			mockService: func() {
				// No service calls expected
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "unauthorized",
			},
		},
		{
			name: "Service error",
			inputQuiz: models.Quiz{
				Question:      "What is the capital of France?",
				QuizType:      models.QuizTypeSingleChoice,
				CorrectAnswer: "Paris",
			},
			setupAuth: func(c *gin.Context) {
				c.Set("userID", uint(1))
			},
			mockService: func() {
				mockQuizService.EXPECT().
					CreateQuiz(gomock.Any()).
					Return(gorm.ErrInvalidDB)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": gorm.ErrInvalidDB.Error(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock service expectations
			tt.mockService()

			// Create a new Gin context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Create request body
			body, _ := json.Marshal(tt.inputQuiz)
			c.Request = httptest.NewRequest(http.MethodPost, "/quizzes", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Setup auth context
			tt.setupAuth(c)

			// Call the handler
			handler.CreateQuiz(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// For successful creation, check specific fields
			if tt.expectedStatus == http.StatusCreated {
				assert.Equal(t, tt.expectedBody["question"], response["question"])
				assert.Equal(t, tt.expectedBody["quiz_type"], response["quiz_type"])
				assert.Equal(t, tt.expectedBody["correct_answer"], response["correct_answer"])
				assert.Equal(t, tt.expectedBody["created_by_id"], response["created_by_id"])
			} else {
				assert.Equal(t, tt.expectedBody["error"], response["error"])
			}
		})
	}
}

func TestGetQuiz(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new controller for mocking
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock quiz service
	mockQuizService := mocks.NewMockQuizService(ctrl)

	// Create quiz handler with mock service
	handler := NewQuizHandler(mockQuizService)

	// Test cases
	tests := []struct {
		name           string
		quizID         string
		setupAuth      func(*gin.Context)
		mockService    func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "Successful quiz retrieval",
			quizID: "1",
			setupAuth: func(c *gin.Context) {
				c.Set("userID", uint(1))
			},
			mockService: func() {
				mockQuizService.EXPECT().
					GetQuizByID(uint(1)).
					Return(&models.Quiz{
						Question:      "What is the capital of France?",
						QuizType:      models.QuizTypeSingleChoice,
						CorrectAnswer: "Paris",
						CreatedByID:   1,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"question":       "What is the capital of France?",
				"quiz_type":      "single_choice",
				"correct_answer": "Paris",
				"created_by_id":  float64(1),
			},
		},
		{
			name:   "Invalid quiz ID",
			quizID: "invalid",
			setupAuth: func(c *gin.Context) {
				c.Set("userID", uint(1))
			},
			mockService: func() {
				// No service calls expected
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid quiz id",
			},
		},
		{
			name:   "Quiz not found",
			quizID: "999",
			setupAuth: func(c *gin.Context) {
				c.Set("userID", uint(1))
			},
			mockService: func() {
				mockQuizService.EXPECT().
					GetQuizByID(uint(999)).
					Return(nil, gorm.ErrRecordNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "quiz not found",
			},
		},
		{
			name:   "Service error",
			quizID: "1",
			setupAuth: func(c *gin.Context) {
				c.Set("userID", uint(1))
			},
			mockService: func() {
				mockQuizService.EXPECT().
					GetQuizByID(uint(1)).
					Return(nil, gorm.ErrInvalidDB)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "quiz not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock service expectations
			tt.mockService()

			// Create a new Gin context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set up the request
			c.Request = httptest.NewRequest(http.MethodGet, "/quizzes/"+tt.quizID, nil)
			c.Params = []gin.Param{{Key: "id", Value: tt.quizID}}

			// Setup auth context
			tt.setupAuth(c)

			// Call the handler
			handler.GetQuiz(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// For successful retrieval, check specific fields
			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, tt.expectedBody["question"], response["question"])
				assert.Equal(t, tt.expectedBody["quiz_type"], response["quiz_type"])
				assert.Equal(t, tt.expectedBody["correct_answer"], response["correct_answer"])
				assert.Equal(t, tt.expectedBody["created_by_id"], response["created_by_id"])
			} else {
				assert.Equal(t, tt.expectedBody["error"], response["error"])
			}
		})
	}
} 