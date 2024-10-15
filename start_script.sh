#!/bin/sh

ulimit -n 4096 &

export DAEMON_LOG_BUFFER_SIZE=700
export DAEMON_HOME=/root/.cyber
export DAEMON_NAME=cyber
export DAEMON_ALLOW_DOWNLOAD_BINARIES=false
export DAEMON_RESTART_AFTER_UPGRADE=true
export UNSAFE_SKIP_BACKUP=true

# Start cyber process

/usr/bin/cosmovisor run start --compute-gpu true --search-api $ALLOW_SEARCH --home /root/.cyber