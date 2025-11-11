SHELL := bash

COMPOSE ?= docker compose

GEN_DIR           := gen
HANDLERS_DIR      := $(GEN_DIR)/openapi
REPO_DIR          := $(GEN_DIR)/mock
GEN_HANDLERS_STAMP:= $(HANDLERS_DIR)/.generated
GEN_REPO_STAMP    := $(REPO_DIR)/.generated
GEN_STAMPS        := $(GEN_HANDLERS_STAMP) $(GEN_REPO_STAMP)

GO_SOURCES := $(shell find cmd internal -type f -name '*.go' 2>/dev/null)
MOD_FILES  := go.mod go.sum
BUILD_CTX  := Dockerfile docker-compose.yml
ALL_INPUTS := $(GO_SOURCES) $(MOD_FILES) $(BUILD_CTX) $(GEN_STAMPS)

BUILD_ARGS ?=
UP_FLAGS   ?=
TEST_FLAGS ?= -race -v
FMT_TOOLS  ?= gofumpt goimports golines gofmt

CACHE_DIR    := .cache
BUILD_STAMP  := $(CACHE_DIR)/build.stamp

.PHONY: all
all: up

.PHONY: gen
gen: $(GEN_STAMPS)

$(GEN_REPO_STAMP): $(shell find internal/repo -type f -name '*.go' 2>/dev/null)
	mkdir -p $(REPO_DIR)
	go generate ./internal/repo/... > /dev/null
	touch $@

$(GEN_HANDLERS_STAMP): $(shell find internal/http -type f -name '*.go' 2>/dev/null)
	mkdir -p $(HANDLERS_DIR)
	go generate ./internal/http/... > /dev/null
	touch $@

$(BUILD_STAMP): $(ALL_INPUTS)
	mkdir -p $(CACHE_DIR)
	$(COMPOSE) build $(BUILD_ARGS)
	touch $@

.PHONY: build
build: $(BUILD_STAMP)

.PHONY: rebuild
rebuild: clean build

.PHONY: up
up: build
	$(COMPOSE) up $(UP_FLAGS)

.PHONY: down
down:
	$(COMPOSE) down --remove-orphans

.PHONY: logs
logs:
	$(COMPOSE) logs -f --tail=200

.PHONY: fmt
fmt:
	gofumpt -w cmd internal || true
	goimports -w cmd internal || true
	golines -w cmd internal || true
	gofmt -s -w cmd internal || true
	go mod tidy

.PHONY: unit_test
unit_test: gen
	go test ./internal/... -cover $(TEST_FLAGS)

.PHONY: clean
clean:
	$(COMPOSE) down -v --remove-orphans
	rm -rf pgdata
	rm -rf $(CACHE_DIR)