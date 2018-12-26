#!/bin/sh

if [ ! -f "/root/.cyberd/config/node_key.json" ]
then
    cyberd init
    cp /genesis.json /root/.cyberd/config/
    cp /config.toml /root/.cyberd/config/
fi

if [ ! -f "/root/.cyberd/config/genesis.json" ]
then
    cp /genesis.json /root/.cyberd/config/
fi

if [ ! -f "/root/.cyberd/config/config.toml" ]
then
    cp /config.toml /root/.cyberd/config/
fi

exec "$@"