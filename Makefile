VERSION = $(shell git describe --tags --dirty --always)
LDFLAGS = "-s -w -X github.com/rapidsai/rvc/pkg/version.version=$(VERSION)"

all: build

build: test build_serverless build_cli

build_serverless:
	@mkdir -p bin
	GOARCH=amd64 GOOS=linux go build -mod=vendor -ldflags=$(LDFLAGS) -o bin/bootstrap ./cmd/rvc_serverless/...

build_cli:
	@mkdir -p bin
	GOARCH=amd64 GOOS=linux go build -mod=vendor -ldflags=$(LDFLAGS) -o bin/rvc ./cmd/rvc_cli/...

test:
	go test ./pkg/...

fmt:
	@gofmt -w -l $$(find pkg/ cmd/ -name '*.go')

coverage:
	go test -coverprofile=coverage.out ./pkg/...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out
	@rm coverage.out

clean:
	@rm -r bin/
