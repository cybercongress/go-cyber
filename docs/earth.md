# Earth cheat sheets

## Copy files to/from earth
```
# Copy from earth
scp -P 33224 earth@earth.cybernode.ai:/path/file /host/path/file
# Copy to earth
scp -P 33324 /host/path/file earth@earth.cybernode.ai:/path/file
```

## Reset cyberd
```
docker stop cyberd
docker rm cyberd
docker run --rm -v /cyberdata/cyberd/data:/root/.cyberd/data -v /cyberdata/cyberd/config:/root/.cyberd/config cybernode/cyberd:master cyberd unsafe-reset-all
docker pull cybernode/cyberd:master
docker run -d --restart always --name=cyberd -p 34656:26656 -p 34657:26657 -p 34660:26660 -v /cyberdata/cyberd/data:/root/.cyberd/data -v /cyberdata/cyberd/config:/root/.cyberd/config cybernode/cyberd:master
```
