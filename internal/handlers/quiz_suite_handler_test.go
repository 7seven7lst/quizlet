package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"quizlet/internal/models"
	"quizlet/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreateQuizSuite(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(services.MockQuizSuiteService)

	testCases := []struct {
		name           string
		userID         uint
		requestBody    map[string]interface{}
		mockSetup      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "Success",
			userID: 1,
			requestBody: map[string]interface{}{
				"title":       "Test Quiz Suite",
				"description": "Test Description",
			},
			mockSetup: func() {
				expectedQuizSuite := &models.QuizSuite{
					ID:          1,
					Title:       "Test Quiz Suite",
					Description: "Test Description",
					CreatedByID: 1,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				mockService.On("CreateQuizSuite", mock.AnythingOfType("*models.QuizSuite")).Return(expectedQuizSuite, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"id":          float64(1),
				"title":       "Test Quiz Suite",
				"description": "Test Description",
				"created_by_id": float64(1),
			},
		},
		{
			name:   "Unauthorized",
			userID: 0,
			requestBody: map[string]interface{}{
				"title":       "Test Quiz Suite",
				"description": "Test Description",
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "unauthorized",
			},
		},
		{
			name:   "Service Error",
			userID: 1,
			requestBody: map[string]interface{}{
				"title":       "Test Quiz Suite",
				"description": "Test Description",
			},
			mockSetup: func() {
				mockService.On("CreateQuizSuite", mock.AnythingOfType("*models.QuizSuite")).Return(nil, gorm.ErrInvalidDB)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": gorm.ErrInvalidDB.Error(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set up request
			body, _ := json.Marshal(tc.requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/quiz-suites", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Set user ID in context
			if tc.userID > 0 {
				c.Set("user_id", tc.userID)
			}

			// Set up mock
			tc.mockSetup()

			// Create handler and execute
			handler := NewQuizSuiteHandler(mockService)
			handler.CreateQuizSuite(c)

			// Assert response
			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedBody, response)
		})
	}
}

func TestGetQuizSuites(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(services.MockQuizSuiteService)

	testCases := []struct {
		name           string
		userID         uint
		mockSetup      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "Success",
			userID: 1,
			mockSetup: func() {
				quizSuites := []models.QuizSuite{
					{
						ID:          1,
						Title:       "Test Quiz Suite 1",
						Description: "Test Description 1",
						CreatedByID: 1,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
					{
						ID:          2,
						Title:       "Test Quiz Suite 2",
						Description: "Test Description 2",
						CreatedByID: 1,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
				}
				mockService.On("GetQuizSuites", uint(1)).Return(quizSuites, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"quiz_suites": []interface{}{
					map[string]interface{}{
						"id":          float64(1),
						"title":       "Test Quiz Suite 1",
						"description": "Test Description 1",
						"created_by_id": float64(1),
					},
					map[string]interface{}{
						"id":          float64(2),
						"title":       "Test Quiz Suite 2",
						"description": "Test Description 2",
						"created_by_id": float64(1),
					},
				},
			},
		},
		{
			name:           "Unauthorized",
			userID:         0,
			mockSetup:      func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "unauthorized",
			},
		},
		{
			name:   "Service Error",
			userID: 1,
			mockSetup: func() {
				mockService.On("GetQuizSuites", uint(1)).Return(nil, gorm.ErrInvalidDB)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": gorm.ErrInvalidDB.Error(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set up request
			c.Request = httptest.NewRequest(http.MethodGet, "/quiz-suites", nil)

			// Set user ID in context
			if tc.userID > 0 {
				c.Set("user_id", tc.userID)
			}

			// Set up mock
			tc.mockSetup()

			// Create handler and execute
			handler := NewQuizSuiteHandler(mockService)
			handler.GetQuizSuites(c)

			// Assert response
			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedBody, response)
		})
	}
}

