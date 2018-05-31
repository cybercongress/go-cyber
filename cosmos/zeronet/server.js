var path = require('path');

var express = require('express'); 
var expressapp = express();

var proxy = require('express-http-proxy');
expressapp.use('/', proxy('http://127.0.0.1:3000'));

expressapp.listen(5000, function () {
  console.log('express server on 5000');
});

var p2p_port = process.env.P2P_PORT || 30091;
var initialPeers = process.env.PEERS ? process.env.PEERS.split(',') : [];
var key = process.env.KEY;

console.log('KEY: ', key);
console.log('PEERS: ', initialPeers);
console.log('P2P_PORT: ', p2p_port);

let app = require('lotion')({
 	genesis: './genesis.json',
 	initialState: { },
 	
 	tendermintPort:30090,
 	p2pPort: p2p_port,

	logTendermint: true,
 	keys: key,
 	p2pPort: 30092,
 	peers: initialPeers
})

app.use(require('./core'))

app.listen(3000).then(({ GCI }) => {
  console.log(GCI)
})


