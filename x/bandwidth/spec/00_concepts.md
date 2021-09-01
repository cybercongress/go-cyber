## Bandwidth model
Bandwidth is used for billing for cyberlinks creation instead of fee billing based on gas.
Using the bandwidth model removes the cognitive gap for cyberlinks creation because not force neurons to pay a fee.
Holding 1 volt represents the possibility to create 1 cyberlinks per ```RecoveryPeriod``` blocks. (example 1 cyberlink per day holding 1 Volt)
Total supply of VOLTs (volts' holdings of all neurons)  it's desirable network bandwidth. By investminting of HYDROGEN to VOLT then neuron increase personal bandwidth and desirable bandwidth of network.
If the network has low bandwidth consumption, then the network provides a discount using ```BasePrice``` multiplier for neurons.

## Personal Bandwidth
The volt's stake of the given neuron are easy to understand as the size of his battery.
The creation of cyberlinks will consume battery charge, and the battery will be fully recharged during ```RecoveryPeriod``` blocks period. If a neuron consumes half of its bandwidth, its battery will be fully charged in the ```RecoveryPeriod/2``` blocks period. If a neuron act when network bandwidth consumption is low, then he will consume less personal bandwidth.

## Bandwidth Load
The network keeps track of total consumed bandwidth during ```RecoveryPeriod``` and weights this window aggregated value to desirable bandwidth. If the current load is more than ```BasePrice```, the price value will be set.

## Bandwidth Price
Bandwidth price it's a multiplier for default bandwidth price. As 1 VOLT allows creating 1 cyberlink per given period, if the price is lower than 1, the network will consume less of bandwidth, allowing neurons to generate more cyberlinks.

## Accounting of bandwidth
Internally 1 Volt represents 1000 millivolts, and 1 cyberlink cost is 1000 bandwidth units. Neurons holdings of 5 Volts means 5000 personal bandwidth units. When the current load is less than BasePrice amount (e.g 0.25), then the network will make the discount for bandwidth bill 4X allowing neurons to create 4X more cyberlinks, 20 cyberlinks in such case.

## Transaction's mempool
For transactions that consist of cyberlinks, a fee check will not apply. But correct required gas amount should be provided.

## Network capacity
The total amount of minted (with investmint) VOLTs represents the demand of bandwidth from neurons. Validators need to keep tracking investments in VOLTs resources to provide great service at scale.
To dynamically adjust available peek load community may adjust maximum bandwidth ```MaxBlockBandwidth``` and gas ```MaxGas``` consumable at block.
