-- Add correct_answer column back to quizzes
ALTER TABLE quizzes ADD COLUMN correct_answer TEXT;

-- Update correct_answer based on correct_selection_id
UPDATE quizzes q
SET correct_answer = qs.selection_text
FROM quiz_selections qs
WHERE qs.id = q.correct_selection_id;

-- Make correct_answer NOT NULL after setting values
ALTER TABLE quizzes ALTER COLUMN correct_answer SET NOT NULL;

-- Drop correct_selection_id column from quizzes
ALTER TABLE quizzes DROP COLUMN correct_selection_id;

-- Drop selection_display_name column from quiz_selections
ALTER TABLE quiz_selections DROP COLUMN selection_display_name; 