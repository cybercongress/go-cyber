#!/usr/bin/env bash

# Run from the project root directory
# This script generates the swagger & openapi.yaml documentation for the rest API on port 1317
#
# Install the following::
# sudo npm install -g swagger2openapi swagger-merger swagger-combine
# go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@v1.16.0

set -eo pipefail

mkdir -p ./tmp-swagger-gen
cd proto

cosmos_sdk_dir=$(go list -f '{{ .Dir }}' -m github.com/cosmos/cosmos-sdk)
ibc_go=$(go list -f '{{ .Dir }}' -m github.com/cosmos/ibc-go/v4)
wasmd=$(go list -f '{{ .Dir }}' -m github.com/CosmWasm/wasmd)

proto_dirs=$(find ./cyber "$cosmos_sdk_dir"/proto "$ibc_go"/proto "$wasmd"/proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  # generate swagger files (filter query files)
  query_file=$(find "${dir}" -maxdepth 1 \( -name 'query.proto' -o -name 'service.proto' \))
  if [[ ! -z "$query_file" ]]; then
    buf generate --template buf.gen.swagger.yaml $query_file
  fi
done

cd ..
# combine swagger files
# uses nodejs package `swagger-combine`.

# Fix circular definition in cosmos/tx/v1beta1/service.swagger.json
jq 'del(.definitions["cosmos.tx.v1beta1.ModeInfo.Multi"].properties.mode_infos.items["$ref"])' ./tmp-swagger-gen/cosmos/tx/v1beta1/service.swagger.json > ./tmp-swagger-gen/cosmos/tx/v1beta1/fixed-service.swagger.json
jq 'del(.definitions["cosmos.autocli."].properties.mode_infos.items["$ref"])' ./tmp-swagger-gen/cosmos/tx/v1beta1/fixed-service.swagger.json > ./tmp-swagger-gen/cosmos/tx/v1beta1/fixed-service2.swagger.json

rm ./tmp-swagger-gen/cosmos/tx/v1beta1/service.swagger.json
rm ./tmp-swagger-gen/ibc/applications/interchain_accounts/host/v1/query.swagger.json
rm ./tmp-swagger-gen/ibc/applications/interchain_accounts/controller/v1/query.swagger.json


# Tag everything as "gRPC Gateway API"
perl -i -pe 's/"(Query|Service)"/"gRPC Gateway API"/' $(find ./tmp-swagger-gen -name '*.swagger.json' -print0 | xargs -0)

# Convert all *.swagger.json files into a single folder _all
files=$(find ./tmp-swagger-gen -name '*.swagger.json' -print0 | xargs -0)
mkdir -p ./tmp-swagger-gen/_all
counter=0
for f in $files; do
  echo "[+] $f"

  if [[ "$f" =~ "cyber" ]]; then
    cp $f ./tmp-swagger-gen/_all/cyber-$counter.json
  elif [[ "$f" =~ "cosmwasm" ]]; then
    cp $f ./tmp-swagger-gen/_all/cosmwasm-$counter.json
  elif [[ "$f" =~ "cosmos" ]]; then
    cp $f ./tmp-swagger-gen/_all/cosmos-$counter.json
  else
    cp $f ./tmp-swagger-gen/_all/other-$counter.json
  fi
  ((counter++))
done

# merges all the above into FINAL.json
python3 ./scripts/merge_protoc.py

# Makes a swagger temp file with reference pointers
swagger-combine ./tmp-swagger-gen/_all/FINAL.json -o ./client/docs/_tmp_swagger.yaml -f yaml --continueOnConflictingPaths --includeDefinitions

# extends out the *ref instances to their full value
swagger-merger --input ./client/docs/_tmp_swagger.yaml -o ./client/docs/swagger.yaml

# Derive openapi from swagger docs
swagger2openapi --patch ./client/docs/swagger.yaml --outfile ./client/docs/static/openapi.yml --yaml

# clean swagger tmp files
rm ./client/docs/_tmp_swagger.yaml
rm -rf ./tmp-swagger-gen