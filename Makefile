.PHONY: buildrun
buildrun:
	docker-compose build
	docker-compose up -d

.PHONY: stop
stop:
	docker-compose down

.PHONY: genMock
genMock:
	mockgen -source=internal/auth/repository.go \
	-destination=internal/auth/mocks/mock_repository.go

.PHONY: genSwagger
genSwagger:
	swag init -parseDependency -g cmd/api/main.go

.PHONY: test
test:
	go test ./test/...
