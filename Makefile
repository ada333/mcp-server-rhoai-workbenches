lint:
	golangci-lint run

test:
	go test -v ./tools/... ./resources/... ./prompts/...

build:
	make lint
	make test
	go build -o mcp-server-rhoai
