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
