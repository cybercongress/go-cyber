# `Resources`

## Abstract

The resources module allows neurons to invest into computer's resources.

### Examples:
I would like to `investmint` 1000000000 HYDROGEN (1 GH) to VOLT resource with lock for 30 DAYS (no spendable) and I will get
newly minted 1 VOLT to my account locked for 30 DAYS (no spendable).
```
(1 GH | 30 DAYS | VOLT) ---investmint---> locked (1 GH | 30 DAYS) + minted and locked (1 VOLT | 30 DAYS)
```

I would like to `investmint` 4000000000 HYDROGEN (4 GH) to AMPERE resource with lock for 7 DAYS (no spendable) and I will get
newly minted 1 AMPERE to my account locked for 7 DAYS (no spendable).
```
(4.2 GH | 7 DAYS | AMPERE) ---investmint---> locked (4.2 GH | 7 DAYS) + minted and locked (1 AMPERE | 7 DAYS)
```

## Contents

1. **[Concepts](00_concepts.md)**
2. **[API](01_api.md)**
3. **[State](02_state.md)**
4. **[State Transitions](03_state_transitions.md)**
5. **[Messages](04_messages.md)**
6. **[Events](05_events.md)**
7. **[Parameters](06_params.md)**
8. **[WASM](07_wasm.md)**
9. **[Errors](08_errors.md)**
10. **[CLI](09_cli.md)**