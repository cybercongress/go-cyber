require('dotenv').config({path:'.env-node3'})
let lotion = require ('lotion')
let app = lotion({
    genesis:'./genesis.json',
    tendermintPort:30092,
    p2pPort: 30094,
    initialState:{ },
    logTendermint: true,
    keys: 'privkey2.json',
    peers: ['192.168.0.102:30091', '192.168.0.130:30093']
})

app.use(require('./core'));


app.listen(3000).then(({GCI}) =>{
    console.log(GCI);
})