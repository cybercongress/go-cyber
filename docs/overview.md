# ://cyber staking and distribution overview

## Params

//todo consider:
// Consensus.BlockTime
// Consensus.MaxBytes
// Staking.UnbondingTime && Slashing.MaxEvidenceAge
// Distr.BaseProposerReward
// Distr.BonusProposerReward
// Slashing.SignedBlocksWindow
// Slashing.MinSignedPerWindow

| Module      | Param         | Value         | Comment                                 |
|-------------|---------------|---------------|-----------------------------------------|
| Staking     | UnbondingTime | 3 weeks       | time duration of unbonding              |
| Staking     | MaxValidators | 100           | maximum number of active validators set |
| Staking     | MaxEntries    | 7             | max entries for either unbonding delegation or redelegation per delegator/validator pair(delegator/validator/validator for redelegation)|
|             |                  |       |                                                    |
| Consensus   | BlockTime        | 5s    | average block generation time                      |
| Consensus   | MaxBytes         | 1mb   | block max bytes limit                              |
| Consensus   | RankCalcWindow   | 600   | full rank recalculation window                     |
|             |                     |     |                                                                        |
| Distr       | CommunityTax        | 0%  | community funding tax, not used                                        |
| Distr       | BaseProposerReward  | 1%  | % of block inflation goes to proposer                                  |
| Distr       | BonusProposerReward | 4%  | addition reward, calculated as % of included votes from validators set |
|             |                          |                   |                                                |
| Slashing    | MaxEvidenceAge           | 3weeks            | misbehaviour evidence max age                  |
| Slashing    | SignedBlocksWindow       | 900               | window to calculate validators liveness        |
| Slashing    | MinSignedPerWindow       | 5%                | min singed block for window to not be jailed   |
| Slashing    | DowntimeJailDuration     | 0                 | unjail delay                                   |
| Slashing    | SlashFractionDoubleSign  | 20%               | % of stake reduction for double sign           |
| Slashing    | SlashFractionDowntime    | 0%                | % of stake reduction for being offline         |
|             |                 |                   |                              |
| Mint        | BlockInflation  | 0.317097920 Gcyb  | validators block rewards     |



## Staking

The cyberd is a public Proof-Of-Stake (PoS) blockchain, 
 meaning that validator's weight is determined by the amount of staking tokens bonded as collateral. 
 These tokens can be staked directly by the validator or delegated to them by token holders.
 The weight (i.e. total stake) of a validator determines whether or not it is an active validator, 
 and also how frequently this node will have to propose a block and how much revenue it will obtain.

### Validator

Any user in the system can declare its intention to become a validator by sending a create-validator transaction. 
 From there, they become validators.
 Validator can set **commission**, that applied on revenue before it is distributed to their delegators.

Each validator holds:
- All bounded tokens(self and delegators). NOTE: not include distribution rewards.
- Own distribution rewards (commission rewards)
- Delegators distribution rewards
- All delegators shares. Share is not mapped 1-to-1 to tokens. 
  In a case a validator being punished for misbehaviour, bounded tokens will be reduced, while shares remain a same. 

### Delegation

Delegators are token holders who cannot, or do not want to run validator operations themselves. 
 A user can delegate tokens to a validator and obtain a part of its revenue in exchange.
 Upon delegation a user converts his tokens to validator shares in a rate `val_tokens/val_shares`. 
  
### Undelegation

A user may want to cancel delegation to specific validator. To do so, he/she send **undelegate** transaction.
 Depending on current validator state, either user receive his revenue proportion and bounded tokens back immediately 
 (for unbonded validator), or just start process of undelegation. 
 If a validator is in unbonding state, than a user will receive tokens at a validator unbonding time. 
 In last case, a user will wait full **UnbondingTime** period.

## Slashing

If validators double sign, are frequently offline or do not participate in governance, 
their staked tokens (including tokens of users that delegated to them) can be destroyed, or 'slashed'.

At the beginning of each block, we update the signing info for each validator and 
 check if they've dipped below the liveness threshold **MinSignedPerWindow** 
 over the tracked window **SignedBlocksWindow**. 
 If so, their stake will be slashed by **SlashFractionDowntime** percentage and 
 will be Jailed for **DowntimeJailDuration**.

## Distribution

All minted tokens goes to fees pool.
 At each **beginblock**, the fees received on previous block are allocated to the proposer, community fund, 
 and previous block active validators set according to next scheme:
 
1. When the validator is the proposer of the round, that validator (and their delegators) 
 receives between **BaseProposerReward** and **BonusProposerReward** of fee rewards. 
 The amount of proposer reward is calculated from pre-commits Tendermint messages 
 in order to incentives validators to wait and include additional pre-commits in the block.
 
2. Community tax is then charged from full fees.

3. The remainder is distributed proportionally by voting power 
 to all bonded validators(and their delegators) independent of whether they voted (social distribution).
 