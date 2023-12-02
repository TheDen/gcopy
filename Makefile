.DEFAULT_GOAL := build
BINARY_NAME=gcopy

build: format
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
	go mod verify

go-update:
	go list -mod=readonly -m -f '{{if not .Indirect}}{{if not .Main}}{{.Path}}{{end}}{{end}}' all | xargs go get -u
	$(MAKE) mod

gosec:
	gosec -severity medium ./...

golines-format:
	# https://github.com/segmentio/golines
	@printf "%s\n" "==== Run golines ====="
	golines --write-output --ignored-dirs=vendor .

go-staticcheck:
	# https://github.com/dominikh/go-tools
	staticcheck ./...

format:
	gofmt -s -w *.go
	prettier -w .

clean:
	go clean
	rm -rf ./bin/*


