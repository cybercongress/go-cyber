#!/bin/sh

cp -n /genesis.json /root/.cyberd/config/
cp -n /config.toml /root/.cyberd/config/

cat /root/.cyberd/config/genesis.json

cyberd init

exec "$@"