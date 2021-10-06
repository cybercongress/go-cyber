#!/usr/bin/make -f
export GO111MODULE = on

CUDA_ENABLED ?= false
LEDGER_ENABLE ?= true

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

TM_VERSION := $(shell go list -m github.com/tendermint/tendermint | sed 's:.* ::') # grab everything after the space in "github.com/tendermint/tendermint v0.34.7"

DOCKER := $(shell which docker)
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf

include contrib/devtools/Makefile

###############################################################################
###                                Build Flags                              ###
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

ifeq ($(WITH_CLEVELDB),yes)
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif

ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))
BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

all: build format lint test

.PHONY: all

###############################################################################
###                                Build                                    ###
###############################################################################

build: go.sum
	go build $(BUILD_FLAGS) -o build/cyber ./cmd/cyber


build-linux: go.sum
	#LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build
	mkdir -p ./build
	docker build --tag cybercongress/cyber ./
	docker create --name temp cybercongress/cyber:latest
	docker cp temp:/usr/bin/cyber ./build/
	docker rm temp

install: go.sum
	go install $(BUILD_FLAGS) ./cmd/cyber

###############################################################################
###                            Tools / Dependencies                         ###
###############################################################################

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download
.PHONY: go-mod-cache

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@#go mod verify
	@go mod tidy

lint:
	$(BINDIR)/golangci-lint run
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "*.pb.go" | xargs gofmt -d -s
	go mod verify
.PHONY: lint

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "*.pb.go" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "*.pb.go" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "*.pb.go" | xargs goimports -w -local github.com/cybercongress/cyber
.PHONY: format

###############################################################################
###                                Localnet                                 ###
###############################################################################

build-docker-cybernode: build-linux
	$(MAKE) -C networks/local

# Run a 4-node testnet locally
localnet-start: localnet-stop
	@if ! [ -f build/node0/cyber/config/genesis.json ]; then docker run --rm -v $(CURDIR)/build:/cyber:Z cybercongress/cyber testnet --v 4 -o . --starting-ip-address 192.168.10.2 --keyring-backend=test ; fi
	docker-compose up -d

# Stop testnet
localnet-stop:
	docker-compose down


###############################################################################
###                                Proto                                    ###
###############################################################################

proto-all: proto-tools proto-format proto-lint proto-gen proto-check-breaking proto-swagger-gen

containerProtoVer=v0.2
containerProtoImage=tendermintdev/sdk-proto-gen:$(containerProtoVer)
containerProtoGen=cosmos-sdk-proto-gen-$(containerProtoVer)
containerProtoGenSwagger=cosmos-sdk-proto-gen-swagger-$(containerProtoVer)
containerProtoFmt=cosmos-sdk-proto-fmt-$(containerProtoVer)

proto-gen:
	@echo "Generating Protobuf files"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoGen}$$"; then docker start -a $(containerProtoGen); else docker run --name $(containerProtoGen) -v $(CURDIR):/workspace --workdir /workspace $(containerProtoImage) \
		sh ./scripts/protocgen.sh; fi

# This generates the SDK's custom wrapper for google.protobuf.Any. It should only be run manually when needed
proto-gen-any:
	@echo "Generating Protobuf Any"
	$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(containerProtoImage) sh ./scripts/protocgen-any.sh

proto-swagger-gen:
	@echo "Generating Protobuf Swagger"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoGenSwagger}$$"; then docker start -a $(containerProtoGenSwagger); else docker run --name $(containerProtoGenSwagger) -v $(CURDIR):/workspace --workdir /workspace $(containerProtoImage) \
		sh ./scripts/protoc-swagger-gen.sh; fi

