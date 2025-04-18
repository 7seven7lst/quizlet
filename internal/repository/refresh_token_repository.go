package repository

import (
	"quizlet/internal/models/user"
	"gorm.io/gorm"
	"time"
)

type RefreshTokenRepository interface {
	Create(token *user.RefreshToken) error
	FindByToken(token string) (*user.RefreshToken, error)
	FindByUserID(userID uint) ([]*user.RefreshToken, error)
	Revoke(token string) error
	DeleteExpired() error
}

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{
		db: db,
	}
}

func (r *refreshTokenRepository) Create(token *user.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *refreshTokenRepository) FindByToken(token string) (*user.RefreshToken, error) {
	var refreshToken user.RefreshToken
	err := r.db.Where("token = ? AND revoked = ?", token, false).First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (r *refreshTokenRepository) FindByUserID(userID uint) ([]*user.RefreshToken, error) {
	var tokens []*user.RefreshToken
	err := r.db.Where("user_id = ? AND revoked = ?", userID, false).Find(&tokens).Error
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (r *refreshTokenRepository) Revoke(token string) error {
	return r.db.Model(&user.RefreshToken{}).Where("token = ?", token).Update("revoked", true).Error
}

func (r *refreshTokenRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ? OR revoked = ?", time.Now(), true).Delete(&user.RefreshToken{}).Error
} 