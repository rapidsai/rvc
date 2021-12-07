all: build

build: test build_serverless build_cli

build_serverless:
	@mkdir -p bin
	GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin ./cmd/rvc_serverless/...

build_cli:
	@mkdir -p bin
	GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin ./cmd/rvc_cli/...

test:
	go test ./pkg/...

fmt:
	gofmt -w ./pkg/... ./cmd/...
