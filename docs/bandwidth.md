# Minimal Bandwidth Spec

**bounded stake**. Stake, that **deducted** from your acc coins and put as **deposit** to take part in consensus. 
Due to passive inflation model and slashing, deposit not match 1-to-1 to final reward. 
So, for example, stakeholders may wish to set up a script,
 that will periodically withdraw and rebound rewards to increase their share(**bounded stake**). 

**active stake**. Currently available for direct transfer, not-bounded stake.

**bandwidth stake**. Sum of active stake and bounded stake for given account.

## Model

**Cyberd** use a very simple bandwidth model. 
Main goal of that model is to reduce daily network growth to given constant, say 3gb per day.

Thus, here we introduce **resource credits**(RS). Each message type have assigned **RS** cost. 
There is constant `DesirableNetworkBandwidthForRecoveryPeriod` determining desirable for `RecoveryPeriod` spent RS value.
`RecoveryPeriod` is defining how fast user can recover their bandwidth from 0 to user max bandwidth.
User has maximum **RS** proportional to his stake by 
 formula `user_max_rc = bandwidth_stake% * DesirableNetworkBandwidthForRecoveryPeriod`.
 
There is period `AdjustPricePeriod` summing how much RS was spent for that period(`AdjustPricePeriodTotalSpent`).
Also, there is constant `AdjustPricePeriodDesiredSpent`, used to calculate network loading. 
`AdjustPricePeriodTotalSpent/AdjustPricePeriodDesiredSpent` ratio defined so called current **price multiplier**.
If network usage is low, **price multiplier** adjust message cost(by simply multiplying) 
 to allow user with lower stake to do more transactions. 
If resource demand increase, **price multiplier** goes `>1` thus increase messages cost 
 and limiting final tx count for some long-term period(RC recovery will be `<` then RC spending).

## Bandwidth stake change

There are only few ways to change acc **bandwidth stake**:

1. Direct coins transfer.
2. When distribution payouts occurs. For example, when validator change his commission rates, 
 all delegations will be automatically unbounded. Another example, delegator itself unbound some part or full share.


**Implementation details.** 
All this cases will change user **active stake**, so we define bank hook `CoinsTransferHook`.
 Due to inconsistent stateupon hook invoking (coins can be substructed, but deposit not added yet), 
 we should update user max RC only after deliver tx occur.