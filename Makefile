PHONY: gen lint test

gen:
	go generate ./...

lint:
	golangci-lint run ./...

test:
	go test -race -cover -v ./...