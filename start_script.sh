#!/bin/sh

ulimit -n 4096 &

export DAEMON_LOG_BUFFER_SIZE=700

# Start cyber process
#// TODO with Cosmovisor v1.x add run command (cosmovisor run)
/usr/bin/cosmovisor run start --compute-gpu true --search-api $ALLOW_SEARCH --home /root/.cyber