proto-format:
	@echo "Formatting Protobuf files"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoFmt}$$"; then docker start -a $(containerProtoFmt); else docker run --name $(containerProtoFmt) -v $(CURDIR):/workspace --workdir /workspace tendermintdev/docker-build-proto \
		find ./ -not -path "./third_party/*" -name *.proto -exec clang-format -i {} \; ; fi

proto-lint:
	@$(DOCKER_BUF) lint --error-format=json

proto-check-breaking:
	@$(DOCKER_BUF) breaking --against $(HTTPS_GIT)#branch=main

TM_URL           = https://raw.githubusercontent.com/tendermint/tendermint/v0.34.x/proto/tendermint
GOGO_PROTO_URL   = https://raw.githubusercontent.com/regen-network/protobuf/cosmos
COSMOS_PROTO_URL = https://raw.githubusercontent.com/regen-network/cosmos-proto/master
COSMOS_SDK_URL   = https://raw.githubusercontent.com/cosmos/cosmos-sdk/release/v0.43.x/proto/cosmos
CONFIO_URL       = https://raw.githubusercontent.com/confio/ics23/v0.6.3

TM_CRYPTO_TYPES     = third_party/proto/tendermint/crypto
TM_ABCI_TYPES       = third_party/proto/tendermint/abci
TM_TYPES            = third_party/proto/tendermint/types
TM_VERSION          = third_party/proto/tendermint/version
TM_LIBS             = third_party/proto/tendermint/libs/bits

GOGO_PROTO_TYPES     = third_party/proto/gogoproto
COSMOS_PROTO_TYPES   = third_party/proto/cosmos_proto
COSMOS_BASE_TYPES    = third_party/proto/cosmos/base
COSMOS_SIGNING_TYPES = third_party/proto/cosmos/tx/signing
COSMOS_CRYPTO_TYPES  = third_party/proto/cosmos/crypto
COSMOS_AUTH_TYPES    = third_party/proto/cosmos/auth
COSMOS_BANK_TYPES    = third_party/proto/cosmos/bank
CONFIO_TYPES         = third_party/proto/confio

proto-update-deps:
	@mkdir -p $(GOGO_PROTO_TYPES)
	@curl -sSL $(GOGO_PROTO_URL)/gogoproto/gogo.proto > $(GOGO_PROTO_TYPES)/gogo.proto

	@mkdir -p $(COSMOS_PROTO_TYPES)
	@curl -sSL $(COSMOS_PROTO_URL)/cosmos.proto > $(COSMOS_PROTO_TYPES)/cosmos.proto

	@mkdir -p $(COSMOS_BASE_TYPES)/v1beta1
	@curl -sSL $(COSMOS_SDK_URL)/base/v1beta1/coin.proto > $(COSMOS_BASE_TYPES)/v1beta1/coin.proto

	@mkdir -p $(COSMOS_BASE_TYPES)/query/v1beta1
	@curl -sSL $(COSMOS_SDK_URL)/base/query/v1beta1/pagination.proto > $(COSMOS_BASE_TYPES)/query/v1beta1/pagination.proto

	@mkdir -p $(COSMOS_SIGNING_TYPES)/v1beta1
	@curl -sSL $(COSMOS_SDK_URL)/tx/signing/v1beta1/signing.proto > $(COSMOS_SIGNING_TYPES)/v1beta1/signing.proto

	@mkdir -p $(COSMOS_CRYPTO_TYPES)/secp256k1
	@curl -sSL $(COSMOS_SDK_URL)/crypto/secp256k1/keys.proto > $(COSMOS_CRYPTO_TYPES)/secp256k1/keys.proto

	@mkdir -p $(COSMOS_CRYPTO_TYPES)/multisig/v1beta1
	@curl -sSL $(COSMOS_SDK_URL)/crypto//multisig/v1beta1/multisig.proto > $(COSMOS_CRYPTO_TYPES)/multisig/v1beta1/multisig.proto

	@mkdir -p $(COSMOS_AUTH_TYPES)/v1beta1
	@curl -sSL $(COSMOS_SDK_URL)/auth/v1beta1/auth.proto > $(COSMOS_AUTH_TYPES)/v1beta1/auth.proto

	@mkdir -p $(COSMOS_BANK_TYPES)/v1beta1
	@curl -sSL $(COSMOS_SDK_URL)/bank/v1beta1/bank.proto > $(COSMOS_BANK_TYPES)/v1beta1/bank.proto

	@mkdir -p $(TM_ABCI_TYPES)
	@curl -sSL $(TM_URL)/abci/types.proto > $(TM_ABCI_TYPES)/types.proto

	@mkdir -p $(TM_VERSION)
	@curl -sSL $(TM_URL)/version/types.proto > $(TM_VERSION)/types.proto

	@mkdir -p $(TM_TYPES)
	@curl -sSL $(TM_URL)/types/types.proto > $(TM_TYPES)/types.proto
	@curl -sSL $(TM_URL)/types/evidence.proto > $(TM_TYPES)/evidence.proto
	@curl -sSL $(TM_URL)/types/params.proto > $(TM_TYPES)/params.proto
	@curl -sSL $(TM_URL)/types/validator.proto > $(TM_TYPES)/validator.proto

	@mkdir -p $(TM_CRYPTO_TYPES)
	@curl -sSL $(TM_URL)/crypto/proof.proto > $(TM_CRYPTO_TYPES)/proof.proto
	@curl -sSL $(TM_URL)/crypto/keys.proto > $(TM_CRYPTO_TYPES)/keys.proto

	@mkdir -p $(TM_LIBS)
	@curl -sSL $(TM_URL)/libs/bits/types.proto > $(TM_LIBS)/types.proto

	@mkdir -p $(CONFIO_TYPES)
	@curl -sSL $(CONFIO_URL)/proofs.proto > $(CONFIO_TYPES)/proofs.proto.orig
## insert go, java package option into proofs.proto file
## Issue link: https://github.com/confio/ics23/issues/32 (instead of a simple sed we need 4 lines cause bsd sed -i is incompatible)
	@head -n3 $(CONFIO_TYPES)/proofs.proto.orig > $(CONFIO_TYPES)/proofs.proto
	@echo 'option go_package = "github.com/confio/ics23/go";' >> $(CONFIO_TYPES)/proofs.proto
	@echo 'option java_package = "tech.confio.ics23";' >> $(CONFIO_TYPES)/proofs.proto
	@echo 'option java_multiple_files = true;' >> $(CONFIO_TYPES)/proofs.proto
	@tail -n+4 $(CONFIO_TYPES)/proofs.proto.orig >> $(CONFIO_TYPES)/proofs.proto
	@rm $(CONFIO_TYPES)/proofs.proto.orig

.PHONY: proto-all proto-gen proto-format proto-gen-any proto-lint proto-check-breaking
.PHONY: proto-update-deps

###############################################################################
###                                Docs                                     ###
###############################################################################

update-swagger-docs: statik proto-swagger-gen
	$(BINDIR)/statik -src=client/docs/swagger-ui -dest=client/docs -f -m

.PHONY: update-swagger-docs

test-rosetta:
	docker build -t rosetta-ci:latest -f client/rosetta/rosetta-ci/Dockerfile .
	docker-compose -f client/rosetta/docker-compose.yaml --project-directory ./ up --abort-on-container-exit --exit-code-from test_rosetta --build
.PHONY: test-rosetta