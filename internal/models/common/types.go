package common

type QuizType string

const (
	QuizTypeSingleChoice QuizType = "single_choice"
	QuizTypeMultiChoice  QuizType = "multi_choice"
	QuizTypeTrueFalse    QuizType = "true_false"
) 