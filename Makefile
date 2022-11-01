GO ?= go
GOPATH ?= $(shell go env GOPATH)

api: install-dev build-api

install-dev:
	go mod download;
	cd scadagobr-client; \
    	yarn

generate:
	go generate pkg/server/api.go

build-api: build-web generate
	CGO_ENABLED=0 $(GO) build -v ./cmd/api/api.go

build-web:
	cd scadagobr-client; \
	VITE_EMBEDDED=true npm run build

clean:
	rm pkg/server/public -r

test:
	go test ./...

