APP_NAME = api

.PHONY: build
build:
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o api main.go
	docker build -t $(APP_NAME) .

.PHONY: run
run:
	docker run $(APP_NAME)
