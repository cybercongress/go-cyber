#!/bin/sh

set -o errexit -o nounset

CHAINID=$1
HMDIR=$2

if [ -z "$1" ]; then
  echo "Need to input chain id"
  exit 1
fi

if [ -z "$2" ]; then
  echo "Need to input home directory"
  exit 1
fi

# validator and dev accounts should be added to .cyber dir in user's root with keyring-backend test
coins="1000000000000boot"
cyber init --chain-id "$CHAINID" "$CHAINID" --home "$HMDIR"
sed -i '' 's#"stake"#"boot"#g' "$HMDIR"/config/genesis.json
cp -R ~/.cyber/keyring-test $HMDIR
cyber genesis add-genesis-account "$(cyber keys show validator -a --keyring-backend="test")" $coins --home "$HMDIR"
cyber genesis add-genesis-account "$(cyber keys show dev -a --keyring-backend="test")" $coins --home "$HMDIR"
cyber genesis gentx validator 5000000000boot --keyring-backend="test" --chain-id $CHAINID --home "$HMDIR"
cyber genesis collect-gentxs --home "$HMDIR"