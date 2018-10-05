!#/bin/bash
docker stop cyberd
docker rm cyberd
docker run --rm -v /cyberdata/cyberd/data:/root/.cyberd/data -v /cyberdata/cyberd/config:/root/.cyberd/config cybernode/cyberd:master cyberd unsafe_reset_all
docker pull cybernode/cyberd:master
docker run -d --restart always --name=cyberd -p 34656:26656 -p 34657:26657 -v /cyberdata/cyberd/data:/root/.cyberd/data -v /cyberdata/cyberd/config:/root/.cyberd/config cybernode/cyberd:master