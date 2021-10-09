# CLI

## Queries

Query all energy routes that made from source account:
```bash
cyber query grid routes-from [source-addr]
```

Query all energy routes that routed to destination account:
```bash
cyber query grid routes-to [destination-addr]
```

Query energy value that routed from source account:
```bash
cyber query grid routed-from [source-addr]
```

Query energy value that routed to destination account:
```bash
cyber query grid routed-to [destination-addr]
```

Query energy route that routes for given source and destination accounts:
```bash
cyber query grid route [source-addr] [destination-addr]
```

Query all energy routes (pagination flags supported):
```bash
cyber query grid routes
```

## Transactions

Create energy route from your address to destination address with provided alias:
```bash
cyber tx grid create-route [destination-addr] [alias]
```

Set value of energy route to destination address:
```bash
cyber tx grid edit-route [destination-addr] [value]
```

Delete your energy route to given destination address:
```bash
cyber tx grid delete-route [destination-addr]
```

Edit alias of energy route to given destination address:
```bash
cyber tx grid edit-route-alias [destination-addr] [new-alias]
```