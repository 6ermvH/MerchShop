GOFILES := ./cmd/ ./internal/

GENERATE_HANDLERS=gen/openapi
GENERATE_REPO=gen/mock

all: up

build: $(GENERATE_REPO) $(GENERATE_HANDLERS)
	docker-compose build

up:
	docker-compose up

fmt:
	gofumpt -w $(GOFILES)
	goimports -w $(GOFILES)
	golines -w $(GOFILES)
	gofmt -s -w $(GOFILES)
.PHONY: fmt

$(GENERATE_REPO):
	go generate ./internal/repo/... > /dev/null

$(GENERATE_HANDLERS):
	go generate ./internal/http/... > /dev/null

clean:
	docker-compose down -v
	rm -rf pgdata
	rm -rf $(GENERATE_HANDLERS) $(GENERATE_REPO)
.PHONY: clean

