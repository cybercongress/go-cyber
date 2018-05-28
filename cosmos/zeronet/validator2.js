require('dotenv').config({path:'.env-node2'})
let lotion = require ('lotion')
let app = lotion({
    genesis:'./genesis.json',
    tendermintPort:30092,
    p2pPort: 30093,
    initialState:{ },
    logTendermint: true,
    keys: 'privkey1.json',
    peers: ['192.168.0.102:30091']
})

app.use(require('./core'));


app.listen(3001).then(({GCI}) =>{
    console.log(GCI);
})