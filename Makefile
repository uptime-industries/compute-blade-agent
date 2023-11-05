FUZZ_TARGETS := ./pkg/smartfanunit/proto

all: lint test

.PHONY: run
run:
	go run cmd/agent/main.go

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test ./... -v


.PHONY: fuzz
fuzz:
	@for target in $(FUZZ_TARGETS); do \
		go test  -fuzz="Fuzz" -fuzztime=5s -fuzzminimizetime=10s  $$target; \
	done


.PHONY: generate
generate: buf
	$(BUF) generate

release:
	goreleaser release --clean

snapshot:
	goreleaser release --snapshot --skip=publish --clean

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
