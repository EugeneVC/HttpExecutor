module httpexecutor

go 1.13

replace (
	common => ./internal/common
	models => ./internal/models
	repository => ./internal/repository
)

require (
	common v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.3.0
	models v0.0.0-00010101000000-000000000000
	repository v0.0.0-00010101000000-000000000000
)
