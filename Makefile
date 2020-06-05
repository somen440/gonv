.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	docker run --rm \
		-v ${PWD}:/app \
		-w /app -it \
		golangci/golangci-lint:v1.26.0 \
		golangci-lint run ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: up
up:
	docker-compose up -d mysql57

.PHONY: down
down:
	docker-compose down

.PHONY: exec-db
exec-db:
	docker-compose exec mysql57 mysql -uroot -ptest test
