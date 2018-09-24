#!/bin/bash

docker build -t cybercongress/cyberd .
docker run -d --name=cyberd -p 26656 -p 26657 -v ~/.cyberd/docker/data:/root/.cyberd/data cybercongress/cyberd:latest