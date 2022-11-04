GO ?= go
GOPATH ?= $(shell go env GOPATH)

VERSION=0.0.1
V_COMMIT := $(shell git rev-parse HEAD)
V_BUILT_BY := $(shell git config user.email)
V_BUILT_AT := $(shell date)

V_LDFLAGS_COMMON := -s -X "github.com/MarcusGoldschmidt/scadagobr/pkg.Version=${VERSION}" \
					-X "github.com/MarcusGoldschmidt/scadagobr/pkg.Commit=${V_COMMIT}" \
					-X "github.com/MarcusGoldschmidt/scadagobr/pkg.BuiltBy=${V_BUILT_BY}" \
					-X "github.com/MarcusGoldschmidt/scadagobr/pkg.BuiltAt=${V_BUILT_AT}"

api: install-dev build-api

install-dev:
	go mod download;
	cd scadagobr-client; \
    	yarn

generate:
	go generate pkg/server/api.go

build-api: build-web generate
	CGO_ENABLED=0 $(GO) build -v -ldflags '$(V_LDFLAGS_COMMON)' ./cmd/api/api.go

build-web:
	cd scadagobr-client; \
	VITE_EMBEDDED=true npm run build

clean:
	rm pkg/server/public -r
	rm api

test:
	mkdir ./pkg/server/public
	go test ./...

