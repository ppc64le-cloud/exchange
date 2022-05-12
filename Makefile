GO_GCFLAGS ?= -gcflags=all='-N -l'
GO=GO111MODULE=on go
GO_BUILD_RECIPE=CGO_ENABLED=0 $(GO) build $(GO_GCFLAGS)

OUT_DIR ?= bin

.PHONY: all
all: build

.PHONY: build
build: pac-server

.PHONY: pac-server
pac-server:
	$(GO_BUILD_RECIPE) -o $(OUT_DIR)/pac-server ./cmd/pac-server
