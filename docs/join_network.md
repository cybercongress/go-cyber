# Run full node
//todo intro


## Prepare your server

First, you have to setup a server. You are supposed to run your validator node all time, so you will need a 
  reliable server to keep it running. Also, you may consider to use any cloud services like AWS or DigitalOcean.

Cyberd is based on Cosmos SDK written in Go. It should work on any platform which can compile and run programs in Go. 
However, I strongly recommend running the validator node on a Linux server.
Here is the current required server specification to run validator node.

1. No. of CPUs: 8
2. RAM: 64GB
3. Card with Nvidia CUDA support(ex 1080ti) and at least 8gb VRAM.
4. Disk: 256GB SSD

## Install Dependencies

Cyberd requires **Go 1.11+**. 
[Here](https://golang.org/doc/install) you can find default installation instructions for Go for various platforms. 

Check Go version:
```bash
go version
// go version go1.11 linux/amd64
```

For installing **CUDA Toolkit**, please, visit [official Nvidia site](https://docs.nvidia.com/cuda/index.html). 

Check toolkit version:
```bash
nvidia-smi
nvcc --version
// ....
// Cuda compilation tools, release 9.1, V9.1.85
``` 

### Compile binaries

Firstly, clone cyberd with required version.
```bash
git clone -b v0.1.0 --depth 1 https://github.com/cybercongress/cyberd.git
cd cyberd/
```

Than, build GPU Kernel (part of cyberd, that computes rank).
```bash
cd app/rank/cuda
nvcc -fmad=false -shared -o libcbdrank.so rank.cu --compiler-options '-fPIC -frounding-math -fsignaling-nans'
sudo cp libcbdrank.so /usr/lib/
sudo cp cbdrank.h /usr/lib/
cd ../../../
```

Compile binaries, copy configs and run daemon.
```bash

go build -tags cuda -o daemon ./cyberd
go build -o cli ./cyberdcli
```

### Running daemon

Create a configuration directory(by default $HOME/.cyberd) in your home and generate some some required files.

```bash
./daemon init
```

We need the genesis and config files in order to connect to the testnet. 
The files are already included in the cyberd repo. Copy them by.
```bash
cp testnet/genesis.json .cyberd/config/genesis.json
cp testnet/config.toml .cyberd/config/config.toml
```

You can run the node in background logging all output to a log file.

```
./daemon start &> cyberd.log &
```

All output will be logged in the gaia.log file. You can keep checking the output using tail -f.

```
tail -f cyberd.log
```

To check if your node is connected to the testnet, you can run this.

```
./cli status
```

You should be seeing a returned JSON with your node status including node_info, sync_info and validator_info. 
If you see this, then you are connected.
