-- Drop the nullable column and its constraints
ALTER TABLE quizzes DROP CONSTRAINT IF EXISTS quizzes_correct_selection_id_fkey;
ALTER TABLE quizzes DROP COLUMN IF EXISTS correct_selection_id;

-- Add back the NOT NULL column with its constraint
ALTER TABLE quizzes ADD COLUMN correct_selection_id INTEGER NOT NULL;
ALTER TABLE quizzes ADD CONSTRAINT quizzes_correct_selection_id_fkey 
    FOREIGN KEY (correct_selection_id) 
    REFERENCES quiz_selections(id); 