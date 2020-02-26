# Testing script

## Linking

1. create link -> link created
2. create dublicated link -> returns error

____________

## Transaction

1. successed -> consumed bandwidth
2. failed -> consumes band
3. non-link msgs -> consumes needed band
4. link msgs -> consumes needed band and increasea personal karma
5. link msgs -> increase global karma
6. send msgs -> change balance
7. send msgs -> change max bandwidth, band resume speed
8. send msgs -> change power of given cyberlinks, affects rank values

____________

## Rank calculation

1. calculation -> every calculation window
2. full balance transfer -> assign default rank values for given CID

____________

## Rank params change

1. change Tolerance -> more accurate rank
2. change CalculationWindow -> different calculation time
3. change DampingFactor -> change rank values

____________

## Search index

1. udpates -> with zero for new values
2. udpates -> calculated values after rank calculation
3. top-1000 -> updates and pagination works
4. node restart -> should available

____________

## Bandwidth

1. exceeded max block band -> returns error
2. not enough personal band -> returns error
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
