CREATE TABLE IF NOT EXISTS quizzes (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    question TEXT NOT NULL,
    quiz_type VARCHAR(20) NOT NULL CHECK (quiz_type IN ('single_choice', 'multi_choice', 'true_false')),
    correct_answer TEXT NOT NULL,
    created_by_id INTEGER REFERENCES users(id)
);

CREATE INDEX idx_quizzes_deleted_at ON quizzes(deleted_at);
CREATE INDEX idx_quizzes_created_by_id ON quizzes(created_by_id); 