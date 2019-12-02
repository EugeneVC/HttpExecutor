module httpexecutor

go 1.13

replace (
	models => ./internal/models
	repository => ./internal/repository
)

require (
	models v0.0.0-00010101000000-000000000000
	repository v0.0.0-00010101000000-000000000000
)
