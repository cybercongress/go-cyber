# Join Cyberd Network As Validator

**Note**. Currently active dev testnet is `euler` (substitute <testnet_chain_id> with that value).

## Prepare your server

First, you have to setup a server.
You are supposed to run your validator node all time, so you will need a reliable server to keep it running.
Also, you may consider to use any cloud services like AWS or DigitalOcean.

Cyberd is based on Cosmos SDK written in Go.
It should work on any platform which can compile and run programs in Go.
However, I strongly recommend running the validator node on a Linux server.

Here is the current required server specification to run validator node:

1. No. of CPUs: 6
2. RAM: 32GB
3. Card with Nvidia CUDA support(ex 1080ti) and at least 8gb VRAM.
4. Disk: 256GB SSD


## Install Dependencies

Our main distribution unit is docker container.
All images are located in default [Dockerhub registry](https://hub.docker.com/r/cyberd/cyberd/).
Docker installation instruction can be found [here](https://docs.docker.com/install/).
In order to access **GPU** from container, nvidia drivers version **410+** and
 [nvidia docker runtime](https://github.com/NVIDIA/nvidia-docker) should be installed on host system.

**Note**: Before installing nvidia docker runtime, reboot pc(nvidia drivers should be loaded into kernel during startup),
 check that drivers loaded correctly by `nvidia-smi` command.

Check both driver and docker runtime installed correctly:
```bash
docker run --runtime=nvidia --rm nvidia/cuda:10.0-base nvidia-smi

# Should be displayed something like this.
Tue Dec 11 18:02:15 2018       
+-----------------------------------------------------------------------------+
| NVIDIA-SMI 410.78       Driver Version: 410.78       CUDA Version: 10.0     |
|-------------------------------+----------------------+----------------------+
| GPU  Name        Persistence-M| Bus-Id        Disp.A | Volatile Uncorr. ECC |
| Fan  Temp  Perf  Pwr:Usage/Cap|         Memory-Usage | GPU-Util  Compute M. |
|===============================+======================+======================|
|   0  GeForce GTX 1050    Off  | 00000000:01:00.0 Off |                  N/A |
| N/A   52C    P0    N/A /  N/A |    302MiB /  4042MiB |      2%      Default |
+-------------------------------+----------------------+----------------------+

+-----------------------------------------------------------------------------+
| Processes:                                                       GPU Memory |
|  GPU       PID   Type   Process name                             Usage      |
|=============================================================================|
+-----------------------------------------------------------------------------+
```

## Start Fullnode

First, create folder, where daemon and cli will store data. Some commands may require admin privileges.
```bash
mkdir /cyberd
```

Run daemon with mounted volumes on created during previous step folder.
```bash
docker run -d --name=cyberd --restart always --runtime=nvidia \
 -p 26656:26656 -p 26657:26657 -p 26660:26660 \
 -v /cyberd/daemon:/root/.cyberd \
 -v /cyberd/cli:/root/.cyberdcli \
 cyberd/cyberd:<testnet_chain_id>
```

To check if your node is connected to the testnet, you can run this:
```
docker exec cyberd cyberdcli status
```
You should be seeing a returned JSON with your node status including node_info, sync_info and validator_info.

## Upgrading

Updating is easy as pulling the new docker container and launching it again

```bash
docker stop cyberd
docker rm cyberd
docker pull cyberd/cyberd:euler

docker run -d --name=cyberd --restart always --runtime=nvidia \
 -p 26656:26656 -p 26657:26657 -p 26660:26660 \
 -v /cyberd/daemon:/root/.cyberd \
 -v /cyberd/cli:/root/.cyberdcli \
 cyberd/cyberd:<testnet_chain_id>
```

Don't forget to unjail if you was jailed during update.

## Get CYB tokens

To be a validator, you will need some CYB to be bounded as your stake.
Top 146 validators by bounded stake will be active validators taking part in consensus and earning validator rewards.

Once the network do not have 146 validators 100M CYB is enougth to join. Don't worry - it is not so much tokens. CYB do not have decimals so 1 CYB is comparable to 1 wei in Ethereum. Cyber network is pure POS [Tendermint](https://tendermint.com/) based network, so there is no ability to mine tokens. But there are [several ways to get tokens](/docs/get_CYB.md).

## Prepare stake address

If you already have address with CYB and know seed phrase just restore it into your local keystore.
```bash
docker exec -ti cyberd cyberdcli keys add <your_key_name> --recover
docker exec cyberd cyberdcli keys show <your_key_name>
```

If you have been lucky enought and your Ethereum address has been included in genesis you can import ethereum private key

> Please, do not import high value Ethereum accounts. This can not be safe! cyberd software is a new software and is not battle tested yet.

```bash
docker exec -ti cyberd cyberdcli keys add import_private <your_key_name>
docker exec cyberd cyberdcli keys show <your_key_name>
```

If you want to create new acccount use the command below.
Also, you should send coins to that address to bound them later during validator submitting.
```
docker exec -ti cyberd cyberdcli keys add <your_key_name>
docker exec cyberd cyberdcli keys show <your_key_name>
```

**<your_key_name>** is any name you pick to represent this key pair.
You have to refer to this <your_key_name> later when you use the keys to sign transactions.
It will ask you to enter your password twice to encrypt the key.
You also need to enter your password when you use your key to sign any transaction.

The command returns the address, public key and a seed phrase which you can use it to
recover your account if you forget your password later.
Keep the seed phrase in a safe place in case you have to use them.

The address showing here is your account address. Let’s call this **<your_account_address>**.
It stores your assets.

## Send create validator transaction

Validators are actors on the network committing new blocks by submitting their votes.
It refers to the node itself, not a single person or a single account.
Therefore, the public key here is referring to the node public key,
not the public key of the address you have just created.

To get the node public key, run the following command.

```bash
docker exec cyberd cyberd tendermint show-validator
```

It will return a bech32 public key. Let’s call it **<your_node_pubkey>**.
The next step you have to declare a validator candidate.
The validator candidate is the account which stake the coins.
So the validator candidate is an account this time.
To declare a validator candidate, run the following command adjusting stake amount and other fields.

```bash
docker exec -ti cyberd cyberdcli tx staking create-validator \
  --amount=100CBD \
  --pubkey=<your_node_pubkey> \
  --moniker=<your_node_nickname> \
  --trust-node \
  --from=<your_key_name> \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --chain-id=<testnet_chain_id>
```

## Verify that you validating

```
docker exec -ti cyberd cyberdcli query stake validators --trust-node=true
```

If you see your valconspub key with status `Bonded` and Jailed `false` everything must be good. You are validating the network.

## Unjailing

If your validator go under slashing conditions it first go to jail. After this event operator must unjail it manually.

```
docker exec -ti cyberd cyberdcli tx slashing unjail --from=<your-wallet> --chain-id=<chain-id>
```
