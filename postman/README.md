# Quizlet API Postman Collection

This folder contains Postman collection and environment files for testing the Quizlet API.

## Files

- `quizlet_api_collection.json`: The Postman collection containing all API endpoints
- `quizlet_api_environment.json`: The Postman environment with variables for the API

## How to Use

### Importing the Collection and Environment

1. Open Postman
2. Click on "Import" in the top left
3. Import both the collection and environment files
4. Select the "Quizlet API Environment" from the environment dropdown in the top right

### Setting Up Authentication

1. After logging in, you'll receive a JWT token
2. Update the `token` variable in the environment with your JWT token
3. All authenticated requests will automatically use this token

### Testing Flow

1. **User Management**:
   - Create a user using the "Create User" request
   - Log in using the "Login" request to get a token
   - Update the `token` variable with the JWT token from the login response

2. **Quiz Management**:
   - Create a quiz using the "Create Quiz" request
   - Add selections to the quiz using the "Add Selection" request
   - Get your quizzes using the "Get User Quizzes" request

3. **Quiz Suite Management**:
   - Create a quiz suite using the "Create Quiz Suite" request
   - Add quizzes to the suite using the "Add Quiz to Suite" request
   - Get your quiz suites using the "Get User Quiz Suites" request

### Example: Creating a Flashcard

1. Create a quiz with type "flashcard":
   ```
   POST /api/quizzes
   {
     "question": "What is the capital of France?",
     "quiz_type": "flashcard",
     "correct_answer": "Paris"
   }
   ```

2. Add a selection for the answer:
   ```
   POST /api/quizzes/{quiz_id}/selections
   {
     "selection_text": "Paris",
     "is_correct": true
   }
   ```

## API Endpoints

### Users
- `POST /api/users` - Create a new user
- `GET /api/users/:id` - Get a user by ID
- `PUT /api/users/:id` - Update a user
- `DELETE /api/users/:id` - Delete a user
- `POST /api/users/login` - Log in a user

### Quizzes
- `POST /api/quizzes` - Create a new quiz
- `GET /api/quizzes/:id` - Get a quiz by ID
- `PUT /api/quizzes/:id` - Update a quiz
- `DELETE /api/quizzes/:id` - Delete a quiz
- `GET /api/quizzes/user` - Get all quizzes for the current user
- `POST /api/quizzes/:id/selections` - Add a selection to a quiz
- `DELETE /api/quizzes/:id/selections/:selectionId` - Remove a selection from a quiz

### Quiz Suites
- `POST /api/quiz-suites` - Create a new quiz suite
- `GET /api/quiz-suites/:id` - Get a quiz suite by ID
- `PUT /api/quiz-suites/:id` - Update a quiz suite
- `DELETE /api/quiz-suites/:id` - Delete a quiz suite
- `GET /api/quiz-suites/user` - Get all quiz suites for the current user
- `POST /api/quiz-suites/:id/quizzes/:quizId` - Add a quiz to a suite
- `DELETE /api/quiz-suites/:id/quizzes/:quizId` - Remove a quiz from a suite 