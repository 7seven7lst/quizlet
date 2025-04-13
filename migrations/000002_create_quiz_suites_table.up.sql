CREATE TABLE IF NOT EXISTS quiz_suites (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    created_by_id INTEGER REFERENCES users(id)
);

CREATE INDEX idx_quiz_suites_deleted_at ON quiz_suites(deleted_at);
CREATE INDEX idx_quiz_suites_created_by_id ON quiz_suites(created_by_id); 