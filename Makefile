GO ?= go
GOPATH ?= $(shell go env GOPATH)

api: install-dev build-api clean

install-dev:
	go mod download;
	cd scadagobr-client; \
    	npm i


build-api: build-web generate
	CGO_ENABLED=0 $(GO) build -v ./cmd/api/api.go

build-web:
	cd scadagobr-client; \
	npm run build

clean:
	rm pkg/server/public -r

generate:
	go generate pkg/server/api.go