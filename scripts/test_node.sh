# Ensure cyber is installed first.

KEY="bostrom1"
CHAINID="bostrom-t1"
MONIKER="localbostrom"
KEYALGO="secp256k1"
KEYRING="test"
LOGLEVEL="info"


cyber config keyring-backend $KEYRING
cyber config chain-id $CHAINID
cyber config output "json"

command -v jq > /dev/null 2>&1 || { echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; exit 1; }

from_scratch () {

  make install

  # remove existing daemon
  rm -rf ~/.cyber/* 

  # if $KEY exists it should be deleted
  # decorate bright ozone fork gallery riot bus exhaust worth way bone indoor calm squirrel merry zero scheme cotton until shop any excess stage laundry
  # bostrom1hj5fveer5cjtn4wd6wstzugjfdxzl0xp9lxpjy
  echo "decorate bright ozone fork gallery riot bus exhaust worth way bone indoor calm squirrel merry zero scheme cotton until shop any excess stage laundry" | cyber keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO --recover
  # Set moniker and chain-id for Craft
  cyber init $MONIKER --chain-id $CHAINID 

  # Function updates the config based on a jq argument as a string
  update_test_genesis () {
    # update_test_genesis '.consensus_params["block"]["max_gas"]="100000000"'
    cat $HOME/.cyber/config/genesis.json | jq "$1" > $HOME/.cyber/config/tmp_genesis.json && mv $HOME/.cyber/config/tmp_genesis.json $HOME/.cyber/config/genesis.json
  }

  # Set gas limit in genesis
  update_test_genesis '.consensus_params["block"]["max_gas"]="100000000"'
  update_test_genesis '.app_state["gov"]["voting_params"]["voting_period"]="15s"'
  
  update_test_genesis '.app_state["staking"]["params"]["bond_denom"]="boot"'

  # v46 only
  # update_test_genesis '.app_state["bank"]["params"]["send_enabled"]=[{"denom": "boot","enabled": true}]'
  # update_test_genesis '.app_state["staking"]["params"]["min_commission_rate"]="0.050000000000000000"'
  
  update_test_genesis '.app_state["mint"]["params"]["mint_denom"]="boot"'  
  update_test_genesis '.app_state["gov"]["deposit_params"]["min_deposit"]=[{"denom": "boot","amount": "1000000"}]' # 1 eve right now
  update_test_genesis '.app_state["crisis"]["constant_fee"]={"denom": "boot","amount": "1000"}'


  # Allocate genesis accounts  
  cyber add-genesis-account $KEY 10000000boot --keyring-backend $KEYRING

  # create gentx with 1 eve
  cyber gentx $KEY 1000000boot --keyring-backend $KEYRING --chain-id $CHAINID

  # Collect genesis tx
  cyber collect-gentxs

  # Run this to ensure everything worked and that the genesis file is setup correctly
  cyber validate-genesis
}

from_scratch

# Opens the RPC endpoint to outside connections
sed -i '/laddr = "tcp:\/\/127.0.0.1:26657"/c\laddr = "tcp:\/\/0.0.0.0:26657"' ~/.cyber/config/config.toml
sed -i 's/cors_allowed_origins = \[\]/cors_allowed_origins = \["\*"\]/g' ~/.cyber/config/config.toml
# cors_allowed_origins = []

# # Start the node (remove the --pruning=nothing flag if historical queries are not needed)
cyber start --pruning=nothing  --minimum-gas-prices=0boot #--mode validator     