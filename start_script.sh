#!/bin/sh

ulimit -n 4096 &

export DAEMON_LOG_BUFFER_SIZE=700
export DAEMON_HOME=/root/.cyber
export DAEMON_NAME=cyber
export DAEMON_ALLOW_DOWNLOAD_BINARIES=false
export DAEMON_RESTART_AFTER_UPGRADE=true
export UNSAFE_SKIP_BACKUP=true

# Start cyber process

/root/.cyber/cosmovisor/upgrades/v4/bin/cyber rollback --home /root/.cyber
