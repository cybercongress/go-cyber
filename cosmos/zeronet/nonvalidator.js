require('dotenv').config({path:'.env-node3'})
let lotion = require ('lotion')
let app = lotion({
    genesis:'./genesis.json',
    tendermintPort:46657,
    initialState:{ },
    logTendermint: true,
    peers:['ws://localhost:30091','ws://localhost:30093']
})

app.use(require('./core'));

app.listen(3002).then(({GCI}) =>{
    console.log(GCI);
})