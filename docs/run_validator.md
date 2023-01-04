
# Join cyber as a Validator

## Prepare your server

First, you should set up a server.
Your node should be online constantly. This means that you will need a reliable server.
You may also consider using any cloud service with a dedicated GPU, like Hetzner, Cherryservers etc. (or use a local machine). Whatever you choose, in order to achieve better stability and consistency we recommend you use a dedicated server for each validator node.

Cyber is based on Cosmos-SDK and is written in Go.
It should work on any platform which can compile and run programs in Go.
We strongly recommend running the validator node on a Linux-based server.

Cyber-rank computations are performed on GPU, so it is required to have it (GPU) on-board your node.

Recommended hardware setup:

```js
CPU: 6 cores
RAM: 32 GB
SSD: 1 TB
Connection: 50+Mbps, Stable and low-latency connection
GPU: Nvidia GeForce (or Tesla/Titan/Quadro) with CUDA-cores; 4+ Gb of video memory*
Software: Ubuntu 18.04 LTS / 20.04 LTS
```

*Cyber runs well on consumer-grade cards like Geforce GTX 1070, but we expect load growth and advise you use Error Correction compatible cards from Tesla or Quadro families. Also, make sure your card is compatible with >= v.410 of NVIDIA driver.*

Of course the hardware is your own choice and technically it might be possible to run the node on *"even - 1 CUDA core GPU"*, but you should be aware of performance drop and rank calculation speed decline.

## Node setup

*To avoid possible misconfiguration issues and simplify the setup of `$ENV`, we recommend to perform all the commands as `root` (here root - is literally root, not just a user with root priveliges). For the case of a dedicated server for cybernode it should be concidered as ok from the security side.*

### Third-party software

