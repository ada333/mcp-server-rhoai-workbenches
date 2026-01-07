lint:
	golangci-lint run

test:
	go test -v ./tools/... ./resources/... ./prompts/...
	


