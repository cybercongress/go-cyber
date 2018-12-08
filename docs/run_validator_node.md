# Run validator node

Firstly, run full node using "[Run full node](run_testnet.md)" guide. 

Alright, you are now connected to the testnet. 
To be a validator, you will need some **CBD**(testnet coin).

## Prepare stake address
If you already have address with **CBD**, just restore it into your local keyspace.

```bash
./cli keys add <your_key_name> --recover
./cli show <your_key_name>
```

If no, create new one and send coins to it. 
You first need to create a pair of private and public keys for sending and 
receiving coins, and bonding stake. 

```
./cli keys add <your_key_name> 
```

<your_key_name> is any name you pick to represent this key pair. 
You have to refer to this <your_key_name> later when you use the keys to sign transactions. 
It will ask you to enter your password twice to encrypt the key. 
You also need to enter your password when you use your key to sign any transaction.

The command returns the address, public key and a seed phrase which you can use it to 
recover your account if you forget your password later.
Keep the seed phrase in a safe place in case you have to use them.

The address showing here is your account address. Let’s call this <your_account_address>. 
It stores your assets.

## Send create validator transaction

Validators are actors on the network committing new blocks by submitting their votes. 
It refers to the node itself, not a single person or a single account.
Therefore, the public key here is referring to the node public key, 
not the public key of the address you have just created.

To get the node public key, run the following command.

```bash
./daemon tendermint show_validator
```

It will return a bech32 public key. Let’s call it <your_node_pubkey>.
The next step you have to declare a validator candidate. 
The validator candidate is the account which stake the coins. 
So the validator candidate is an account this time.
To declare a validator candidate, run the following command.

```bash
./cli tx stake create-validator \
  --amount=100CBD \
  --pubkey=<your_node_pubkey> \
  --moniker="hlb" \
  --trust-node \
  --from=<your_key_name> \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01"
```