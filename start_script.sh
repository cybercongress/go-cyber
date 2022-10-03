#!/bin/sh

ulimit -n 4096 &

export DAEMON_LOG_BUFFER_SIZE=700

# Start pussy process
#// TODO with Cosmovisor v1.x add run command (cosmovisor run)
/usr/bin/cosmovisor run start --compute-gpu=$COMPUTE_GPU --search-api $ALLOW_SEARCH --home /root/.pussy