# API reference 

Cyberd provides a [JSON-RPC](http://json-rpc.org/wiki/specification) API. Http endpoint is served under 
 `localhost:20657`. WebSockets are the preferred transport for cyberd RPC and are used by applications such as cyb. 
 Default WebSocket connection endpoint for cyberd is `ws://localhost:20657/websocket`. There are test endpoints 
 available at `http://earth.cybernode.ai:34657` and `ws://earth.cybernode.ai:34657/websocket`.

<br />

## Standard Methods

### Query Example

Query http endpoint using curl:
```bash
curl --data '{"method":"cyberd_mostResentBlockNumber","params":[],"id":"1","jsonrpc":"2.0"}' \
-H "Content-Type: application/json" -X POST earth.cybernode.ai:34657
```

Query ws endpoint from js:
```js
let websocket = new WebSocket("ws://earth.cybernode.ai:34657/websocket");
websocket.send(JSON.stringify({
  "method":"cyberd_mostResentBlockNumber",
  "params":[],
  "id":"1",
  "jsonrpc":"2.0"
}));
```


### Method Overview

The following is an overview of the RPC methods and their current status.  Click
the method name for further details such as parameter and return information.

|#|Method|Description|
|---|------|-----------|
|1|[mostResentBlockNumber](#mostResentBlockNumber)|Returns the number of most recent block.|

### Method Details

***
<a name="mostResentBlockNumber"/>

|   |   |
|---|---|
|Method|cyberd_mostResentBlockNumber|
|Parameters|None|
|Description|Returns the number of most recent block.|
|Returns|`int`|
|Example Return |`42`|
|[Return to Overview](#method-overview)<br />

***
<br />
<br />

## Notifications (WebSocket-specific)

Cyberd uses standard JSON-RPC notifications to notify clients of changes, rather than requiring clients to poll cyberd
 for updates. JSON-RPC notifications are a subset of requests, but do not contain an ID. The notification type 
 is categorized by the `query` params field.
 
### Subscribe Example 
Subscribe for new blocks header from js:
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
|1|[NewBlockHeader](#NewBlockHeader)|Sends block header notification when a new block is committed.|

### Events Details

#### NewBlockHeader    
|   |   |
|---|---|
|Event|NewBlockHeader|
|Description|Sends block header notification when a new block is committed.|
|Query|'tm.event=\'NewBlockHeader\''|
|[Return to Overview](#events-overview)<br />
```json
{
  "jsonrpc": "2.0",
  "id": "1#event",
  "result": {
    "query": "tm.event='NewBlockHeader'",
    "data": {
      "type": "tendermint/event/NewBlockHeader",
      "value": {
        "header": {
          "chain_id": "test-chain-gRXWCL",
          "height": "40561",
          "time": "2018-11-08T14:40:11.820674115Z",
          "num_txs": "0",
          "total_txs": "510",
          "last_block_id": {
            "hash": "CC1693981B0353C907EBC2EBEB1578ABD21C955D",
            "parts": {
              "total": "1",
              "hash": "E854DA23981283464484FEA41D8AB5FF7697728C"
            }
          },
          "last_commit_hash": "51E1DEBA90591C532E8C58BF99F29AE4B2264644",
          "data_hash": "",
          "validators_hash": "985037CEBE01051C949F01278621449712C85715",
          "next_validators_hash": "985037CEBE01051C949F01278621449712C85715",
          "consensus_hash": "0E520AF30D47BE28F293E040E418D0361BFB5370",
          "app_hash": "2DBED732547DD6A3223027EF8C80FD3C3A1AD11420CA899533DCFCA260F8D170",
          "last_results_hash": "",
          "evidence_hash": "",
          "proposer_address": "A41C2ED742E47A9BDC4E71BE4B879F949045EC97"
        }
      }
    }
  }
}
```




