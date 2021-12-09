all: build

build: test build_serverless build_cli

build_serverless:
	@mkdir -p bin
	GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/rvc_serverless ./cmd/rvc_serverless/...

build_cli:
	@mkdir -p bin
	GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/rvc ./cmd/rvc_cli/...

test:
	go test ./pkg/...

fmt:
	@gofmt -w -l $$(find . -name '*.go')

coverage:
	go test -coverprofile=coverage.out ./pkg/...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out
	@rm coverage.out
