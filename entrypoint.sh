#!/bin/sh

if [ ! -f "/root/.cyberd/config/genesis.json" ]
then
    cp /genesis.json /root/.cyberd/config/
    cat /root/.cyberd/config/genesis.json
fi

if [ ! -f "/root/.cyberd/config/config.toml" ]
then
    cp /config.toml /root/.cyberd/config/
fi


if [ ! -f "/root/.cyberd/config/node_key.json" ]
then
    cyberd init
fi

exec "$@"