package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"quizlet/internal/models/user"
	"quizlet/internal/auth"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(user *user.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) GetUserByID(id uint) (*user.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserService) GetUserByEmail(email string) (*user.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(user *user.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) DeleteUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) ValidatePassword(email, password string) (*user.User, error) {
	args := m.Called(email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserService) CreateRefreshToken(userID uint) (*user.RefreshToken, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.RefreshToken), args.Error(1)
}

func (m *MockUserService) ValidateRefreshToken(token string) (*user.RefreshToken, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.RefreshToken), args.Error(1)
}

func (m *MockUserService) RevokeRefreshToken(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func TestCreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockUserService)

	testCases := []struct {
		name           string
		requestBody    map[string]interface{}
		mockSetup      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "Success",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"email":    "test@example.com",
				"password": "password123",
			},
			mockSetup: func() {
				mockService.On("CreateUser", mock.MatchedBy(func(u *user.User) bool {
					u.ID = 1
					return u.Username == "testuser" &&
						u.Email == "test@example.com" &&
						u.Password == "password123"
				})).Return(nil).Once()
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"id":         float64(1),
				"username":   "testuser",
				"email":      "test@example.com",
				"created_at": "0001-01-01T00:00:00Z",
				"updated_at": "0001-01-01T00:00:00Z",
			},
		},
		{
			name:   "Invalid Request",
			requestBody: map[string]interface{}{
				"username": "",
				"email":    "invalid-email",
				"password": "",
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Key: 'CreateUserRequest.Username' Error:Field validation for 'Username' failed on the 'required' tag\nKey: 'CreateUserRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag\nKey: 'CreateUserRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag",
			},
		},
		{
			name:   "Service Error",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"email":    "test@example.com",
				"password": "password123",
			},
			mockSetup: func() {
				mockService.On("CreateUser", mock.AnythingOfType("*user.User")).Return(gorm.ErrInvalidDB).Once()
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
			body, _ := json.Marshal(tc.requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Set up mock
			tc.mockSetup()

			// Create handler and execute
			handler := NewUserHandler(mockService)
			handler.CreateUser(c)

			// Assert response
			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedBody, response)
		})
	}
}

func TestGetUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockUserService)

	testCases := []struct {
		name           string
		userID         string
		mockSetup      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "Success",
			userID: "1",
			mockSetup: func() {
				mockUser := &user.User{
					ID:        1,
					Username:  "testuser",
					Email:     "test@example.com",
					Password:  "hashedpassword",
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				}
				mockService.On("GetUserByID", uint(1)).Return(mockUser, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":         float64(1),
				"username":   "testuser",
				"email":      "test@example.com",
				"created_at": "0001-01-01T00:00:00Z",
				"updated_at": "0001-01-01T00:00:00Z",
			},
		},
		{
			name:           "Invalid ID",
			userID:         "invalid",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid user id",
			},
		},
		{
			name:   "User Not Found",
			userID: "1",
			mockSetup: func() {
				mockService.On("GetUserByID", uint(1)).Return(nil, gorm.ErrRecordNotFound).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "user not found",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/users/"+tc.userID, nil)
			c.Params = []gin.Param{{Key: "id", Value: tc.userID}}

			tc.mockSetup()

			handler := NewUserHandler(mockService)
			handler.GetUser(c)

			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedBody, response)
		})
	}
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockUserService)

	testCases := []struct {
		name           string
		requestBody    map[string]interface{}
		mockSetup      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "Success",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
			mockSetup: func() {
				mockService.On("ValidatePassword", "test@example.com", "password123").Return(&user.User{
					ID:        1,
					Username:  "testuser",
					Email:     "test@example.com",
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				}, nil).Once()
				
				mockService.On("CreateRefreshToken", uint(1)).Return(&user.RefreshToken{
					Token:     "refresh-token-123",
					UserID:    1,
					ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
				}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"access_token": mock.Anything,
				"refresh_token": "refresh-token-123",
				"expires_in": float64(900),
				"user": map[string]interface{}{
					"id":         float64(1),
					"username":   "testuser",
					"email":      "test@example.com",
					"created_at": "0001-01-01T00:00:00Z",
					"updated_at": "0001-01-01T00:00:00Z",
				},
			},
		},
		{
			name:   "Invalid Request",
			requestBody: map[string]interface{}{
				"email":    "",
				"password": "",
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Key: 'LoginRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag\nKey: 'LoginRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag",
			},
		},
		{
			name:   "Invalid Credentials",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "wrongpassword",
			},
			mockSetup: func() {
				mockService.On("ValidatePassword", "test@example.com", "wrongpassword").Return(nil, gorm.ErrRecordNotFound).Once()
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "record not found",
			},
		},
		{
			name:   "Service Error",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
			mockSetup: func() {
				mockService.On("ValidatePassword", "test@example.com", "password123").Return(nil, gorm.ErrInvalidDB).Once()
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
			body, _ := json.Marshal(tc.requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Set up mock
			tc.mockSetup()

			// Create handler and execute
			handler := NewUserHandler(mockService)
			handler.Login(c)

			// Assert response
			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tc.expectedStatus == http.StatusOK {
				// For successful login, verify token exists but don't check its exact value
				assert.Contains(t, response, "access_token")
				assert.NotEmpty(t, response["access_token"])
				assert.Equal(t, tc.expectedBody["refresh_token"], response["refresh_token"])
				assert.Equal(t, tc.expectedBody["expires_in"], response["expires_in"])
				assert.Equal(t, tc.expectedBody["user"], response["user"])
			} else {
				assert.Equal(t, tc.expectedBody, response)
			}
		})
	}
}

func TestRefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockUserService)

	testCases := []struct {
		name           string
		requestBody    map[string]interface{}
		mockSetup      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "Success",
			requestBody: map[string]interface{}{
				"refresh_token": "valid-refresh-token",
			},
			mockSetup: func() {
				mockService.On("ValidateRefreshToken", "valid-refresh-token").Return(&user.RefreshToken{
					UserID: 1,
				}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"access_token": mock.Anything,
				"expires_in": float64(900),
			},
		},
		{
			name:   "Invalid Request",
			requestBody: map[string]interface{}{
				"refresh_token": "",
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Key: 'RefreshTokenRequest.RefreshToken' Error:Field validation for 'RefreshToken' failed on the 'required' tag",
			},
		},
		{
			name:   "Invalid Token",
			requestBody: map[string]interface{}{
				"refresh_token": "invalid-token",
			},
			mockSetup: func() {
				mockService.On("ValidateRefreshToken", "invalid-token").Return(nil, gorm.ErrRecordNotFound).Once()
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "invalid refresh token",
			},
		},
		{
			name:   "Expired Token",
			requestBody: map[string]interface{}{
				"refresh_token": "expired-token",
			},
			mockSetup: func() {
				mockService.On("ValidateRefreshToken", "expired-token").Return(nil, auth.ErrExpiredToken).Once()
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "refresh token has expired",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set up request
			body, _ := json.Marshal(tc.requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Set up mock
			tc.mockSetup()

			// Create handler and execute
			handler := NewUserHandler(mockService)
			handler.RefreshToken(c)

			// Assert response
			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tc.expectedStatus == http.StatusOK {
				// For successful refresh, verify token exists but don't check its exact value
				assert.Contains(t, response, "access_token")
				assert.NotEmpty(t, response["access_token"])
				assert.Equal(t, tc.expectedBody["expires_in"], response["expires_in"])
			} else {
				assert.Equal(t, tc.expectedBody, response)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockUserService)

	testCases := []struct {
		name           string
		requestBody    map[string]interface{}
		mockSetup      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "Success",
			requestBody: map[string]interface{}{
				"refresh_token": "valid-refresh-token",
			},
			mockSetup: func() {
				mockService.On("RevokeRefreshToken", "valid-refresh-token").Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "successfully logged out",
			},
		},
		{
			name:   "Invalid Request",
			requestBody: map[string]interface{}{
				"refresh_token": "",
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Key: 'RefreshTokenRequest.RefreshToken' Error:Field validation for 'RefreshToken' failed on the 'required' tag",
			},
		},
		{
			name:   "Service Error",
			requestBody: map[string]interface{}{
				"refresh_token": "invalid-token",
			},
			mockSetup: func() {
				mockService.On("RevokeRefreshToken", "invalid-token").Return(gorm.ErrInvalidDB).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "failed to revoke refresh token",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set up request
			body, _ := json.Marshal(tc.requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/logout", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Set up mock
			tc.mockSetup()

			// Create handler and execute
			handler := NewUserHandler(mockService)
			handler.Logout(c)

			// Assert response
			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedBody, response)
		})
	}
} 