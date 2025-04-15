-- Drop the foreign key constraint first
ALTER TABLE quizzes DROP CONSTRAINT IF EXISTS quizzes_correct_selection_id_fkey;

-- Drop the correct_selection_id column
ALTER TABLE quizzes DROP COLUMN IF EXISTS correct_selection_id; 