# State Transitions

These messages (Msg) in the energy module trigger state transitions.

```
Value may represents {V}, {A} and {A,V}

Route Key-Value storage is `0x00 | Source | Destination -> ProtocolBuffer(Route)`

Value Key-Value storage is `0x01 | Destination -> ProtocolBuffer(Value)`
```

## Route creation
1. Save destination account if account not already exist in storage (sdk's account keeper)
2. Save route with given source, destination, alias and empty route value if route not already exist
3. Call `OnCoinsTransfer` cyberbank's module hook to trigger index of balances to update with destination account (not storage)

## Route set
1. If route have zero value than escrow new value into `EnergyGrid` account
2. If route have old value more than new value to set than send `SendCoinsFromModuleToAccount` from `EnergyGrid` account to source account difference between old and new values (old-new)
3. If route have old value less than new value to set than escrow with `SendCoinsFromAccountToModule` from source account to `EnergyGrid` account difference between new and old values (new-old)
4. Update routed to destination energy value (add of sub value)
5. Update route with updated value
6. Call `OnCoinsTransfer` cyberbank's module hook to trigger index of balances to update with source and destination accounts (not storage)

## Route's alias edit
1. Update route with new alias

## Route remove
1. Send `SendCoinsFromModuleToAccount` from `EnergyGrid` account to source account route value (volts and amperes)
2. Update routed to destination energy value (sub value)
3. Remove route
4. Call `OnCoinsTransfer` cyberbank's module hook to trigger index of balances to update with source and destination accounts (not storage)

## Import routes on genesis
1. Initialize or update (add) routed to destination value
2. Set route
3. Call `OnCoinsTransfer` cyberbank's module hook to trigger index of balances to update with source and destination accounts (not storage)