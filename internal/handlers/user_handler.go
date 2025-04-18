package handlers

import (
	"net/http"
	"quizlet/internal/models/user"
	"quizlet/internal/service"
	"strconv"
	"log"

	"github.com/gin-gonic/gin"
	"quizlet/internal/auth"
	"gorm.io/gorm"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	User         user.User `json:"user"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int64     `json:"expires_in"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// @Summary Create a new user
// @Description Create a new user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param user body user.CreateUserRequest true "User information"
// @Success 201 {object} user.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req user.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := &user.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.userService.CreateUser(u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, u)
}

// @Summary Get a user by ID
// @Description Get user information by user ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} user.User
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	u, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Don't send password back in response
	u.Password = ""
	c.JSON(http.StatusOK, u)
}

// @Summary Update a user
// @Description Update user information by user ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body user.User true "User information"
// @Success 200 {object} user.User
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var u user.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u.ID = uint(id)
	if err := h.userService.UpdateUser(&u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Don't send password back in response
	u.Password = ""
	c.JSON(http.StatusOK, u)
}

// @Summary Delete a user
// @Description Delete a user by user ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.userService.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

// @Summary Login user
// @Description Authenticate user with email and password
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Login attempt for email: %s", req.Email)

	u, err := h.userService.ValidatePassword(req.Email, req.Password)
	if err != nil {
		log.Printf("Login failed for email %s: %v", req.Email, err)
		if err == gorm.ErrRecordNotFound || err.Error() == "invalid password" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	log.Printf("Login successful for user: %s", u.Email)

	// Generate access token
	accessToken, err := auth.GenerateAccessToken(u.ID)
	if err != nil {
		log.Printf("Failed to generate access token for user %s: %v", u.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	// Generate refresh token
	refreshToken, err := h.userService.CreateRefreshToken(u.ID)
	if err != nil {
		log.Printf("Failed to generate refresh token for user %s: %v", u.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	// Don't send password back in response
	u.Password = ""
	
	response := LoginResponse{
		User:         *u,
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
		ExpiresIn:    900, // 15 minutes in seconds
	}
	
	c.JSON(http.StatusOK, response)
}

// @Summary Refresh access token
// @Description Get a new access token using a refresh token
// @Tags users
// @Accept json
// @Produce json
// @Param refresh_token body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} RefreshTokenResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /users/refresh [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate refresh token
	refreshToken, err := h.userService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		if err == auth.ErrExpiredToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token has expired"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		}
		return
	}

	// Generate new access token
	accessToken, err := auth.GenerateAccessToken(refreshToken.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	response := RefreshTokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   900, // 15 minutes in seconds
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Logout user
// @Description Revoke the refresh token from cookie
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security BearerAuth
// @Router /users/logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.RevokeRefreshToken(req.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to revoke refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully logged out"})
} 