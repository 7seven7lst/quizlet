package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"quizlet/internal/models/quiz_suite"
	"quizlet/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"github.com/golang/mock/gomock"
	"quizlet/internal/models/user"
	"quizlet/internal/models/quiz"
	"quizlet/tests/mocks"
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

			// For success case, only compare non-timestamp fields
			if tc.expectedStatus == http.StatusOK {
				assert.Equal(t, tc.expectedBody["id"], response["id"])
				assert.Equal(t, tc.expectedBody["title"], response["title"])
				assert.Equal(t, tc.expectedBody["description"], response["description"])
				assert.Equal(t, tc.expectedBody["created_by_id"], response["created_by_id"])
				assert.Equal(t, tc.expectedBody["deleted_at"], response["deleted_at"])
			} else {
				assert.Equal(t, tc.expectedBody, response)
			}
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

			// For success case, only compare non-timestamp fields
			if tc.expectedStatus == http.StatusOK {
				assert.Equal(t, tc.expectedBody["quiz_suites"], response["quiz_suites"])
			} else {
				assert.Equal(t, tc.expectedBody, response)
			}
		})
	}
}

func TestGetQuizSuite(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new controller for mocking
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock quiz suite service
	mockQuizSuiteService := mocks.NewMockQuizSuiteService(ctrl)

	// Create quiz suite handler with mock service
	handler := NewQuizSuiteHandler(mockQuizSuiteService)

	// Test cases
	tests := []struct {
		name           string
		suiteID        string
		setupAuth      func(*gin.Context)
		mockService    func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:    "Successful quiz suite retrieval",
			suiteID: "1",
			setupAuth: func(c *gin.Context) {
				c.Set("userID", uint(1))
			},
			mockService: func() {
				mockQuizSuiteService.EXPECT().
					GetQuizSuite(uint(1)).
					Return(&quiz_suite.QuizSuite{
						Title:       "Test Suite",
						Description: "Test Description",
						CreatedByID: 1,
						CreatedBy: &user.User{
							ID:       1,
							Username: "testuser",
							Email:    "test@example.com",
							Password: "hashedpassword", // This should not be exposed in response
						},
						Quizzes: []*quiz.Quiz{
							{
								Question:    "What is the capital of France?",
								QuizType:    quiz.QuizTypeSingleChoice,
								CreatedByID: 1,
								CreatedBy: &user.User{
									ID:       1,
									Username: "testuser",
									Email:    "test@example.com",
									Password: "hashedpassword", // This should not be exposed in response
								},
							},
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"title":        "Test Suite",
				"description":  "Test Description",
				"created_by_id": float64(1),
				"created_by": map[string]interface{}{
					"id":       float64(1),
					"username": "testuser",
					"email":    "test@example.com",
				},
				"quizzes": []interface{}{
					map[string]interface{}{
						"question":      "What is the capital of France?",
						"quiz_type":     "single_choice",
						"created_by_id": float64(1),
						"created_by": map[string]interface{}{
							"id":       float64(1),
							"username": "testuser",
							"email":    "test@example.com",
						},
					},
				},
			},
		},
		{
			name:   "Service Error",
			suiteID: "1",
			setupAuth: func(c *gin.Context) {
				c.Set("userID", uint(1))
			},
			mockService: func() {
				mockQuizSuiteService.EXPECT().
					GetQuizSuite(uint(1)).
					Return(nil, gorm.ErrInvalidDB).
					Times(1)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "gorm: invalid db",
			},
		},
		{
			name:           "Invalid ID",
			suiteID:        "invalid",
			setupAuth:      func(c *gin.Context) {},
			mockService:    func() {}, // Add empty mock service function
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid quiz suite id",
			},
		},
		{
			name:        "Not Found",
			suiteID:     "1",
			setupAuth: func(c *gin.Context) {
				c.Set("userID", uint(1))
			},
			mockService: func() {
				mockQuizSuiteService.EXPECT().
					GetQuizSuite(uint(1)).
					Return(nil, gorm.ErrRecordNotFound).
					Times(1)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "quiz suite not found",
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
			c.Request = httptest.NewRequest(http.MethodGet, "/quiz-suites/"+tt.suiteID, nil)
			c.Params = []gin.Param{{Key: "id", Value: tt.suiteID}}

			// Setup auth context
			tt.setupAuth(c)

			// Call the handler
			handler.GetQuizSuite(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// For successful retrieval, check specific fields
			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, tt.expectedBody["title"], response["title"])
				assert.Equal(t, tt.expectedBody["description"], response["description"])
				assert.Equal(t, tt.expectedBody["created_by_id"], response["created_by_id"])
				
				// Verify user data is present but password is not exposed
				createdBy, ok := response["created_by"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, tt.expectedBody["created_by"].(map[string]interface{})["id"], createdBy["id"])
				assert.Equal(t, tt.expectedBody["created_by"].(map[string]interface{})["username"], createdBy["username"])
				assert.Equal(t, tt.expectedBody["created_by"].(map[string]interface{})["email"], createdBy["email"])
				_, hasPassword := createdBy["password"]
				assert.False(t, hasPassword, "Password should not be exposed in response")

				// Verify quizzes array
				quizzes, ok := response["quizzes"].([]interface{})
				assert.True(t, ok)
				assert.Len(t, quizzes, 1)

				// Verify first quiz
				firstQuiz, ok := quizzes[0].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, "What is the capital of France?", firstQuiz["question"])
				assert.Equal(t, "single_choice", firstQuiz["quiz_type"])

				// Verify quiz creator data
				quizCreator, ok := firstQuiz["created_by"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, float64(1), quizCreator["id"])
				assert.Equal(t, "testuser", quizCreator["username"])
				assert.Equal(t, "test@example.com", quizCreator["email"])
				_, hasQuizCreatorPassword := quizCreator["password"]
				assert.False(t, hasQuizCreatorPassword, "Password should not be exposed in quiz creator data")
			} else {
				assert.Equal(t, tt.expectedBody["error"], response["error"])
			}
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
		requestBody    quiz_suite.UpdateQuizSuiteRequest
		mockSetup      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:        "Not Found",
			userID:      1,
			quizSuiteID: "999",
			requestBody: quiz_suite.UpdateQuizSuiteRequest{
				Title:       "Updated Title",
				Description: "Updated Description",
			},
			mockSetup: func() {
				mockService.On("GetQuizSuite", uint(999)).Return(nil, gorm.ErrRecordNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "quiz suite not found",
			},
		},
		{
			name:        "Success",
			userID:      1,
			quizSuiteID: "1",
			requestBody: quiz_suite.UpdateQuizSuiteRequest{
				Title:       "Updated Title",
				Description: "Updated Description",
			},
			mockSetup: func() {
				existingSuite := &quiz_suite.QuizSuite{
					ID:           1,
					Title:        "Original Title",
					Description:  "Original Description",
					CreatedByID:  1,
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				}
				mockService.On("GetQuizSuite", uint(1)).Return(existingSuite, nil)
				mockService.On("UpdateQuizSuite", mock.AnythingOfType("*quiz_suite.QuizSuite")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":           float64(1),
				"title":        "Updated Title",
				"description":  "Updated Description",
				"created_by_id": float64(1),
				"created_at":   "2025-04-18T00:20:18.128464635Z",
				"updated_at":   "2025-04-18T00:20:18.128464688Z",
				"deleted_at":   nil,
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

			// For success case, only compare non-timestamp fields
			if tc.expectedStatus == http.StatusOK {
				assert.Equal(t, tc.expectedBody["id"], response["id"])
				assert.Equal(t, tc.expectedBody["title"], response["title"])
				assert.Equal(t, tc.expectedBody["description"], response["description"])
				assert.Equal(t, tc.expectedBody["created_by_id"], response["created_by_id"])
				assert.Equal(t, tc.expectedBody["deleted_at"], response["deleted_at"])
			} else {
				assert.Equal(t, tc.expectedBody, response)
			}
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

			// For success case, only compare non-timestamp fields
			if tc.expectedStatus == http.StatusOK {
				assert.Equal(t, tc.expectedBody["message"], response["message"])
			} else {
				assert.Equal(t, tc.expectedBody, response)
			}
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