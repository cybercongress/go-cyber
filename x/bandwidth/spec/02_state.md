# State

## Account Bandwidth

AccountBandwidth is used for tracking bandwidth of accounts in the network.

AccountBandwidth type has the following structure:
```
type AccountBandwidth struct {
    address          string     // address of neuron
	remainedValue    uint64     // current bandwidth value 
	lastUpdatedBlock uint64     // last block when last time updated
	maxValue         uint64     // max current bandwidth value of neuron
}
```


## Last price
Value is used to store up-to-date price of bandwidth.

```
type Price struct {
    price          sdk.Dec   // current multiplier for bandwidth billing
}
```

## Block bandwidth
Storing used bandwidth for each block. Used for calculation of load using sum of used bandwidth in blocks at recovery period window.
Used for reverting transactions with cyberlinks if rise more than ```MaxBlockBandwidth```

```
sdk.Uint64ToBigEndian(value) // where value is amount of bandwidth used by all neurons in given block
```

## Desirable bandwidth
Desirable bandwidth represents amount of cyberlinks that network would like to process.

```
sdk.Uint64ToBigEndian(value) // where value is total current supply of mvolt (uint64)
```

-------

## Keys

- Account bandwidth: `0x01 | []byte(address) -> ProtocolBuffer(AccountBandwidth)`
- Block bandwidth: `0x02 | sdk.Uint64ToBigEndian(blockNumber) -> sdk.Uint64ToBigEndian(value)`
- Last bandwidth price: `0x00 | []byte("lastBandwidthPrice") -> ProtocolBuffer(Price)`
- Desirable bandwidth: `0x00 | []byte("desirableBandwidth") -> sdk.Uint64ToBigEndian(value)`
- ModuleName, StoreKey, QuerierRoute: `bandwidth`