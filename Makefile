GOFILES := ./cmd/ ./internal/

GENERATE=gen/openapi

all: up

build: $(GENERATE)
	docker-compose build

up:
	docker-compose up

fmt:
	gofumpt -w $(GOFILES)
	goimports -w $(GOFILES)
	golines -w $(GOFILES)
	gofmt -s -w $(GOFILES)
.PHONY: fmt

$(GENERATE):
	go generate ./internal/... > /dev/null

clean:
	docker-compose down -v
	rm -rf pgdata
	rm -rf gen/openapi gen/.openapi-generator
.PHONY: clean

