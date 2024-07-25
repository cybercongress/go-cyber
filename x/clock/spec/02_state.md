<!--
order: 2
-->

# State

## State Objects

The `x/clock` module only manages the following object in state: ClockContract. This object is used to store the address of the contract and its jail status. The jail status is used to determine if the contract should be executed at the end of every block. If the contract is jailed, it will not be executed.

```go
// This object is used to store the contract address and the
// jail status of the contract.
message ClockContract {
    // The address of the contract.
    string contract_address = 1;
    // The jail status of the contract.
    bool is_jailed = 2;
}
```

## Genesis & Params

The `x/clock` module's `GenesisState` defines the state necessary for initializing the chain from a previously exported height. It simply contains the gas limit parameter which is used to determine the maximum amount of gas that can be used by a contract. This value can be modified with a governance proposal.

```go
// GenesisState - initial state of module
message GenesisState {
  // Params of this module
  Params params = 1 [
    (gogoproto.nullable) = false,
    (gogoproto.jsontag) = "params,omitempty"
  ];
}

// Params defines the set of module parameters.
message Params {
  // contract_gas_limit defines the maximum amount of gas that can be used by a contract.
  uint64 contract_gas_limit = 1 [
    (gogoproto.jsontag) = "contract_gas_limit,omitempty",
    (gogoproto.moretags) = "yaml:\"contract_gas_limit\""
  ];
}
```

## State Transitions

The following state transitions are possible:

- Register a contract creates a new ClockContract object in state.
- Jailing a contract updates the is_jailed field of a ClockContract object in state.
- Unjailing a contract updates the is_jailed field of a ClockContract object in state.
- Unregister a contract deletes a ClockContract object from state.