# State

## Periodic Vesting Account (auth)

A vesting account implementation that vests coins according to a custom vesting schedule.

```
type PeriodicVestingAccount struct {
  base_vesting_account BaseVestingAccount 
  start_time           uint64
  vesting_periods      []Period
}
```

## Period (auth)

```
type Period struct {
	length int64                                   
	amount []sdk.Coin
}
```

## Keys

- ModuleName, RouterKey: `energy`