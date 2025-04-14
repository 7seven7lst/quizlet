//go:generate mockgen -source=../../internal/service/quiz_service.go -destination=quiz_service_mock.go
//go:generate mockgen -source=../../internal/service/user_service.go -destination=user_service_mock.go
//go:generate mockgen -source=../../internal/service/quiz_suite_service.go -destination=quiz_suite_service_mock.go
//go:generate mockgen -source=../../internal/repository/quiz_repository.go -destination=quiz_repository_mock.go
//go:generate mockgen -source=../../internal/repository/user_repository.go -destination=user_repository_mock.go
//go:generate mockgen -source=../../internal/repository/quiz_suite_repository.go -destination=quiz_suite_repository_mock.go

package mocks 