# Messages

Messages (Msg) are objects that trigger state transitions. 
Msgs are wrapped in transactions (Txs) that clients submit to the network. 
The Cosmos SDK wraps and unwraps energy module messages from transactions.

## MsgCreateRoute

A route is created with zero initial value with the `MsgCreateRoute` message.

```go
type MsgCreateRoute struct {
    source       string     // account address of the origin of this message
    destination  string     // account address routing to
    alias        string     // related alias or tag
}
```

### Validity Checks

Validity checks are performed for MsgCreatePool messages. The transaction that is triggered with `MsgCreateRoute` fails if:

- Stateless
    - if msg.Source is not valid address
    - if msg.Destination is not valid address
    - if msg.Alias is equal 0 or more than 32 bytes
- Stateful 
    - if route with given msg.Source and msg.Destination exist
    - if total amount of routes from msg.Source in store more than or equal max_routes param.

## MsgEditRoute

SET value (volts or amperes) to given route with the `MsgEditRoute` message.

```go
type MsgEditRoute struct {
    source       string       // account address of the origin of this message
    destination  string       // account address editing to
    value        sdk.Coin     // value volts are amperes to set
}
```

Note: SET means that you are setting value of volts and amperes to be routed to. 
Depending of the state of the given route msg.Source will need to provide more coins (if value increased) to `EnergyGrid` or will receive coins (if value decreased) from `EnergyGrid`.

### Validity Checks

Validity checks are performed for MsgEditRoute messages. The transaction that is triggered with `MsgEditRoute` fails if:

- Stateless
    - if msg.Source is not valid address
    - if msg.Destination is not valid address
    - if msg.Value denom is not equal denoms of volts or amperes
- Stateful
    - if Route with given msg.Source and msg.Destination not exist
    - if the balance of `Source` does not have enough amount of coins for `SendCoinsFromAccountToModule` to `EnergyGrid`

## MsgDeleteRoute

Delete the given route with the `MsgDeleteRoute` message.

```go
type MsgDeleteRoute struct {
    source       string       // account address of the origin of this message
    destination  string       // account address deleting to
}
```

### Validity Checks

Validity checks are performed for MsgDeleteRoute messages. The transaction that is triggered with `MsgDeleteRoute` fails if:

- Stateless
    - if msg.Source is not valid address
    - if msg.Destination is not valid address
- Stateful
    - if Route with given msg.Source and msg.Destination not exist

## MsgEditRouteAlias

SET value (volts or amperes) to given route with the `MsgEditRoute` message.

```go
type MsgEditRouteAlias struct {
    source       string       // account address of the origin of this message
    destination  string       // account address editing to 
    alias        string       // new related alias or tag
}
```

### Validity Checks

Validity checks are performed for MsgEditRoute messages. The transaction that is triggered with `MsgEditRoute` fails if:

- Stateless
    - if msg.Source is not valid address
    - if msg.Destination is not valid address
    - if msg.Alias is equal 0 or more than 32 bytes
- Stateful
    - if Route with given msg.Source and msg.Destination not exist
