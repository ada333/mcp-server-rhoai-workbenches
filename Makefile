lint:
	golangci-lint run

test:
	go test -v ./tools/... ./resources/... ./prompts/...

build:
	make lint
	make test
	go build -o mcp-server-rhoai

eval:
	npx promptfoo@latest eval -c promptfoo.yaml

eval-view:
	npx promptfoo@latest eval -c promptfoo.yaml && npx promptfoo@latest view
