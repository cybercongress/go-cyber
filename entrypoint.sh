#!/bin/sh

if [ ! -d "/root/.pussy/" ]
then
  mkdir /root/.pussy/
  mkdir /root/.pussy/config/
  mkdir /root/.pussy/data/
  /pussy/cosmovisor/genesis/bin/pussy init ${NODE_MONIKER}
  cp -r /pussy/cosmovisor/  /root/.pussy
fi

if [ ! -f "/root/.pussy/config/node_key.json" ]
then
  /pussy/cosmovisor/genesis/bin/pussy init ${NODE_MONIKER}
  cp /genesis.json /root/.pussy/config/
  cp -r /pussy/cosmovisor/  /root/.pussy
fi

if [ ! -d "/root/.pussy/cosmovisor" ]
then
  cp -r /pussy/cosmovisor/  /root/.pussy/
fi

if [ ! -d "/root/.pussy/cosmovisor/updgrades" ]
then
  cp -r /pussy/cosmovisor/upgrades  /root/.pussy/cosmovisor/
fi

if [ ! -d "/root/.pussy/cosmovisor/updgrades/v0.0.3" ]
then
  cp -r /pussy/cosmovisor/upgrades/v0.0.3  /root/.pussy/cosmovisor/upgrades
fi


if [  -f "/root/.pussy/cosmovisor/genesis/bin/pussy" ]
then
  cp /pussy/cosmovisor/genesis/bin/pussy  /root/.pussy/cosmovisor/genesis/bin/pussy
fi

if [ ! -f "/root/.pussy/config/genesis.json" ]
then
  cp /genesis.json /root/.pussy/config/
fi

if [ "$2" = 'init' ]
then
  return 0
else
  exec "$@"
fi
