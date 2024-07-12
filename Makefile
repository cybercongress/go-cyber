#!/usr/bin/make -f
CUDA_ENABLED ?= false
LEDGER_ENABLED ?= true
VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
TM_VERSION := $(shell go list -m github.com/tendermint/tendermint | sed 's:.* ::')
LINK_STATICALLY ?= false

BINDIR ?= $(GOPATH)/bin
BUILDDIR ?= $(CURDIR)/build/

# for dockerized protobuf tools
DOCKER := $(shell which docker)
BUF_IMAGE=bufbuild/buf@sha256:3cb1f8a4b48bd5ad8f09168f10f607ddc318af202f5c057d52a45216793d85e5 #v1.4.0
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(BUF_IMAGE)
HTTPS_GIT := https://github.com/cybercongress/go-cyber.git

export GO111MODULE = on

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
empty = $(whitespace) $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(empty),$(comma),$(build_tags))

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=cyber \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=cyber \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)" \
		  -X github.com/cometbft/cometbft/version.TMCoreSemVer=$(TM_VERSION)

ifeq ($(LINK_STATICALLY),true)
	ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
endif

ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))
BUILD_FLAGS := -tags "$(build_tags_comma_sep)" -ldflags '$(ldflags)' -trimpath

include contrib/devtools/Makefile

all: build format lint test
.PHONY: all

###############################################################################
###                                Build                                    ###
###############################################################################

build: go.sum
	go build $(BUILD_FLAGS) -o $(BUILDDIR) ./cmd/cyber

install: go.sum
	go install $(BUILD_FLAGS) ./cmd/cyber

.PHONY: build install

###############################################################################
###                           Tools / Dependencies                          ###
###############################################################################

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	go mod verify
	go mod tidy
.PHONY: go.sum

.PHONY: go.sum go-mod-cache

###############################################################################
###                                Localnet                                 ###
###############################################################################

# TODO update localnet flow
#build-docker-cybernode: build-linux
#	$(MAKE) -C networks/local
#
## Run a 4-node testnet locally
#localnet-start: localnet-stop
#	@if ! [ -f build/node0/cyber/config/genesis.json ]; then docker run --rm -v $(CURDIR)/build:/cyber:Z cybercongress/cyber testnet --v 4 -o . --starting-ip-address 192.168.10.2 --keyring-backend=test --chain-id=chain-local ; fi
#	docker-compose up -d
#
## Stop testnet
#localnet-stop:
#	docker-compose down

###############################################################################
###                                Linting                                  ###
###############################################################################

format-tools:
	go install mvdan.cc/gofumpt@v0.4.0
	go install github.com/client9/misspell/cmd/misspell@v0.3.4
	go install golang.org/x/tools/cmd/goimports@latest

lint: format-tools
	golangci-lint run --tests=false
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "*_test.go" | xargs gofumpt -d

format: format-tools
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/docs/statik/statik.go" | xargs gofumpt -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/docs/statik/statik.go" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/docs/statik/statik.go" | xargs goimports -w -local github.com/cybercongress/go-cyber

###############################################################################
###                                Protobuf                                 ###
###############################################################################

PROTO_BUILDER_IMAGE=tendermintdev/sdk-proto-gen:v0.7
PROTO_FORMATTER_IMAGE=tendermintdev/docker-build-proto@sha256:aabcfe2fc19c31c0f198d4cd26393f5e5ca9502d7ea3feafbfe972448fee7cae

proto-all: proto-format proto-lint proto-gen format

proto-gen:
	@echo "Generating Protobuf files"
	$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(PROTO_BUILDER_IMAGE) sh ./scripts/protocgen.sh

proto-format:
	@echo "Formatting Protobuf files"
	$(DOCKER) run --rm -v $(CURDIR):/workspace \
	--workdir /workspace $(PROTO_FORMATTER_IMAGE) \
	find ./ -name *.proto -exec clang-format -i {} \;

proto-swagger-gen:
	@echo "Generating Protobuf Swagger OpenAPI"
	@./scripts/protoc_swagger_openapi_gen.sh

proto-lint:
	@$(DOCKER_BUF) lint --error-format=json

proto-check-breaking:
	@$(DOCKER_BUF) breaking --against $(HTTPS_GIT)#branch=main

###############################################################################
###                                Docs                                     ###
###############################################################################

# TODO update statik and swagger flow
#update-swagger-docs: statik proto-swagger-gen
#	$(BINDIR)/statik -src=client/docs/swagger-ui -dest=client/docs -f -m
#
#.PHONY: update-swagger-docs