CREATE TABLE IF NOT EXISTS quiz_suite_quizzes (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    quiz_suite_id INTEGER REFERENCES quiz_suites(id) ON DELETE CASCADE,
    quiz_id INTEGER REFERENCES quizzes(id) ON DELETE CASCADE,
    UNIQUE(quiz_suite_id, quiz_id)
);

CREATE INDEX idx_quiz_suite_quizzes_quiz_suite_id ON quiz_suite_quizzes(quiz_suite_id);
CREATE INDEX idx_quiz_suite_quizzes_quiz_id ON quiz_suite_quizzes(quiz_id); 