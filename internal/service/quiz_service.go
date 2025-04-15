package service

import (
	"quizlet/internal/models/quiz"
	"quizlet/internal/repository"
)

// QuizSelection is a type alias for quiz.QuizSelection to ensure type compatibility
type QuizSelection = quiz.QuizSelection

type QuizService interface {
	CreateQuiz(quiz *quiz.Quiz) error
	GetQuizByID(id uint) (*quiz.Quiz, error)
	GetQuizzesByUserID(userID uint) ([]*quiz.Quiz, error)
	UpdateQuiz(quiz *quiz.Quiz) error
	DeleteQuiz(id uint) error
	AddSelection(quizID uint, selection QuizSelection) error
	RemoveSelection(quizID uint, selectionID uint) error
}

type quizService struct {
	quizRepo repository.QuizRepository
}

func NewQuizService(quizRepo repository.QuizRepository) QuizService {
	return &quizService{
		quizRepo: quizRepo,
	}
}

func (s *quizService) CreateQuiz(quiz *quiz.Quiz) error {
	return s.quizRepo.Create(quiz)
}

func (s *quizService) GetQuizByID(id uint) (*quiz.Quiz, error) {
	return s.quizRepo.FindByID(id)
}

func (s *quizService) GetQuizzesByUserID(userID uint) ([]*quiz.Quiz, error) {
	return s.quizRepo.FindByUserID(userID)
}

func (s *quizService) UpdateQuiz(quiz *quiz.Quiz) error {
	existing, err := s.quizRepo.FindByID(quiz.ID)
	if err != nil {
		return err
	}

	existing.Question = quiz.Question
	existing.QuizType = quiz.QuizType
	return s.quizRepo.Update(existing)
}

func (s *quizService) DeleteQuiz(id uint) error {
	return s.quizRepo.Delete(id)
}

func (s *quizService) AddSelection(quizID uint, selection QuizSelection) error {
	// Verify the quiz exists
	quiz, err := s.quizRepo.FindByID(quizID)
	if err != nil {
		return err
	}

	// Add selection to quiz
	if quiz.Selections == nil {
		quiz.Selections = make([]QuizSelection, 0)
	}
	quiz.Selections = append(quiz.Selections, selection)
	return s.quizRepo.Update(quiz)
}

func (s *quizService) RemoveSelection(quizID uint, selectionID uint) error {
	// Verify the quiz exists
	quiz, err := s.quizRepo.FindByID(quizID)
	if err != nil {
		return err
	}

	// Remove selection from quiz
	for i, selection := range quiz.Selections {
		if selection.ID == selectionID {
			quiz.Selections = append(quiz.Selections[:i], quiz.Selections[i+1:]...)
			break
		}
	}

	return s.quizRepo.Update(quiz)
} 