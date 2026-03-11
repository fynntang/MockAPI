BINARY=mockapi
VERSION?=1.0.0

.PHONY: build run clean test docker

build:
	go build -ldflags="-s -w" -o $(BINARY) ./cmd/mockapi

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY)

test:
	go test ./...

# Cross compilation
build-all:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BINARY)-linux-amd64 ./cmd/mockapi
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o $(BINARY)-linux-arm64 ./cmd/mockapi
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(BINARY)-darwin-amd64 ./cmd/mockapi
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $(BINARY)-darwin-arm64 ./cmd/mockapi
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(BINARY)-windows-amd64.exe ./cmd/mockapi

docker:
	docker build -t mockapi:$(VERSION) .

docker-run:
	docker run -p 8088:8088 mockapi:$(VERSION)
