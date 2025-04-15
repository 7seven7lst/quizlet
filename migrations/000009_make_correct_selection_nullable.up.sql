-- Drop the existing column and its constraints
ALTER TABLE quizzes DROP COLUMN IF EXISTS correct_selection_id;

-- Add the column back as nullable with the foreign key constraint
ALTER TABLE quizzes ADD COLUMN correct_selection_id INTEGER;
ALTER TABLE quizzes ADD CONSTRAINT quizzes_correct_selection_id_fkey 
    FOREIGN KEY (correct_selection_id) 
    REFERENCES quiz_selections(id) 
    ON DELETE SET NULL; 