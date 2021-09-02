# `cyberbank`

## Abstract

The cyberbank module is the module that wraps the original Cosmos-SDK's bank module while providing extra functionality:
1. tracks of neuron's VOLTs balance change for [bandwidth](../../bandwidth/spec/README.md) module to adjust neuron's personal bandwidth
2. keeps two in-memory indexes of AMPEREs balances for the [rank](../../rank/spec/README.md) module (one index for last calculation and one for the next one)
3. tracking of routed to neurons resources (VOLTs and AMPEREs) from [energy](../../energy/spec/README.md) module

Note: cyberbank wrapper keep all his state completely in-memory and not introduce any new modification to application storage  
  

## Contents

1. **[Concepts](00_concepts.md)**
3. **[State](02_state.md)**
4. **[State Transitions](03_state_transitions.md)**