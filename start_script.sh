#!/bin/sh

ulimit -n 4096 &

export DAEMON_LOG_BUFFER_SIZE=700

# Start cyber process

/usr/bin/cosmovisor start --compute-gpu true --search-api true --home /root/.cyber 