func TestGetQuizSuite(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(services.MockQuizSuiteService)

	testCases := []struct {
		name           string
		quizSuiteID    string
		mockSetup      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:        "Success",
			quizSuiteID: "1",
			mockSetup: func() {
				quizSuite := &models.QuizSuite{
					ID:          1,
					Title:       "Test Quiz Suite",
					Description: "Test Description",
					CreatedByID: 1,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				mockService.On("GetQuizSuite", uint(1)).Return(quizSuite, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":          float64(1),
				"title":       "Test Quiz Suite",
				"description": "Test Description",
				"created_by_id": float64(1),
			},
		},
		{
			name:           "Invalid ID",
			quizSuiteID:    "invalid",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid quiz suite id",
			},
		},
		{
			name:        "Not Found",
			quizSuiteID: "1",
			mockSetup: func() {
				mockService.On("GetQuizSuite", uint(1)).Return(nil, gorm.ErrRecordNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "quiz suite not found",
			},
		},
		{
			name:        "Service Error",
			quizSuiteID: "1",
			mockSetup: func() {
				mockService.On("GetQuizSuite", uint(1)).Return(nil, gorm.ErrInvalidDB)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": gorm.ErrInvalidDB.Error(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set up request
			c.Request = httptest.NewRequest(http.MethodGet, "/quiz-suites/"+tc.quizSuiteID, nil)
			c.Params = []gin.Param{{Key: "id", Value: tc.quizSuiteID}}

			// Set up mock
			tc.mockSetup()

			// Create handler and execute
			handler := NewQuizSuiteHandler(mockService)
			handler.GetQuizSuite(c)

			// Assert response
			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedBody, response)
		})
	}
}

func TestUpdateQuizSuite(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(services.MockQuizSuiteService)

	testCases := []struct {
		name           string
		quizSuiteID    string
		requestBody    map[string]interface{}
		mockSetup      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:        "Success",
			quizSuiteID: "1",
			requestBody: map[string]interface{}{
				"title":       "Updated Quiz Suite",
				"description": "Updated Description",
			},
			mockSetup: func() {
				updatedQuizSuite := &models.QuizSuite{
					ID:          1,
					Title:       "Updated Quiz Suite",
					Description: "Updated Description",
					CreatedByID: 1,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				mockService.On("UpdateQuizSuite", mock.AnythingOfType("*models.QuizSuite")).Return(updatedQuizSuite, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":          float64(1),
				"title":       "Updated Quiz Suite",
				"description": "Updated Description",
				"created_by_id": float64(1),
			},
		},
		{
			name:           "Invalid ID",
			quizSuiteID:    "invalid",
			requestBody:    map[string]interface{}{},
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid quiz suite id",
			},
		},
		{
			name:        "Not Found",
			quizSuiteID: "1",
			requestBody: map[string]interface{}{
				"title": "Updated Quiz Suite",
			},
			mockSetup: func() {
				mockService.On("UpdateQuizSuite", mock.AnythingOfType("*models.QuizSuite")).Return(nil, gorm.ErrRecordNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "quiz suite not found",
			},
		},
		{
			name:        "Service Error",
			quizSuiteID: "1",
			requestBody: map[string]interface{}{
				"title": "Updated Quiz Suite",
			},
			mockSetup: func() {
				mockService.On("UpdateQuizSuite", mock.AnythingOfType("*models.QuizSuite")).Return(nil, gorm.ErrInvalidDB)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": gorm.ErrInvalidDB.Error(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set up request
			body, _ := json.Marshal(tc.requestBody)
			c.Request = httptest.NewRequest(http.MethodPut, "/quiz-suites/"+tc.quizSuiteID, bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Params = []gin.Param{{Key: "id", Value: tc.quizSuiteID}}

			// Set up mock
			tc.mockSetup()

			// Create handler and execute
			handler := NewQuizSuiteHandler(mockService)
			handler.UpdateQuizSuite(c)

			// Assert response
			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedBody, response)
		})
	}
}

func TestDeleteQuizSuite(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(services.MockQuizSuiteService)

	testCases := []struct {
		name           string
		quizSuiteID    string
		mockSetup      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:        "Success",
			quizSuiteID: "1",
			mockSetup: func() {
				mockService.On("DeleteQuizSuite", uint(1)).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "quiz suite deleted successfully",
			},
		},
		{
			name:           "Invalid ID",
			quizSuiteID:    "invalid",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid quiz suite id",
			},
		},
		{
			name:        "Not Found",
			quizSuiteID: "1",
			mockSetup: func() {
				mockService.On("DeleteQuizSuite", uint(1)).Return(gorm.ErrRecordNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "quiz suite not found",
			},
		},
		{
			name:        "Service Error",
			quizSuiteID: "1",
			mockSetup: func() {
				mockService.On("DeleteQuizSuite", uint(1)).Return(gorm.ErrInvalidDB)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": gorm.ErrInvalidDB.Error(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set up request
			c.Request = httptest.NewRequest(http.MethodDelete, "/quiz-suites/"+tc.quizSuiteID, nil)
			c.Params = []gin.Param{{Key: "id", Value: tc.quizSuiteID}}

			// Set up mock
			tc.mockSetup()

			// Create handler and execute
			handler := NewQuizSuiteHandler(mockService)
			handler.DeleteQuizSuite(c)

			// Assert response
			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedBody, response)
		})
	}
} 