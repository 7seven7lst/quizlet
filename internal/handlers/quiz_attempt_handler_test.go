package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"quizlet/internal/models/quiz_attempt"
	"quizlet/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockQuizAttemptService is a mock implementation of the QuizAttemptService interface
type MockQuizAttemptService struct {
	mock.Mock
}

// Ensure MockQuizAttemptService implements the QuizAttemptService interface
var _ service.QuizAttemptService = (*MockQuizAttemptService)(nil)

func (m *MockQuizAttemptService) ListByQuizSuite(ctx context.Context, quizSuiteID, userID int64) ([]quiz_attempt.QuizAttempt, error) {
	args := m.Called(ctx, quizSuiteID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]quiz_attempt.QuizAttempt), args.Error(1)
}

func (m *MockQuizAttemptService) Create(ctx context.Context, quizSuiteID, userID int64, req quiz_attempt.CreateQuizAttemptRequest) (*quiz_attempt.QuizAttempt, error) {
	args := m.Called(ctx, quizSuiteID, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*quiz_attempt.QuizAttempt), args.Error(1)
}

func (m *MockQuizAttemptService) Get(ctx context.Context, id, userID int64) (*quiz_attempt.QuizAttempt, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*quiz_attempt.QuizAttempt), args.Error(1)
}

func (m *MockQuizAttemptService) Update(ctx context.Context, id, userID int64, req quiz_attempt.UpdateQuizAttemptRequest) (*quiz_attempt.QuizAttempt, error) {
	args := m.Called(ctx, id, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*quiz_attempt.QuizAttempt), args.Error(1)
}

func (m *MockQuizAttemptService) Delete(ctx context.Context, id, userID int64) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func setupTestRouter() (*gin.Engine, *MockQuizAttemptService) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(gin.Recovery())

	mockService := new(MockQuizAttemptService)
	
	return router, mockService
}

// Helper function to format time for testing
func formatTimeForTest(t time.Time) string {
	// Format without nanoseconds for consistent testing
	return t.Format("2006-01-02T15:04:05Z")
}

