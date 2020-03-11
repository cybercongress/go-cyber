# Testing script

## Linking

- [x] create link -> link created **@savetheales**
- [x] create existed link by the same address -> returns error **@savetheales**
- [x] broadcast tx with multiple links messages -> links created **@savetheales**
- [x] linkchain in the one cyberlink message -> links created **@savetheales**

____________

## Transaction

- [x] successed -> consumed bandwidth **@savetheales**
- [x] failed -> consumes band **@savetheales**
3. non-link msgs -> consumes needed band
- [x] link msgs -> consumes needed band and increasea personal karma **@savetheales**
- [x] link msgs -> increase global karma **@savetheales**
- [x] send msgs -> change balance **@savetheales**
- [x] send msgs -> change max bandwidth, band resume speed **@savetheales**
- [x] send msgs -> change power of given cyberlinks, affects rank values **@savetheales**

____________

## Rank calculation

- [x] calculation -> every calculation window (if rank changed) **@savetheales**
- [x] full balance transfer -> assign close to default rank values for given CID **@savetheales**

____________

## Rank params change

1. change Tolerance -> more accurate rank
2. change CalculationWindow -> different calculation time
3. change DampingFactor -> change rank values

____________

## Search index

1. udpates -> with zero for new values
- [x] udpates -> calculated values after rank calculation **@savetheales**
- [x] top-1000 -> updates and pagination works **@savetheales**
4. node restart -> should available

____________

## Bandwidth

- [x] exceeded max block band -> returns error **@savetheales**
- [x] not enough personal band -> returns error **@savetheales**
3. personal band -> resumes based on DesirableBandwidth, personal balance and Recovery Period

____________

## Bandwidth params

1. change DesirableBandwidth ->
2. change BaseCreditPrice ->
3. change TxCost ->
4. change LinkMsgCost ->
5. change NonLinkMsgCost ->
6. change AdjustPricePeriod ->
7. change RecoveryPeriod ->

____________

## Price multiplier

1. txs -> change RC
2. time passed -> decrease

____________

## CLI

1. modules -> returns their params
2. linking and all other txs -> works
3. all queries -> available (also with --trust-node=true)
4. ethereum key import -> works
5. ethereum key sign -> works

____________

## Node

1. restart -> works
2. search index -> works
3. rank -> computes
4. cyberlinks -> restored
5. sync from zero
6. sync from backups
7. backups -> works

____________

## Export

1. generates -> graph file and genesis

____________

## Import

1. chain start from exported state -> graph and state applied
2. chain start from exported state -> search index works

____________

## RPC

1. all endpoints -> works

____________

## LCD

1. all endpoints -> works
2. swagger -> works

____________

## Websockets

1. all endpoints -> works

____________

## Consensus

1. fail nodes -> chain operates
2. change tendermint timings -> changed block time

____________

## Crawler

1. wiki indexation -> index to chain works

____________

## Staking

1. delegation -> works
2. redelegation -> works
3. undelegation -> works
4. rewards -> works

____________

## Validators

1. downtime -> jail, slash tokens
2. double sign -> jail, slash tokens
3. unjail -> works
4. change comission -> unbonding

____________

## Community pool

1. Community pool spend -> works

____________

## Online upgrades

1. proposal passed -> node updates to new version

____________

## Cosmwasm

1. contract deploy -> works
2. contract init -> works
3. contract calls -> works
4. contract query -> works
