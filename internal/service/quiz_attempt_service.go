package service

import (
	"context"
	"errors"
	"time"

	"quizlet/internal/models/quiz_attempt"
	"quizlet/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrQuizAttemptNotFound = errors.New("quiz attempt not found")
	ErrUnauthorized        = errors.New("unauthorized access")
)

// QuizAttemptService defines the interface for quiz attempt operations
type QuizAttemptService interface {
	ListByQuizSuite(ctx context.Context, quizSuiteID, userID int64) ([]quiz_attempt.QuizAttempt, error)
	Create(ctx context.Context, quizSuiteID, userID int64, req quiz_attempt.CreateQuizAttemptRequest) (*quiz_attempt.QuizAttempt, error)
	Get(ctx context.Context, id, userID int64) (*quiz_attempt.QuizAttempt, error)
	Update(ctx context.Context, id, userID int64, req quiz_attempt.UpdateQuizAttemptRequest) (*quiz_attempt.QuizAttempt, error)
	Delete(ctx context.Context, id, userID int64) error
}

// QuizAttemptServiceImpl is the concrete implementation of QuizAttemptService
type QuizAttemptServiceImpl struct {
	repo *repository.QuizAttemptRepository
}

// Ensure QuizAttemptServiceImpl implements QuizAttemptService
var _ QuizAttemptService = (*QuizAttemptServiceImpl)(nil)

func NewQuizAttemptService(repo *repository.QuizAttemptRepository) QuizAttemptService {
	return &QuizAttemptServiceImpl{
		repo: repo,
	}
}

func (s *QuizAttemptServiceImpl) ListByQuizSuite(ctx context.Context, quizSuiteID, userID int64) ([]quiz_attempt.QuizAttempt, error) {
	return s.repo.ListByQuizSuite(ctx, quizSuiteID, userID)
}

func (s *QuizAttemptServiceImpl) Create(ctx context.Context, quizSuiteID, userID int64, req quiz_attempt.CreateQuizAttemptRequest) (*quiz_attempt.QuizAttempt, error) {
	now := time.Now()
	attempt := &quiz_attempt.QuizAttempt{
		UserID:      userID,
		QuizSuiteID: quizSuiteID,
		Score:       req.Score,
		Completed:   req.Completed,
		StartedAt:   now,
	}

	if req.Completed {
		attempt.CompletedAt = &now
	}

	return s.repo.Create(ctx, attempt)
}

func (s *QuizAttemptServiceImpl) Get(ctx context.Context, id, userID int64) (*quiz_attempt.QuizAttempt, error) {
	attempt, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrQuizAttemptNotFound
		}
		return nil, err
	}

	if attempt.UserID != userID {
		return nil, ErrUnauthorized
	}

	return attempt, nil
}

func (s *QuizAttemptServiceImpl) Update(ctx context.Context, id, userID int64, req quiz_attempt.UpdateQuizAttemptRequest) (*quiz_attempt.QuizAttempt, error) {
	attempt, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrQuizAttemptNotFound
		}
		return nil, err
	}

	if attempt.UserID != userID {
		return nil, ErrUnauthorized
	}

	if req.Score != nil {
		attempt.Score = *req.Score
	}

	if req.Completed != nil {
		attempt.Completed = *req.Completed
		if *req.Completed && attempt.CompletedAt == nil {
			now := time.Now()
			attempt.CompletedAt = &now
		}
	}

	return s.repo.Update(ctx, attempt)
}

func (s *QuizAttemptServiceImpl) Delete(ctx context.Context, id, userID int64) error {
	attempt, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrQuizAttemptNotFound
		}
		return err
	}

	if attempt.UserID != userID {
		return ErrUnauthorized
	}

	return s.repo.Delete(ctx, id)
} 