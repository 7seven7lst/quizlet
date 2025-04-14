package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"quizlet/internal/models/quiz_suite"
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
				mockService.On("CreateQuizSuite", mock.MatchedBy(func(qs *quiz_suite.QuizSuite) bool {
					return qs.Title == "Test Quiz Suite" &&
						qs.Description == "Test Description" &&
						qs.CreatedByID == uint(1)
				})).Return(nil).Once()
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"id":          float64(0),
				"title":       "Test Quiz Suite",
				"description": "Test Description",
				"created_by_id": float64(1),
				"created_at":   "0001-01-01T00:00:00Z",
				"updated_at":   "0001-01-01T00:00:00Z",
				"deleted_at":   nil,
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
				mockService.On("CreateQuizSuite", mock.MatchedBy(func(qs *quiz_suite.QuizSuite) bool {
					return qs.Title == "Test Quiz Suite" &&
						qs.Description == "Test Description" &&
						qs.CreatedByID == uint(1)
				})).Return(gorm.ErrInvalidDB).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "gorm: invalid db",
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
				c.Set("userID", tc.userID)
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
				quizSuites := []*quiz_suite.QuizSuite{
					{
						ID:          1,
						Title:       "Test Quiz Suite 1",
						Description: "Test Description 1",
						CreatedByID: 1,
					},
					{
						ID:          2,
						Title:       "Test Quiz Suite 2",
						Description: "Test Description 2",
						CreatedByID: 1,
					},
				}
				mockService.On("GetUserQuizSuites", uint(1)).Return(quizSuites, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"quiz_suites": []interface{}{
					map[string]interface{}{
						"id":           float64(1),
						"title":        "Test Quiz Suite 1",
						"description":  "Test Description 1",
						"created_by_id": float64(1),
						"created_at":   "0001-01-01T00:00:00Z",
						"updated_at":   "0001-01-01T00:00:00Z",
						"deleted_at":   nil,
					},
					map[string]interface{}{
						"id":           float64(2),
						"title":        "Test Quiz Suite 2",
						"description":  "Test Description 2",
						"created_by_id": float64(1),
						"created_at":   "0001-01-01T00:00:00Z",
						"updated_at":   "0001-01-01T00:00:00Z",
						"deleted_at":   nil,
					},
				},
			},
		},
		{
			name:   "Service Error",
			userID: 1,
			mockSetup: func() {
				mockService.On("GetUserQuizSuites", uint(1)).Return(nil, gorm.ErrInvalidDB).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "gorm: invalid db",
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
				c.Set("userID", tc.userID)
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
		userID         uint
		quizSuiteID    string
		mockSetup      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "Success",
			userID: 1,
			quizSuiteID: "1",
			mockSetup: func() {
				quizSuite := &quiz_suite.QuizSuite{
					ID:          1,
					Title:       "Test Quiz Suite",
					Description: "Test Description",
					CreatedByID: 1,
				}
				mockService.On("GetQuizSuite", uint(1)).Return(quizSuite, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":           float64(1),
				"title":        "Test Quiz Suite",
				"description":  "Test Description",
				"created_by_id": float64(1),
				"created_at":   "0001-01-01T00:00:00Z",
				"updated_at":   "0001-01-01T00:00:00Z",
				"deleted_at":   nil,
			},
		},
		{
			name:   "Service Error",
			userID: 1,
			quizSuiteID: "1",
			mockSetup: func() {
				mockService.On("GetQuizSuite", uint(1)).Return(nil, gorm.ErrInvalidDB).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "gorm: invalid db",
			},
		},
		{
			name:           "Invalid ID",
			userID:         1,
			quizSuiteID:    "invalid",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid quiz suite id",
			},
		},
		{
			name:        "Not Found",
			userID:      1,
			quizSuiteID: "1",
			mockSetup: func() {
				mockService.On("GetQuizSuite", uint(1)).Return(nil, gorm.ErrRecordNotFound).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "quiz suite not found",
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

			// Set user ID in context
			if tc.userID > 0 {
				c.Set("userID", tc.userID)
			}

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
		userID         uint
		quizSuiteID    string
		requestBody    map[string]interface{}
		mockSetup      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "Success",
			userID: 1,
			quizSuiteID: "1",
			requestBody: map[string]interface{}{
				"title":       "Updated Quiz Suite",
				"description": "Updated Description",
			},
			mockSetup: func() {
				mockService.On("UpdateQuizSuite", mock.MatchedBy(func(qs *quiz_suite.QuizSuite) bool {
					return qs.ID == uint(1) &&
						qs.Title == "Updated Quiz Suite" &&
						qs.Description == "Updated Description"
				})).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":           float64(1),
				"title":        "Updated Quiz Suite",
				"description":  "Updated Description",
				"created_by_id": float64(1),
				"created_at":   "0001-01-01T00:00:00Z",
				"updated_at":   "0001-01-01T00:00:00Z",
				"deleted_at":   nil,
			},
		},
		{
			name:   "Service Error",
			userID: 1,
			quizSuiteID: "1",
			requestBody: map[string]interface{}{
				"title":       "Updated Quiz Suite",
				"description": "Updated Description",
			},
			mockSetup: func() {
				mockService.On("UpdateQuizSuite", mock.MatchedBy(func(qs *quiz_suite.QuizSuite) bool {
					return qs.ID == uint(1) &&
						qs.Title == "Updated Quiz Suite" &&
						qs.Description == "Updated Description"
				})).Return(gorm.ErrInvalidDB).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "gorm: invalid db",
			},
		},
		{
			name:           "Invalid ID",
			userID:         1,
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
			userID:      1,
			quizSuiteID: "1",
			requestBody: map[string]interface{}{
				"title": "Updated Quiz Suite",
			},
			mockSetup: func() {
				mockService.On("UpdateQuizSuite", mock.AnythingOfType("*quiz_suite.QuizSuite")).Return(gorm.ErrRecordNotFound).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "quiz suite not found",
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

			// Set user ID in context
			if tc.userID > 0 {
				c.Set("userID", tc.userID)
			}

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
		userID         uint
		quizSuiteID    string
		mockSetup      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "Success",
			userID: 1,
			quizSuiteID: "1",
			mockSetup: func() {
				mockService.On("DeleteQuizSuite", uint(1)).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "quiz suite deleted successfully",
			},
		},
		{
			name:   "Service Error",
			userID: 1,
			quizSuiteID: "1",
			mockSetup: func() {
				mockService.On("DeleteQuizSuite", uint(1)).Return(gorm.ErrInvalidDB).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "gorm: invalid db",
			},
		},
		{
			name:           "Invalid ID",
			userID:         1,
			quizSuiteID:    "invalid",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid quiz suite id",
			},
		},
		{
			name:        "Not Found",
			userID:      1,
			quizSuiteID: "1",
			mockSetup: func() {
				mockService.On("DeleteQuizSuite", uint(1)).Return(gorm.ErrRecordNotFound).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "quiz suite not found",
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

			// Set user ID in context
			if tc.userID > 0 {
				c.Set("userID", tc.userID)
			}

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

func TestAddQuizToSuite(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(services.MockQuizSuiteService)

	testCases := []struct {
		name           string
		userID         uint
		quizSuiteID    string
		quizID         string
		mockSetup      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:        "Success",
			userID:      1,
			quizSuiteID: "1",
			quizID:      "2",
			mockSetup: func() {
				mockService.On("GetQuizSuite", uint(1)).Return(&quiz_suite.QuizSuite{
					ID:          1,
					CreatedByID: 1,
				}, nil).Once()
				mockService.On("AddQuizToSuite", uint(1), uint(2)).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "quiz added to suite successfully",
			},
		},
		{
			name:        "Invalid Quiz Suite ID",
			userID:      1,
			quizSuiteID: "invalid",
			quizID:      "2",
			mockSetup:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid quiz suite id",
			},
		},
		{
			name:        "Invalid Quiz ID",
			userID:      1,
			quizSuiteID: "1",
			quizID:      "invalid",
			mockSetup:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid quiz id",
			},
		},
		{
			name:        "Unauthorized - No User ID",
			userID:      0,
			quizSuiteID: "1",
			quizID:      "2",
			mockSetup:   func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "unauthorized",
			},
		},
		{
			name:        "Quiz Suite Not Found",
			userID:      1,
			quizSuiteID: "1",
			quizID:      "2",
			mockSetup: func() {
				mockService.On("GetQuizSuite", uint(1)).Return(nil, gorm.ErrRecordNotFound).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "quiz suite not found",
			},
		},
		{
			name:        "Unauthorized - Different User",
			userID:      2,
			quizSuiteID: "1",
			quizID:      "2",
			mockSetup: func() {
				mockService.On("GetQuizSuite", uint(1)).Return(&quiz_suite.QuizSuite{
					ID:          1,
					CreatedByID: 1,
				}, nil).Once()
			},
			expectedStatus: http.StatusForbidden,
			expectedBody: map[string]interface{}{
				"error": "unauthorized: you don't have permission to access this quiz suite",
			},
		},
		{
			name:        "Service Error",
			userID:      1,
			quizSuiteID: "1",
			quizID:      "2",
			mockSetup: func() {
				mockService.On("GetQuizSuite", uint(1)).Return(&quiz_suite.QuizSuite{
					ID:          1,
					CreatedByID: 1,
				}, nil).Once()
				mockService.On("AddQuizToSuite", uint(1), uint(2)).Return(gorm.ErrInvalidDB).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "gorm: invalid db",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set up request
			c.Request = httptest.NewRequest(http.MethodPost, "/quiz-suites/"+tc.quizSuiteID+"/quizzes/"+tc.quizID, nil)
			c.Params = []gin.Param{
				{Key: "id", Value: tc.quizSuiteID},
				{Key: "quizId", Value: tc.quizID},
			}

			// Set user ID in context
			if tc.userID > 0 {
				c.Set("userID", tc.userID)
			}

			// Set up mock
			tc.mockSetup()

			// Create handler and execute
			handler := NewQuizSuiteHandler(mockService)
			handler.AddQuizToSuite(c)

			// Assert response
			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedBody, response)
		})
	}
}

