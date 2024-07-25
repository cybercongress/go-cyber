# `clock`

```text
Thanks to Juno and Tgrade teams for the initial development of this module.
This module is based on Juno clock's module implementation with few modifications:
1. Added support for ABCI BeginBlocker
2. TODO: Add support for gas limit for each contract
```

## Abstract

This document specifies the internal `x/clock` module of cyber-sdk based networks.

The `x/clock` module allows specific contracts to be executed at the end of every block. This allows the smart contract to perform actions that may need to happen every block or at set block intervals.

By using this module, your application can remove the headache of external whitelisted bots and instead depend on the chain itself for constant executions.

## Contents

1. **[Concepts](01_concepts.md)**
2. **[State](02_state.md)**
3. **[Contract Integration](03_integration.md)**
4. **[Clients](04_clients.md)**
