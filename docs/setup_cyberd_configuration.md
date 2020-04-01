# Setup config.toml

Correct configuration is one of the main keys to consistent and proper functioning of your node no matter is it a validator or sentinel/service one.

Hereby we'll check all the key points of `config.toml` file and explain hot to configure them for all use-cases.

We will operate basic configuration file (line numbers from it), generated after initialization of cyber daemon and typically located inside `$HOME/.cyberd/config directory`.

> All changes to `config.toml` file require cyberd restart to take effect!

## Ports / Addresses configuration

### RPC port

First of all lets look through the ports cyberd using to operate with outside world. At the *line 84* specified port used for RPC server (TCP and UNIX websocket connections):

```bash
# TCP or UNIX socket address for the RPC server to listen on
laddr = "tcp://127.0.0.1:26657"
```

After the node start RPC server providing some endpoints to check chain/node parameters, accepting $POST transactions and so on, could be opened locally at browser at `http://localhost:26657`.

- For Validator node it's not recommended to open this port to outside world, as it may allow anyone to produce transactions using your node and allows DOS attacks ( you didn't want your validator attacked, right? ). So lets leave it like this:

```bash
laddr = "tcp://127.0.0.1:26657"
```

- For Sentinel node this shold be kept same as for vaidator:

```bash
laddr = "tcp://127.0.0.1:26657"
```

- For Service node, when usecases include remote access to RPC for yourself, or your great service it is allowable to expose it otside by using following values:

```bash
laddr = "tcp://0.0.0.0:26657"
```

Also if you would like to make RPC server respond on different port you may change this value to whatever you'd like (just make sure it will not cross with any other services), for exapmle to change it to `9588` use:

```bash
# TCP or UNIX socket address for the RPC server to listen on
laddr = "tcp://0.0.0.0:9588"
```

### Cyberd communication port

At the line 163 we could find following:

```bash
# Address to listen for incoming connections
laddr = "tcp://0.0.0.0:26656"
```

This is way node communicating with other nodes in the chain. For all cases(**Validator, Sentinel, Service**) leave it at default binded to  `0.0.0.0`. And if you need to change port nubmer to something different like `35622` just use:

```bash
laddr = "tcp://0.0.0.0:35622"
```

And, if changed, your node peer address would be changed accordingly: `75e8f44072b0dd598dfa95aaf9b5f2c60f956819@your_external_ip:35622`.

### Prometheus collectors port

At the line 325 located port for Prometheus monitoring service:

```bash
# Address to listen for Prometheus collector(s) connections
prometheus_listen_addr = ":26660"
```

