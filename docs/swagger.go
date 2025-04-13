package docs

import "github.com/swaggo/swag"

// @title           Quizlet API
// @version         1.0
// @description     API for managing quizzes, quiz suites, and users
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @tag.name users
// @tag.description User management endpoints

// @tag.name quizzes
// @tag.description Quiz management endpoints

// @tag.name quiz-suites
// @tag.description Quiz suite management endpoints 