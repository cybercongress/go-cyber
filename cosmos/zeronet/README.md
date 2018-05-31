### Quick start
```
docker-compose up --build
```

### test
```
curl -H "Content-type:application/json" --data '{"type" : "search", "keyword": "test" }' http://localhost:3002/txs

curl -H "Content-type:application/json" --data '{"type" : "link", "keyword": "test", hash": "hash" }' http://localhost:3002/txs

http://localhost:3002/state
http://localhost:3003/state
http://localhost:3004/state
```

### generate priv key
```
./node_modules/lotion/bin/tendermint gen_validator > privkey1.json
```