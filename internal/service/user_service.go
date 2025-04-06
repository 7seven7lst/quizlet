package service

import (
	"errors"
	"quizlet/internal/models"
	"quizlet/internal/repository"
)

type UserService interface {
	CreateUser(user *models.User) error
	GetUserByID(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(user *models.User) error {
	// Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(user.Email)
	if err == nil && existingUser != nil {
		return errors.New("user with this email already exists")
	}

	return s.userRepo.Create(user)
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.FindByEmail(email)
}

func (s *userService) UpdateUser(user *models.User) error {
	return s.userRepo.Update(user)
}

func (s *userService) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
} 