package service

import (
	"errors"
	"log"
	"quizlet/internal/models/user"
	"quizlet/internal/repository"
	"quizlet/internal/auth"
	"time"
)

type UserService interface {
	CreateUser(user *user.User) error
	GetUserByID(id uint) (*user.User, error)
	GetUserByEmail(email string) (*user.User, error)
	UpdateUser(user *user.User) error
	DeleteUser(id uint) error
	ValidatePassword(email, password string) (*user.User, error)
	CreateRefreshToken(userID uint) (*user.RefreshToken, error)
	ValidateRefreshToken(token string) (*user.RefreshToken, error)
	RevokeRefreshToken(token string) error
}

type userService struct {
	userRepo repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
}

func NewUserService(userRepo repository.UserRepository, refreshTokenRepo repository.RefreshTokenRepository) UserService {
	return &userService{
		userRepo: userRepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (s *userService) CreateUser(user *user.User) error {
	// Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(user.Email)
	if err == nil && existingUser != nil {
		return errors.New("user with this email already exists")
	}

	// Hash the password before saving
	if err := user.HashPassword(); err != nil {
		return err
	}

	return s.userRepo.Create(user)
}

func (s *userService) GetUserByID(id uint) (*user.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *userService) GetUserByEmail(email string) (*user.User, error) {
	return s.userRepo.FindByEmail(email)
}

func (s *userService) UpdateUser(user *user.User) error {
	// If password is being updated, hash it
	if user.Password != "" {
		if err := user.HashPassword(); err != nil {
			return err
		}
	}
	return s.userRepo.Update(user)
}

func (s *userService) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}

func (s *userService) ValidatePassword(email, password string) (*user.User, error) {
	log.Printf("Validating password for email: %s (password length: %d)", email, len(password))
	
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		log.Printf("Error finding user by email %s: %v", email, err)
		return nil, err
	}
	if user == nil {
		log.Printf("User not found for email: %s", email)
		return nil, errors.New("user not found")
	}

	log.Printf("Found user with email %s (ID: %d), checking password", email, user.ID)
	log.Printf("Stored password hash length: %d", len(user.Password))
	
	if !user.CheckPassword(password) {
		log.Printf("Invalid password for user: %s (ID: %d)", email, user.ID)
		return nil, errors.New("invalid password")
	}

	log.Printf("Password validation successful for user: %s (ID: %d)", email, user.ID)
	return user, nil
}

func (s *userService) CreateRefreshToken(userID uint) (*user.RefreshToken, error) {
	token, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	refreshToken := &user.RefreshToken{
		Token:     token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour), // 30 days
	}

	if err := s.refreshTokenRepo.Create(refreshToken); err != nil {
		return nil, err
	}

	return refreshToken, nil
}

func (s *userService) ValidateRefreshToken(token string) (*user.RefreshToken, error) {
	refreshToken, err := s.refreshTokenRepo.FindByToken(token)
	if err != nil {
		return nil, err
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return nil, auth.ErrExpiredToken
	}

	return refreshToken, nil
}

func (s *userService) RevokeRefreshToken(token string) error {
	return s.refreshTokenRepo.Revoke(token)
} 