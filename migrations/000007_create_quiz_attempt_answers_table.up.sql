CREATE TABLE IF NOT EXISTS quiz_attempt_answers (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    quiz_attempt_id INTEGER REFERENCES quiz_attempts(id) ON DELETE CASCADE,
    quiz_id INTEGER REFERENCES quizzes(id) ON DELETE CASCADE,
    user_answer TEXT NOT NULL,
    is_correct BOOLEAN NOT NULL
);

CREATE INDEX idx_quiz_attempt_answers_quiz_attempt_id ON quiz_attempt_answers(quiz_attempt_id);
CREATE INDEX idx_quiz_attempt_answers_quiz_id ON quiz_attempt_answers(quiz_id); 