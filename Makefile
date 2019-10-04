APP_NAME = api

.PHONY: build
build:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api main.go
	docker build -t $(APP_NAME) .

.PHONY: run
run:
	docker run -p 8080:8080 $(APP_NAME)
