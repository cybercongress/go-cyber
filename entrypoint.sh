#!/bin/sh

if [ ! -d "/root/.deepchain/" ]
then
  mkdir /root/.deepchain/
  mkdir /root/.deepchain/config/
  /deepchain/cosmovisor/genesis/bin/deepchain init ${NODE_MONIKER}
  cp -r /deepchain/cosmovisor/  /root/.deepchain
fi

if [ ! -f "/root/.deepchain/config/node_key.json" ]
then
  /deepchain/cosmovisor/genesis/bin/deepchain init ${NODE_MONIKER}
  cp /genesis.json /root/.deepchain/config/
  cp -r /deepchain/cosmovisor/  /root/.deepchain
fi

if [ ! -d "/root/.deepchain/cosmovisor" ]
then
  cp -r /deepchain/cosmovisor/  /root/.deepchain/
fi

if [ ! -d "/root/.deepchain/cosmovisor/updgrades" ]
then
  cp -r /deepchain/cosmovisor/upgrades  /root/.deepchain/cosmovisor/
fi

if [  -f "/root/.deepchain/cosmovisor/genesis/bin/deepchain" ]
then
  cp /deepchain/cosmovisor/genesis/bin/deepchain  /root/.deepchain/cosmovisor/genesis/bin/deepchain
fi

if [ ! -f "/root/.deepchain/config/genesis.json" ]
then
  cp /genesis.json /root/.deepchain/config/
fi

if [ "$2" = 'init' ]
then
  return 0
else
  exec "$@"
fi
