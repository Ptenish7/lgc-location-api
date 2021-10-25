.PHONY: build
build:
	go build -o bin/lgc-location-api cmd/lgc-location-api/main.go

.PHONY: run
run:
	go run cmd/lgc-location-api/main.go

.PHONY: test
test:
	go test -v ./...
