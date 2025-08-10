GOPATH_BIN := $(shell go env GOPATH)/bin
PROTO_SRC_DIR := internal/adapter/left/grpc
PROTO_OUT_DIR := $(PROTO_SRC_DIR)/pb
PROTOS := $(wildcard $(PROTO_SRC_DIR)/*.proto)

.PHONY: all proto tidy build run test clean

all: build

proto: $(GOPATH_BIN)/protoc-gen-go $(GOPATH_BIN)/protoc-gen-go-grpc
	@echo "Generating protobuf files into $(PROTO_OUT_DIR)"
	@mkdir -p $(PROTO_OUT_DIR)
	@PATH="$(GOPATH_BIN):$$PATH" protoc -I=$(PROTO_SRC_DIR) \
		--go_out=paths=source_relative:$(PROTO_OUT_DIR) \
		--go-grpc_out=paths=source_relative:$(PROTO_OUT_DIR) \
		$(PROTOS)

$(GOPATH_BIN)/protoc-gen-go:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

$(GOPATH_BIN)/protoc-gen-go-grpc:
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

tidy:
	go mod tidy

build:
	go build ./...

run:
	go run ./cmd/toq_server.go

test:
	go test ./...

clean:
	rm -f $(PROTO_OUT_DIR)/*_pb.go
