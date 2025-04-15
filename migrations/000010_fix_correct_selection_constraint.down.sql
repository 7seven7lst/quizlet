-- Drop the foreign key constraint
ALTER TABLE quizzes DROP CONSTRAINT IF EXISTS quizzes_correct_selection_id_fkey;

-- Add back the original foreign key constraint without ON DELETE SET NULL
ALTER TABLE quizzes ADD CONSTRAINT quizzes_correct_selection_id_fkey 
    FOREIGN KEY (correct_selection_id) 
    REFERENCES quiz_selections(id); 