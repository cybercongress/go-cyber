#!/bin/sh

set -o errexit -o nounset

CHAINID=$1
GENACCT=$2

if [ -z "$1" ]; then
  echo "Need to input chain id..."
  exit 1
fi

if [ -z "$2" ]; then
  echo "Need to input genesis account address..."
  exit 1
fi

# Build genesis file incl account for passed address
coins="10000000000nick"
cyber init --chain-id "$CHAINID" "$CHAINID"
cyber keys add validator --keyring-backend="test"
cyber add-genesis-account "$(cyber keys show validator -a --keyring-backend="test")" $coins
cyber add-genesis-account "$GENACCT" $coins
cyber gentx validator 5000000000nick --keyring-backend="test" --chain-id $CHAINID
cyber collect-gentxs

# Set proper defaults and change ports
sed -i '' 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ~/.cyber/config/config.toml
sed -i '' 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ~/.cyber/config/config.toml
sed -i '' 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ~/.cyber/config/config.toml
sed -i '' 's/index_all_keys = false/index_all_keys = true/g' ~/.cyber/config/config.toml

# Start the akash
cyber start --pruning=nothing