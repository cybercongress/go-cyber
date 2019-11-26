# Join Cyberd devnet as a Validator

**Note** The current active dev testnet is `euler-dev` (substitute <testnet_chain_id> with that value, do not forget to remove the `<` and the `>` symbols).

## Prepare your server

First, you have to setup a server.
You should to run your validator node all time. This means that you will need a reliable server to keep it running.
Also, you may consider to use any cloud service with dedicated GPU, like Hetzner (or a local machine).

Cyberd is based on Cosmos SDK and written in Go.
It should work on any platform which can compile and run programs in Go.
However, we strongly recommend running the validator node on a Linux server.

Rank calculation in cyberd is benefit to GPU computation. 
They are easy to parallelize. This is why it is best to use GPU.

Recommended requirements:

```js
CPU: 6 cores
RAM: 32 GB
SSD: 256 GB
Connection: 100Mb, Fiber, Stable and low-latency connection
GPU: nvidia GeForce(or Tesla/Titan/Quadro) with CUDA-cores; at least 6gb of memory*
Software: Docker, Ubuntu 16.04/18.04 LTS
```
*Cyberd runs well on comsumer grade cards like Geforce GTX 1070, but expecting load growth we advise to use Error Correction compatible cards from Tesla or Quadro families. 

But, of cource, hardware is your onw choise and tecnically it migth be possible to run the chain on "even - 1 CUDA core gpu", but, you should be aware of stabilty and a decline in calculation speed.

## Validator setup

### Third-party software

