BINARY_NAME=gcopy

build:
	go generate
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -ldflags="-s -w" -v -o ./bin/${BINARY_NAME}-darwin-amd64 main.go
	CGO_ENABLED=0 GOARCH=arm64 GOOS=darwin go build -ldflags="-s -w" -v -o ./bin/${BINARY_NAME}-darwin-arm64 main.go
	lipo -create -output bin/${BINARY_NAME} bin/${BINARY_NAME}-darwin-amd64 bin/${BINARY_NAME}-darwin-arm64

run:
	./bin/${BINARY_NAME}

build_and_run: build run

mod:
	go mod tidy
	go mod vendor

gosec:
	gosec -severity medium ./...

lint:
	gofmt -s -w .
	prettier -w .

clean:
	go clean
	rm -rf ./bin/*
