-- Add back the correct_selection_id column
ALTER TABLE quizzes ADD COLUMN correct_selection_id BIGINT;

-- Add back the foreign key constraint
ALTER TABLE quizzes ADD CONSTRAINT quizzes_correct_selection_id_fkey 
    FOREIGN KEY (correct_selection_id) 
    REFERENCES quiz_selections(id) 
    ON DELETE SET NULL; 