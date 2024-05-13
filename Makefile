.PHONY: buildrun
buildrun:
	docker-compose build
	docker-compose up -d

.PHONY: stop
stop:
	docker-compose down

.PHONY: test
test:
	go test ./test/...
