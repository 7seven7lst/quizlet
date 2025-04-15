-- Add selection_display_name column to quiz_selections
ALTER TABLE quiz_selections ADD COLUMN selection_display_name VARCHAR(10);

-- Add correct_selection_id column to quizzes
ALTER TABLE quizzes ADD COLUMN correct_selection_id INTEGER REFERENCES quiz_selections(id);

-- Update correct_selection_id based on existing correct_answer and is_correct
UPDATE quizzes q
SET correct_selection_id = qs.id
FROM quiz_selections qs
WHERE qs.quiz_id = q.id AND qs.is_correct = TRUE;

-- Drop correct_answer column from quizzes
ALTER TABLE quizzes DROP COLUMN correct_answer; 