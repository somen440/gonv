dc := docker-compose -f build/docker-compose.yaml -p gonv

.PHONY: up
up:
	${dc} up -d mysql57

.PHONY: down
down:
	${dc} down

.PHONY: clean
clean:
	${dc} down --volumes

.PHONY: ps
ps:
	${dc} ps

.PHONY: wait
wait:
	${dc} run --rm wait

.PHONY: exec-db
exec-db:
	${dc} exec mysql57 mysql -uroot -ptest test

.PHONY: logs-db
logs-db:
	${dc} logs mysql57

.PHONY: reflect
reflect:
	go run . reflect -u root -p test -P 33066 -o build/mysql/schema test

.PHONY: diff
diff:
	go run . diff -u root -p test -P 33066 test build/mysql/schema

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