The main distribution unit for Cyberd is a [docker](https://www.docker.com/) container. All images are located in the default [Dockerhub registry](https://hub.docker.com/r/cyberd/cyberd/). In order to access the GPU from the container, Nvidia drivers version **410+** and [Nvidia docker runtime](https://github.com/NVIDIA/nvidia-docker) should be installed on the host system. For better user experience, we propose you use [portainer](https://portainer.io) - a docker containers manager. You can skip any subsection of this guide if you already have any of the necessary software configured.

#### Docker installation
Simply, copy the commands below into your CLI.

1. Update the apt package index:

```bash
sudo apt-get update
```

2. Install packages to allow apt to use a repository over HTTPS:

```bash
sudo apt-get install \
     apt-transport-https \
     ca-certificates \
     curl \
     gnupg-agent \
     software-properties-common
```

> You may require installing `curl` - `apt-get install curl`

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
sudo apt-get update
```

5. Install the latest version of Docker CE and containerd or skip to the next step to install a specific version (as of Nov 2019 version 19.03 is required):

```bash
sudo apt-get install docker-ce docker-ce-cli containerd.io
```

If you don’t want to preface docker commands with sudo create a Unix group called docker and add users to that group. When the Docker daemon starts it creates a Unix socket accessible by members of the docker group.

6. Create the docker group:

```bash
sudo groupadd docker
```

7. Add your user to the docker group:

```bash
sudo usermod -aG docker $YOUR-USER-NAME
```

8. Reboot the system for the changes to take effect.

#### Portainer installation (optional)

1. Before installing Portainer, download the Portainer image from the DockerHub using the docker pull command below:

```bash
docker pull portainer/portainer
```

2. Now, run Portainer by using the simple docker command from below:

```bash
docker run -d --restart always -p 9000:9000 -v /var/run/docker.sock:/var/run/docker.sock portainer/portainer
```

3. Open your browser and go to:

```bash
localhost:9000
```

![](https://ipfs.io/ipfs/QmS42MJxjUB7Cu1GoJeE6eBmWkjHTZdgiAUcX4Qqy9NR3M)

4. Create a username and set a password. Chose the `local` tab and click `connect`. 
All the containers will be available in the `containers` tab on your dashboard.

#### Nvidia drivers installation

1. To proceed, first add the `ppa:graphics-drivers/ppa` repository into your system (you might see some warnings - press `enter`):

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

3. Next, identify your graphic card model and the recommended drivers:

```bash
ubuntu-drivers devices
```
You should see something similar to this this:

```bash
== /sys/devices/pci0000:00/0000:00:01.0/0000:01:00.0 ==
modalias : pci:v000010DEd00001BA1sv00001462sd000011E4bc03sc00i00
vendor   : NVIDIA Corporation
model    : GP104M [GeForce GTX 1070 Mobile]
driver   : nvidia-driver-418 - third-party free
driver   : nvidia-driver-410 - third-party free
driver   : nvidia-driver-430 - third-party free
driver   : nvidia-driver-440 - third-party free recommended
driver   : xserver-xorg-video-nouveau - distro free builtin
```

4. We need the
**410+**
drivers release. As you can see v440 is recommended. The command below will install the recommended version of drivers:

```bash
sudo ubuntu-drivers autoinstall
```

Drivers will install for approximately 10 minutes.

```bash
DKMS: install completed.
Setting up libxdamage1:i386 (1:1.1.4-3) ...
Setting up libxext6:i386 (2:1.3.3-1) ...
Setting up libxfixes3:i386 (1:5.0.3-1) ...
Setting up libnvidia-decode-415:i386 (415.27-0ubuntu0~gpu18.04.1) ...
Setting up build-essential (12.4ubuntu1) ...
Setting up libnvidia-gl-415:i386 (415.27-0ubuntu0~gpu18.04.1) ...
Setting up libnvidia-encode-415:i386 (415.27-0ubuntu0~gpu18.04.1) ...
Setting up nvidia-driver-415 (415.27-0ubuntu0~gpu18.04.1) ...
Setting up libxxf86vm1:i386 (1:1.1.4-1) ...
Setting up libglx-mesa0:i386 (18.0.5-0ubuntu0~18.04.1) ...
Setting up libglx0:i386 (1.0.0-2ubuntu2.2) ...
Setting up libgl1:i386 (1.0.0-2ubuntu2.2) ...
Setting up libnvidia-ifr1-415:i386 (415.27-0ubuntu0~gpu18.04.1) ...
Setting up libnvidia-fbc1-415:i386 (415.27-0ubuntu0~gpu18.04.1) ...
Processing triggers for libc-bin (2.27-3ubuntu1) ...
Processing triggers for initramfs-tools (0.130ubuntu3.1) ...
update-initramfs: Generating /boot/initrd.img-4.15.0-45-generic
```

5. Reboot the system for the changes to take effect.

6. Check the installed drivers:

```bash
nvidia-smi
```
You should see this:
(Some version/driver numbers migth differ. You also might have some processes already running)

```
+-----------------------------------------------------------------------------+
| NVIDIA-SMI 440.26       Driver Version: 430.14       CUDA Version: 10.2     |
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

You should see this:

```bash
deb https://nvidia.github.io/libnvidia-container/ubuntu18.04/$(ARCH) /
deb https://nvidia.github.io/nvidia-container-runtime/ubuntu18.04/$(ARCH) /
deb https://nvidia.github.io/nvidia-docker/ubuntu18.04/$(ARCH) /
```

3. Install nvidia-docker2 and reload the Docker daemon configuration

```bash
sudo apt-get update && sudo apt-get install -y nvidia-container-toolkit
```

```bash
sudo systemctl restart docker
```

4. Test nvidia-smi with the latest official CUDA image

```bash
docker run --runtime=nvidia --rm nvidia/cuda:10.0-base nvidia-smi
```

Output logs should coincide as earlier:

```
Unable to find image 'nvidia/cuda:10.0-base' locally
10.0-base: Pulling from nvidia/cuda
38e2e6cd5626: Pull complete
705054bc3f5b: Pull complete
c7051e069564: Pull complete
7308e914506c: Pull complete
5260e5fce42c: Pull complete
8e2b19e62adb: Pull complete
Digest: sha256:625491db7e15efcc78a529d3a2e41b77ffb5b002015983fdf90bf28955277d68
Status: Downloaded newer image for nvidia/cuda:10.0-base
Fri Feb  1 05:41:12 2019       
+-----------------------------------------------------------------------------+
| NVIDIA-SMI 440.26      Driver Version: 440.26       CUDA Version: 10.0     |
|-------------------------------+----------------------+----------------------+
| GPU  Name        Persistence-M| Bus-Id        Disp.A | Volatile Uncorr. ECC |
| Fan  Temp  Perf  Pwr:Usage/Cap|         Memory-Usage | GPU-Util  Compute M. |
|===============================+======================+======================|
|   0  GeForce GTX 1070    Off  | 00000000:01:00.0  On |                  N/A |
| N/A   55C    P0    31W /  N/A |    445MiB /  8117MiB |     38%      Default |
+-------------------------------+----------------------+----------------------+

+-----------------------------------------------------------------------------+
| Processes:                                                       GPU Memory |
|  GPU       PID   Type   Process name                             Usage      |
|=============================================================================|
+-----------------------------------------------------------------------------+
```

Your machine is ready to launch the fullnode.

### Cyberd fullnode launching

1. Create folders for keys and data storing where you want, for example:

``` bash
mkdir /cyberd-dev/cyberd
mkdir /cyberd-dev/cyberdcli
```

2. Run the fullnode:
(This will pull and extract the image from cyberd/cyberd)

```bash
docker run -d --name=euler-dev --restart always -p 26656:26656 -p 26657:26657 -p 1317:1317 -e ALLOW_SEARCH=true -v /cyberd-dev/cyberd:/root/.cyberd  -v /cyberd-dev/cyberdcli:/root/.cyberdcli  cyberd/cyberd:euler-dev
```

3. After successful pulling of the container and launching, run to check if your node is connected to the testnet:

```bash
docker exec euler-dev cyberdcli status
```

A possible output looks like this:

```bash
{"node_info":{"protocol_version":{"p2p":"6","block":"9","app":"0"},"id":"93b776d3eb3f3ce9d9bda7164bc8af3acacff7b6","listen_addr":"tcp://0.0.0.0:26656","network":"euler-dev","version":"0.32.7","channels":"4020212223303800","moniker":"anon","other":{"tx_index":"on","rpc_address":"tcp://0.0.0.0:26657"}},"sync_info":{"latest_block_hash":"686B4E65415D4E56D3B406153C965C0897D0CE27004E9CABF65064B6A0ED4240","latest_app_hash":"0A1F6D260945FD6E926785F07D41049B8060C60A132F5BA49DD54F7B1C5B2522","latest_block_height":"4553","latest_block_time":"2019-11-24T09:49:19.771375108Z","catching_up":false},"validator_info":{"address":"66098853CF3B61C4313DD487BA21EDF8DECACDF0","pub_key":{"type":"tendermint/PubKeyEd25519","value":"uZrCCdZTJoHE1/v+EvhtZufJgA3zAm1bN4uZA3RyvoY="},"voting_power":"0"}}
```

Your node has started to sync. If that didn't happen, check your config.toml file located at /<your euler-dev directory>/cyberd/config/config.toml and add at least a couple of addresses to <persistent_peers = ""> and <seeds = "">, some of those you can fing on our [forum](https://ai.cybercongress.ai/t/euler-dev-testnet/32).

You can follow the syncing process in the terminal. Open a new tab and run the following command:

```bash
docker logs euler-dev --follow
```

Additional information about the chain is available via an API endpoint at: `localhost:26657` (access via your browser) 

e.i. the number of active validators is available at: `localhost:26657/validators`

## Validator start
After your node has successfully synced, you can run a validator.

#### Prepare the staking address

We included 1 million Ethereum addresses, over 8000 Cosmos addresses and all of `euler-4` validators addresses into  genesis, so there's a huge chance that you alredy have some EUL tokens. Here are 3 ways to check this: 

If you already have a cyberd address with EUL and know the seed phrase or your private key, just restore it into your local keystore:
```bash
docker exec -ti euler-dev cyberdcli keys add <your_key_name> --recover
docker exec euler-dev cyberdcli keys show <your_key_name>
```

If you have an Ethereum address that had ~0.2Eth or more at block 8080808 (on the ETH network), you can import your Ethereum private key. To do this, please check out our Ethereum [gift tool](qhttps://github.com/cybercongress/launch-kit/tree/0.1.0/ethereum_gift_tool)

> Please do not import high value Ethereum accounts. This is not safe! cyberd software is a new and has not been battle tested yet.

```bash
docker exec -ti euler-dev cyberdcli keys add import_private <your_key_name>
docker exec euler-dev cyberdcli keys show <your_key_name>
```

If you want to create a new acccount, use the command below:
(You should send coins to that address to bound them later during the submitting of the validator)

```bash
docker exec -ti euler-dev cyberdcli keys add <your_key_name>
docker exec euler-dev cyberdcli keys show <your_key_name>
```

You could use your ledger device with a Cosmos app installed on it to sign and store cyber addresses: [guide here](https://github.com/cybercongress/cyberd/blob/0.1.5/docs/cyberd_Ledger_guide.md). 
In common case use the --ledger flag, with your commands:

```bash
docker exec -ti euler-dev cyberdcli keys add <your_key_name> --ledger
```

**<your_key_name>** is any name you pick to represent this key pair.
You have to refer to this parameter <your_key_name> later, when you use the keys to sign transactions.
It will ask you to enter your password twice to encrypt the key.
You will also need to enter your password when you use your key to sign any transaction.

The command returns the address, a public key and a seed phrase, which you can use to
recover your account if you forget your password later.
Keep the seed phrase at a safe place (preferably, not hot storage) in case you have to use it.

The address shown here is your account address. Let’s call this **<your_account_address>**.
It stores your assets.

#### Send the create validator transaction

Validators are actors on the network committing to new blocks by submitting their votes.
This refers to the node itself, not a single person or a single account.
Therefore, the public key here is referring to the nodes public key,
not the public key of the address you have just created.

To get the nodes public key, run the following command:

```bash
docker exec euler-dev cyberd tendermint show-validator
```

It will return a bech32 public key. Let’s call it **<your_node_pubkey>**.
The next step is to to declare a validator candidate.
The validator candidate is the account which stakes the coins.
So the validator candidate is an account this time.
To declare a validator candidate, run the following command adjusting the stake amount and the other fields:

```bash
docker exec -ti euler-dev cyberdcli tx staking create-validator \
  --amount=10000000eul \
  --min-self-delegation "1000000" \
  --pubkey=<your_node_pubkey> \
  --moniker=<your_node_nickname> \
  --trust-node \
  --from=<your_key_name> \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --chain-id=euler-dev
```

#### Verify that you are validating

```bash
docker exec -ti euler-dev cyberdcli query staking validators --trust-node=true
```

If you see your `<your_node_nickname>` with status `Bonded` and Jailed `false` everything is good. 
You are validating the network.

## Maintenance of the validator

#### Jailing

If your validator got under slashing conditions, it will be jailed. 
After such event, an operator must unjail the validator manually:

```bash
docker exec -ti euler-dev cyberdcli tx slashing unjail --from=<your_key_name> --chain-id euler-dev
```
