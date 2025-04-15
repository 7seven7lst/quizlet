-- Drop the existing foreign key constraint if it exists
ALTER TABLE quizzes DROP CONSTRAINT IF EXISTS quizzes_correct_selection_id_fkey;

-- Add the foreign key constraint with ON DELETE SET NULL
ALTER TABLE quizzes ADD CONSTRAINT quizzes_correct_selection_id_fkey 
    FOREIGN KEY (correct_selection_id) 
    REFERENCES quiz_selections(id) 
    ON DELETE SET NULL; 