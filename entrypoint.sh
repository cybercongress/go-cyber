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
#  cp /config.toml  /root/.cyber/config/
  cp -r /cyber/cosmovisor/  /root/.cyber
fi

if [ ! -f "/root/.cyber/config/genesis.json" ]
then
  cp /genesis.json /root/.cyber/config/
fi

if [ ! -f "/root/.cyber/config/config.toml" ]
then
  cp /config.toml /root/.cyber/config/
fi

if [ "$2" = 'init' ]
then
  return 0
else
  exec "$@"
fi