It is usefull if you want to monitor remotely your node condition using [Prometheus](https://prometheus.io/) metrics collector service and could be changed to whatever like `23456`:

```bash
prometheus_listen_addr = ":23456"
```

Dont forget to enable Prometheus metrics by changing to `true` value at line 322, if needed:

```bash
# When true, Prometheus metrics are served under /metrics on
# PrometheusListenAddr.
# Check out the documentation for the list of available metrics.
prometheus = true
```

### External address

At the line 169 you may fing the following:

```bash
# Address to advertise to peers for them to dial
# If empty, will use the same port as the laddr,
# and will introspect on the listener or use UPnP
# to figure out the address.
external_address = ""
```

That line implies specifying your external ip address, so it means presense of static external address at your network connection. If you dont have one - just skip it.

- For Validator node you may scip it, until you have **enough private peers** to get synced with. Othewise you have to specify your external statis IP to impove peers dicsovery for your node. Also don't forget to change the port according to line [163](#cyberd-communication-port):

```bash
external_address = "tcp://<your_external_static_ip>:26656"
```

- For Sentinel node it is good idea to specify IP for better peer discovery.

```bash
external_address = "tcp://<your_external_static_ip>:26656"
```

- For Servie node this setting could be same as for Sentinel;

```bash
external_address = "tcp://<your_external_static_ip>:26656"
```

And again, all above settings apply to the cases when **STATIC EXTERNAL IP** is available.

### Allow duplicated ip's

Line 224 of config.toml holds following:

```bash
# Toggle to disable guard against peers connecting from the same ip.
allow_duplicate_ip = false
```

That variable configures allowance for different peers to be connected from the same ip. Lets imagine situation that you run 2 nodes (lets say node ID of the first is `75e8f44072b0dd598dfa95aaf9b5f2c60f956819` and second is `d0518ce9881a4b0c5872e5e9b7c4ea8d760dad3f`) on one internet provider with external IP of 92.23.45.123. In this case all other nodes in the network with `allow_duplicate_ip = false` would see  attemtps to connect from peers `d0518ce9881a4b0c5872e5e9b7c4ea8d760dad3f@92.23.45.123:26656` and `75e8f44072b0dd598dfa95aaf9b5f2c60f956819@92.23.45.123:36656` and will block one which comes later, because originating IP is same for both. If this case applies to you, change this setting to following:

```bash
# Toggle to disable guard against peers connecting from the same ip.
allow_duplicate_ip = true
```

## P2P configuration

### Seed nodes

At the line 172 of the config.toml we see following:

```bash
# Comma separated list of seed nodes to connect to
seeds = ""
```

This line is dedicated to list of seed nodes you want to establish connection with. To get seed nodes addresses take a llok at our [forum](https://ai.cybercongress.ai/) or ask at our [chat in Telegram](https://t.me/fuckgoogle).

- For validators with the sentinel nodes, or decent quantity of peers conected it is not required to fill this up:

```bash
seeds = ""
```

- For Sentinel nodes and Service nodes it's a good idea to fill up with couple of seed nodes addresses, separated with commas: 

```bash
seeds = "<seed_node1_ID>@<seed_node1_ip>:<port>,<seed_node2_ID>@<seed_node2_ip>:<port>"
```

### Persistent peers

Place to put in persistent peers locatet at the line 175. Persistent peers presense is very important for correct node functioning.

```bash
# Comma separated list of nodes to keep persistent connections to
persistent_peers = ""
```

- For Validator node you have to fill this line with decent ammount of peers you **trust** to, othervise you validator node address will be exposed to. In the perfect case you should put to this section only with addresses of your sentinel nodes.

```bash
persistent_peers ="<sentinel_node1_ID>@<sentinel_node1_ip>:<port>,<sentinel_node2_ID>@<sentinel_node2_ip>:<port>"
```

- For Sentinel node and Service node add as much peers as possiple to keep persistent connection and network stability, but **DO NOT** put here you'r validator nodes ID's:

```bash
persistent_peers ="<node1_ID>@<node1_ip>:<port>,<node2_ID>@<node2_ip>:<port>,...,<node_n_ID>@<node_n_ip>:<port>"
```

### Peer Exchange Reactor

Line 212 shows by default:

```bash
# Set true to enable the peer-exchange reactor
pex = true
```

This is peer exchange module, which responsible for exchanging node ID's across the network.

- For validator node with sentinel architecture set this have to be disabled:

```bash
pex = false
```

- For Sentinel node and Service node leave as default:

```bash
pex = true
```

### Private peers ID's

At the line 221 we see:

```bash
# Comma separated list of peer IDs to keep private (will not be gossiped to other peers)
private_peer_ids = ""
```

This is the list of the ones, which ID's should not be gossiped to others.

- For Validator nodes - leave as default.

```bash
private_peer_ids = ""
```

Or you may add your 2nd validator ID here (if you running more than 1 validator).

- For Sentinel node add here your validator/s address:

```bash
private_peer_ids = "<validator_node_ID>@<validator_node_ip>:<port>"
```

- For Service node leave blank.

```bash
private_peer_ids = ""
```

## Node Index, Naming

### Indexed tags

Node is able to index and store decent ammount of keyas and values regarding transactions, accounts etc. Lines 306 and 304 restonsible for that:

```bash
# You can also index transactions by height by adding "tx.height" key here.
#
# It's recommended to index only a subset of keys due to possible memory
# bloat. This is, of course, depends on the indexer's DB and the volume of
# transactions.
index_keys = ""

# When set to true, tells indexer to index all compositeKeys (predefined keys:
# "tx.hash", "tx.height" and all keys from DeliverTx responses).
#
# Note this may be not desirable (see the comment above). IndexKeys has a
# precedence over IndexAllKeys (i.e. when given both, IndexKeys will be
# indexed).
index_all_keys = false
```

- For Validato and Sentinel node this is not necessary, so leave as default: 

```bash
index_keys = ""

index_all_keys = false
```

- For Service node you must specify subset of keys you want to idex:

```bash
index_keys = "tx.hash,tx.height,...etc.."

index_all_keys = all
```

### Naming 

To setup up your node moniker please refer to line 16 and type in whatever you want to have as moniker:

```bash
# A custom human readable name for this node
moniker = "god_damn_node"
```
