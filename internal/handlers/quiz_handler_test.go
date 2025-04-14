package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"quizlet/internal/models/quiz"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"quizlet/tests/mocks"
	"gorm.io/gorm"
)

type MockQuizService struct {
	mock.Mock
}

func (m *MockQuizService) CreateQuiz(quiz *quiz.Quiz) error {
	args := m.Called(quiz)
	return args.Error(0)
}

func (m *MockQuizService) GetQuizByID(id uint) (*quiz.Quiz, error) {
	args := m.Called(id)
	return args.Get(0).(*quiz.Quiz), args.Error(1)
}

func (m *MockQuizService) GetQuizzesByUserID(userID uint) ([]*quiz.Quiz, error) {
	args := m.Called(userID)
	return args.Get(0).([]*quiz.Quiz), args.Error(1)
}

func (m *MockQuizService) UpdateQuiz(quiz *quiz.Quiz) error {
	args := m.Called(quiz)
	return args.Error(0)
}

func (m *MockQuizService) DeleteQuiz(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockQuizService) AddSelection(quizID uint, selection quiz.QuizSelection) error {
	args := m.Called(quizID, selection)
	return args.Error(0)
}

func (m *MockQuizService) RemoveSelection(quizID uint, selectionID uint) error {
	args := m.Called(quizID, selectionID)
	return args.Error(0)
}

func TestCreateQuiz(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockQuizService := new(MockQuizService)
	handler := NewQuizHandler(mockQuizService)

	testCases := []struct {
		name           string
		input          quiz.Quiz
		userID         uint
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			input: quiz.Quiz{
				Question:      "Test Question",
				QuizType:      quiz.QuizTypeMultiChoice,
				CorrectAnswer: "Test Answer",
			},
			userID: 1,
			mockSetup: func() {
				mockQuizService.On("CreateQuiz", mock.MatchedBy(func(q *quiz.Quiz) bool {
					return q.Question == "Test Question" &&
						q.QuizType == quiz.QuizTypeMultiChoice &&
						q.CorrectAnswer == "Test Answer" &&
						q.CreatedByID == uint(1)
				})).Return(nil).Once()
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":0,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","question":"Test Question","quiz_type":"multi_choice","correct_answer":"Test Answer","created_by_id":1}`,
		},
		{
			name:           "Unauthorized",
			input:          quiz.Quiz{},
			userID:         0,
			mockSetup:      func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"unauthorized"}`,
		},
		{
			name:   "Service Error",
			input: quiz.Quiz{
				Question:      "Test Question",
				QuizType:      quiz.QuizTypeMultiChoice,
				CorrectAnswer: "Test Answer",
			},
			userID: 1,
			mockSetup: func() {
				mockQuizService.On("CreateQuiz", mock.MatchedBy(func(q *quiz.Quiz) bool {
					return q.Question == "Test Question" &&
						q.QuizType == quiz.QuizTypeMultiChoice &&
						q.CorrectAnswer == "Test Answer" &&
						q.CreatedByID == uint(1)
				})).Return(gorm.ErrInvalidDB).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"invalid db"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tc.input)
			c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer(body))

			if tc.userID > 0 {
				c.Set("userID", tc.userID)
			}

			tc.mockSetup()

			handler.CreateQuiz(c)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())

			mockQuizService.AssertExpectations(t)
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
					Return(&quiz.Quiz{
						Question:      "What is the capital of France?",
						QuizType:      quiz.QuizTypeSingleChoice,
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