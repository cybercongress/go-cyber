#!/usr/bin/make -f
export GO111MODULE = on

CUDA_ENABLED ?= false
LEDGER_ENABLED ?= true

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

TM_VERSION := $(shell go list -m github.com/tendermint/tendermint | sed 's:.* ::') # grab everything after the space in "github.com/tendermint/tendermint v0.34.7"

BINDIR ?= $(GOPATH)/bin
BUILDDIR ?= $(CURDIR)/build/

include contrib/devtools/Makefile

###############################################################################
###                              Build Flags/Tags                           ###
###############################################################################

build_tags = netgo

ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

ifeq ($(CUDA_ENABLED),true)
    NVCC_RESULT := $(shell which nvcc 2> NULL)
    NVCC_TEST := $(notdir $(NVCC_RESULT))
    ifeq ($(NVCC_TEST),nvcc)
        build_tags += cuda
    else
        $(error CUDA not installed for GPU support, please install or set CUDA_ENABLED=false)
    endif
endif

build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))
whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=cyber \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=cyber \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)" \
		  -X github.com/tendermint/tendermint/version.TMCoreSemVer=$(TM_VERSION)

ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))
BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

all: build format lint test

.PHONY: all

###############################################################################
###                                Build                                    ###
###############################################################################

build: go.sum
	go build $(BUILD_FLAGS) -o $(BUILDDIR) ./cmd/cyber


build-linux: go.sum
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build
#	mkdir -p ./build
#	docker build --tag cybercongress/cyber ./
#	docker create --name temp cybercongress/cyber:latest
#	docker cp temp:/usr/bin/cyber ./build/
#	docker rm temp

install: go.sum
	go install $(BUILD_FLAGS) ./cmd/cyber

run:
	$(BUILDDIR)/cyber --home $(BUILDDIR)/bostrom-dev start

###############################################################################
###                           Tools / Dependencies                          ###
###############################################################################

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download
.PHONY: go-mod-cache

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify # TODO uncomment on release
	go mod tidy -compat=1.17
.PHONY: go.sum

lint:
	$(BINDIR)/golangci-lint run
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "*.pb.go" | xargs gofmt -d -s
	go mod verify
.PHONY: lint

statik:
	$(GO) get -u github.com/rakyll/statik
	$(GO) generate ./api/...
.PHONY: statik

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "*.pb.go" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "*.pb.go" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "*.pb.go" | xargs goimports -w -local github.com/cybercongress/go-cyber
.PHONY: format

###############################################################################
###                                Localnet                                 ###
###############################################################################

build-docker-cybernode: build-linux
	$(MAKE) -C networks/local

# Run a 4-node testnet locally
localnet-start: localnet-stop
	@if ! [ -f build/node0/cyber/config/genesis.json ]; then docker run --rm -v $(CURDIR)/build:/cyber:Z cybercongress/cyber testnet --v 4 -o . --starting-ip-address 192.168.10.2 --keyring-backend=test --chain-id=chain-local ; fi
	docker-compose up -d

# Stop testnet
localnet-stop:
	docker-compose down


###############################################################################
###                                Protobuf                                 ###
###############################################################################

protoVer=0.13.1
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	@$(protoImage) sh ./scripts/protocgen.sh

proto-swagger-gen:
	@echo "Generating Protobuf Swagger"
	@$(protoImage) sh ./scripts/protoc-swagger-gen.sh

proto-format:
	@$(protoImage) find ./ -name "*.proto" -exec clang-format -i {} \;

proto-lint:
	@$(protoImage) buf lint --error-format=json

proto-check-breaking:
	@$(protoImage) buf breaking --against $(HTTPS_GIT)#branch=main

.PHONY: proto-all proto-gen proto-gen-any proto-swagger-gen proto-format proto-lint proto-check-breaking proto-update-deps


###############################################################################
###                                Docs                                     ###
###############################################################################

update-swagger-docs: statik proto-swagger-gen
	$(BINDIR)/statik -src=client/docs/swagger-ui -dest=client/docs -f -m

.PHONY: update-swagger-docs