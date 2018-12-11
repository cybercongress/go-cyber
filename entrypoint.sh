#!/bin/sh

if [ ! -f "/root/.cyberd/config/node_key.json" ]
then
    cyberd init
    cp /genesis.json /root/.cyberd/config/
    cp /config.toml /root/.cyberd/config/
    cat /root/.cyberd/config/genesis.json
fi

exec "$@"