require('dotenv').config({path:'.env-node1'})
let lotion = require ('lotion')
let app = lotion({
    genesis:'./genesis.json',
    tendermintPort:30090,
    p2pPort: 30091,
    initialState:{ },
    logTendermint: true,
    keys: 'privkey0.json',
    peers: ['192.168.0.130:30093', '192.168.0.103:30094']
})

app.use(require('./core'));

app.listen(3000).then(({GCI}) =>{
    console.log(GCI);
})