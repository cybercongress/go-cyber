# Parameters

The energy module contains the following parameters:

| Key                     | Type           | Example                |
| ----------------------- | -------------- | ---------------------- |
| RecoveryPeriod          | uint64         | 16000                  |
| AdjustPricePeriod       | uint64         | 5                      |
| BasePrice               | sdk.Dec        | 0.25                   |
| MaxBlockBandwidth       | uint64         | 100000                 |

## Recovery Period
Recovery period is amount of blocks that bandwidth of any neuron will be fully recovered (from zero to max value).

## Adjust Price Period
Adjust price period is amount of blocks that form period of bandwidth price recalculation adapting to current load.

## Base Price
Base price is multiplier for bandwidth billing, bandwidth discount for moments of low neurons' activity.
If load rise more than value of base price than current price will be applied.

## Max Block Bandwidth
Max block bandwidth is amount of bandwidth from neurons that network can process for one block.