
.PHONY: run
run:
	go run cmd/agent/main.go

lint:
	golangci-lint run