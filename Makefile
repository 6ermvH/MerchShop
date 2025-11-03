SPEC       ?= schema.json
OUT        ?= gen/api
GENERATOR  ?= go-gin-server
PKG        ?= api
MODULE     ?= github.com/6ermvH/MerchShop/gen/api
GO_VERSION ?= 1.22

PROPS = interfaceOnly=true,packageName=$(PKG),moduleName=$(MODULE),hideGenerationTimestamp=true

.PHONY: gen modfix

gen:
	openapi-generator generate \
		-i $(SPEC) \
		-g $(GENERATOR) \
		-o $(OUT) \
		-p $(PROPS)
	cd $(OUT) && rm go.mod main.go Dockerfile
