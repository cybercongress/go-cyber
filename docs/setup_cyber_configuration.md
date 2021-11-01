# Setup config.toml

Correct configuration is one of the main keys to consistent and proper functioning of your node no matter if it is a validator or a sentinel/service node.

Throughout this document, we will check all the key points of the `config.toml` file and explain how to configure them for all use-cases.

We will operate a basic configuration file (according to the number of the line in the actual file), generated after the initialization of cyber daemon and typically located inside `$HOME/.cyber/config directory`.

> All changes made to the `config.toml` file, require to restart cyber to take effect!

## Port / Address configuration

### RPC port

First of all, let's look through the ports cyber uses to communicate with the outside world. On *line 84* the specified port is used for an RPC server (TCP and UNIX websocket connections):

```bash
# TCP or UNIX socket address for the RPC server to listen on
laddr = "tcp://127.0.0.1:26657"
```

After the node starts the RPC server provides endpoints to check chain/node parameters, accepts $POST transactions and so on. It can be opened locally using your favourite browser via: `http://localhost:26657`.

- We do not recommend a validator node to open this port to the outside world, as it may allow anyone to produce transactions using your node and allows DOS attacks (you don't want your validator attacked, right?). So let's leave it like this:

```bash
laddr = "tcp://127.0.0.1:26657"
```

- For Sentinel nodes this should be kept the same as for validators:

```bash
laddr = "tcp://127.0.0.1:26657"
```

- For Service nodes, when use cases include remote access to the RPC for yourself or for your great service, it is allowable to expose it to the outside by using the following values:

```bash
laddr = "tcp://0.0.0.0:26657"
```

If you would like to make the RPC server respond on a different port, you may change this value to whatever you'd like (just make sure it will not cross with any of the other services), for example, to change it to `9588` use:

```bash
# TCP or UNIX socket address for the RPC server to listen on
laddr = "tcp://0.0.0.0:9588"
```

### Cyberd communication port

On *line 163* we can find the following:

```bash
# Address to listen for incoming connections
laddr = "tcp://0.0.0.0:26656"
```

This is the way the node communicates with other nodes in the chain. For all possible cases(**Validator, Sentinel, Service**) leave it as default, bound to  `0.0.0.0`. And if you need to change the port number to something different like `35622` just use:

```bash
laddr = "tcp://0.0.0.0:35622"
```

If changed, your node peer address would be changed accordingly: `75e8f44072b0dd598dfa95aaf9b5f2c60f956819@your_external_ip:35622`.

### Prometheus collectors port

On *line 325* the port for Prometheus monitoring service is located: 

```bash
# Address to listen for Prometheus collector(s) connections
prometheus_listen_addr = ":26660"
```

It is useful if you want to monitor remotely the condition of your node using the [Prometheus](https://prometheus.io/) metrics collector service and could be changed to whatever you like `23456`:

```bash
prometheus_listen_addr = ":23456"
```

Don't forget to enable Prometheus metrics by changing to `true` on *line 322*, if needed:

```bash
# When true, Prometheus metrics are served under /metrics on
# PrometheusListenAddr.
# Check out the documentation for the list of available metrics.
prometheus = true
```

### External address

On *line 169* you should find the following:

```bash
# Address to advertise to peers for them to dial
# If empty, will use the same port as the laddr,
# and will introspect on the listener or use UPnP
# to figure out the address.
external_address = ""
```

This line implies specifying your external IP address, which means the presence of a static external address at your network connection. If you don't have one, just skip it.

- For Validator nodes, you may skip it, until you have **enough private peers** to get synced with. Otherwise, you have to specify your external static IP to improve peer discovery for your node. Also, don't forget to change the port according to *line [163](#cyber-communication-port)*:

```bash
external_address = "tcp://<your_external_static_ip>:26656"
```

- For Sentinel nodes it is a good idea to specify the IP for better peer discovery:

```bash
external_address = "tcp://<your_external_static_ip>:26656"
```

- For Servie nodes this setting can be the same as for Sentinel nodes:

```bash
external_address = "tcp://<your_external_static_ip>:26656"
```

And again, all of the above settings apply to the cases when **STATIC EXTERNAL IP** is available.

### Allow duplicated IP's

*Line 224* of the config.toml holds the following:

```bash
# Toggle to disable guard against peers connecting from the same ip.
allow_duplicate_ip = false
```

This variable configures the possibility for different peers to be connected from the same IP. Lets imagine a situation where you run 2 nodes (lets say that the node ID of the first one is: `75e8f44072b0dd598dfa95aaf9b5f2c60f956819` and the second one is: `d0518ce9881a4b0c5872e5e9b7c4ea8d760dad3f`) on one internet provider, with an external IP of 92.23.45.123. In this case, all other nodes in the network with `allow_duplicate_ip = false` will see attempts to connect from peers `d0518ce9881a4b0c5872e5e9b7c4ea8d760dad3f@92.23.45.123:26656` and `75e8f44072b0dd598dfa95aaf9b5f2c60f956819@92.23.45.123:36656` and will block the one which comes last because the originating IP address is the same for both nodes. If this case applies to you, change this setting to the following:

```bash
# Toggle to disable guard against peers connecting from the same ip.
allow_duplicate_ip = true
```

## P2P configuration

### Seed nodes

On *line 172* of the config.toml we see the following:

```bash
# Comma separated list of seed nodes to connect to
seeds = ""
```

This line is dedicated to the list of seed nodes you want to establish a connection with. To get some seed and peers addresses take a look at networks [repository](https://github.com/cybercongress/networks).

- For validators with sentinel nodes or with a decent quantity of peers connected it is not required to fill it out:

```bash
seeds = ""
```

- For Sentinel nodes and Service nodes it's a good idea to fill it out with a couple of seed node addresses, separated with commas: 

```bash
seeds = "<seed_node1_ID>@<seed_node1_ip>:<port>,<seed_node2_ID>@<seed_node2_ip>:<port>"
```

### Persistent peers

The place to add persistent peers is located on *line 175.* Presence of persistent peers is very important for the correct functioning of the node:

```bash
# Comma separated list of nodes to keep persistent connections to
persistent_peers = ""
```

- For Validator nodes you have to fill out this line with a decent amount of peers you **trust**, otherwise, your validator node address will be exposed. In the perfect case scenario, you should add to this section only the addresses of your sentinel nodes:

```bash
persistent_peers ="<sentinel_node1_ID>@<sentinel_node1_ip>:<port>,<sentinel_node2_ID>@<sentinel_node2_ip>:<port>"
```

- For Sentinel nodes and Service nodes add as many peers as possible to keep a persistent connection and network stability, but **DO NOT** put here you'r validator nodes ID's:

```bash
persistent_peers ="<node1_ID>@<node1_ip>:<port>,<node2_ID>@<node2_ip>:<port>,...,<node_n_ID>@<node_n_ip>:<port>"
```

### Peer Exchange Reactor

*Line 212* shows by default:

```bash
# Set true to enable the peer-exchange reactor
pex = true
```

This is a peer exchange module, which is responsible for exchanging node IDs across the network.

- For Validator nodes with Sentinel architecture set this to be disabled:

```bash
pex = false
```

- For Sentinel nodes and Service nodes leave as default:

```bash
pex = true
```

### Private peers ID's

On *line 221* we see:

```bash
# Comma separated list of peer IDs to keep private (will not be gossiped to other peers)
private_peer_ids = ""
```

This is the list of peers which IDs should not gossip to others.

- For Validator nodes, leave as default:

```bash
private_peer_ids = ""
```

- For Sentinel nodes, add your validator/s address  here:

```bash
private_peer_ids = "<validator_node_ID>@<validator_node_ip>:<port>"
```

- For Service nodes leave blank:

```bash
private_peer_ids = ""
```

## Node Index, Naming

### Indexed tags

A node can index and store a decent amount of keys and values with regards to transactions, accounts etc. *Lines 306 and 304* are responsible for this:

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

- For Validator and Sentinel nodes this is not necessary, so leave as default: 

```bash
index_keys = ""

index_all_keys = false
```

- For Service nodes, you should specify a subset of keys you want to index:

```bash
index_keys = "tx.hash,tx.height,...etc.."

index_all_keys = true
```

### Naming

To setup up your node moniker please refer to *line 16* and type in whatever you want to have as moniker:

```bash
# A custom human readable name for this node
moniker = "rocket_node"
```
