
all: lint

.PHONY: run
run:
	go run cmd/agent/main.go

.PHONY: lint
lint:
	golangci-lint run

.PHONY: generate
generate: buf
	$(BUF) generate

# Dependencies
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

BUF ?= $(LOCALBIN)/buf
BUF_VERSION ?= v1.25.0

.PHONY: buf
buf: $(BUF)
$(BUF): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION)
