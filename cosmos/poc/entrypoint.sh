#!/bin/sh

cp /genesis.json /root/.cyberd/config/
cp /config.toml /root/.cyberd/config/

cat /root/.cyberd/config/genesis.json

cyberd init

exec "$@"