# State Transitions [WIP]

Bandwidth module doesn't have own messages that trigger state transition.
State transition is happen in such cases:

## Processing of transaction with cyberlinks messages in transaction middleware (ante handler)
1. calculate total bandwidth amount for all cyberlinks messages in transaction using current price and consume neuron's bandwidth
2. add consumed bandwidth to block bandwidth (in-memory)

## Processing of cyberlink message created by VM contract (graph module)
1. calculate bandwidth for message using current price and consume neuron's bandwidth
2. add consumed bandwidth to block bandwidth (in-memory)

Note: billing happens in the graph module for contracts because contracts creates messages not grouped into transactions (ante handler are not processing them)

## Transfers of Volts (EndBlocker)
1. Update account's bandwidth for an account with changed stake collected by ```CollectAddressesWithStakeChange``` hook (e.g transfer of investmint).

Note: minting of new volts (using investmint) will trigger the account's bandwidth update with an increased max bandwidth value

## Save consumed bandwidth by block (EndBlocker)
1. Save the total amount (sum aggregated in-memory before) of consumed bandwidth by all neurons on a given block (to storage & in-memory).
2. Remove value for a block that is out of ```RecoveryWindow``` block period and not perform in bandwidth load calculation (to storage & in-memory).

## Adjust bandwidth price (EndBlocker)
1. If block height number's remainder of division by ```AdjustPrice``` parameter is equal to zero
   - calculate and save price based on current load (or apply ```BasePrice``` if load less than ```BasePrice```).

## Genesis
1. If neuron have volts in genesis 
   - initialize and save account bandwidth with max value
