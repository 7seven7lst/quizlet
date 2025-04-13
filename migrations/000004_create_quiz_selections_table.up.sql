CREATE TABLE IF NOT EXISTS quiz_selections (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    quiz_id INTEGER REFERENCES quizzes(id) ON DELETE CASCADE,
    selection_text TEXT NOT NULL,
    is_correct BOOLEAN NOT NULL
);

CREATE INDEX idx_quiz_selections_deleted_at ON quiz_selections(deleted_at);
CREATE INDEX idx_quiz_selections_quiz_id ON quiz_selections(quiz_id); 