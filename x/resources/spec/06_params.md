# Parameters

The resources module contains the following parameters:

| Key                         | Type           | Example                |
| --------------------------- | -------------- | ---------------------- |
| MaxSlots                    | uint32         | 8                      | 
| BaseHalvingPeriodVolt       | uint32         | 12000000 (blocks)      |
| BaseHalvingPeriodAmpere     | uint32         | 12000000 (blocks)      |
| BaseInvestmintPeriodVolt    | uint32         | 3000000  (seconds)     |
| BaseInvestmintPeriodAmpere  | uint32         | 3000000  (seconds)     |
| BaseInvestmintAmountVolt    | sdk.Coin       | 1000000000 BOOT        |
| BaseInvestmintAmountAmpere  | sdk.Coin       | 1000000000 BOOT        |
| MinInvestmintPeriodSec      | uint32         | 86400 (seconds)        |


## Max Slots
Maximum amount of active slots at same amount of time

## Base Halving Period Volt
Period of blocks to mint rate' halving for VOLT resource

## Base Halving Period Ampere
Period of blocks to mint rate' halving for AMPERE resource

## Base Investmint Period Volt
Base amount of basic resource (BOOT) to invesmint  to get VOLT resource

## Base Investmint Period Ampere
Base amount of basic resource (BOOT) to invesmint  to get AMPERE resource

## Base Investmint Amount Volt
Base length of lock period in seconds of basic resource (BOOT) to get VOLT resource

## Base Investmint Amount Ampere
Base length of lock period in seconds of basic resource (BOOT) to get AMPERE resource

## Min Investmint Period Sec
Minimum length of lock period for all investmint operations