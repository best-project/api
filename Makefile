APP_NAME = api

.PHONY: build
build: fmt vet
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api main.go
	docker build -t $(APP_NAME) .

.PHONY: run
run:
	docker run -p 8080:8080 $(APP_NAME)

# Run go fmt against code
.PHONY: format
format: fmt vet

# Run go fmt against code
.PHONY: fmt
fmt:
	go fmt ./...

# Run go vet against code
.PHONY: vet
vet:
	go vet ./...
