
# Join Cyber testnet as a Validator

## Prepare your server

First of all, you should set up a server.
Your node should be running (online) constantly. This means that you will need a reliable server to keep it running.
You may also consider using any cloud service with a dedicated GPU, like Hetzner (or use a local machine). Whatever you'll choose, for better stability and consistency we recommend to use a dedicated server for each validator node.

Cyberd is based on Cosmos-SDK and is written in Go.
It should work on any platform which can compile and run programs in Go.
However, we strongly recommend running the validator node on a Linux-based server.

The rank calculations are done via GPU computations.
They are easy to parallelize. This is why we recommended using a GPU.

Recommended requirements:

```js
CPU: 6 cores
RAM: 32 GB
SSD: 1 TB
Connection: 50+Mbps, Stable and low-latency connection
GPU: Nvidia GeForce (or Tesla/Titan/Quadro) with CUDA-cores; 4+ Gb of video memory*
Software: Ubuntu 18.04 LTS / 20.04 LTS
```

*Cyber runs well on consumer-grade cards like Geforce GTX 1070, but we expect load growth and advise to use Error Correction compatible cards from Tesla or Quadro families. Also, make sure your card is compatible with >=v.410 of NVIDIA drivers.*

Of course, the hardware is your own choice and technically it might be possible to run the node on "even - 1 CUDA core GPU", but you should be aware of stability drop and rank calculation speed decline .

## Node setup

*To avoid possible misconfiguration issues and simplify the setup of `$ENV`, we recommend to perform all the commands as `root` (here root - is literally root, not just a user with root priveliges)*

### Third-party software

You may skip any sections of the guide if you already have any of the necessary software configured. 

