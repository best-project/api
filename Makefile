APP_NAME = api

.PHONY: build
build:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api main.go
	docker build -t $(APP_NAME) .

.PHONY: run
run:
	docker run -p $(PORT):$(PORT) --env PORT=$(PORT) --env-file=.env $(APP_NAME)

.PHONY: start
start: format build run

.PHONY: format
format: fmt vet

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...
