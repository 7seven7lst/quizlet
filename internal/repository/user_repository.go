package repository

import (
	"quizlet/internal/models/user"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *user.User) error
	FindByID(id uint) (*user.User, error)
	FindByEmail(email string) (*user.User, error)
	Update(user *user.User) error
	Delete(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *user.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uint) (*user.User, error) {
	var user user.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*user.User, error) {
	var user user.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *user.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&user.User{}, id).Error
} 