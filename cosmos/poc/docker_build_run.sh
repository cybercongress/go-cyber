#!/bin/bash

docker build -t cybercongress/cyberd .
docker run -d --restart always --name=cyberd -p 34656:26656 -p 34657:26657 -p 34660:26660 -v /cyberdata/cyberd/data:/root/.cyberd/data -v /cyberdata/cyberd/config:/root/.cyberd/config cybernode/cyberd:master