func TestListQuizAttempts(t *testing.T) {
	router, mockService := setupTestRouter()
	
	// Create a handler with the mock service
	handler := NewQuizAttemptHandler(mockService)
	
	router.GET("/quiz-suites/:id/attempts", func(c *gin.Context) {
		c.Set("userID", int64(1))
		handler.ListQuizAttempts(c)
	})

	tests := []struct {
		name           string
		quizSuiteID    string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "Success",
			quizSuiteID: "1",
			setupMock: func() {
				now := time.Now()
				attempts := []quiz_attempt.QuizAttempt{
					{
						ID:          1,
						UserID:      1,
						QuizSuiteID: 1,
						Score:       80,
						Completed:   true,
						StartedAt:   now,
						CompletedAt: &now,
					},
				}
				mockService.On("ListByQuizSuite", mock.Anything, int64(1), int64(1)).Return(attempts, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":1,"user_id":1,"quiz_suite_id":1,"score":80,"completed":true,"started_at":"` + formatTimeForTest(time.Now()) + `","completed_at":"` + formatTimeForTest(time.Now()) + `"}]`,
		},
		{
			name:        "Invalid Quiz Suite ID",
			quizSuiteID: "invalid",
			setupMock:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid quiz suite id"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/quiz-suites/"+tt.quizSuiteID+"/attempts", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			
			// For success case, we need to check that the response contains the expected fields
			// but ignore additional fields like created_at and updated_at
			if tt.expectedStatus == http.StatusOK {
				var expectedData, actualData []map[string]interface{}
				err := json.Unmarshal([]byte(tt.expectedBody), &expectedData)
				assert.NoError(t, err)
				
				err = json.Unmarshal(w.Body.Bytes(), &actualData)
				assert.NoError(t, err)
				
				// Check that all expected fields are present with correct values
				for i, expected := range expectedData {
					actual := actualData[i]
					for key, value := range expected {
						// For timestamp fields, we need to compare only the date part
						if key == "started_at" || key == "completed_at" {
							expectedTime, _ := time.Parse(time.RFC3339, value.(string))
							actualTime, _ := time.Parse(time.RFC3339, actual[key].(string))
							assert.Equal(t, formatTimeForTest(expectedTime), formatTimeForTest(actualTime), "Field %s should match", key)
						} else {
							assert.Equal(t, value, actual[key], "Field %s should match", key)
						}
					}
				}
			} else {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestCreateQuizAttempt(t *testing.T) {
	router, mockService := setupTestRouter()
	
	// Create a handler with the mock service
	handler := NewQuizAttemptHandler(mockService)
	
	router.POST("/quiz-suites/:id/attempts", func(c *gin.Context) {
		c.Set("userID", int64(1))
		handler.CreateQuizAttempt(c)
	})

	tests := []struct {
		name           string
		quizSuiteID    string
		requestBody    string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "Success",
			quizSuiteID: "1",
			requestBody: `{"score":80,"completed":true}`,
			setupMock: func() {
				now := time.Now()
				attempt := &quiz_attempt.QuizAttempt{
					ID:          1,
					UserID:      1,
					QuizSuiteID: 1,
					Score:       80,
					Completed:   true,
					StartedAt:   now,
					CompletedAt: &now,
				}
				mockService.On("Create", mock.Anything, int64(1), int64(1), quiz_attempt.CreateQuizAttemptRequest{
					Score:     80,
					Completed: true,
				}).Return(attempt, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":1,"user_id":1,"quiz_suite_id":1,"score":80,"completed":true,"started_at":"` + formatTimeForTest(time.Now()) + `","completed_at":"` + formatTimeForTest(time.Now()) + `"}`,
		},
		{
			name:           "Invalid Request Body",
			quizSuiteID:    "1",
			requestBody:    `{"score":"invalid"}`,
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"json: cannot unmarshal string into Go struct field CreateQuizAttemptRequest.score of type int"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/quiz-suites/"+tt.quizSuiteID+"/attempts", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			
			// For success case, we need to check that the response contains the expected fields
			// but ignore additional fields like created_at and updated_at
			if tt.expectedStatus == http.StatusCreated {
				var expectedData, actualData map[string]interface{}
				err := json.Unmarshal([]byte(tt.expectedBody), &expectedData)
				assert.NoError(t, err)
				
				err = json.Unmarshal(w.Body.Bytes(), &actualData)
				assert.NoError(t, err)
				
				// Check that all expected fields are present with correct values
				for key, value := range expectedData {
					// For timestamp fields, we need to compare only the date part
					if key == "started_at" || key == "completed_at" {
						expectedTime, _ := time.Parse(time.RFC3339, value.(string))
						actualTime, _ := time.Parse(time.RFC3339, actualData[key].(string))
						assert.Equal(t, formatTimeForTest(expectedTime), formatTimeForTest(actualTime), "Field %s should match", key)
					} else {
						assert.Equal(t, value, actualData[key], "Field %s should match", key)
					}
				}
			} else {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestGetQuizAttempt(t *testing.T) {
	router, mockService := setupTestRouter()
	
	// Create a handler with the mock service
	handler := NewQuizAttemptHandler(mockService)
	
	router.GET("/quiz-suites/:id/attempts/:attemptId", func(c *gin.Context) {
		c.Set("userID", int64(1))
		handler.GetQuizAttempt(c)
	})

	tests := []struct {
		name           string
		quizSuiteID    string
		attemptID      string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "Success",
			quizSuiteID: "1",
			attemptID:   "1",
			setupMock: func() {
				now := time.Now()
				attempt := &quiz_attempt.QuizAttempt{
					ID:          1,
					UserID:      1,
					QuizSuiteID: 1,
					Score:       80,
					Completed:   true,
					StartedAt:   now,
					CompletedAt: &now,
				}
				mockService.On("Get", mock.Anything, int64(1), int64(1)).Return(attempt, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"user_id":1,"quiz_suite_id":1,"score":80,"completed":true,"started_at":"` + formatTimeForTest(time.Now()) + `","completed_at":"` + formatTimeForTest(time.Now()) + `"}`,
		},
		{
			name:           "Not Found",
			quizSuiteID:    "1",
			attemptID:      "1",
			setupMock:      func() {
				mockService.On("Get", mock.Anything, int64(1), int64(1)).Return(nil, service.ErrQuizAttemptNotFound).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"quiz attempt not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/quiz-suites/"+tt.quizSuiteID+"/attempts/"+tt.attemptID, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			
			// For success case, we need to check that the response contains the expected fields
			// but ignore additional fields like created_at and updated_at
			if tt.expectedStatus == http.StatusOK {
				var expectedData, actualData map[string]interface{}
				err := json.Unmarshal([]byte(tt.expectedBody), &expectedData)
				assert.NoError(t, err)
				
				err = json.Unmarshal(w.Body.Bytes(), &actualData)
				assert.NoError(t, err)
				
				// Check that all expected fields are present with correct values
				for key, value := range expectedData {
					// For timestamp fields, we need to compare only the date part
					if key == "started_at" || key == "completed_at" {
						expectedTime, _ := time.Parse(time.RFC3339, value.(string))
						actualTime, _ := time.Parse(time.RFC3339, actualData[key].(string))
						assert.Equal(t, formatTimeForTest(expectedTime), formatTimeForTest(actualTime), "Field %s should match", key)
					} else {
						assert.Equal(t, value, actualData[key], "Field %s should match", key)
					}
				}
			} else {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestUpdateQuizAttempt(t *testing.T) {
	router, mockService := setupTestRouter()
	
	// Create a handler with the mock service
	handler := NewQuizAttemptHandler(mockService)
	
	router.PUT("/quiz-suites/:id/attempts/:attemptId", func(c *gin.Context) {
		c.Set("userID", int64(1))
		handler.UpdateQuizAttempt(c)
	})

	tests := []struct {
		name           string
		quizSuiteID    string
		attemptID      string
		requestBody    string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "Success",
			quizSuiteID: "1",
			attemptID:   "1",
			requestBody: `{"score":90,"completed":true}`,
			setupMock: func() {
				now := time.Now()
				attempt := &quiz_attempt.QuizAttempt{
					ID:          1,
					UserID:      1,
					QuizSuiteID: 1,
					Score:       90,
					Completed:   true,
					StartedAt:   now,
					CompletedAt: &now,
				}
				score := 90
				completed := true
				mockService.On("Update", mock.Anything, int64(1), int64(1), quiz_attempt.UpdateQuizAttemptRequest{
					Score:     &score,
					Completed: &completed,
				}).Return(attempt, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"user_id":1,"quiz_suite_id":1,"score":90,"completed":true,"started_at":"` + formatTimeForTest(time.Now()) + `","completed_at":"` + formatTimeForTest(time.Now()) + `"}`,
		},
		{
			name:           "Not Found",
			quizSuiteID:    "1",
			attemptID:      "1",
			requestBody:    `{"score":90,"completed":true}`,
			setupMock:      func() {
				score := 90
				completed := true
				mockService.On("Update", mock.Anything, int64(1), int64(1), quiz_attempt.UpdateQuizAttemptRequest{
					Score:     &score,
					Completed: &completed,
				}).Return(nil, service.ErrQuizAttemptNotFound).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"quiz attempt not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/quiz-suites/"+tt.quizSuiteID+"/attempts/"+tt.attemptID, bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			
			// For success case, we need to check that the response contains the expected fields
			// but ignore additional fields like created_at and updated_at
			if tt.expectedStatus == http.StatusOK {
				var expectedData, actualData map[string]interface{}
				err := json.Unmarshal([]byte(tt.expectedBody), &expectedData)
				assert.NoError(t, err)
				
				err = json.Unmarshal(w.Body.Bytes(), &actualData)
				assert.NoError(t, err)
				
				// Check that all expected fields are present with correct values
				for key, value := range expectedData {
					// For timestamp fields, we need to compare only the date part
					if key == "started_at" || key == "completed_at" {
						expectedTime, _ := time.Parse(time.RFC3339, value.(string))
						actualTime, _ := time.Parse(time.RFC3339, actualData[key].(string))
						assert.Equal(t, formatTimeForTest(expectedTime), formatTimeForTest(actualTime), "Field %s should match", key)
					} else {
						assert.Equal(t, value, actualData[key], "Field %s should match", key)
					}
				}
			} else {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestDeleteQuizAttempt(t *testing.T) {
	router, mockService := setupTestRouter()
	
	// Create a handler with the mock service
	handler := NewQuizAttemptHandler(mockService)
	
	router.DELETE("/quiz-suites/:id/attempts/:attemptId", func(c *gin.Context) {
		c.Set("userID", int64(1))
		handler.DeleteQuizAttempt(c)
	})

	tests := []struct {
		name           string
		quizSuiteID    string
		attemptID      string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Success",
			quizSuiteID:    "1",
			attemptID:      "1",
			setupMock:      func() { mockService.On("Delete", mock.Anything, int64(1), int64(1)).Return(nil).Once() },
			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
		},
		{
			name:           "Not Found",
			quizSuiteID:    "1",
			attemptID:      "1",
			setupMock:      func() { mockService.On("Delete", mock.Anything, int64(1), int64(1)).Return(service.ErrQuizAttemptNotFound).Once() },
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"quiz attempt not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/quiz-suites/"+tt.quizSuiteID+"/attempts/"+tt.attemptID, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			} else {
				assert.Empty(t, w.Body.String())
			}
		})
	}
} 