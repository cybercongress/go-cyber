# State Transitions

These messages (Msg) in the resources module trigger state transitions.

## Investmint
1. Check that the desirable period of lock is less than the currently available maximum period of lock
2. Check that the spendable balance is more than or equal to the desired amount to lock
3. Check that the desirable period of lock is more than the current minimum period of lock
4. Initialize periodic vesting account on the first investmint
6. Put provided neuron's basis resources (BOOT) into vesting schedule to desirable period to empty slot
    - add vesting period and update periodic vesting account (auth module)
7. Calculate the amount of resources to mint (VOLT/AMPERE/...)
8. Check that calculated amount more than 1000 units (1000 milliampere or 1000 millivolt for example)
9. Mint calculated amount of given resource to neuron and lock that resource to same lock period to the same slot
    - mint tokens to module `resources` account (bank module)
    - transfer from module account `resource` to neuron's account (bank module)
    - add and save vesting period to same slot and update periodic vesting account (auth module)
10. Check if it's first neuron's investmint operation and if so then charge personal neuron's bandwidth with 1000 bandwidth units (bandwidth module)
11. Increase desirable bandwidth if there was investmint operation to VOLT resource (bandwidth module)

Note: neuron will be initially charged only if it will be first investmint operation to VOLT resource. This removes the gap to start neuron activity as soon as possible.

Note: total supply of VOLT resource is desirable bandwidth of the network.

## Slots logic (simplified)
1. Basic resource (BOOT) and desirable resource (VOLT/AMPERE/...) goes to the same vesting slot
2. If there are empty slots and all current slots are active then add a new one as active and reorder slots
3. If there are some expired slots then clean them and reorder currently active slots and add a new one as an active slot
4. If there are all slots active (amount of active slots is equal to max_slots parameter) then return an error
5. If all slots are passed then clean all of them and put new as active slot

Note: active slot => unlock time at future, expired slot => unlock time in past

Note: for fully understanding of slots logic do research of `addCoinsToVestingSchedule` function

## Investmints calculation

```
cycle = neuron's desirable period of lock / base period for given resource
base = neurons's desirable amount to lock / base amount for given resource
halving = 2^(current block height / halving period)  

=>

mint and lock (cycles * base * halving) VOLT or AMPERES to given neuron 
```

Note: VOLT and AMPERE as basic native computer's resources have separated base period/amount/halving parameters  that
adjustable with governance and dynamic cybernetic feedback loops (in a future release)

Note: for fully understanding of resources economy do research of `Mint` function




