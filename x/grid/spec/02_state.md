# State

## Route

Route is used for tracking amount of volts and amperes routed to given account.

Route type has the following structure:
```
type Route struct {
    source         string       // source account address, route from
    destination    string       // destination account address, route to
    alias          string       // releated alias or tag
    value          []sdk.Coin   // amount of volts and amperes in this route
}
```

## Value
Value is used to store up-to-date summarized value of all routes routed to given destination.

```
type Value struct {
    value          []sdk.Coin   // amount of volts and amperes in this route
}
```

-------

## Keys

- Route: `0x00 | Source | Destination -> ProtocolBuffer(Route)`
- Value: `0x01 | Destination -> ProtocolBuffer(Value)`
- ModuleName, RouterKey, StoreKey: `energy`
- EnergyPoolName: `energy_grid`