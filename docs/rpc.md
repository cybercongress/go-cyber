# API reference

Cyberd provides a [JSON-RPC](http://json-rpc.org/wiki/specification) type API. The HTTP endpoint is served under
 `localhost:20657`. WebSockets are the preferred transport for cyberd RPC and are used by applications, such as cyb.
 The default WebSocket connection endpoint for cyberd is `ws://localhost:26657/websocket`. There are test endpoints
 available at `https://titan.cybernode.ai/api/` and `ws://titan.cybernode.ai/websocket`.

<br />

## Standard Methods

### Query Example

Query the HTTP endpoint using curl:

```bash
curl --data '{"method":"status","params":[],"id":"1","jsonrpc":"2.0"}' \
-H "Content-Type: application/json" -X POST earth.cybernode.ai:34657
```

Query ws endpoint from JS:

```js
let websocket = new WebSocket("ws://earth.cybernode.ai:34657/websocket");
websocket.send(JSON.stringify({
  "method":"status",
  "params":[],
  "id":"1",
  "jsonrpc":"2.0"
}));
```

### Method Overview

The following is an overview of the RPC methods and their current status.  Click
the method name for further detail, such as parameter, and this will return information.

|#|Method|Description|
|---|------|-----------|
|1|[status](#status)|Get node info, pubkey, latest block hash, app hash, block height and time.|
|2|[account](#account)|Get account nonce, pubkey, number, and coins.|
|3|[account_bandwidth](#account-bandwidth)|Get account bandwidth info for current height.|
|4|[is_link_exist](#link-exist)|Return true, if given link exist.|
|5|[current_bandwidth_price](#current-bandwidth-price)|Returns current bandwidth credit price.|
|6|[index_stats](#index-stats)|Returns current index entities count.|

### Method Details

***
<a name="status"/>

|   |   |
|---|---|
|Method|status|
|Parameters|None|
|Description|Get node info, pubkey, latest block hash, app hash, block height and time.|
|[Return to Overview](#method-overview)<br />

<a name="account"/>

|   |   |
|---|---|
|Method|account|
|Parameters|1. address (string, required)<br />|
|Description|Get account nonce, pubkey, number, and coins.|
|[Return to Overview](#method-overview)<br />

<a name="account-bandwidth"/>

|   |   |
|---|---|
|Method|account_bandwidth|
|Parameters|1. address (string, required)<br />|
|Description|Get account bandwidth info for current height.|
|[Return to Overview](#method-overview)<br />

<a name="link-exist"/>

|   |   |
|---|---|
|Method|is_link_exist|
|Parameters|1. from (cid, required)<br />2. to (cid, required)<br />3. address (string, required)<br />|
|Description|Return true, if given link exist.|
|[Return to Overview](#method-overview)<br />

<a name="current-bandwidth-price"/>

|   |   |
|---|---|
|Method|current_bandwidth_price|
|Parameters|None|
|Description|Returns current bandwidth credit price.|
|[Return to Overview](#method-overview)<br />

<a name="index-stats"/>

|   |   |
|---|---|
|Method|index_stats|
|Parameters|None|
|Description|Returns current index entities count.|
|[Return to Overview](#method-overview)<br />

***
<br />
<br />

## Notifications (WebSocket-specific)

Cyberd uses standard JSON-RPC notifications to notify clients of changes rather than requiring clients to poll cyberd
 for updates. JSON-RPC notifications are a subset of requests, but do not contain an ID. The notification type 
 is categorized by the `query` params field.

### Subscribe Example 

Subscribe for new block headers from JS:

 ```js
 let websocket = new WebSocket("ws://earth.cybernode.ai:34657/websocket");
 websocket.send(JSON.stringify({
   "method": "subscribe",
   "params": ["tm.event='NewBlockHeader'"],
   "id": "1",
   "jsonrpc": "2.0"
 }));
 ```

### Events Overview

|#|Event|Description|
|---|------|-----------|
|1|[NewBlockHeader](#NewBlockHeader)|Sends a block header notification when a new block is committed.|
|2|[CoinsReceived](#CoinsReceived)|Sends a notification when new coinc arrive to a given address.|
|3|[CoinsSend](#CoinsSend)|Sends a notification when new coins are sent from a given address.|
|4|[СidsLinked](#СidsLinked)|Notification of links created by a given address.|
|5|[SignedTxCommitted](#SignedTxCommitted)|Notify when any tx for a given signer is committed.|

### Events Details

#### NewBlockHeader

|   |   |
|---|---|
|Event|NewBlockHeader|
|Description|Sends block header notification when a new block is committed.|
|Query|`tm.event='NewBlockHeader'`|
|[Return to Overview](#events-overview)<br />

#### CoinsReceived

|   |   |
|---|---|
|Event|CoinsReceived|
|Description|Sends a notification when new coins have arrived to a given address.|
|Query|`tm.event='EventTx' AND recipient='cbd1sk3uvpacpjm2t3389caqk4gd9n9gkzq2054yds'`|
|[Return to Overview](#events-overview)<br />

#### CoinsSend

|   |   |
|---|---|
|Event|CoinsSend|
|Description|Sends a notification when new coins are sent from a given address.|
|Query|`tm.event='EventTx' AND sender='cbd1sk3uvpacpjm2t3389caqk4gd9n9gkzq2054yds'`|
|[Return to Overview](#events-overview)<br />

#### СidsLinked
|   |   |
|---|---|
|Event|СidsLinked|
|Description|Notification of links created by a given address.|
|Query|`tm.event='EventTx' AND signer='cbd1sk3uvpacpjm2t3389caqk4gd9n9gkzq2054yds' AND action='link'`|
|[Return to Overview](#events-overview)<br />

#### SignedTxCommitted

|   |   |
|---|---|
|Event|SignedTxCommitted|
|Description|Notify when any tx for a given signer is committed.|
|Query|`tm.event='EventTx' AND signer='cbd1sk3uvpacpjm2t3389caqk4gd9n9gkzq2054yds'`|
|[Return to Overview](#events-overview)<br />