func TestRemoveQuizFromSuite(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(services.MockQuizSuiteService)

	testCases := []struct {
		name           string
		quizSuiteID    string
		quizID         string
		mockSetup      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:        "Success",
			quizSuiteID: "1",
			quizID:      "2",
			mockSetup: func() {
				mockService.On("RemoveQuizFromSuite", uint(1), uint(2)).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "quiz removed from suite successfully",
			},
		},
		{
			name:        "Invalid Quiz Suite ID",
			quizSuiteID: "invalid",
			quizID:      "2",
			mockSetup:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid quiz suite id",
			},
		},
		{
			name:        "Invalid Quiz ID",
			quizSuiteID: "1",
			quizID:      "invalid",
			mockSetup:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid quiz id",
			},
		},
		{
			name:        "Service Error",
			quizSuiteID: "1",
			quizID:      "2",
			mockSetup: func() {
				mockService.On("RemoveQuizFromSuite", uint(1), uint(2)).Return(gorm.ErrInvalidDB).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "invalid db",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set up request
			c.Request = httptest.NewRequest(http.MethodDelete, "/quiz-suites/"+tc.quizSuiteID+"/quizzes/"+tc.quizID, nil)
			c.Params = []gin.Param{
				{Key: "id", Value: tc.quizSuiteID},
				{Key: "quizId", Value: tc.quizID},
			}

			// Set up mock
			tc.mockSetup()

			// Create handler and execute
			handler := NewQuizSuiteHandler(mockService)
			handler.RemoveQuizFromSuite(c)

			// Assert response
			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedBody, response)
		})
	}
} 