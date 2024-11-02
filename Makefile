.PHONY: mock
mock:
	@mockgen -source=webook/internal/service/user.go -package=svcmocks -destination=webook/internal/service/mocks/user.mock.go
	@mockgen -source=webook/internal/repository/user.go -package=repomocks -destination=webook/internal/repository/mocks/user.mock.go
	@go mod tidy