package repository

import (
	"context"
	"time"

	"quizlet/internal/models/quiz_attempt"
	"gorm.io/gorm"
)

type QuizAttemptRepository struct {
	db *gorm.DB
}

func NewQuizAttemptRepository(db *gorm.DB) *QuizAttemptRepository {
	return &QuizAttemptRepository{db: db}
}

func (r *QuizAttemptRepository) ListByQuizSuite(ctx context.Context, quizSuiteID, userID int64) ([]quiz_attempt.QuizAttempt, error) {
	var attempts []quiz_attempt.QuizAttempt
	err := r.db.WithContext(ctx).
		Where("quiz_suite_id = ? AND user_id = ?", quizSuiteID, userID).
		Find(&attempts).Error
	return attempts, err
}

func (r *QuizAttemptRepository) Create(ctx context.Context, attempt *quiz_attempt.QuizAttempt) (*quiz_attempt.QuizAttempt, error) {
	err := r.db.WithContext(ctx).Create(attempt).Error
	if err != nil {
		return nil, err
	}
	return attempt, nil
}

func (r *QuizAttemptRepository) Get(ctx context.Context, id int64) (*quiz_attempt.QuizAttempt, error) {
	var attempt quiz_attempt.QuizAttempt
	err := r.db.WithContext(ctx).First(&attempt, id).Error
	if err != nil {
		return nil, err
	}
	return &attempt, nil
}

func (r *QuizAttemptRepository) Update(ctx context.Context, attempt *quiz_attempt.QuizAttempt) (*quiz_attempt.QuizAttempt, error) {
	attempt.UpdatedAt = time.Now()
	err := r.db.WithContext(ctx).Save(attempt).Error
	if err != nil {
		return nil, err
	}
	return attempt, nil
}

func (r *QuizAttemptRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&quiz_attempt.QuizAttempt{}, id).Error
} 