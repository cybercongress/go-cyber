require('dotenv').config({path:'.env-node1'})
let lotion = require ('lotion');

var http_port = process.env.HTTP_PORT || 3001;
var p2p_port = process.env.P2P_PORT || 30091;
var initialPeers = process.env.PEERS ? process.env.PEERS.split(',') : [];
var key = process.env.KEY;

console.log('HTTP_PORT: ', http_port);
console.log('P2P_PORT: ', p2p_port);
console.log('PEERS: ', initialPeers);
console.log('LOTION_HOME: ', process.env.LOTION_HOME);
console.log('KEY: ', key);


let app = lotion({
    genesis:'./genesis.json',
    // tendermintPort: 30090,
    // p2pPort: p2p_port,
    // initialState:{ },
    // logTendermint: true,
    // keys: key,
    // peers: initialPeers
})

app.use(require('./core'));

app.listen(http_port).then(({GCI}) =>{
    console.log(GCI);
})