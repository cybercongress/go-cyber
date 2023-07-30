#!/bin/sh

if [ ! -d "/root/.cyber/" ]
then
  mkdir /root/.cyber/
  mkdir /root/.cyber/config/
  /cyber/cosmovisor/genesis/bin/cyber init ${NODE_MONIKER}
  cp -r /cyber/cosmovisor/  /root/.cyber
fi

if [ ! -f "/root/.cyber/config/node_key.json" ]
then
  /cyber/cosmovisor/genesis/bin/cyber init ${NODE_MONIKER}
  cp /genesis.json /root/.cyber/config/
  cp -r /cyber/cosmovisor/  /root/.cyber
fi

if [ ! -d "/root/.cyber/cosmovisor" ]
then
  cp -r /cyber/cosmovisor/  /root/.cyber/
fi

if [ ! -d "/root/.cyber/cosmovisor/updgrades" ]
then
  cp -r /cyber/cosmovisor/upgrades  /root/.cyber/cosmovisor/
fi

if [  -f "/root/.cyber/cosmovisor/genesis/bin/cyber" ]
then
  cp /cyber/cosmovisor/genesis/bin/cyber  /root/.cyber/cosmovisor/genesis/bin/cyber
fi

if [ ! -f "/root/.cyber/config/genesis.json" ]
then
  cp /genesis.json /root/.cyber/config/
fi

if [ "$2" = 'init' ]
then
  return 0
else
  exec "$@"
fi
