require('dotenv').config({path:'.env-node0'})
let lotion = require ('lotion')
let app = lotion({
    genesis:'./genesis.json',
    tendermintPort:46657,
    initialState:{ },
    logTendermint: true,
    peers:[
	    'ws://192.168.0.102:30091', 
	    'ws://192.168.0.103:30094', 
	    'ws://192.168.0.130:30093'
    ]
})

app.use(require('./core'));

app.listen(3002).then(({GCI}) =>{
    console.log(GCI);
})