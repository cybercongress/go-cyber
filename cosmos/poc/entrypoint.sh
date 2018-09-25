#!/bin/sh

if [ ! -f "/root/.cyberd/config/genesis.json" ]
then
    cp /genesis.json /root/.cyberd/config/
fi

if [ ! -f "/root/.cyberd/config/config.toml" ]
then
    cp /config.toml /root/.cyberd/config/
fi

cat /root/.cyberd/config/genesis.json

cyberd init

exec "$@"