#!/bin/sh

for OUTPUT in ../build/node*
do
  cyber unsafe-reset-all --home $OUTPUT/cyber
done