The main distribution unit for Cyber is a [docker](https://www.docker.com/) container. All images are located in the default [Dockerhub registry](https://hub.docker.com/r/cyberd/cyber). In order to access the GPU from the container, Nvidia driver version **410+** and [Nvidia docker runtime](https://github.com/NVIDIA/nvidia-docker) should be installed on the host system.

All commands below suppose `amd64` architecture, as the different architectures commands may differ accordingly.

### Docker installation

Simply copy the commands below and paste into CLI.

1. Update the `apt` package index:

```bash
sudo apt-get update
```

2. Install packages to allow apt to use a repository over HTTPS:

```bash
sudo apt install -y \
     apt-transport-https \
     ca-certificates \
     curl \
     gnupg-agent \
     software-properties-common
```

3. Add Docker’s official GPG key:

```bash
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
```

```bash
sudo add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
```

4. Update the apt package index:

```bash
sudo apt update
```

5. Install the latest version of Docker CE and containerd:

```bash
sudo apt-get install docker-ce docker-ce-cli containerd.io
```

### Installing Nvidia drivers

1. To proceed, first add the `ppa:graphics-drivers/ppa` repository:

```bash
sudo add-apt-repository ppa:graphics-drivers/ppa
```

```bash
sudo apt update
```

2. Install Ubuntu-drivers:

```bash
sudo apt install -y ubuntu-drivers-common
```

3. Next check what are recommended drivers for your card:

```bash
ubuntu-drivers devices
```

You should see something similar to this:

```bash
== /sys/devices/pci0000:00/0000:00:01.0/0000:01:00.0 ==
modalias : pci:v000010DEd00001BA1sv00001462sd000011E4bc03sc00i00
vendor   : NVIDIA Corporation
model    : GP104M [GeForce GTX 1070 Mobile]
driver   : nvidia-driver-418 - third-party free
driver   : nvidia-driver-430 - third-party free
driver   : nvidia-driver-440 - third-party free
driver   : nvidia-driver-460 - third-party free recommended
driver   : xserver-xorg-video-nouveau - distro free builtin
```

4. We need the **410+** drivers release. As you can see the v460 is recommended. The command below will install the recommended version of the drivers:

```bash
sudo ubuntu-drivers autoinstall
```

To install specific version of a driver use `sudo apt install nvidia-driver-460`


The driver installation takes approximately 10 minutes.

```bash
DKMS: install completed.
Setting up libxdamage1:i386 (1:1.1.4-3) ...
Setting up libxext6:i386 (2:1.3.3-1) ...
Setting up libxfixes3:i386 (1:5.0.3-1) ...
Setting up libnvidia-decode-415:i386 (460.84-0ubuntu0~gpu18.04.1) ...
Setting up build-essential (12.4ubuntu1) ...
Setting up libnvidia-gl-415:i386 (460.84-0ubuntu0~gpu18.04.1) ...
Setting up libnvidia-encode-415:i386 (460.84-0ubuntu0~gpu18.04.1) ...
Setting up nvidia-driver-415 (460.84-0ubuntu0~gpu18.04.1) ...
Setting up libxxf86vm1:i386 (1:1.1.4-1) ...
Setting up libglx-mesa0:i386 (18.0.5-0ubuntu0~18.04.1) ...
Setting up libglx0:i386 (1.0.0-2ubuntu2.2) ...
Setting up libgl1:i386 (1.0.0-2ubuntu2.2) ...
Setting up libnvidia-ifr1-415:i386 (460.84-0ubuntu0~gpu18.04.1) ...
Setting up libnvidia-fbc1-415:i386 (460.84-0ubuntu0~gpu18.04.1) ...
Processing triggers for libc-bin (2.27-3ubuntu1) ...
Processing triggers for initramfs-tools (0.130ubuntu3.1) ...
update-initramfs: Generating /boot/initrd.img-4.15.0-45-generic
```

5. **Reboot** the system for the changes to take effect.

```bash
sudo reboot
```

6. Check the installed drivers:

```bash
nvidia-smi
```

You should see this:
(Some version/driver numbers might differ. You might also have some processes already running)

```bash
+-----------------------------------------------------------------------------+
| NVIDIA-SMI 460.84       Driver Version: 460.84       CUDA Version: 11.2     |
|-------------------------------+----------------------+----------------------+
| GPU  Name        Persistence-M| Bus-Id        Disp.A | Volatile Uncorr. ECC |
| Fan  Temp  Perf  Pwr:Usage/Cap|         Memory-Usage | GPU-Util  Compute M. |
|===============================+======================+======================|
|   0  GeForce GTX 1070    Off  | 00000000:01:00.0 Off |                  N/A |
| 26%   36C    P5    26W / 180W |      0MiB /  8119MiB |      2%      Default |
+-------------------------------+----------------------+----------------------+

+-----------------------------------------------------------------------------+
| Processes:                                                       GPU Memory |
|  GPU       PID   Type   Process name                             Usage      |
|=============================================================================|
|  No running processes found                                                 |
+-----------------------------------------------------------------------------+
```

#### Install Nvidia container runtime for docker

1. Add package repositories:

```bash
distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
```

```bash
curl -s -L https://nvidia.github.io/nvidia-docker/gpgkey | sudo apt-key add -
```

```bash
curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.list | sudo tee /etc/apt/sources.list.d/nvidia-docker.list
```

You should see something like this:

```bash
deb https://nvidia.github.io/libnvidia-container/ubuntu20.04/$(ARCH) /
deb https://nvidia.github.io/nvidia-container-runtime/ubuntu20.04/$(ARCH) /
deb https://nvidia.github.io/nvidia-docker/ubuntu20.04/$(ARCH) /
```

2. Install nvidia-container toolkit and reload the Docker daemon configuration

```bash
sudo apt-get update && sudo apt-get install -y nvidia-container-toolkit
```

```bash
sudo systemctl restart docker
```

3. Test nvidia-smi with the latest official CUDA image

```bash
docker run --gpus all nvidia/cuda:11.1-base nvidia-smi
```

Output logs should coincide as earlier:

```bash
Unable to find image 'nvidia/cuda:11.1-base' locally
11.1-base: Pulling from nvidia/cuda
54ee1f796a1e: Pull complete 
f7bfea53ad12: Pull complete 
46d371e02073: Pull complete 
b66c17bbf772: Pull complete 
3642f1a6dfb3: Pull complete 
e5ce55b8b4b9: Pull complete 
155bc0332b0a: Pull complete 
Digest: sha256:774ca3d612de15213102c2dbbba55df44dc5cf9870ca2be6c6e9c627fa63d67a
Status: Downloaded newer image for nvidia/cuda:11.1-base
Mon Jun 21 14:07:52 2021 
+------------------------------------------------------------------------+
|NVIDIA-SMI 460.84      Driver Version:460.84      CUDA Version: 11.2    |
|-----------------------------+--------------------+---------------------+
|GPU  Name       Persistence-M| Bus-Id       Disp.A| Volatile Uncorr. ECC|
|Fan  Temp  Perf Pwr:Usage/Cap|        Memory-Usage| GPU-Util  Compute M.|
|                             |                    |               MIG M.|
|=============================+====================+=====================|
|  0  GeForce GTX165...  Off  |00000000:01:00.0 Off|                  N/A|
|N/A   47C    P0   16W /  N/A |      0MB /  3914MiB|      0%      Default|
|                             |                    |                  N/A|
+-----------------------------+--------------------+---------------------+
                                                                       
+------------------------------------------------------------------------+
|Processes:                                                              |
| GPU   GI   CI       PID   Type   Process name                GPU Memory|
|       ID   ID                                                Usage     |
|========================================================================|
| No running processes found                                             |
+------------------------------------------------------------------------+
```

Your machine is ready to launch the node.


### Launching cyber fullnode

Make a directory tree for storing your daemon:

```bash
mkdir $HOME/.cyber
mkdir $HOME/.cyber/data
mkdir $HOME/.cyber/config
```


2. Run the full node:
(This will pull and extract the image from cyberd/cyber)

```bash
docker run -d --gpus all --name=bostrom --restart always -p 26656:26656 -p 26657:26657 -p 1317:1317 -e ALLOW_SEARCH=true -v $HOME/.cyber:/root/.cyber  cyberd/cyber:bostrom-1
```

3. Setup some peers to `persistent_peers` and `seeds` to $HOME/.cyber/config/config.toml line 184:


```bash
# Comma separated list of seed nodes to connect to
seeds = ""

# Comma separated list of nodes to keep persistent connections to
persistent_peers = ""
```

For peers addresses please refer to appropriate section of the [networks](https://github.com/cybercongress/networks) repo.
When done, please restart container using:

4. To apply config changes restart the container:

```bash
docker restart bostrom
```

5. Then check the status of your node:

```bash
docker exec bostrom cyber status
```

A possible output may look like this:

```bash
{"NodeInfo":{"protocol_version":{"p2p":"8","block":"11","app":"0"},"id":"808a3773d8adabc78bca6ef8d6b2ee20456bfbcb","listen_addr":"tcp://86.57.207.105:26656","network":"bostrom","version":"","channels":"40202122233038606100","moniker":"node1234","other":
{"tx_index":"on","rpc_address":"tcp://0.0.0.0:26657"}},"SyncInfo":{"latest_block_hash":"241BA3E744A9024A2D04BDF4CE7CF4985D7922054B38AF258712027D0854E930","latest_app_hash":"5BF4B64508A95984F017BD6C29012FE5E66ADCB367D06345EE1EB2ED18314437","latest_block_height":"521",
"latest_block_time":"2021-06-21T14:21:41.817756021Z","earliest_block_hash":"98DD3065543108F5EBEBC45FAAAEA868B3C84426572BE9FDA2E3F1C49A2C0CE8","earliest_app_hash":"E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855","earliest_block_height":"1",
"earliest_block_time":"2021-06-10T00:00:00Z","catching_up":false},"ValidatorInfo":{"Address":"611C9F3568E7341155DBF0546795BF673FD84EB1","PubKey":{"type":"tendermint/PubKeyEd25519","value":"0iGriT3gyRJXQXR/98c+2MTAhChIo5F5v7FfPmOAH5o="},"VotingPower":"0"}}

```

To check container logs use:

```bash
docker logs bostrom -f --tail 10
```

## Validator start

After your node has successfully synced, you can run a validator.

### Prepare the staking address

1. To proceed further you need to add your existing address to the node or generate one and fund it. 

To **create** a new one use:

```bash
docker exec -ti bostrom cyber keys add <your_key_name>
```

The above command returns the address, the public key and the seed phrase, which you can use to
recover your account if you forget your password later.

**Keep you seed phrase safe. Your keys are only your responsibility!**

To **import** existing address use: 

```bash
docker exec -ti bostrom cyber keys add <your_key_name> --recover
```

You can use your **ledger** device with the Cosmos app installed on it to sign transactions. Add address from Ledger:

```bash
docker exec -ti bostrom cyber keys add <your_key_name> --ledger
```

**<your_key_name>** is any name you pick to represent this key pair.
You have to refer to that name later when you use cli to sign transactions.

### Send the create validator transaction

Validators are actors on the network committing new blocks by submitting their votes.
This refers to the node itself, not a single person or a single account.

The next step is to to declare a validator candidate.
To declare a validator candidate, run the following command adjusting the stake amount and the other fields:

```bash
docker exec -ti bostrom cyber tx staking create-validator \
  --amount=10000000boot \
  --min-self-delegation "1000000" \
  --pubkey=$(docker exec -ti bostrom cyber tendermint show-validator) \
  --moniker=<your_node_nickname> \
  --from=<your_key_name> \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --chain-id=bostrom \
  --gas-prices=0.01boot \
  --gas=600000
```

### Verify that you are validating

```bash
docker exec -ti bostrom cyber query staking validators
```

If you see your `<your_node_nickname>` with the status `Bonded` and Jailed `false` everything is good.
You are validating the network.

## Maintenance of the validator

### Jailing

If your validator got under slashing conditions, it will be jailed.
After such an event an operator must unjail the validator manually:

```bash
docker exec -ti bostrom cyber tx slashing unjail --from=<your_key_name> --chain-id=bostrom --gas-prices=0.01boot --gas=300000
```

### Back-up validator keys (!)

Your identity as a validator consists of two things: 

- your account (to sign transactions)
- your validator private key (to sign stuff on the chain consensus layer)

Please back up `$HOME/.cyber/config/priv_validator_key.json` along with your seed phrase. In case of occasional node loss you would be able to restore your validator operation with this file and another full node.

Finally, in case you want to keep your cyber node ID consistent during networks please backup `$HOME/.cyber/config/node_key.json`.

### Rebase to snapshot

You can use a snapshot to start a node from scratch, as well as to reduce its disk space.  
Please check the latest snapshot [here](https://jupiter.cybernode.ai/shared/).  

Download snapshot
```bash
wget https://jupiter.cybernode.ai/shared/bostrom_pruned_<*>.tar
```

Stop your node
```bash
docker stop bostrom
```
Replace data folder to snapshot
```bash 
sudo -H rm -r ~/.cyber/data/
tar -xf bostrom_pruned_<*>.tar -C ~/.cyber/
```
Start node
```bash
docker start bostrom
```
Check node status and logs after startup
```bash
docker exec -ti bostrom cyber status
docker logs bostrom -f --tail 20
```
