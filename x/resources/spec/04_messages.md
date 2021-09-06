# Messages

Messages (Msg) are objects that trigger state transitions. 
Msgs are wrapped in transactions (Txs) that clients submit to the network. 
The Cosmos SDK wraps and unwraps resources module messages from transactions.

## MsgInvestmint

Neuron make investmint operation and lock base resource amount to given length with return of desired advanced resource locked to same length.

```go
type MsgInvestmint struct {
    agent    string    // agent's address
    amount   sdk.Coin  // amount of basic resource to invesmint 
	resource string    // desirable resource
	length   uint32    // desirable lock period in seconds
}
```

The message will fail under the following conditions:

- Stateless
    - if msg.Agent is not valid address
    - if msg.Amount is invalid (not positive or wrong denom)
    - if msg.Resource is invalid (resource not exist)
    - if msg.Length is equal to 0
- Stateful 
    - if Msg.Length is more than available maximum lock period
    - if Msg.Length is less than available minimum lock period
    - if account have less spendable balance that Msg.Amount
    - if there are no empty or expired slots
    - if amount of resource in return less than 1000 milli{volt,amper}