The main distribution unit for Cyber is a [docker](https://www.docker.com/) container. All images are located in the default [Dockerhub registry](https://hub.docker.com/r/cyberd/cyber). In order to access the GPU from the container, Nvidia drivers version **410+** and [Nvidia docker runtime](https://github.com/NVIDIA/nvidia-docker) should be installed on the host system.

### Docker installation

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

### Installing Nvidia drivers

1. To proceed, first, add the `ppa:graphics-drivers/ppa` repository into your system (you might see some warnings - press `enter`):

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

You should see something similar to this:

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

4. We need the **410+** drivers release. As you can see that v440 is recommended. The command below will install the recommended version of the drivers:

```bash
sudo ubuntu-drivers autoinstall
```

The driver installation takes approximately 10 minutes.

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

**Reboot** the system for the changes to take effect.

5. Check the installed drivers:

```bash
nvidia-smi
```

You should see this:
(Some version/driver numbers might differ. You might also have some processes already running)

```bash
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

2. Install nvidia-container toolkit and reload the Docker daemon configuration

```bash
sudo apt-get update && sudo apt-get install -y nvidia-container-toolkit
```

```bash
sudo systemctl restart docker

```

3. Test nvidia-smi with the latest official CUDA image

```bash
docker run --gpus all nvidia/cuda:11.0-base nvidia-smi
```

Output logs should coincide as earlier:

```bash
Unable to find image 'nvidia/cuda:11.0-base' locally
11.0-base: Pulling from nvidia/cuda
54ee1f796a1e: Pull complete 
f7bfea53ad12: Pull complete 
46d371e02073: Pull complete 
b66c17bbf772: Pull complete 
3642f1a6dfb3: Pull complete 
e5ce55b8b4b9: Pull complete 
155bc0332b0a: Pull complete 
Digest: sha256:774ca3d612de15213102c2dbbba55df44dc5cf9870ca2be6c6e9c627fa63d67a
Status: Downloaded newer image for nvidia/cuda:11.0-base
Mon Jun 21 14:07:52 2021 
+-----------------------------------------------------------------------------+
| NVIDIA-SMI 460.84       Driver Version: 460.84       CUDA Version: 11.2     |
|-------------------------------+----------------------+----------------------+
| GPU  Name        Persistence-M| Bus-Id        Disp.A | Volatile Uncorr. ECC |
| Fan  Temp  Perf  Pwr:Usage/Cap|         Memory-Usage | GPU-Util  Compute M. |
|                               |                      |               MIG M. |
|===============================+======================+======================|
|   0  GeForce GTX 165...  Off  | 00000000:01:00.0 Off |                  N/A |
| N/A   47C    P0    16W /  N/A |      0MiB /  3914MiB |      0%      Default |
|                               |                      |                  N/A |
+-------------------------------+----------------------+----------------------+
                                                                               
+-----------------------------------------------------------------------------+
| Processes:                                                                  |
|  GPU   GI   CI        PID   Type   Process name                  GPU Memory |
|        ID   ID                                                   Usage      |
|=============================================================================|
|  No running processes found                                                 |
+-----------------------------------------------------------------------------+
```

Your machine is ready to launch the fullnode.


### Launching Cyberd fullnode

Make a directory tree for storing your daemon:

```bash
mkdir $HOME/.cyber
```


2. Run the fullnode:
(This will pull and extract the image from cyberd/cyberd)

```bash
docker run -d --gpus all --name=cyber-bostromdev --restart always -p 26656:26656 -p 26657:26657 -p 1317:1317 -e ALLOW_SEARCH=true -v $HOME/.cyber:/root/.cyber  cyberd/cyber:bostromdev-2
```
**TODO update image name for bostrom-dev**

3. After container successfully pulled and launched, check the status of your node:

```bash
docker exec bostrom-dev cyber status
```

A possible output may look like this:

```bash
{"NodeInfo":{"protocol_version":{"p2p":"8","block":"11","app":"0"},"id":"808a3773d8adabc78bca6ef8d6b2ee20456bfbcb","listen_addr":"tcp://86.57.207.105:26656","network":"bostromdev-2","version":"","channels":"40202122233038606100","moniker":"node1234","other":
{"tx_index":"on","rpc_address":"tcp://0.0.0.0:26657"}},"SyncInfo":{"latest_block_hash":"241BA3E744A9024A2D04BDF4CE7CF4985D7922054B38AF258712027D0854E930","latest_app_hash":"5BF4B64508A95984F017BD6C29012FE5E66ADCB367D06345EE1EB2ED18314437","latest_block_height":"52829",
"latest_block_time":"2021-06-21T14:21:41.817756021Z","earliest_block_hash":"98DD3065543108F5EBEBC45FAAAEA868B3C84426572BE9FDA2E3F1C49A2C0CE8","earliest_app_hash":"E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855","earliest_block_height":"1",
"earliest_block_time":"2021-06-10T00:00:00Z","catching_up":false},"ValidatorInfo":{"Address":"611C9F3568E7341155DBF0546795BF673FD84EB1","PubKey":{"type":"tendermint/PubKeyEd25519","value":"0iGriT3gyRJXQXR/98c+2MTAhChIo5F5v7FfPmOAH5o="},"VotingPower":"0"}}

```

4 Setup some peers to $HOME/.cyber/config/config.toml:

```bash
# Comma separated list of nodes to keep persistent connections to
persistent_peers = "d0518ce9881a4b0c5872e5e9b7c4ea8d760dad3f@85.10.207.173:26656,0f7d8d5bb8e831a67d29d5950cff0f0ecafbab54@195.201.105.229:36656,30d949f592baf210dd2fc500c83f087f7ce95a84@86.57.254.202:36656"
```
**TODO update image name for bostrom-dev**

You can follow the syncing process in the terminal:

```bash
docker logs bostrom-dev --f --tail 10
```

## Validator start

After your node has successfully synced, you can run a validator.

### Prepare the staking address

1. To proceed further you need to add your address to node, or generete one and fund it. 
Please checkout if you're eligible for cyber Gift: **TODO add link to gift**
Or use our [port](https://cyber.page/brain) to enter cyber.

To create a new one use:

```bash
docker exec -ti bostrom-dev cyber keys add <your_key_name>
```

To import existing address use: 

```bash
docker exec -ti bostrom-dev cyber keys add <your_key_name> --recover
```

It also possible to import address from Ethereum private key:

```bash
docker exec -ti bostrom-dev cyber keys add import_private <your_key_name>
```

You could use your ledger device with the Cosmos app installed on it to sign and store cyber addresses:

```bash
docker exec -ti bostrom-dev cyber keys add <your_key_name> --ledger
```

**<your_key_name>** is any name you pick to represent this key pair.
You have to refer to this parameter <your_key_name> later, when you use the keys to sign transactions.

The command returns the address, a public key and a seed phrase, which you can use to
recover your account if you forget your password later.

**Keep you seed phrase safe. Your keys is only your responsibility!**


#### Send the create validator transaction

Validators are actors on the network committing to new blocks by submitting their votes.
This refers to the node itself, not a single person or a single account.
Therefore, the public key here is referring to the nodes public key,
not the public key of the address you have just created.

To get the nodes public key, run the following command:

```bash
docker exec bostrom-dev cyberd tendermint show-validator
```

It will return a bech32 public key. Let’s call it **<your_node_pubkey>**.
The next step is to to declare a validator candidate.
The validator candidate is the account which stakes the coins.
So the validator candidate is an account this time.
To declare a validator candidate, run the following command adjusting the stake amount and the other fields:

```bash
docker exec -ti bostrom-dev cyber tx staking create-validator \
  --amount=10000000eul \
  --min-self-delegation "1000000" \
  --pubkey=<your_node_pubkey> \
  --moniker=<your_node_nickname> \
  --trust-node \
  --from=<your_key_name> \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --chain-id=bostrom-dev
```

#### Verify that you are validating

```bash
docker exec -ti bostrom-dev cyber query staking validators --trust-node=true
```

If you see your `<your_node_nickname>` with status `Bonded` and Jailed `false` everything is good.
You are validating the network.

## Maintenance of the validator

### Jailing

If your validator got under slashing conditions, it will be jailed.
After such event, an operator must unjail the validator manually:

```bash
docker exec -ti bostrom-dev cyber tx slashing unjail --from=<your_key_name> --chain-id bostrom-dev
```



















#### Launch cyberd

Reload `systemd` after the creation of the new service:

```bash
systemctl daemon-reload
```

Start the node:

```bash
sudo systemctl start cyberd
```

Check node status:

```bash
sudo systemctl status cyberd
```

Enable the service to start with the system:

```bash
sudo systemctl enable cyberd
```

Check logs:

```bash
journalctl -u cyberd -f --lines 20
```

If you need to stop the node:

```bash
sudo systemctl stop cyberd
```

All commands in this section are also applicable to `cyber-rest.service`.

At this point your cyberd should be running in the backgroud and you should be able to call `cyber` to operate with the client. Try calling `cyber status`. A possible output looks like this:

```bash
{"node_info":{"protocol_version":{"p2p":"6","block":"9","app":"0"},"id":"93b776d3eb3f3ce9d9bda7164bc8af3acacff7b6","listen_addr":"tcp://0.0.0.0:26656","network":"euler-6","version":"0.32.7","channels":"4020212223303800","moniker":"anon","other":{"tx_index":"on","rpc_address":"tcp://0.0.0.0:26657"}},"sync_info":{"latest_block_hash":"686B4E65415D4E56D3B406153C965C0897D0CE27004E9CABF65064B6A0ED4240","latest_app_hash":"0A1F6D260945FD6E926785F07D41049B8060C60A132F5BA49DD54F7B1C5B2522","latest_block_height":"4553","latest_block_time":"2019-11-24T09:49:19.771375108Z","catching_up":false},"validator_info":{"address":"66098853CF3B61C4313DD487BA21EDF8DECACDF0","pub_key":{"type":"tendermint/PubKeyEd25519","value":"uZrCCdZTJoHE1/v+EvhtZufJgA3zAm1bN4uZA3RyvoY="},"voting_power":"0"}}
```

Your node has started to sync. If that didn't happen, check your config.toml file located at `$DAEMON_HOME/config/config.toml` and add at least a couple of addresses to <persistent_peers = ""> and <seeds = "">, some of those you can find on our [forum](https://ai.cybercongress.ai/).

Additional information about the chain is available via an API endpoint at: `localhost:26657` (access via your browser)

E.G. the number of active validators is available at: `localhost:26657/validators`

If your node did not launch correctly from the genesis, you need to set the current link to cosmosd for cyber daemon:

```bash
ln -s $DAEMON_HOME/upgrade_manager/genesis current
```

If you joined the testnet **after** a chain upgrade happened, you must point your current link to a new location (with an approptiatly upgraded binary file):

```bash
mkdir $DAEMON_HOME/upgrade_manager/upgrades
cp <path_to_upgraded_cyberd> $DAEMON_HOME/upgrade_manager/upgrades
ln -s $DAEMON_HOME/upgrade_manager/upgrades current
```

## Validator start

After your node has successfully sycned, you can run your validator.

### Prepare a staking address

We included ~1 million Ethereum addresses, over 10000 Cosmos addresses and all of `euler-4` validators addresses into the genesis file. This means that there's a huge chance that you already have some EUL tokens. Here are 3 ways to check this:

If you already have a cyberd address with EUL and know the seed phrase or your private key, just restore it into your local Keystore:

```bash
cyber keys add <your_key_name> --recover
cyber keys show <your_key_name>
```

If you have an Ethereum address that had ~0.2Eth or more at block 8080808 (on the ETH network), you probably received a gift and may import your Ethereum private key. To check your gift balance, paste your Ethereum address on [cyber.page](https://cyber.page).

> Please do not import high-value Ethereum accounts. This is not safe! cyberd software is new and has not been audited yet.

```bash
cyber keys add private <your_key_name>
cyber keys show <your_key_name>
```

If you want to create a new account, use the command below:
(You should send coins to that address to bound them later during the launch of the validator)

```bash
cyber keys add <your_key_name>
cyber keys show <your_key_name>
```

You could use your Ledger device, with the Cosmos app installed on it to sign and store cyber addresses: [guide here](https://github.com/cybercongress/cyberd/blob/0.1.5/docs/cyberd_Ledger_guide.md).
In most cases use the --ledger flag, with your commands:

```bash
cyber keys add <your_key_name> --ledger
```

**<your_key_name>** is any name you pick to represent this key pair.
You have to refer to this parameter <your_key_name> later when you use the keys to sign transactions.
It will ask you to enter your password twice to encrypt the key.
You will also need to enter your password when you use your key to sign any transaction.

The command returns the address, a public key and a seed phrase, which you can use to
recover your account if you forget your password later.
Keep the seed phrase at a safe place (not in hot storage) in case you have to use it.

The address shown here is your account address. Let’s call this **<your_account_address>**.
It stores your assets.

**Important note**: Starting with v.38 cosmos-SDK uses os-native keyring to store all your keys. We've noticed that in certain cases it does not work well by default (for example if you don't have any GUI installed on your machine). If during the execution of the `cyber keys add` command, you are getting this type of error:

```bash
panic: No such  interface 'org.freedesktop.DBus.Properties' on object at path /

goroutine 1 [running]:
github.com/cosmos/cosmos-sdk/crypto/keys.keyringKeybase.writeInfo(0x1307a18, 0x1307a10, 0xc000b37160, 0x1, 0x1, 0xc000b37170, 0x1, 0x1, 0x147a6c0, 0xc000f1c780, ...)
    /root/go/pkg/mod/github.com/cosmos/cosmos-sdk@v0.38.1/crypto/keys/keyring.go:479 +0x38c
github.com/cosmos/cosmos-sdk/crypto/keys.keyringKeybase.writeLocalKey(0x1307a18, 0x1307a10, 0xc000b37160, 0x1, 0x1, 0xc000b37170, 0x1, 0x1, 0x147a6c0, 0xc000f1c780, ...)
    /root/go/pkg/mod/github.com/cosmos/cosmos-sdk@v0.38.1/crypto/keys/keyring.go:465 +0x189
github.com/cosmos/cosmos-sdk/crypto/keys.baseKeybase.CreateAccount(0x1307a18, 0x1307a10, 0xc000b37160, 0x1, 0x1, 0xc000b37170, 0x1, 0x1, 0x146aa00, 0xc000b15630, ...)
    /root/go/pkg/mod/github.com/cosmos/cosmos-sdk@v0.38.1/crypto/keys/keybase_base.go:171 +0x192
github.com/cosmos/cosmos-sdk/crypto/keys.keyringKeybase.CreateAccount(...)
    /root/go/pkg/mod/github.com/cosmos/cosmos-sdk@v0.38.1/crypto/keys/keyring.go:107
github.com/cosmos/cosmos-sdk/client/keys.RunAddCmd(0xc000f0b400, 0xc000f125f0, 0x1, 0x1, 0x148dcc0, 0xc000aca550, 0xc000ea75c0, 0xc000ae1c08, 0x5e93b7)
    /root/go/pkg/mod/github.com/cosmos/cosmos-sdk@v0.38.1/client/keys/add.go:273 +0xa8b
... etc
```

You will have to use another keyring backend to keep your keys. Here are 2 options: store the files within the cli folder or a `pass` manager.

Setting keyring backend to a **local file**:

Execute:

```bash
cyber config keyring-backend file
```

As a result you migth see following: `configuration saved to /root/.cybercli/config/config.toml`

Execute:

```bash
cyber config --get keyring-backend
```

The result should be the following:

```bash
user@node:~# cyber config --get keyring-backend
file
```

This means that you've set your keyring-backend to a local file. *Note*, in this case, all the keys in your keyring will be encrypted using the same password. If you would like to set up a unique password for each key, you should set a unique `--home` folder for each key. To do that, just use `--home=/<unique_path_to_key_folder>/` with setup keyring backend and at all interactions with keys when using cyber:

```bash
cyber config keyring-backend file --home=/<unique_path_to_key_folder>/
cyber keys add <your_second_key_name> --home=/<unique_path_to_key_folder>/
cyber keys list --home=/<unique_path_to_key_folder>/
```

Set keyring backend to [**pass manager**](https://github.com/cosmos/cosmos-sdk/blob/9cce836c08d14dc6836d07164dd964b2b7226f36/crypto/keyring/doc.go#L30):

Pass utility uses a GPG key to encrypt your keys (but again, it uses the same GPG for all the keys). To install and generate your GPG key you should follow [this guide](https://www.passwordstore.org/) or this very [detailed guide](http://tuxlabs.com/?p=450). When you'll get your `pass` set, configure `cyber` to use it as a keyring backend:

```bash
cyber config keyring-backend pass
```

And verify that all has been set as planned:

```bash
cyber config --get keyring-backend
pass
```

#### Send the create validator transaction

Validators are actors on the network committing to new blocks by submitting their votes.
This refers to the node itself, not a single person or a single account.
Therefore, the public key here is referring to the nodes public key,
not the public key of the address you have just created.

To get the nodes public key run the following command:

```bash
cyberd tendermint show-validator
```

It will return a bech32 public key. Let’s call it **<your_node_pubkey>**.
The next step is to declare a validator candidate.
The validator candidate is the account which stakes the coins.
So the validator candidate is the account this time.
To declare a validator candidate, run the following command adjusting the staked amount and the other fields:

```bash
cyber tx staking create-validator \
  --amount=10000000eul \
  --min-self-delegation "1000000" \
  --pubkey=<your_node_pubkey> \
  --moniker=<your_node_nickname> \
  --trust-node \
  --from=<your_key_name> \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --chain-id=euler-6
```

#### Verify that you are validating

```bash
cyber query staking validators --trust-node=true
```

If you see your `<your_node_nickname>` with status `Bonded` and Jailed `false`, everything is good.
You are validating the network.

## Maintenance of the validator

### Jailing

If your validator got slashed, it will get jailed.
If it happens the operator must unjail the validator manually:

```bash
cyber tx slashing unjail --from=<your_key_name> --chain-id euler-6
```
