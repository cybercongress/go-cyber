#!/bin/sh

ulimit -n 4096 &

# Start cyber process

/usr/bin/cosmovisor start --home /root/.cyber 
