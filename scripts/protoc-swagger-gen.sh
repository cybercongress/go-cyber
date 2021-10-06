#!/usr/bin/env bash

verbose=''
if [[ "$1" == '-v' || "$1" == '--verbose' ]]; then
  verbose='--verbose'
fi

set -eo pipefail

mkdir -p ./tmp-swagger-gen
proto_dirs=$(find ./proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do

  # generate swagger files for the queries.
  query_file=$(find "${dir}" -maxdepth 1 -name 'query.proto')
  if [[ -n "$query_file" ]]; then
    [[ -n "$verbose" ]] && printf 'Generating swagger file for [%s].\n' "$query_file"
    buf protoc  \
    -I "proto" \
    -I "third_party/proto" \
    "$query_file" \
    --swagger_out=./tmp-swagger-gen \
    --swagger_opt=logtostderr=true --swagger_opt=fqn_for_swagger_name=true --swagger_opt=simple_operation_ids=true
  fi
  # generate swagger files for the transactions.
  tx_file=$(find "${dir}" -maxdepth 1 -name 'tx.proto')
  if [[ -n "$tx_file" ]]; then
    [[ -n "$verbose" ]] && printf 'Generating swagger file for [%s].\n' "$tx_file"
    buf protoc  \
    -I "proto" \
    -I "third_party/proto" \
    "$tx_file" \
    --swagger_out=./tmp-swagger-gen \
    --swagger_opt=logtostderr=true --swagger_opt=fqn_for_swagger_name=true --swagger_opt=simple_operation_ids=true
  fi
done

[[ -n "$verbose" ]] && printf 'Combining swagger files.\n'
# combine swagger files
# uses nodejs package `swagger-combine`.
# all the individual swagger files need to be configured in `config.json` for merging
swagger-combine ./client/docs/config.json -o ./client/docs/swagger-ui/swagger.yaml -f yaml --continueOnConflictingPaths true --includeDefinitions true

# clean swagger files
[[ -n "$verbose" ]] && printf 'Deleting ./tmp-swagger-gen\n'
rm -rf ./tmp-swagger-gen
