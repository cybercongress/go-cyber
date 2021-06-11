#!/bin/sh

ulimit -n 4096 &

# Start cyber process

/usr/bin/cosmovisor start --compute-gpu true --search-api true --home /root/.cyber 