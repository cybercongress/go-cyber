<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [cyber/bandwidth/v1beta1/types.proto](#cyber/bandwidth/v1beta1/types.proto)
    - [NeuronBandwidth](#cyber.bandwidth.v1beta1.NeuronBandwidth)
    - [Params](#cyber.bandwidth.v1beta1.Params)
    - [Price](#cyber.bandwidth.v1beta1.Price)
  
- [cyber/bandwidth/v1beta1/genesis.proto](#cyber/bandwidth/v1beta1/genesis.proto)
    - [GenesisState](#cyber.bandwidth.v1beta1.GenesisState)
  
- [cyber/bandwidth/v1beta1/query.proto](#cyber/bandwidth/v1beta1/query.proto)
    - [QueryLoadRequest](#cyber.bandwidth.v1beta1.QueryLoadRequest)
    - [QueryLoadResponse](#cyber.bandwidth.v1beta1.QueryLoadResponse)
    - [QueryNeuronBandwidthRequest](#cyber.bandwidth.v1beta1.QueryNeuronBandwidthRequest)
    - [QueryNeuronBandwidthResponse](#cyber.bandwidth.v1beta1.QueryNeuronBandwidthResponse)
    - [QueryParamsRequest](#cyber.bandwidth.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#cyber.bandwidth.v1beta1.QueryParamsResponse)
    - [QueryPriceRequest](#cyber.bandwidth.v1beta1.QueryPriceRequest)
    - [QueryPriceResponse](#cyber.bandwidth.v1beta1.QueryPriceResponse)
    - [QueryTotalBandwidthRequest](#cyber.bandwidth.v1beta1.QueryTotalBandwidthRequest)
    - [QueryTotalBandwidthResponse](#cyber.bandwidth.v1beta1.QueryTotalBandwidthResponse)
  
    - [Query](#cyber.bandwidth.v1beta1.Query)
  
- [cyber/bandwidth/v1beta1/tx.proto](#cyber/bandwidth/v1beta1/tx.proto)
    - [MsgUpdateParams](#cyber.bandwidth.v1beta1.MsgUpdateParams)
    - [MsgUpdateParamsResponse](#cyber.bandwidth.v1beta1.MsgUpdateParamsResponse)
  
    - [Msg](#cyber.bandwidth.v1beta1.Msg)
  
- [cyber/clock/v1/clock.proto](#cyber/clock/v1/clock.proto)
    - [ClockContract](#cyber.clock.v1.ClockContract)
  
- [cyber/clock/v1/genesis.proto](#cyber/clock/v1/genesis.proto)
    - [GenesisState](#cyber.clock.v1.GenesisState)
    - [Params](#cyber.clock.v1.Params)
  
- [cyber/clock/v1/query.proto](#cyber/clock/v1/query.proto)
    - [QueryClockContract](#cyber.clock.v1.QueryClockContract)
    - [QueryClockContractResponse](#cyber.clock.v1.QueryClockContractResponse)
    - [QueryClockContracts](#cyber.clock.v1.QueryClockContracts)
    - [QueryClockContractsResponse](#cyber.clock.v1.QueryClockContractsResponse)
    - [QueryParamsRequest](#cyber.clock.v1.QueryParamsRequest)
    - [QueryParamsResponse](#cyber.clock.v1.QueryParamsResponse)
  
    - [Query](#cyber.clock.v1.Query)
  
- [cyber/clock/v1/tx.proto](#cyber/clock/v1/tx.proto)
    - [MsgRegisterClockContract](#cyber.clock.v1.MsgRegisterClockContract)
    - [MsgRegisterClockContractResponse](#cyber.clock.v1.MsgRegisterClockContractResponse)
    - [MsgUnjailClockContract](#cyber.clock.v1.MsgUnjailClockContract)
    - [MsgUnjailClockContractResponse](#cyber.clock.v1.MsgUnjailClockContractResponse)
    - [MsgUnregisterClockContract](#cyber.clock.v1.MsgUnregisterClockContract)
    - [MsgUnregisterClockContractResponse](#cyber.clock.v1.MsgUnregisterClockContractResponse)
    - [MsgUpdateParams](#cyber.clock.v1.MsgUpdateParams)
    - [MsgUpdateParamsResponse](#cyber.clock.v1.MsgUpdateParamsResponse)
  
    - [Msg](#cyber.clock.v1.Msg)
  
- [cyber/dmn/v1beta1/types.proto](#cyber/dmn/v1beta1/types.proto)
    - [Load](#cyber.dmn.v1beta1.Load)
    - [Params](#cyber.dmn.v1beta1.Params)
    - [Thought](#cyber.dmn.v1beta1.Thought)
    - [ThoughtStats](#cyber.dmn.v1beta1.ThoughtStats)
    - [Trigger](#cyber.dmn.v1beta1.Trigger)
  
- [cyber/dmn/v1beta1/genesis.proto](#cyber/dmn/v1beta1/genesis.proto)
    - [GenesisState](#cyber.dmn.v1beta1.GenesisState)
  
- [cyber/dmn/v1beta1/query.proto](#cyber/dmn/v1beta1/query.proto)
    - [QueryParamsRequest](#cyber.dmn.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#cyber.dmn.v1beta1.QueryParamsResponse)
    - [QueryThoughtParamsRequest](#cyber.dmn.v1beta1.QueryThoughtParamsRequest)
    - [QueryThoughtResponse](#cyber.dmn.v1beta1.QueryThoughtResponse)
    - [QueryThoughtStatsResponse](#cyber.dmn.v1beta1.QueryThoughtStatsResponse)
    - [QueryThoughtsFeesRequest](#cyber.dmn.v1beta1.QueryThoughtsFeesRequest)
    - [QueryThoughtsFeesResponse](#cyber.dmn.v1beta1.QueryThoughtsFeesResponse)
    - [QueryThoughtsRequest](#cyber.dmn.v1beta1.QueryThoughtsRequest)
    - [QueryThoughtsResponse](#cyber.dmn.v1beta1.QueryThoughtsResponse)
    - [QueryThoughtsStatsRequest](#cyber.dmn.v1beta1.QueryThoughtsStatsRequest)
    - [QueryThoughtsStatsResponse](#cyber.dmn.v1beta1.QueryThoughtsStatsResponse)
  
    - [Query](#cyber.dmn.v1beta1.Query)
  
- [cyber/dmn/v1beta1/tx.proto](#cyber/dmn/v1beta1/tx.proto)
    - [MsgChangeThoughtBlock](#cyber.dmn.v1beta1.MsgChangeThoughtBlock)
    - [MsgChangeThoughtBlockResponse](#cyber.dmn.v1beta1.MsgChangeThoughtBlockResponse)
    - [MsgChangeThoughtGasPrice](#cyber.dmn.v1beta1.MsgChangeThoughtGasPrice)
    - [MsgChangeThoughtGasPriceResponse](#cyber.dmn.v1beta1.MsgChangeThoughtGasPriceResponse)
    - [MsgChangeThoughtInput](#cyber.dmn.v1beta1.MsgChangeThoughtInput)
    - [MsgChangeThoughtInputResponse](#cyber.dmn.v1beta1.MsgChangeThoughtInputResponse)
    - [MsgChangeThoughtName](#cyber.dmn.v1beta1.MsgChangeThoughtName)
    - [MsgChangeThoughtNameResponse](#cyber.dmn.v1beta1.MsgChangeThoughtNameResponse)
    - [MsgChangeThoughtParticle](#cyber.dmn.v1beta1.MsgChangeThoughtParticle)
    - [MsgChangeThoughtParticleResponse](#cyber.dmn.v1beta1.MsgChangeThoughtParticleResponse)
    - [MsgChangeThoughtPeriod](#cyber.dmn.v1beta1.MsgChangeThoughtPeriod)
    - [MsgChangeThoughtPeriodResponse](#cyber.dmn.v1beta1.MsgChangeThoughtPeriodResponse)
    - [MsgCreateThought](#cyber.dmn.v1beta1.MsgCreateThought)
    - [MsgCreateThoughtResponse](#cyber.dmn.v1beta1.MsgCreateThoughtResponse)
    - [MsgForgetThought](#cyber.dmn.v1beta1.MsgForgetThought)
    - [MsgForgetThoughtResponse](#cyber.dmn.v1beta1.MsgForgetThoughtResponse)
    - [MsgUpdateParams](#cyber.dmn.v1beta1.MsgUpdateParams)
    - [MsgUpdateParamsResponse](#cyber.dmn.v1beta1.MsgUpdateParamsResponse)
  
    - [Msg](#cyber.dmn.v1beta1.Msg)
  
- [cyber/graph/v1beta1/query.proto](#cyber/graph/v1beta1/query.proto)
    - [QueryGraphStatsRequest](#cyber.graph.v1beta1.QueryGraphStatsRequest)
    - [QueryGraphStatsResponse](#cyber.graph.v1beta1.QueryGraphStatsResponse)
  
    - [Query](#cyber.graph.v1beta1.Query)
  
- [cyber/graph/v1beta1/types.proto](#cyber/graph/v1beta1/types.proto)
    - [Link](#cyber.graph.v1beta1.Link)
  
- [cyber/graph/v1beta1/tx.proto](#cyber/graph/v1beta1/tx.proto)
    - [MsgCyberlink](#cyber.graph.v1beta1.MsgCyberlink)
    - [MsgCyberlinkResponse](#cyber.graph.v1beta1.MsgCyberlinkResponse)
  
    - [Msg](#cyber.graph.v1beta1.Msg)
  
- [cyber/grid/v1beta1/types.proto](#cyber/grid/v1beta1/types.proto)
    - [Params](#cyber.grid.v1beta1.Params)
    - [Route](#cyber.grid.v1beta1.Route)
    - [Value](#cyber.grid.v1beta1.Value)
  
- [cyber/grid/v1beta1/genesis.proto](#cyber/grid/v1beta1/genesis.proto)
    - [GenesisState](#cyber.grid.v1beta1.GenesisState)
  
- [cyber/grid/v1beta1/query.proto](#cyber/grid/v1beta1/query.proto)
    - [QueryDestinationRequest](#cyber.grid.v1beta1.QueryDestinationRequest)
    - [QueryParamsRequest](#cyber.grid.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#cyber.grid.v1beta1.QueryParamsResponse)
    - [QueryRouteRequest](#cyber.grid.v1beta1.QueryRouteRequest)
    - [QueryRouteResponse](#cyber.grid.v1beta1.QueryRouteResponse)
    - [QueryRoutedEnergyResponse](#cyber.grid.v1beta1.QueryRoutedEnergyResponse)
    - [QueryRoutesRequest](#cyber.grid.v1beta1.QueryRoutesRequest)
    - [QueryRoutesResponse](#cyber.grid.v1beta1.QueryRoutesResponse)
    - [QuerySourceRequest](#cyber.grid.v1beta1.QuerySourceRequest)
  
    - [Query](#cyber.grid.v1beta1.Query)
  
- [cyber/grid/v1beta1/tx.proto](#cyber/grid/v1beta1/tx.proto)
    - [MsgCreateRoute](#cyber.grid.v1beta1.MsgCreateRoute)
    - [MsgCreateRouteResponse](#cyber.grid.v1beta1.MsgCreateRouteResponse)
    - [MsgDeleteRoute](#cyber.grid.v1beta1.MsgDeleteRoute)
    - [MsgDeleteRouteResponse](#cyber.grid.v1beta1.MsgDeleteRouteResponse)
    - [MsgEditRoute](#cyber.grid.v1beta1.MsgEditRoute)
    - [MsgEditRouteName](#cyber.grid.v1beta1.MsgEditRouteName)
    - [MsgEditRouteNameResponse](#cyber.grid.v1beta1.MsgEditRouteNameResponse)
    - [MsgEditRouteResponse](#cyber.grid.v1beta1.MsgEditRouteResponse)
    - [MsgUpdateParams](#cyber.grid.v1beta1.MsgUpdateParams)
    - [MsgUpdateParamsResponse](#cyber.grid.v1beta1.MsgUpdateParamsResponse)
  
    - [Msg](#cyber.grid.v1beta1.Msg)
  
- [cyber/liquidity/v1beta1/tx.proto](#cyber/liquidity/v1beta1/tx.proto)
    - [MsgCreatePool](#cyber.liquidity.v1beta1.MsgCreatePool)
    - [MsgCreatePoolResponse](#cyber.liquidity.v1beta1.MsgCreatePoolResponse)
    - [MsgDepositWithinBatch](#cyber.liquidity.v1beta1.MsgDepositWithinBatch)
    - [MsgDepositWithinBatchResponse](#cyber.liquidity.v1beta1.MsgDepositWithinBatchResponse)
    - [MsgSwapWithinBatch](#cyber.liquidity.v1beta1.MsgSwapWithinBatch)
    - [MsgSwapWithinBatchResponse](#cyber.liquidity.v1beta1.MsgSwapWithinBatchResponse)
    - [MsgWithdrawWithinBatch](#cyber.liquidity.v1beta1.MsgWithdrawWithinBatch)
    - [MsgWithdrawWithinBatchResponse](#cyber.liquidity.v1beta1.MsgWithdrawWithinBatchResponse)
  
    - [Msg](#cyber.liquidity.v1beta1.Msg)
  
- [cyber/liquidity/v1beta1/liquidity.proto](#cyber/liquidity/v1beta1/liquidity.proto)
    - [DepositMsgState](#cyber.liquidity.v1beta1.DepositMsgState)
    - [Params](#cyber.liquidity.v1beta1.Params)
    - [Pool](#cyber.liquidity.v1beta1.Pool)
    - [PoolBatch](#cyber.liquidity.v1beta1.PoolBatch)
    - [PoolMetadata](#cyber.liquidity.v1beta1.PoolMetadata)
    - [PoolType](#cyber.liquidity.v1beta1.PoolType)
    - [SwapMsgState](#cyber.liquidity.v1beta1.SwapMsgState)
    - [WithdrawMsgState](#cyber.liquidity.v1beta1.WithdrawMsgState)
  
- [cyber/liquidity/v1beta1/genesis.proto](#cyber/liquidity/v1beta1/genesis.proto)
    - [GenesisState](#cyber.liquidity.v1beta1.GenesisState)
    - [PoolRecord](#cyber.liquidity.v1beta1.PoolRecord)
  
- [cyber/liquidity/v1beta1/query.proto](#cyber/liquidity/v1beta1/query.proto)
    - [QueryLiquidityPoolBatchRequest](#cyber.liquidity.v1beta1.QueryLiquidityPoolBatchRequest)
    - [QueryLiquidityPoolBatchResponse](#cyber.liquidity.v1beta1.QueryLiquidityPoolBatchResponse)
    - [QueryLiquidityPoolByPoolCoinDenomRequest](#cyber.liquidity.v1beta1.QueryLiquidityPoolByPoolCoinDenomRequest)
    - [QueryLiquidityPoolByReserveAccRequest](#cyber.liquidity.v1beta1.QueryLiquidityPoolByReserveAccRequest)
    - [QueryLiquidityPoolRequest](#cyber.liquidity.v1beta1.QueryLiquidityPoolRequest)
    - [QueryLiquidityPoolResponse](#cyber.liquidity.v1beta1.QueryLiquidityPoolResponse)
    - [QueryLiquidityPoolsRequest](#cyber.liquidity.v1beta1.QueryLiquidityPoolsRequest)
    - [QueryLiquidityPoolsResponse](#cyber.liquidity.v1beta1.QueryLiquidityPoolsResponse)
    - [QueryParamsRequest](#cyber.liquidity.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#cyber.liquidity.v1beta1.QueryParamsResponse)
    - [QueryPoolBatchDepositMsgRequest](#cyber.liquidity.v1beta1.QueryPoolBatchDepositMsgRequest)
    - [QueryPoolBatchDepositMsgResponse](#cyber.liquidity.v1beta1.QueryPoolBatchDepositMsgResponse)
    - [QueryPoolBatchDepositMsgsRequest](#cyber.liquidity.v1beta1.QueryPoolBatchDepositMsgsRequest)
    - [QueryPoolBatchDepositMsgsResponse](#cyber.liquidity.v1beta1.QueryPoolBatchDepositMsgsResponse)
    - [QueryPoolBatchSwapMsgRequest](#cyber.liquidity.v1beta1.QueryPoolBatchSwapMsgRequest)
    - [QueryPoolBatchSwapMsgResponse](#cyber.liquidity.v1beta1.QueryPoolBatchSwapMsgResponse)
    - [QueryPoolBatchSwapMsgsRequest](#cyber.liquidity.v1beta1.QueryPoolBatchSwapMsgsRequest)
    - [QueryPoolBatchSwapMsgsResponse](#cyber.liquidity.v1beta1.QueryPoolBatchSwapMsgsResponse)
    - [QueryPoolBatchWithdrawMsgRequest](#cyber.liquidity.v1beta1.QueryPoolBatchWithdrawMsgRequest)
    - [QueryPoolBatchWithdrawMsgResponse](#cyber.liquidity.v1beta1.QueryPoolBatchWithdrawMsgResponse)
    - [QueryPoolBatchWithdrawMsgsRequest](#cyber.liquidity.v1beta1.QueryPoolBatchWithdrawMsgsRequest)
    - [QueryPoolBatchWithdrawMsgsResponse](#cyber.liquidity.v1beta1.QueryPoolBatchWithdrawMsgsResponse)
  
    - [Query](#cyber.liquidity.v1beta1.Query)
  
- [cyber/rank/v1beta1/types.proto](#cyber/rank/v1beta1/types.proto)
    - [Params](#cyber.rank.v1beta1.Params)
    - [RankedParticle](#cyber.rank.v1beta1.RankedParticle)
  
- [cyber/rank/v1beta1/genesis.proto](#cyber/rank/v1beta1/genesis.proto)
    - [GenesisState](#cyber.rank.v1beta1.GenesisState)
  
- [cyber/rank/v1beta1/pagination.proto](#cyber/rank/v1beta1/pagination.proto)
    - [PageRequest](#cyber.rank.v1beta1.PageRequest)
    - [PageResponse](#cyber.rank.v1beta1.PageResponse)
  
- [cyber/rank/v1beta1/query.proto](#cyber/rank/v1beta1/query.proto)
    - [QueryIsAnyLinkExistRequest](#cyber.rank.v1beta1.QueryIsAnyLinkExistRequest)
    - [QueryIsLinkExistRequest](#cyber.rank.v1beta1.QueryIsLinkExistRequest)
    - [QueryKarmaRequest](#cyber.rank.v1beta1.QueryKarmaRequest)
    - [QueryKarmaResponse](#cyber.rank.v1beta1.QueryKarmaResponse)
    - [QueryLinkExistResponse](#cyber.rank.v1beta1.QueryLinkExistResponse)
    - [QueryNegentropyParticleResponse](#cyber.rank.v1beta1.QueryNegentropyParticleResponse)
    - [QueryNegentropyPartilceRequest](#cyber.rank.v1beta1.QueryNegentropyPartilceRequest)
    - [QueryNegentropyRequest](#cyber.rank.v1beta1.QueryNegentropyRequest)
    - [QueryNegentropyResponse](#cyber.rank.v1beta1.QueryNegentropyResponse)
    - [QueryParamsRequest](#cyber.rank.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#cyber.rank.v1beta1.QueryParamsResponse)
    - [QueryRankRequest](#cyber.rank.v1beta1.QueryRankRequest)
    - [QueryRankResponse](#cyber.rank.v1beta1.QueryRankResponse)
    - [QuerySearchRequest](#cyber.rank.v1beta1.QuerySearchRequest)
    - [QuerySearchResponse](#cyber.rank.v1beta1.QuerySearchResponse)
    - [QueryTopRequest](#cyber.rank.v1beta1.QueryTopRequest)
  
    - [Query](#cyber.rank.v1beta1.Query)
  
- [cyber/rank/v1beta1/tx.proto](#cyber/rank/v1beta1/tx.proto)
    - [MsgUpdateParams](#cyber.rank.v1beta1.MsgUpdateParams)
    - [MsgUpdateParamsResponse](#cyber.rank.v1beta1.MsgUpdateParamsResponse)
  
    - [Msg](#cyber.rank.v1beta1.Msg)
  
- [cyber/resources/v1beta1/types.proto](#cyber/resources/v1beta1/types.proto)
    - [Params](#cyber.resources.v1beta1.Params)
  
- [cyber/resources/v1beta1/genesis.proto](#cyber/resources/v1beta1/genesis.proto)
    - [GenesisState](#cyber.resources.v1beta1.GenesisState)
  
- [cyber/resources/v1beta1/query.proto](#cyber/resources/v1beta1/query.proto)
    - [QueryInvestmintRequest](#cyber.resources.v1beta1.QueryInvestmintRequest)
    - [QueryInvestmintResponse](#cyber.resources.v1beta1.QueryInvestmintResponse)
    - [QueryParamsRequest](#cyber.resources.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#cyber.resources.v1beta1.QueryParamsResponse)
  
    - [Query](#cyber.resources.v1beta1.Query)
  
- [cyber/resources/v1beta1/tx.proto](#cyber/resources/v1beta1/tx.proto)
    - [MsgInvestmint](#cyber.resources.v1beta1.MsgInvestmint)
    - [MsgInvestmintResponse](#cyber.resources.v1beta1.MsgInvestmintResponse)
    - [MsgUpdateParams](#cyber.resources.v1beta1.MsgUpdateParams)
    - [MsgUpdateParamsResponse](#cyber.resources.v1beta1.MsgUpdateParamsResponse)
  
    - [Msg](#cyber.resources.v1beta1.Msg)
  
- [osmosis/tokenfactory/v1beta1/authority_metadata.proto](#osmosis/tokenfactory/v1beta1/authority_metadata.proto)
    - [DenomAuthorityMetadata](#osmosis.tokenfactory.v1beta1.DenomAuthorityMetadata)
  
- [osmosis/tokenfactory/v1beta1/params.proto](#osmosis/tokenfactory/v1beta1/params.proto)
    - [Params](#osmosis.tokenfactory.v1beta1.Params)
  
- [osmosis/tokenfactory/v1beta1/genesis.proto](#osmosis/tokenfactory/v1beta1/genesis.proto)
    - [GenesisDenom](#osmosis.tokenfactory.v1beta1.GenesisDenom)
    - [GenesisState](#osmosis.tokenfactory.v1beta1.GenesisState)
  
- [osmosis/tokenfactory/v1beta1/query.proto](#osmosis/tokenfactory/v1beta1/query.proto)
    - [QueryDenomAuthorityMetadataRequest](#osmosis.tokenfactory.v1beta1.QueryDenomAuthorityMetadataRequest)
    - [QueryDenomAuthorityMetadataResponse](#osmosis.tokenfactory.v1beta1.QueryDenomAuthorityMetadataResponse)
    - [QueryDenomsFromCreatorRequest](#osmosis.tokenfactory.v1beta1.QueryDenomsFromCreatorRequest)
    - [QueryDenomsFromCreatorResponse](#osmosis.tokenfactory.v1beta1.QueryDenomsFromCreatorResponse)
    - [QueryParamsRequest](#osmosis.tokenfactory.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#osmosis.tokenfactory.v1beta1.QueryParamsResponse)
  
    - [Query](#osmosis.tokenfactory.v1beta1.Query)
  
- [osmosis/tokenfactory/v1beta1/tx.proto](#osmosis/tokenfactory/v1beta1/tx.proto)
    - [MsgBurn](#osmosis.tokenfactory.v1beta1.MsgBurn)
    - [MsgBurnResponse](#osmosis.tokenfactory.v1beta1.MsgBurnResponse)
    - [MsgChangeAdmin](#osmosis.tokenfactory.v1beta1.MsgChangeAdmin)
    - [MsgChangeAdminResponse](#osmosis.tokenfactory.v1beta1.MsgChangeAdminResponse)
    - [MsgCreateDenom](#osmosis.tokenfactory.v1beta1.MsgCreateDenom)
    - [MsgCreateDenomResponse](#osmosis.tokenfactory.v1beta1.MsgCreateDenomResponse)
    - [MsgForceTransfer](#osmosis.tokenfactory.v1beta1.MsgForceTransfer)
    - [MsgForceTransferResponse](#osmosis.tokenfactory.v1beta1.MsgForceTransferResponse)
    - [MsgMint](#osmosis.tokenfactory.v1beta1.MsgMint)
    - [MsgMintResponse](#osmosis.tokenfactory.v1beta1.MsgMintResponse)
    - [MsgSetDenomMetadata](#osmosis.tokenfactory.v1beta1.MsgSetDenomMetadata)
    - [MsgSetDenomMetadataResponse](#osmosis.tokenfactory.v1beta1.MsgSetDenomMetadataResponse)
    - [MsgUpdateParams](#osmosis.tokenfactory.v1beta1.MsgUpdateParams)
    - [MsgUpdateParamsResponse](#osmosis.tokenfactory.v1beta1.MsgUpdateParamsResponse)
  
    - [Msg](#osmosis.tokenfactory.v1beta1.Msg)
  
- [Scalar Value Types](#scalar-value-types)



<a name="cyber/bandwidth/v1beta1/types.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/bandwidth/v1beta1/types.proto



<a name="cyber.bandwidth.v1beta1.NeuronBandwidth"></a>

### NeuronBandwidth



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `neuron` | [string](#string) |  |  |
| `remained_value` | [uint64](#uint64) |  |  |
| `last_updated_block` | [uint64](#uint64) |  |  |
| `max_value` | [uint64](#uint64) |  |  |






<a name="cyber.bandwidth.v1beta1.Params"></a>

### Params



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `recovery_period` | [uint64](#uint64) |  |  |
| `adjust_price_period` | [uint64](#uint64) |  |  |
| `base_price` | [string](#string) |  |  |
| `base_load` | [string](#string) |  |  |
| `max_block_bandwidth` | [uint64](#uint64) |  |  |






<a name="cyber.bandwidth.v1beta1.Price"></a>

### Price



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `price` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cyber/bandwidth/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/bandwidth/v1beta1/genesis.proto



<a name="cyber.bandwidth.v1beta1.GenesisState"></a>

### GenesisState



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cyber.bandwidth.v1beta1.Params) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cyber/bandwidth/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/bandwidth/v1beta1/query.proto



<a name="cyber.bandwidth.v1beta1.QueryLoadRequest"></a>

### QueryLoadRequest







<a name="cyber.bandwidth.v1beta1.QueryLoadResponse"></a>

### QueryLoadResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `load` | [string](#string) |  |  |






<a name="cyber.bandwidth.v1beta1.QueryNeuronBandwidthRequest"></a>

### QueryNeuronBandwidthRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `neuron` | [string](#string) |  |  |






<a name="cyber.bandwidth.v1beta1.QueryNeuronBandwidthResponse"></a>

### QueryNeuronBandwidthResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `neuron_bandwidth` | [NeuronBandwidth](#cyber.bandwidth.v1beta1.NeuronBandwidth) |  |  |






<a name="cyber.bandwidth.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest







<a name="cyber.bandwidth.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cyber.bandwidth.v1beta1.Params) |  |  |






<a name="cyber.bandwidth.v1beta1.QueryPriceRequest"></a>

### QueryPriceRequest







<a name="cyber.bandwidth.v1beta1.QueryPriceResponse"></a>

### QueryPriceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `price` | [string](#string) |  |  |






<a name="cyber.bandwidth.v1beta1.QueryTotalBandwidthRequest"></a>

### QueryTotalBandwidthRequest







<a name="cyber.bandwidth.v1beta1.QueryTotalBandwidthResponse"></a>

### QueryTotalBandwidthResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `total_bandwidth` | [uint64](#uint64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cyber.bandwidth.v1beta1.Query"></a>

### Query


| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Load` | [QueryLoadRequest](#cyber.bandwidth.v1beta1.QueryLoadRequest) | [QueryLoadResponse](#cyber.bandwidth.v1beta1.QueryLoadResponse) |  | GET|/cyber/bandwidth/v1beta1/bandwidth/load|
| `Price` | [QueryPriceRequest](#cyber.bandwidth.v1beta1.QueryPriceRequest) | [QueryPriceResponse](#cyber.bandwidth.v1beta1.QueryPriceResponse) |  | GET|/cyber/bandwidth/v1beta1/bandwidth/price|
| `TotalBandwidth` | [QueryTotalBandwidthRequest](#cyber.bandwidth.v1beta1.QueryTotalBandwidthRequest) | [QueryTotalBandwidthResponse](#cyber.bandwidth.v1beta1.QueryTotalBandwidthResponse) |  | GET|/cyber/bandwidth/v1beta1/bandwidth/total|
| `NeuronBandwidth` | [QueryNeuronBandwidthRequest](#cyber.bandwidth.v1beta1.QueryNeuronBandwidthRequest) | [QueryNeuronBandwidthResponse](#cyber.bandwidth.v1beta1.QueryNeuronBandwidthResponse) |  | GET|/cyber/bandwidth/v1beta1/bandwidth/neuron/{neuron}|
| `Params` | [QueryParamsRequest](#cyber.bandwidth.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#cyber.bandwidth.v1beta1.QueryParamsResponse) |  | GET|/cyber/bandwidth/v1beta1/bandwidth/params|

 <!-- end services -->



<a name="cyber/bandwidth/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/bandwidth/v1beta1/tx.proto



<a name="cyber.bandwidth.v1beta1.MsgUpdateParams"></a>

### MsgUpdateParams



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authority` | [string](#string) |  |  |
| `params` | [Params](#cyber.bandwidth.v1beta1.Params) |  |  |






<a name="cyber.bandwidth.v1beta1.MsgUpdateParamsResponse"></a>

### MsgUpdateParamsResponse






 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cyber.bandwidth.v1beta1.Msg"></a>

### Msg


| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `UpdateParams` | [MsgUpdateParams](#cyber.bandwidth.v1beta1.MsgUpdateParams) | [MsgUpdateParamsResponse](#cyber.bandwidth.v1beta1.MsgUpdateParamsResponse) |  | |

 <!-- end services -->



<a name="cyber/clock/v1/clock.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/clock/v1/clock.proto



<a name="cyber.clock.v1.ClockContract"></a>

### ClockContract
This object is used to store the contract address and the
jail status of the contract.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_address` | [string](#string) |  | The address of the contract. |
| `is_jailed` | [bool](#bool) |  | The jail status of the contract. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cyber/clock/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/clock/v1/genesis.proto



<a name="cyber.clock.v1.GenesisState"></a>

### GenesisState
GenesisState - initial state of module


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cyber.clock.v1.Params) |  | Params of this module |






<a name="cyber.clock.v1.Params"></a>

### Params
Params defines the set of module parameters.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_gas_limit` | [uint64](#uint64) |  | contract_gas_limit defines the maximum amount of gas that can be used by a contract. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cyber/clock/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/clock/v1/query.proto



<a name="cyber.clock.v1.QueryClockContract"></a>

### QueryClockContract
QueryClockContract is the request type to get a single contract.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_address` | [string](#string) |  | contract_address is the address of the contract to query. |






<a name="cyber.clock.v1.QueryClockContractResponse"></a>

### QueryClockContractResponse
QueryClockContractResponse is the response type for the Query/ClockContract
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `clock_contract` | [ClockContract](#cyber.clock.v1.ClockContract) |  | contract is the clock contract. |






<a name="cyber.clock.v1.QueryClockContracts"></a>

### QueryClockContracts
QueryClockContracts is the request type to get all contracts.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cyber.clock.v1.QueryClockContractsResponse"></a>

### QueryClockContractsResponse
QueryClockContractsResponse is the response type for the Query/ClockContracts
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `clock_contracts` | [ClockContract](#cyber.clock.v1.ClockContract) | repeated | clock_contracts are the clock contracts. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="cyber.clock.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParams is the request type to get all module params.






<a name="cyber.clock.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryClockContractsResponse is the response type for the Query/ClockContracts
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cyber.clock.v1.Params) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cyber.clock.v1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `ClockContracts` | [QueryClockContracts](#cyber.clock.v1.QueryClockContracts) | [QueryClockContractsResponse](#cyber.clock.v1.QueryClockContractsResponse) | ClockContracts | GET|/cyber/clock/v1/contracts|
| `ClockContract` | [QueryClockContract](#cyber.clock.v1.QueryClockContract) | [QueryClockContractResponse](#cyber.clock.v1.QueryClockContractResponse) | ClockContract | GET|/cyber/clock/v1/contracts/{contract_address}|
| `Params` | [QueryParamsRequest](#cyber.clock.v1.QueryParamsRequest) | [QueryParamsResponse](#cyber.clock.v1.QueryParamsResponse) | Params | GET|/cyber/clock/v1/params|

 <!-- end services -->



<a name="cyber/clock/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/clock/v1/tx.proto



<a name="cyber.clock.v1.MsgRegisterClockContract"></a>

### MsgRegisterClockContract
MsgRegisterClockContract is the Msg/RegisterClockContract request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender_address` | [string](#string) |  | The address of the sender. |
| `contract_address` | [string](#string) |  | The address of the contract to register. |






<a name="cyber.clock.v1.MsgRegisterClockContractResponse"></a>

### MsgRegisterClockContractResponse
MsgRegisterClockContractResponse defines the response structure for executing
a MsgRegisterClockContract message.






<a name="cyber.clock.v1.MsgUnjailClockContract"></a>

### MsgUnjailClockContract
MsgUnjailClockContract is the Msg/UnjailClockContract request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender_address` | [string](#string) |  | The address of the sender. |
| `contract_address` | [string](#string) |  | The address of the contract to unjail. |






<a name="cyber.clock.v1.MsgUnjailClockContractResponse"></a>

### MsgUnjailClockContractResponse
MsgUnjailClockContractResponse defines the response structure for executing a
MsgUnjailClockContract message.






<a name="cyber.clock.v1.MsgUnregisterClockContract"></a>

### MsgUnregisterClockContract
MsgUnregisterClockContract is the Msg/UnregisterClockContract request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender_address` | [string](#string) |  | The address of the sender. |
| `contract_address` | [string](#string) |  | The address of the contract to unregister. |






<a name="cyber.clock.v1.MsgUnregisterClockContractResponse"></a>

### MsgUnregisterClockContractResponse
MsgUnregisterClockContractResponse defines the response structure for
executing a MsgUnregisterClockContract message.






<a name="cyber.clock.v1.MsgUpdateParams"></a>

### MsgUpdateParams
MsgUpdateParams is the Msg/UpdateParams request type.

Since: cosmos-sdk 0.47


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authority` | [string](#string) |  | authority is the address of the governance account. |
| `params` | [Params](#cyber.clock.v1.Params) |  | params defines the x/clock parameters to update.

NOTE: All parameters must be supplied. |






<a name="cyber.clock.v1.MsgUpdateParamsResponse"></a>

### MsgUpdateParamsResponse
MsgUpdateParamsResponse defines the response structure for executing a
MsgUpdateParams message.

Since: cosmos-sdk 0.47





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cyber.clock.v1.Msg"></a>

### Msg
Msg defines the Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `RegisterClockContract` | [MsgRegisterClockContract](#cyber.clock.v1.MsgRegisterClockContract) | [MsgRegisterClockContractResponse](#cyber.clock.v1.MsgRegisterClockContractResponse) | RegisterClockContract defines the endpoint for registering a new clock contract. | POST|/cyber/clock/v1/tx/register|
| `UnregisterClockContract` | [MsgUnregisterClockContract](#cyber.clock.v1.MsgUnregisterClockContract) | [MsgUnregisterClockContractResponse](#cyber.clock.v1.MsgUnregisterClockContractResponse) | UnregisterClockContract defines the endpoint for unregistering a clock contract. | POST|/cyber/clock/v1/tx/unregister|
| `UnjailClockContract` | [MsgUnjailClockContract](#cyber.clock.v1.MsgUnjailClockContract) | [MsgUnjailClockContractResponse](#cyber.clock.v1.MsgUnjailClockContractResponse) | UnjailClockContract defines the endpoint for unjailing a clock contract. | POST|/cyber/clock/v1/tx/unjail|
| `UpdateParams` | [MsgUpdateParams](#cyber.clock.v1.MsgUpdateParams) | [MsgUpdateParamsResponse](#cyber.clock.v1.MsgUpdateParamsResponse) | UpdateParams defines a governance operation for updating the x/clock module parameters. The authority is hard-coded to the x/gov module account.

Since: cosmos-sdk 0.47 | |

 <!-- end services -->



<a name="cyber/dmn/v1beta1/types.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/dmn/v1beta1/types.proto



<a name="cyber.dmn.v1beta1.Load"></a>

### Load



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `input` | [string](#string) |  |  |
| `gas_price` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="cyber.dmn.v1beta1.Params"></a>

### Params



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `max_slots` | [uint32](#uint32) |  |  |
| `max_gas` | [uint32](#uint32) |  |  |
| `fee_ttl` | [uint32](#uint32) |  |  |






<a name="cyber.dmn.v1beta1.Thought"></a>

### Thought



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `program` | [string](#string) |  |  |
| `trigger` | [Trigger](#cyber.dmn.v1beta1.Trigger) |  |  |
| `load` | [Load](#cyber.dmn.v1beta1.Load) |  |  |
| `name` | [string](#string) |  |  |
| `particle` | [string](#string) |  |  |






<a name="cyber.dmn.v1beta1.ThoughtStats"></a>

### ThoughtStats



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `program` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |
| `calls` | [uint64](#uint64) |  |  |
| `fees` | [uint64](#uint64) |  |  |
| `gas` | [uint64](#uint64) |  |  |
| `last_block` | [uint64](#uint64) |  |  |






<a name="cyber.dmn.v1beta1.Trigger"></a>

### Trigger



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `period` | [uint64](#uint64) |  |  |
| `block` | [uint64](#uint64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cyber/dmn/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/dmn/v1beta1/genesis.proto



<a name="cyber.dmn.v1beta1.GenesisState"></a>

### GenesisState



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cyber.dmn.v1beta1.Params) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cyber/dmn/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/dmn/v1beta1/query.proto



<a name="cyber.dmn.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest







<a name="cyber.dmn.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cyber.dmn.v1beta1.Params) |  |  |






<a name="cyber.dmn.v1beta1.QueryThoughtParamsRequest"></a>

### QueryThoughtParamsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `program` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |






<a name="cyber.dmn.v1beta1.QueryThoughtResponse"></a>

### QueryThoughtResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `thought` | [Thought](#cyber.dmn.v1beta1.Thought) |  |  |






<a name="cyber.dmn.v1beta1.QueryThoughtStatsResponse"></a>

### QueryThoughtStatsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `thought_stats` | [ThoughtStats](#cyber.dmn.v1beta1.ThoughtStats) |  |  |






<a name="cyber.dmn.v1beta1.QueryThoughtsFeesRequest"></a>

### QueryThoughtsFeesRequest







<a name="cyber.dmn.v1beta1.QueryThoughtsFeesResponse"></a>

### QueryThoughtsFeesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `fees` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="cyber.dmn.v1beta1.QueryThoughtsRequest"></a>

### QueryThoughtsRequest







<a name="cyber.dmn.v1beta1.QueryThoughtsResponse"></a>

### QueryThoughtsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `thoughts` | [Thought](#cyber.dmn.v1beta1.Thought) | repeated |  |






<a name="cyber.dmn.v1beta1.QueryThoughtsStatsRequest"></a>

### QueryThoughtsStatsRequest







<a name="cyber.dmn.v1beta1.QueryThoughtsStatsResponse"></a>

### QueryThoughtsStatsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `thoughts_stats` | [ThoughtStats](#cyber.dmn.v1beta1.ThoughtStats) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cyber.dmn.v1beta1.Query"></a>

### Query


| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#cyber.dmn.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#cyber.dmn.v1beta1.QueryParamsResponse) |  | GET|/cyber/dmn/v1beta1/dmn/params|
| `Thought` | [QueryThoughtParamsRequest](#cyber.dmn.v1beta1.QueryThoughtParamsRequest) | [QueryThoughtResponse](#cyber.dmn.v1beta1.QueryThoughtResponse) |  | GET|/cyber/dmn/v1beta1/dmn/thought|
| `ThoughtStats` | [QueryThoughtParamsRequest](#cyber.dmn.v1beta1.QueryThoughtParamsRequest) | [QueryThoughtStatsResponse](#cyber.dmn.v1beta1.QueryThoughtStatsResponse) |  | GET|/cyber/dmn/v1beta1/dmn/thought_stats|
| `Thoughts` | [QueryThoughtsRequest](#cyber.dmn.v1beta1.QueryThoughtsRequest) | [QueryThoughtsResponse](#cyber.dmn.v1beta1.QueryThoughtsResponse) |  | GET|/cyber/dmn/v1beta1/dmn/thoughts|
| `ThoughtsStats` | [QueryThoughtsStatsRequest](#cyber.dmn.v1beta1.QueryThoughtsStatsRequest) | [QueryThoughtsStatsResponse](#cyber.dmn.v1beta1.QueryThoughtsStatsResponse) |  | GET|/cyber/dmn/v1beta1/dmn/thoughts_stats|
| `ThoughtsFees` | [QueryThoughtsFeesRequest](#cyber.dmn.v1beta1.QueryThoughtsFeesRequest) | [QueryThoughtsFeesResponse](#cyber.dmn.v1beta1.QueryThoughtsFeesResponse) |  | GET|/cyber/dmn/v1beta1/dmn/thoughts_fees|

 <!-- end services -->



<a name="cyber/dmn/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/dmn/v1beta1/tx.proto



<a name="cyber.dmn.v1beta1.MsgChangeThoughtBlock"></a>

### MsgChangeThoughtBlock



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `program` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |
| `block` | [uint64](#uint64) |  |  |






<a name="cyber.dmn.v1beta1.MsgChangeThoughtBlockResponse"></a>

### MsgChangeThoughtBlockResponse







<a name="cyber.dmn.v1beta1.MsgChangeThoughtGasPrice"></a>

### MsgChangeThoughtGasPrice



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `program` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |
| `gas_price` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="cyber.dmn.v1beta1.MsgChangeThoughtGasPriceResponse"></a>

### MsgChangeThoughtGasPriceResponse







<a name="cyber.dmn.v1beta1.MsgChangeThoughtInput"></a>

### MsgChangeThoughtInput



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `program` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |
| `input` | [string](#string) |  |  |






<a name="cyber.dmn.v1beta1.MsgChangeThoughtInputResponse"></a>

### MsgChangeThoughtInputResponse







<a name="cyber.dmn.v1beta1.MsgChangeThoughtName"></a>

### MsgChangeThoughtName



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `program` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |
| `new_name` | [string](#string) |  |  |






<a name="cyber.dmn.v1beta1.MsgChangeThoughtNameResponse"></a>

### MsgChangeThoughtNameResponse







<a name="cyber.dmn.v1beta1.MsgChangeThoughtParticle"></a>

### MsgChangeThoughtParticle



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `program` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |
| `particle` | [string](#string) |  |  |






<a name="cyber.dmn.v1beta1.MsgChangeThoughtParticleResponse"></a>

### MsgChangeThoughtParticleResponse







<a name="cyber.dmn.v1beta1.MsgChangeThoughtPeriod"></a>

### MsgChangeThoughtPeriod



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `program` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |
| `period` | [uint64](#uint64) |  |  |






<a name="cyber.dmn.v1beta1.MsgChangeThoughtPeriodResponse"></a>

### MsgChangeThoughtPeriodResponse







<a name="cyber.dmn.v1beta1.MsgCreateThought"></a>

### MsgCreateThought



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `program` | [string](#string) |  |  |
| `trigger` | [Trigger](#cyber.dmn.v1beta1.Trigger) |  |  |
| `load` | [Load](#cyber.dmn.v1beta1.Load) |  |  |
| `name` | [string](#string) |  |  |
| `particle` | [string](#string) |  |  |






<a name="cyber.dmn.v1beta1.MsgCreateThoughtResponse"></a>

### MsgCreateThoughtResponse







<a name="cyber.dmn.v1beta1.MsgForgetThought"></a>

### MsgForgetThought



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `program` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |






<a name="cyber.dmn.v1beta1.MsgForgetThoughtResponse"></a>

### MsgForgetThoughtResponse







<a name="cyber.dmn.v1beta1.MsgUpdateParams"></a>

### MsgUpdateParams



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authority` | [string](#string) |  |  |
| `params` | [Params](#cyber.dmn.v1beta1.Params) |  |  |






<a name="cyber.dmn.v1beta1.MsgUpdateParamsResponse"></a>

### MsgUpdateParamsResponse






 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cyber.dmn.v1beta1.Msg"></a>

### Msg


| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `CreateThought` | [MsgCreateThought](#cyber.dmn.v1beta1.MsgCreateThought) | [MsgCreateThoughtResponse](#cyber.dmn.v1beta1.MsgCreateThoughtResponse) |  | |
| `ForgetThought` | [MsgForgetThought](#cyber.dmn.v1beta1.MsgForgetThought) | [MsgForgetThoughtResponse](#cyber.dmn.v1beta1.MsgForgetThoughtResponse) |  | |
| `ChangeThoughtParticle` | [MsgChangeThoughtParticle](#cyber.dmn.v1beta1.MsgChangeThoughtParticle) | [MsgChangeThoughtParticleResponse](#cyber.dmn.v1beta1.MsgChangeThoughtParticleResponse) |  | |
| `ChangeThoughtName` | [MsgChangeThoughtName](#cyber.dmn.v1beta1.MsgChangeThoughtName) | [MsgChangeThoughtNameResponse](#cyber.dmn.v1beta1.MsgChangeThoughtNameResponse) |  | |
| `ChangeThoughtInput` | [MsgChangeThoughtInput](#cyber.dmn.v1beta1.MsgChangeThoughtInput) | [MsgChangeThoughtInputResponse](#cyber.dmn.v1beta1.MsgChangeThoughtInputResponse) |  | |
| `ChangeThoughtGasPrice` | [MsgChangeThoughtGasPrice](#cyber.dmn.v1beta1.MsgChangeThoughtGasPrice) | [MsgChangeThoughtGasPriceResponse](#cyber.dmn.v1beta1.MsgChangeThoughtGasPriceResponse) |  | |
| `ChangeThoughtPeriod` | [MsgChangeThoughtPeriod](#cyber.dmn.v1beta1.MsgChangeThoughtPeriod) | [MsgChangeThoughtPeriodResponse](#cyber.dmn.v1beta1.MsgChangeThoughtPeriodResponse) |  | |
| `ChangeThoughtBlock` | [MsgChangeThoughtBlock](#cyber.dmn.v1beta1.MsgChangeThoughtBlock) | [MsgChangeThoughtBlockResponse](#cyber.dmn.v1beta1.MsgChangeThoughtBlockResponse) |  | |
| `UpdateParams` | [MsgUpdateParams](#cyber.dmn.v1beta1.MsgUpdateParams) | [MsgUpdateParamsResponse](#cyber.dmn.v1beta1.MsgUpdateParamsResponse) |  | |

 <!-- end services -->



<a name="cyber/graph/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/graph/v1beta1/query.proto



<a name="cyber.graph.v1beta1.QueryGraphStatsRequest"></a>

### QueryGraphStatsRequest







<a name="cyber.graph.v1beta1.QueryGraphStatsResponse"></a>

### QueryGraphStatsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cyberlinks` | [uint64](#uint64) |  |  |
| `particles` | [uint64](#uint64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cyber.graph.v1beta1.Query"></a>

### Query


| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `GraphStats` | [QueryGraphStatsRequest](#cyber.graph.v1beta1.QueryGraphStatsRequest) | [QueryGraphStatsResponse](#cyber.graph.v1beta1.QueryGraphStatsResponse) |  | GET|/cyber/graph/v1beta1/graph_stats|

 <!-- end services -->



<a name="cyber/graph/v1beta1/types.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/graph/v1beta1/types.proto



<a name="cyber.graph.v1beta1.Link"></a>

### Link



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `from` | [string](#string) |  |  |
| `to` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cyber/graph/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/graph/v1beta1/tx.proto



<a name="cyber.graph.v1beta1.MsgCyberlink"></a>

### MsgCyberlink



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `neuron` | [string](#string) |  |  |
| `links` | [Link](#cyber.graph.v1beta1.Link) | repeated |  |






<a name="cyber.graph.v1beta1.MsgCyberlinkResponse"></a>

### MsgCyberlinkResponse






 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cyber.graph.v1beta1.Msg"></a>

### Msg


| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Cyberlink` | [MsgCyberlink](#cyber.graph.v1beta1.MsgCyberlink) | [MsgCyberlinkResponse](#cyber.graph.v1beta1.MsgCyberlinkResponse) |  | |

 <!-- end services -->



<a name="cyber/grid/v1beta1/types.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/grid/v1beta1/types.proto



<a name="cyber.grid.v1beta1.Params"></a>

### Params



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `max_routes` | [uint32](#uint32) |  |  |






<a name="cyber.grid.v1beta1.Route"></a>

### Route



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `source` | [string](#string) |  |  |
| `destination` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |
| `value` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="cyber.grid.v1beta1.Value"></a>

### Value



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `value` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cyber/grid/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/grid/v1beta1/genesis.proto



<a name="cyber.grid.v1beta1.GenesisState"></a>

### GenesisState



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cyber.grid.v1beta1.Params) |  |  |
| `routes` | [Route](#cyber.grid.v1beta1.Route) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cyber/grid/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/grid/v1beta1/query.proto



<a name="cyber.grid.v1beta1.QueryDestinationRequest"></a>

### QueryDestinationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `destination` | [string](#string) |  |  |






<a name="cyber.grid.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest







<a name="cyber.grid.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cyber.grid.v1beta1.Params) |  |  |






<a name="cyber.grid.v1beta1.QueryRouteRequest"></a>

### QueryRouteRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `source` | [string](#string) |  |  |
| `destination` | [string](#string) |  |  |






<a name="cyber.grid.v1beta1.QueryRouteResponse"></a>

### QueryRouteResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `route` | [Route](#cyber.grid.v1beta1.Route) |  |  |






<a name="cyber.grid.v1beta1.QueryRoutedEnergyResponse"></a>

### QueryRoutedEnergyResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `value` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="cyber.grid.v1beta1.QueryRoutesRequest"></a>

### QueryRoutesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="cyber.grid.v1beta1.QueryRoutesResponse"></a>

### QueryRoutesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `routes` | [Route](#cyber.grid.v1beta1.Route) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |






<a name="cyber.grid.v1beta1.QuerySourceRequest"></a>

### QuerySourceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `source` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cyber.grid.v1beta1.Query"></a>

### Query


| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#cyber.grid.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#cyber.grid.v1beta1.QueryParamsResponse) |  | GET|/cyber/grid/v1beta1/grid/params|
| `SourceRoutes` | [QuerySourceRequest](#cyber.grid.v1beta1.QuerySourceRequest) | [QueryRoutesResponse](#cyber.grid.v1beta1.QueryRoutesResponse) |  | GET|/cyber/grid/v1beta1/grid/source_routes|
| `DestinationRoutes` | [QueryDestinationRequest](#cyber.grid.v1beta1.QueryDestinationRequest) | [QueryRoutesResponse](#cyber.grid.v1beta1.QueryRoutesResponse) |  | GET|/cyber/grid/v1beta1/grid/destination_routes|
| `DestinationRoutedEnergy` | [QueryDestinationRequest](#cyber.grid.v1beta1.QueryDestinationRequest) | [QueryRoutedEnergyResponse](#cyber.grid.v1beta1.QueryRoutedEnergyResponse) |  | GET|/cyber/grid/v1beta1/grid/destination_routed_energy|
| `SourceRoutedEnergy` | [QuerySourceRequest](#cyber.grid.v1beta1.QuerySourceRequest) | [QueryRoutedEnergyResponse](#cyber.grid.v1beta1.QueryRoutedEnergyResponse) |  | GET|/cyber/grid/v1beta1/grid/source_routed_energy|
| `Route` | [QueryRouteRequest](#cyber.grid.v1beta1.QueryRouteRequest) | [QueryRouteResponse](#cyber.grid.v1beta1.QueryRouteResponse) |  | GET|/cyber/grid/v1beta1/grid/route|
| `Routes` | [QueryRoutesRequest](#cyber.grid.v1beta1.QueryRoutesRequest) | [QueryRoutesResponse](#cyber.grid.v1beta1.QueryRoutesResponse) |  | GET|/cyber/grid/v1beta1/grid/routes|

 <!-- end services -->



<a name="cyber/grid/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/grid/v1beta1/tx.proto



<a name="cyber.grid.v1beta1.MsgCreateRoute"></a>

### MsgCreateRoute



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `source` | [string](#string) |  |  |
| `destination` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |






<a name="cyber.grid.v1beta1.MsgCreateRouteResponse"></a>

### MsgCreateRouteResponse







<a name="cyber.grid.v1beta1.MsgDeleteRoute"></a>

### MsgDeleteRoute



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `source` | [string](#string) |  |  |
| `destination` | [string](#string) |  |  |






<a name="cyber.grid.v1beta1.MsgDeleteRouteResponse"></a>

### MsgDeleteRouteResponse







<a name="cyber.grid.v1beta1.MsgEditRoute"></a>

### MsgEditRoute



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `source` | [string](#string) |  |  |
| `destination` | [string](#string) |  |  |
| `value` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="cyber.grid.v1beta1.MsgEditRouteName"></a>

### MsgEditRouteName



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `source` | [string](#string) |  |  |
| `destination` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |






<a name="cyber.grid.v1beta1.MsgEditRouteNameResponse"></a>

### MsgEditRouteNameResponse







<a name="cyber.grid.v1beta1.MsgEditRouteResponse"></a>

### MsgEditRouteResponse







<a name="cyber.grid.v1beta1.MsgUpdateParams"></a>

### MsgUpdateParams



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authority` | [string](#string) |  |  |
| `params` | [Params](#cyber.grid.v1beta1.Params) |  |  |






<a name="cyber.grid.v1beta1.MsgUpdateParamsResponse"></a>

### MsgUpdateParamsResponse






 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cyber.grid.v1beta1.Msg"></a>

### Msg


| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `CreateRoute` | [MsgCreateRoute](#cyber.grid.v1beta1.MsgCreateRoute) | [MsgCreateRouteResponse](#cyber.grid.v1beta1.MsgCreateRouteResponse) |  | |
| `EditRoute` | [MsgEditRoute](#cyber.grid.v1beta1.MsgEditRoute) | [MsgEditRouteResponse](#cyber.grid.v1beta1.MsgEditRouteResponse) |  | |
| `DeleteRoute` | [MsgDeleteRoute](#cyber.grid.v1beta1.MsgDeleteRoute) | [MsgDeleteRouteResponse](#cyber.grid.v1beta1.MsgDeleteRouteResponse) |  | |
| `EditRouteName` | [MsgEditRouteName](#cyber.grid.v1beta1.MsgEditRouteName) | [MsgEditRouteNameResponse](#cyber.grid.v1beta1.MsgEditRouteNameResponse) |  | |
| `UpdateParams` | [MsgUpdateParams](#cyber.grid.v1beta1.MsgUpdateParams) | [MsgUpdateParamsResponse](#cyber.grid.v1beta1.MsgUpdateParamsResponse) |  | |

 <!-- end services -->



<a name="cyber/liquidity/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/liquidity/v1beta1/tx.proto



<a name="cyber.liquidity.v1beta1.MsgCreatePool"></a>

### MsgCreatePool
MsgCreatePool defines an sdk.Msg type that supports submitting a create
liquidity pool tx.

See:
https://github.com/gravity-devs/liquidity/blob/develop/x/liquidity/spec/04_messages.md


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pool_creator_address` | [string](#string) |  |  |
| `pool_type_id` | [uint32](#uint32) |  | id of the target pool type, must match the value in the pool. Only pool-type-id 1 is supported. |
| `deposit_coins` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | reserve coin pair of the pool to deposit. |






<a name="cyber.liquidity.v1beta1.MsgCreatePoolResponse"></a>

### MsgCreatePoolResponse
MsgCreatePoolResponse defines the Msg/CreatePool response type.






<a name="cyber.liquidity.v1beta1.MsgDepositWithinBatch"></a>

### MsgDepositWithinBatch
`MsgDepositWithinBatch defines` an `sdk.Msg` type that supports submitting
a deposit request to the batch of the liquidity pool.
Deposit is submitted to the batch of the Liquidity pool with the specified
`pool_id`, `deposit_coins` for reserve.
This request is stacked in the batch of the liquidity pool, is not processed
immediately, and is processed in the `endblock` at the same time as other
requests.

See:
https://github.com/gravity-devs/liquidity/blob/develop/x/liquidity/spec/04_messages.md


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `depositor_address` | [string](#string) |  |  |
| `pool_id` | [uint64](#uint64) |  | id of the target pool |
| `deposit_coins` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | reserve coin pair of the pool to deposit |






<a name="cyber.liquidity.v1beta1.MsgDepositWithinBatchResponse"></a>

### MsgDepositWithinBatchResponse
MsgDepositWithinBatchResponse defines the Msg/DepositWithinBatch response
type.






<a name="cyber.liquidity.v1beta1.MsgSwapWithinBatch"></a>

### MsgSwapWithinBatch
`MsgSwapWithinBatch` defines an sdk.Msg type that supports submitting a swap
offer request to the batch of the liquidity pool. Submit swap offer to the
liquidity pool batch with the specified the `pool_id`, `swap_type_id`,
`demand_coin_denom` with the coin and the price you're offering
and `offer_coin_fee` must be half of offer coin amount * current
`params.swap_fee_rate` and ceil for reservation to pay fees. This request is
stacked in the batch of the liquidity pool, is not processed immediately, and
is processed in the `endblock` at the same time as other requests. You must
request the same fields as the pool. Only the default `swap_type_id` 1 is
supported.

See: https://github.com/gravity-devs/liquidity/tree/develop/doc
https://github.com/gravity-devs/liquidity/blob/develop/x/liquidity/spec/04_messages.md


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `swap_requester_address` | [string](#string) |  | address of swap requester |
| `pool_id` | [uint64](#uint64) |  | id of swap type, must match the value in the pool. Only `swap_type_id` 1 is supported. |
| `swap_type_id` | [uint32](#uint32) |  | id of swap type. Must match the value in the pool. |
| `offer_coin` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | offer sdk.coin for the swap request, must match the denom in the pool. |
| `demand_coin_denom` | [string](#string) |  | denom of demand coin to be exchanged on the swap request, must match the denom in the pool. |
| `offer_coin_fee` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | half of offer coin amount * params.swap_fee_rate and ceil for reservation to pay fees. |
| `order_price` | [string](#string) |  | limit order price for the order, the price is the exchange ratio of X/Y where X is the amount of the first coin and Y is the amount of the second coin when their denoms are sorted alphabetically. |






<a name="cyber.liquidity.v1beta1.MsgSwapWithinBatchResponse"></a>

### MsgSwapWithinBatchResponse
MsgSwapWithinBatchResponse defines the Msg/Swap response type.






<a name="cyber.liquidity.v1beta1.MsgWithdrawWithinBatch"></a>

### MsgWithdrawWithinBatch
`MsgWithdrawWithinBatch` defines an `sdk.Msg` type that supports submitting
a withdraw request to the batch of the liquidity pool.
Withdraw is submitted to the batch from the Liquidity pool with the
specified `pool_id`, `pool_coin` of the pool.
This request is stacked in the batch of the liquidity pool, is not processed
immediately, and is processed in the `endblock` at the same time as other
requests.

See:
https://github.com/gravity-devs/liquidity/blob/develop/x/liquidity/spec/04_messages.md


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `withdrawer_address` | [string](#string) |  |  |
| `pool_id` | [uint64](#uint64) |  | id of the target pool |
| `pool_coin` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="cyber.liquidity.v1beta1.MsgWithdrawWithinBatchResponse"></a>

### MsgWithdrawWithinBatchResponse
MsgWithdrawWithinBatchResponse defines the Msg/WithdrawWithinBatch response
type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cyber.liquidity.v1beta1.Msg"></a>

### Msg
Msg defines the liquidity Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `CreatePool` | [MsgCreatePool](#cyber.liquidity.v1beta1.MsgCreatePool) | [MsgCreatePoolResponse](#cyber.liquidity.v1beta1.MsgCreatePoolResponse) | Submit a create liquidity pool message. | |
| `DepositWithinBatch` | [MsgDepositWithinBatch](#cyber.liquidity.v1beta1.MsgDepositWithinBatch) | [MsgDepositWithinBatchResponse](#cyber.liquidity.v1beta1.MsgDepositWithinBatchResponse) | Submit a deposit to the liquidity pool batch. | |
| `WithdrawWithinBatch` | [MsgWithdrawWithinBatch](#cyber.liquidity.v1beta1.MsgWithdrawWithinBatch) | [MsgWithdrawWithinBatchResponse](#cyber.liquidity.v1beta1.MsgWithdrawWithinBatchResponse) | Submit a withdraw from the liquidity pool batch. | |
| `Swap` | [MsgSwapWithinBatch](#cyber.liquidity.v1beta1.MsgSwapWithinBatch) | [MsgSwapWithinBatchResponse](#cyber.liquidity.v1beta1.MsgSwapWithinBatchResponse) | Submit a swap to the liquidity pool batch. | |

 <!-- end services -->



<a name="cyber/liquidity/v1beta1/liquidity.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/liquidity/v1beta1/liquidity.proto



<a name="cyber.liquidity.v1beta1.DepositMsgState"></a>

### DepositMsgState
DepositMsgState defines the state of deposit message that contains state
information as it is processed in the next batch or batches.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `msg_height` | [int64](#int64) |  | height where this message is appended to the batch |
| `msg_index` | [uint64](#uint64) |  | index of this deposit message in this liquidity pool |
| `executed` | [bool](#bool) |  | true if executed on this batch, false if not executed |
| `succeeded` | [bool](#bool) |  | true if executed successfully on this batch, false if failed |
| `to_be_deleted` | [bool](#bool) |  | true if ready to be deleted on kvstore, false if not ready to be deleted |
| `msg` | [MsgDepositWithinBatch](#cyber.liquidity.v1beta1.MsgDepositWithinBatch) |  | MsgDepositWithinBatch |






<a name="cyber.liquidity.v1beta1.Params"></a>

### Params
Params defines the parameters for the liquidity module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pool_types` | [PoolType](#cyber.liquidity.v1beta1.PoolType) | repeated | list of available pool types |
| `min_init_deposit_amount` | [string](#string) |  | Minimum number of coins to be deposited to the liquidity pool on pool creation. |
| `init_pool_coin_mint_amount` | [string](#string) |  | Initial mint amount of pool coins upon pool creation. |
| `max_reserve_coin_amount` | [string](#string) |  | Limit the size of each liquidity pool to minimize risk. In development, set to 0 for no limit. In production, set a limit. |
| `pool_creation_fee` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | Fee paid to create a Liquidity Pool. Set a fee to prevent spamming. |
| `swap_fee_rate` | [string](#string) |  | Swap fee rate for every executed swap. |
| `withdraw_fee_rate` | [string](#string) |  | Reserve coin withdrawal with less proportion by withdrawFeeRate. |
| `max_order_amount_ratio` | [string](#string) |  | Maximum ratio of reserve coins that can be ordered at a swap order. |
| `unit_batch_height` | [uint32](#uint32) |  | The smallest unit batch height for every liquidity pool. |
| `circuit_breaker_enabled` | [bool](#bool) |  | Circuit breaker enables or disables transaction messages in liquidity module. |






<a name="cyber.liquidity.v1beta1.Pool"></a>

### Pool
Pool defines the liquidity pool that contains pool information.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  | id of the pool |
| `type_id` | [uint32](#uint32) |  | id of the pool_type |
| `reserve_coin_denoms` | [string](#string) | repeated | denoms of reserve coin pair of the pool |
| `reserve_account_address` | [string](#string) |  | reserve account address of the pool |
| `pool_coin_denom` | [string](#string) |  | denom of pool coin of the pool |






<a name="cyber.liquidity.v1beta1.PoolBatch"></a>

### PoolBatch
PoolBatch defines the batch or batches of a given liquidity pool that
contains indexes of deposit, withdraw, and swap messages. Index param
increments by 1 if the pool id is same.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pool_id` | [uint64](#uint64) |  | id of the pool |
| `index` | [uint64](#uint64) |  | index of this batch |
| `begin_height` | [int64](#int64) |  | height where this batch is started |
| `deposit_msg_index` | [uint64](#uint64) |  | last index of DepositMsgStates |
| `withdraw_msg_index` | [uint64](#uint64) |  | last index of WithdrawMsgStates |
| `swap_msg_index` | [uint64](#uint64) |  | last index of SwapMsgStates |
| `executed` | [bool](#bool) |  | true if executed, false if not executed |






<a name="cyber.liquidity.v1beta1.PoolMetadata"></a>

### PoolMetadata
Metadata for the state of each pool for invariant checking after genesis
export or import.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pool_id` | [uint64](#uint64) |  | id of the pool |
| `pool_coin_total_supply` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | pool coin issued at the pool |
| `reserve_coins` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | reserve coins deposited in the pool |






<a name="cyber.liquidity.v1beta1.PoolType"></a>

### PoolType
Structure for the pool type to distinguish the characteristics of the reserve
pools.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint32](#uint32) |  | This is the id of the pool_type that is used as pool_type_id for pool creation. In this version, only pool-type-id 1 is supported. {"id":1,"name":"ConstantProductLiquidityPool","min_reserve_coin_num":2,"max_reserve_coin_num":2,"description":""} |
| `name` | [string](#string) |  | name of the pool type. |
| `min_reserve_coin_num` | [uint32](#uint32) |  | minimum number of reserveCoins for LiquidityPoolType, only 2 reserve coins are supported. |
| `max_reserve_coin_num` | [uint32](#uint32) |  | maximum number of reserveCoins for LiquidityPoolType, only 2 reserve coins are supported. |
| `description` | [string](#string) |  | description of the pool type. |






<a name="cyber.liquidity.v1beta1.SwapMsgState"></a>

### SwapMsgState
SwapMsgState defines the state of the swap message that contains state
information as the message is processed in the next batch or batches.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `msg_height` | [int64](#int64) |  | height where this message is appended to the batch |
| `msg_index` | [uint64](#uint64) |  | index of this swap message in this liquidity pool |
| `executed` | [bool](#bool) |  | true if executed on this batch, false if not executed |
| `succeeded` | [bool](#bool) |  | true if executed successfully on this batch, false if failed |
| `to_be_deleted` | [bool](#bool) |  | true if ready to be deleted on kvstore, false if not ready to be deleted |
| `order_expiry_height` | [int64](#int64) |  | swap orders are cancelled when current height is equal to or higher than ExpiryHeight |
| `exchanged_offer_coin` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | offer coin exchanged until now |
| `remaining_offer_coin` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | offer coin currently remaining to be exchanged |
| `reserved_offer_coin_fee` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | reserve fee for pays fee in half offer coin |
| `msg` | [MsgSwapWithinBatch](#cyber.liquidity.v1beta1.MsgSwapWithinBatch) |  | MsgSwapWithinBatch |






<a name="cyber.liquidity.v1beta1.WithdrawMsgState"></a>

### WithdrawMsgState
WithdrawMsgState defines the state of the withdraw message that contains
state information as the message is processed in the next batch or batches.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `msg_height` | [int64](#int64) |  | height where this message is appended to the batch |
| `msg_index` | [uint64](#uint64) |  | index of this withdraw message in this liquidity pool |
| `executed` | [bool](#bool) |  | true if executed on this batch, false if not executed |
| `succeeded` | [bool](#bool) |  | true if executed successfully on this batch, false if failed |
| `to_be_deleted` | [bool](#bool) |  | true if ready to be deleted on kvstore, false if not ready to be deleted |
| `msg` | [MsgWithdrawWithinBatch](#cyber.liquidity.v1beta1.MsgWithdrawWithinBatch) |  | MsgWithdrawWithinBatch |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cyber/liquidity/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/liquidity/v1beta1/genesis.proto



<a name="cyber.liquidity.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the liquidity module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cyber.liquidity.v1beta1.Params) |  | params defines all the parameters for the liquidity module. |
| `pool_records` | [PoolRecord](#cyber.liquidity.v1beta1.PoolRecord) | repeated |  |






<a name="cyber.liquidity.v1beta1.PoolRecord"></a>

### PoolRecord
records the state of each pool after genesis export or import, used to check
variables


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pool` | [Pool](#cyber.liquidity.v1beta1.Pool) |  |  |
| `pool_metadata` | [PoolMetadata](#cyber.liquidity.v1beta1.PoolMetadata) |  |  |
| `pool_batch` | [PoolBatch](#cyber.liquidity.v1beta1.PoolBatch) |  |  |
| `deposit_msg_states` | [DepositMsgState](#cyber.liquidity.v1beta1.DepositMsgState) | repeated |  |
| `withdraw_msg_states` | [WithdrawMsgState](#cyber.liquidity.v1beta1.WithdrawMsgState) | repeated |  |
| `swap_msg_states` | [SwapMsgState](#cyber.liquidity.v1beta1.SwapMsgState) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cyber/liquidity/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/liquidity/v1beta1/query.proto



<a name="cyber.liquidity.v1beta1.QueryLiquidityPoolBatchRequest"></a>

### QueryLiquidityPoolBatchRequest
the request type for the QueryLiquidityPoolBatch RPC method. requestable
including specified pool_id.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pool_id` | [uint64](#uint64) |  | id of the target pool for query |






<a name="cyber.liquidity.v1beta1.QueryLiquidityPoolBatchResponse"></a>

### QueryLiquidityPoolBatchResponse
the response type for the QueryLiquidityPoolBatchResponse RPC method. Returns
the liquidity pool batch that corresponds to the requested pool_id.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `batch` | [PoolBatch](#cyber.liquidity.v1beta1.PoolBatch) |  |  |






<a name="cyber.liquidity.v1beta1.QueryLiquidityPoolByPoolCoinDenomRequest"></a>

### QueryLiquidityPoolByPoolCoinDenomRequest
the request type for the QueryLiquidityByPoolCoinDenomPool RPC method.
Requestable specified pool_coin_denom.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pool_coin_denom` | [string](#string) |  |  |






<a name="cyber.liquidity.v1beta1.QueryLiquidityPoolByReserveAccRequest"></a>

### QueryLiquidityPoolByReserveAccRequest
the request type for the QueryLiquidityByReserveAcc RPC method. Requestable
specified reserve_acc.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `reserve_acc` | [string](#string) |  |  |






<a name="cyber.liquidity.v1beta1.QueryLiquidityPoolRequest"></a>

### QueryLiquidityPoolRequest
the request type for the QueryLiquidityPool RPC method. requestable specified
pool_id.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pool_id` | [uint64](#uint64) |  |  |






<a name="cyber.liquidity.v1beta1.QueryLiquidityPoolResponse"></a>

### QueryLiquidityPoolResponse
the response type for the QueryLiquidityPoolResponse RPC method. Returns the
liquidity pool that corresponds to the requested pool_id.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pool` | [Pool](#cyber.liquidity.v1beta1.Pool) |  |  |






<a name="cyber.liquidity.v1beta1.QueryLiquidityPoolsRequest"></a>

### QueryLiquidityPoolsRequest
the request type for the QueryLiquidityPools RPC method. Requestable
including pagination offset, limit, key.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cyber.liquidity.v1beta1.QueryLiquidityPoolsResponse"></a>

### QueryLiquidityPoolsResponse
the response type for the QueryLiquidityPoolsResponse RPC method. This
includes a list of all existing liquidity pools and paging results that
contain next_key and total count.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pools` | [Pool](#cyber.liquidity.v1beta1.Pool) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. not working on this version. |






<a name="cyber.liquidity.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is request type for the QueryParams RPC method.






<a name="cyber.liquidity.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse
the response type for the QueryParamsResponse RPC method. This includes
current parameter of the liquidity module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cyber.liquidity.v1beta1.Params) |  | params holds all the parameters of this module. |






<a name="cyber.liquidity.v1beta1.QueryPoolBatchDepositMsgRequest"></a>

### QueryPoolBatchDepositMsgRequest
the request type for the QueryPoolBatchDeposit RPC method. requestable
including specified pool_id and msg_index.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pool_id` | [uint64](#uint64) |  | id of the target pool for query |
| `msg_index` | [uint64](#uint64) |  | target msg_index of the pool |






<a name="cyber.liquidity.v1beta1.QueryPoolBatchDepositMsgResponse"></a>

### QueryPoolBatchDepositMsgResponse
the response type for the QueryPoolBatchDepositMsg RPC method. This includes
a batch swap message of the batch.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `deposit` | [DepositMsgState](#cyber.liquidity.v1beta1.DepositMsgState) |  |  |






<a name="cyber.liquidity.v1beta1.QueryPoolBatchDepositMsgsRequest"></a>

### QueryPoolBatchDepositMsgsRequest
the request type for the QueryPoolBatchDeposit RPC method. Requestable
including specified pool_id and pagination offset, limit, key.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pool_id` | [uint64](#uint64) |  | id of the target pool for query |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cyber.liquidity.v1beta1.QueryPoolBatchDepositMsgsResponse"></a>

### QueryPoolBatchDepositMsgsResponse
the response type for the QueryPoolBatchDeposit RPC method. This includes a
list of all currently existing deposit messages of the batch and paging
results that contain next_key and total count.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `deposits` | [DepositMsgState](#cyber.liquidity.v1beta1.DepositMsgState) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. not working on this version. |






<a name="cyber.liquidity.v1beta1.QueryPoolBatchSwapMsgRequest"></a>

### QueryPoolBatchSwapMsgRequest
the request type for the QueryPoolBatchSwap RPC method. Requestable including
specified pool_id and msg_index.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pool_id` | [uint64](#uint64) |  | id of the target pool for query |
| `msg_index` | [uint64](#uint64) |  | target msg_index of the pool |






<a name="cyber.liquidity.v1beta1.QueryPoolBatchSwapMsgResponse"></a>

### QueryPoolBatchSwapMsgResponse
the response type for the QueryPoolBatchSwapMsg RPC method. This includes a
batch swap message of the batch.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `swap` | [SwapMsgState](#cyber.liquidity.v1beta1.SwapMsgState) |  |  |






<a name="cyber.liquidity.v1beta1.QueryPoolBatchSwapMsgsRequest"></a>

### QueryPoolBatchSwapMsgsRequest
the request type for the QueryPoolBatchSwapMsgs RPC method. Requestable
including specified pool_id and pagination offset, limit, key.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pool_id` | [uint64](#uint64) |  | id of the target pool for query |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cyber.liquidity.v1beta1.QueryPoolBatchSwapMsgsResponse"></a>

### QueryPoolBatchSwapMsgsResponse
the response type for the QueryPoolBatchSwapMsgs RPC method. This includes
list of all currently existing swap messages of the batch and paging results
that contain next_key and total count.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `swaps` | [SwapMsgState](#cyber.liquidity.v1beta1.SwapMsgState) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. not working on this version. |






<a name="cyber.liquidity.v1beta1.QueryPoolBatchWithdrawMsgRequest"></a>

### QueryPoolBatchWithdrawMsgRequest
the request type for the QueryPoolBatchWithdraw RPC method. requestable
including specified pool_id and msg_index.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pool_id` | [uint64](#uint64) |  | id of the target pool for query |
| `msg_index` | [uint64](#uint64) |  | target msg_index of the pool |






<a name="cyber.liquidity.v1beta1.QueryPoolBatchWithdrawMsgResponse"></a>

### QueryPoolBatchWithdrawMsgResponse
the response type for the QueryPoolBatchWithdrawMsg RPC method. This includes
a batch swap message of the batch.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `withdraw` | [WithdrawMsgState](#cyber.liquidity.v1beta1.WithdrawMsgState) |  |  |






<a name="cyber.liquidity.v1beta1.QueryPoolBatchWithdrawMsgsRequest"></a>

### QueryPoolBatchWithdrawMsgsRequest
the request type for the QueryPoolBatchWithdraw RPC method. Requestable
including specified pool_id and pagination offset, limit, key.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pool_id` | [uint64](#uint64) |  | id of the target pool for query |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cyber.liquidity.v1beta1.QueryPoolBatchWithdrawMsgsResponse"></a>

### QueryPoolBatchWithdrawMsgsResponse
the response type for the QueryPoolBatchWithdraw RPC method. This includes a
list of all currently existing withdraw messages of the batch and paging
results that contain next_key and total count.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `withdraws` | [WithdrawMsgState](#cyber.liquidity.v1beta1.WithdrawMsgState) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. Not supported on this version. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cyber.liquidity.v1beta1.Query"></a>

### Query
Query defines the gRPC query service for the liquidity module.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `LiquidityPools` | [QueryLiquidityPoolsRequest](#cyber.liquidity.v1beta1.QueryLiquidityPoolsRequest) | [QueryLiquidityPoolsResponse](#cyber.liquidity.v1beta1.QueryLiquidityPoolsResponse) | Get existing liquidity pools. | GET|/cosmos/liquidity/v1beta1/pools|
| `LiquidityPool` | [QueryLiquidityPoolRequest](#cyber.liquidity.v1beta1.QueryLiquidityPoolRequest) | [QueryLiquidityPoolResponse](#cyber.liquidity.v1beta1.QueryLiquidityPoolResponse) | Get specific liquidity pool. | GET|/cosmos/liquidity/v1beta1/pools/{pool_id}|
| `LiquidityPoolByPoolCoinDenom` | [QueryLiquidityPoolByPoolCoinDenomRequest](#cyber.liquidity.v1beta1.QueryLiquidityPoolByPoolCoinDenomRequest) | [QueryLiquidityPoolResponse](#cyber.liquidity.v1beta1.QueryLiquidityPoolResponse) | Get specific liquidity pool corresponding to the pool_coin_denom. | GET|/cosmos/liquidity/v1beta1/pools/pool_coin_denom/{pool_coin_denom}|
| `LiquidityPoolByReserveAcc` | [QueryLiquidityPoolByReserveAccRequest](#cyber.liquidity.v1beta1.QueryLiquidityPoolByReserveAccRequest) | [QueryLiquidityPoolResponse](#cyber.liquidity.v1beta1.QueryLiquidityPoolResponse) | Get specific liquidity pool corresponding to the reserve account. | GET|/cosmos/liquidity/v1beta1/pools/reserve_acc/{reserve_acc}|
| `LiquidityPoolBatch` | [QueryLiquidityPoolBatchRequest](#cyber.liquidity.v1beta1.QueryLiquidityPoolBatchRequest) | [QueryLiquidityPoolBatchResponse](#cyber.liquidity.v1beta1.QueryLiquidityPoolBatchResponse) | Get the pool's current batch. | GET|/cosmos/liquidity/v1beta1/pools/{pool_id}/batch|
| `PoolBatchSwapMsgs` | [QueryPoolBatchSwapMsgsRequest](#cyber.liquidity.v1beta1.QueryPoolBatchSwapMsgsRequest) | [QueryPoolBatchSwapMsgsResponse](#cyber.liquidity.v1beta1.QueryPoolBatchSwapMsgsResponse) | Get all swap messages in the pool's current batch. | GET|/cosmos/liquidity/v1beta1/pools/{pool_id}/batch/swaps|
| `PoolBatchSwapMsg` | [QueryPoolBatchSwapMsgRequest](#cyber.liquidity.v1beta1.QueryPoolBatchSwapMsgRequest) | [QueryPoolBatchSwapMsgResponse](#cyber.liquidity.v1beta1.QueryPoolBatchSwapMsgResponse) | Get a specific swap message in the pool's current batch. | GET|/cosmos/liquidity/v1beta1/pools/{pool_id}/batch/swaps/{msg_index}|
| `PoolBatchDepositMsgs` | [QueryPoolBatchDepositMsgsRequest](#cyber.liquidity.v1beta1.QueryPoolBatchDepositMsgsRequest) | [QueryPoolBatchDepositMsgsResponse](#cyber.liquidity.v1beta1.QueryPoolBatchDepositMsgsResponse) | Get all deposit messages in the pool's current batch. | GET|/cosmos/liquidity/v1beta1/pools/{pool_id}/batch/deposits|
| `PoolBatchDepositMsg` | [QueryPoolBatchDepositMsgRequest](#cyber.liquidity.v1beta1.QueryPoolBatchDepositMsgRequest) | [QueryPoolBatchDepositMsgResponse](#cyber.liquidity.v1beta1.QueryPoolBatchDepositMsgResponse) | Get a specific deposit message in the pool's current batch. | GET|/cosmos/liquidity/v1beta1/pools/{pool_id}/batch/deposits/{msg_index}|
| `PoolBatchWithdrawMsgs` | [QueryPoolBatchWithdrawMsgsRequest](#cyber.liquidity.v1beta1.QueryPoolBatchWithdrawMsgsRequest) | [QueryPoolBatchWithdrawMsgsResponse](#cyber.liquidity.v1beta1.QueryPoolBatchWithdrawMsgsResponse) | Get all withdraw messages in the pool's current batch. | GET|/cosmos/liquidity/v1beta1/pools/{pool_id}/batch/withdraws|
| `PoolBatchWithdrawMsg` | [QueryPoolBatchWithdrawMsgRequest](#cyber.liquidity.v1beta1.QueryPoolBatchWithdrawMsgRequest) | [QueryPoolBatchWithdrawMsgResponse](#cyber.liquidity.v1beta1.QueryPoolBatchWithdrawMsgResponse) | Get a specific withdraw message in the pool's current batch. | GET|/cosmos/liquidity/v1beta1/pools/{pool_id}/batch/withdraws/{msg_index}|
| `Params` | [QueryParamsRequest](#cyber.liquidity.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#cyber.liquidity.v1beta1.QueryParamsResponse) | Get all parameters of the liquidity module. | GET|/cosmos/liquidity/v1beta1/params|

 <!-- end services -->



<a name="cyber/rank/v1beta1/types.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/rank/v1beta1/types.proto



<a name="cyber.rank.v1beta1.Params"></a>

### Params



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `calculation_period` | [int64](#int64) |  |  |
| `damping_factor` | [string](#string) |  |  |
| `tolerance` | [string](#string) |  |  |






<a name="cyber.rank.v1beta1.RankedParticle"></a>

### RankedParticle



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `particle` | [string](#string) |  |  |
| `rank` | [uint64](#uint64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cyber/rank/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/rank/v1beta1/genesis.proto



<a name="cyber.rank.v1beta1.GenesisState"></a>

### GenesisState



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cyber.rank.v1beta1.Params) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cyber/rank/v1beta1/pagination.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/rank/v1beta1/pagination.proto



<a name="cyber.rank.v1beta1.PageRequest"></a>

### PageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `page` | [uint32](#uint32) |  |  |
| `per_page` | [uint32](#uint32) |  |  |






<a name="cyber.rank.v1beta1.PageResponse"></a>

### PageResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `total` | [uint32](#uint32) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cyber/rank/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/rank/v1beta1/query.proto



<a name="cyber.rank.v1beta1.QueryIsAnyLinkExistRequest"></a>

### QueryIsAnyLinkExistRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `from` | [string](#string) |  |  |
| `to` | [string](#string) |  |  |






<a name="cyber.rank.v1beta1.QueryIsLinkExistRequest"></a>

### QueryIsLinkExistRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `from` | [string](#string) |  |  |
| `to` | [string](#string) |  |  |
| `address` | [string](#string) |  |  |






<a name="cyber.rank.v1beta1.QueryKarmaRequest"></a>

### QueryKarmaRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `neuron` | [string](#string) |  |  |






<a name="cyber.rank.v1beta1.QueryKarmaResponse"></a>

### QueryKarmaResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `karma` | [uint64](#uint64) |  |  |






<a name="cyber.rank.v1beta1.QueryLinkExistResponse"></a>

### QueryLinkExistResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `exist` | [bool](#bool) |  |  |






<a name="cyber.rank.v1beta1.QueryNegentropyParticleResponse"></a>

### QueryNegentropyParticleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entropy` | [uint64](#uint64) |  |  |






<a name="cyber.rank.v1beta1.QueryNegentropyPartilceRequest"></a>

### QueryNegentropyPartilceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `particle` | [string](#string) |  |  |






<a name="cyber.rank.v1beta1.QueryNegentropyRequest"></a>

### QueryNegentropyRequest







<a name="cyber.rank.v1beta1.QueryNegentropyResponse"></a>

### QueryNegentropyResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `negentropy` | [uint64](#uint64) |  |  |






<a name="cyber.rank.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest







<a name="cyber.rank.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cyber.rank.v1beta1.Params) |  |  |






<a name="cyber.rank.v1beta1.QueryRankRequest"></a>

### QueryRankRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `particle` | [string](#string) |  |  |






<a name="cyber.rank.v1beta1.QueryRankResponse"></a>

### QueryRankResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `rank` | [uint64](#uint64) |  |  |






<a name="cyber.rank.v1beta1.QuerySearchRequest"></a>

### QuerySearchRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `particle` | [string](#string) |  |  |
| `pagination` | [PageRequest](#cyber.rank.v1beta1.PageRequest) |  |  |






<a name="cyber.rank.v1beta1.QuerySearchResponse"></a>

### QuerySearchResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `result` | [RankedParticle](#cyber.rank.v1beta1.RankedParticle) | repeated |  |
| `pagination` | [PageResponse](#cyber.rank.v1beta1.PageResponse) |  |  |






<a name="cyber.rank.v1beta1.QueryTopRequest"></a>

### QueryTopRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [PageRequest](#cyber.rank.v1beta1.PageRequest) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cyber.rank.v1beta1.Query"></a>

### Query


| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#cyber.rank.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#cyber.rank.v1beta1.QueryParamsResponse) |  | GET|/cyber/rank/v1beta1/rank/params|
| `Rank` | [QueryRankRequest](#cyber.rank.v1beta1.QueryRankRequest) | [QueryRankResponse](#cyber.rank.v1beta1.QueryRankResponse) |  | GET|/cyber/rank/v1beta1/rank/rank/{particle}|
| `Search` | [QuerySearchRequest](#cyber.rank.v1beta1.QuerySearchRequest) | [QuerySearchResponse](#cyber.rank.v1beta1.QuerySearchResponse) |  | GET|/cyber/rank/v1beta1/rank/search/{particle}|
| `Backlinks` | [QuerySearchRequest](#cyber.rank.v1beta1.QuerySearchRequest) | [QuerySearchResponse](#cyber.rank.v1beta1.QuerySearchResponse) |  | GET|/cyber/rank/v1beta1/rank/backlinks/{particle}|
| `Top` | [QueryTopRequest](#cyber.rank.v1beta1.QueryTopRequest) | [QuerySearchResponse](#cyber.rank.v1beta1.QuerySearchResponse) |  | GET|/cyber/rank/v1beta1/rank/top|
| `IsLinkExist` | [QueryIsLinkExistRequest](#cyber.rank.v1beta1.QueryIsLinkExistRequest) | [QueryLinkExistResponse](#cyber.rank.v1beta1.QueryLinkExistResponse) |  | GET|/cyber/rank/v1beta1/is_link_exist|
| `IsAnyLinkExist` | [QueryIsAnyLinkExistRequest](#cyber.rank.v1beta1.QueryIsAnyLinkExistRequest) | [QueryLinkExistResponse](#cyber.rank.v1beta1.QueryLinkExistResponse) |  | GET|/cyber/rank/v1beta1/is_any_link_exist|
| `ParticleNegentropy` | [QueryNegentropyPartilceRequest](#cyber.rank.v1beta1.QueryNegentropyPartilceRequest) | [QueryNegentropyParticleResponse](#cyber.rank.v1beta1.QueryNegentropyParticleResponse) |  | GET|/cyber/rank/v1beta1/negentropy/{particle}|
| `Negentropy` | [QueryNegentropyRequest](#cyber.rank.v1beta1.QueryNegentropyRequest) | [QueryNegentropyResponse](#cyber.rank.v1beta1.QueryNegentropyResponse) |  | GET|/cyber/rank/v1beta1/negentropy|
| `Karma` | [QueryKarmaRequest](#cyber.rank.v1beta1.QueryKarmaRequest) | [QueryKarmaResponse](#cyber.rank.v1beta1.QueryKarmaResponse) |  | GET|/cyber/rank/v1beta1/karma/{neuron}|

 <!-- end services -->



<a name="cyber/rank/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/rank/v1beta1/tx.proto



<a name="cyber.rank.v1beta1.MsgUpdateParams"></a>

### MsgUpdateParams



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authority` | [string](#string) |  |  |
| `params` | [Params](#cyber.rank.v1beta1.Params) |  |  |






<a name="cyber.rank.v1beta1.MsgUpdateParamsResponse"></a>

### MsgUpdateParamsResponse






 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cyber.rank.v1beta1.Msg"></a>

### Msg


| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `UpdateParams` | [MsgUpdateParams](#cyber.rank.v1beta1.MsgUpdateParams) | [MsgUpdateParamsResponse](#cyber.rank.v1beta1.MsgUpdateParamsResponse) |  | |

 <!-- end services -->



<a name="cyber/resources/v1beta1/types.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/resources/v1beta1/types.proto



<a name="cyber.resources.v1beta1.Params"></a>

### Params



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `max_slots` | [uint32](#uint32) |  |  |
| `halving_period_volt_blocks` | [uint32](#uint32) |  |  |
| `halving_period_ampere_blocks` | [uint32](#uint32) |  |  |
| `base_investmint_period_volt` | [uint32](#uint32) |  |  |
| `base_investmint_period_ampere` | [uint32](#uint32) |  |  |
| `min_investmint_period` | [uint32](#uint32) |  |  |
| `base_investmint_amount_volt` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `base_investmint_amount_ampere` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cyber/resources/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/resources/v1beta1/genesis.proto



<a name="cyber.resources.v1beta1.GenesisState"></a>

### GenesisState



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cyber.resources.v1beta1.Params) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cyber/resources/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/resources/v1beta1/query.proto



<a name="cyber.resources.v1beta1.QueryInvestmintRequest"></a>

### QueryInvestmintRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `resource` | [string](#string) |  |  |
| `length` | [uint64](#uint64) |  |  |






<a name="cyber.resources.v1beta1.QueryInvestmintResponse"></a>

### QueryInvestmintResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="cyber.resources.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest







<a name="cyber.resources.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cyber.resources.v1beta1.Params) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cyber.resources.v1beta1.Query"></a>

### Query


| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#cyber.resources.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#cyber.resources.v1beta1.QueryParamsResponse) |  | GET|/cyber/resources/v1beta1/resources/params|
| `Investmint` | [QueryInvestmintRequest](#cyber.resources.v1beta1.QueryInvestmintRequest) | [QueryInvestmintResponse](#cyber.resources.v1beta1.QueryInvestmintResponse) |  | GET|/cyber/resources/v1beta1/resources/investmint|

 <!-- end services -->



<a name="cyber/resources/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cyber/resources/v1beta1/tx.proto



<a name="cyber.resources.v1beta1.MsgInvestmint"></a>

### MsgInvestmint



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `neuron` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `resource` | [string](#string) |  |  |
| `length` | [uint64](#uint64) |  |  |






<a name="cyber.resources.v1beta1.MsgInvestmintResponse"></a>

### MsgInvestmintResponse







<a name="cyber.resources.v1beta1.MsgUpdateParams"></a>

### MsgUpdateParams



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authority` | [string](#string) |  |  |
| `params` | [Params](#cyber.resources.v1beta1.Params) |  |  |






<a name="cyber.resources.v1beta1.MsgUpdateParamsResponse"></a>

### MsgUpdateParamsResponse






 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cyber.resources.v1beta1.Msg"></a>

### Msg


| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Investmint` | [MsgInvestmint](#cyber.resources.v1beta1.MsgInvestmint) | [MsgInvestmintResponse](#cyber.resources.v1beta1.MsgInvestmintResponse) |  | |
| `UpdateParams` | [MsgUpdateParams](#cyber.resources.v1beta1.MsgUpdateParams) | [MsgUpdateParamsResponse](#cyber.resources.v1beta1.MsgUpdateParamsResponse) |  | |

 <!-- end services -->



<a name="osmosis/tokenfactory/v1beta1/authority_metadata.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## osmosis/tokenfactory/v1beta1/authority_metadata.proto



<a name="osmosis.tokenfactory.v1beta1.DenomAuthorityMetadata"></a>

### DenomAuthorityMetadata
DenomAuthorityMetadata specifies metadata for addresses that have specific
capabilities over a token factory denom. Right now there is only one Admin
permission, but is planned to be extended to the future.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `admin` | [string](#string) |  | Can be empty for no admin, or a valid osmosis address |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="osmosis/tokenfactory/v1beta1/params.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## osmosis/tokenfactory/v1beta1/params.proto



<a name="osmosis.tokenfactory.v1beta1.Params"></a>

### Params
Params defines the parameters for the tokenfactory module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom_creation_fee` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
| `denom_creation_gas_consume` | [uint64](#uint64) |  | if denom_creation_fee is an empty array, then this field is used to add more gas consumption to the base cost. https://github.com/CosmWasm/token-factory/issues/11 |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="osmosis/tokenfactory/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## osmosis/tokenfactory/v1beta1/genesis.proto



<a name="osmosis.tokenfactory.v1beta1.GenesisDenom"></a>

### GenesisDenom
GenesisDenom defines a tokenfactory denom that is defined within genesis
state. The structure contains DenomAuthorityMetadata which defines the
denom's admin.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |
| `authority_metadata` | [DenomAuthorityMetadata](#osmosis.tokenfactory.v1beta1.DenomAuthorityMetadata) |  |  |






<a name="osmosis.tokenfactory.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the tokenfactory module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#osmosis.tokenfactory.v1beta1.Params) |  | params defines the parameters of the module. |
| `factory_denoms` | [GenesisDenom](#osmosis.tokenfactory.v1beta1.GenesisDenom) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="osmosis/tokenfactory/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## osmosis/tokenfactory/v1beta1/query.proto



<a name="osmosis.tokenfactory.v1beta1.QueryDenomAuthorityMetadataRequest"></a>

### QueryDenomAuthorityMetadataRequest
QueryDenomAuthorityMetadataRequest defines the request structure for the
DenomAuthorityMetadata gRPC query.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |






<a name="osmosis.tokenfactory.v1beta1.QueryDenomAuthorityMetadataResponse"></a>

### QueryDenomAuthorityMetadataResponse
QueryDenomAuthorityMetadataResponse defines the response structure for the
DenomAuthorityMetadata gRPC query.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authority_metadata` | [DenomAuthorityMetadata](#osmosis.tokenfactory.v1beta1.DenomAuthorityMetadata) |  |  |






<a name="osmosis.tokenfactory.v1beta1.QueryDenomsFromCreatorRequest"></a>

### QueryDenomsFromCreatorRequest
QueryDenomsFromCreatorRequest defines the request structure for the
DenomsFromCreator gRPC query.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `creator` | [string](#string) |  |  |






<a name="osmosis.tokenfactory.v1beta1.QueryDenomsFromCreatorResponse"></a>

### QueryDenomsFromCreatorResponse
QueryDenomsFromCreatorRequest defines the response structure for the
DenomsFromCreator gRPC query.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denoms` | [string](#string) | repeated |  |






<a name="osmosis.tokenfactory.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.






<a name="osmosis.tokenfactory.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#osmosis.tokenfactory.v1beta1.Params) |  | params defines the parameters of the module. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="osmosis.tokenfactory.v1beta1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#osmosis.tokenfactory.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#osmosis.tokenfactory.v1beta1.QueryParamsResponse) | Params defines a gRPC query method that returns the tokenfactory module's parameters. | GET|/osmosis/tokenfactory/v1beta1/params|
| `DenomAuthorityMetadata` | [QueryDenomAuthorityMetadataRequest](#osmosis.tokenfactory.v1beta1.QueryDenomAuthorityMetadataRequest) | [QueryDenomAuthorityMetadataResponse](#osmosis.tokenfactory.v1beta1.QueryDenomAuthorityMetadataResponse) | DenomAuthorityMetadata defines a gRPC query method for fetching DenomAuthorityMetadata for a particular denom. | GET|/osmosis/tokenfactory/v1beta1/denoms/{denom}/authority_metadata|
| `DenomsFromCreator` | [QueryDenomsFromCreatorRequest](#osmosis.tokenfactory.v1beta1.QueryDenomsFromCreatorRequest) | [QueryDenomsFromCreatorResponse](#osmosis.tokenfactory.v1beta1.QueryDenomsFromCreatorResponse) | DenomsFromCreator defines a gRPC query method for fetching all denominations created by a specific admin/creator. | GET|/osmosis/tokenfactory/v1beta1/denoms_from_creator/{creator}|

 <!-- end services -->



<a name="osmosis/tokenfactory/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## osmosis/tokenfactory/v1beta1/tx.proto



<a name="osmosis.tokenfactory.v1beta1.MsgBurn"></a>

### MsgBurn
MsgBurn is the sdk.Msg type for allowing an admin account to burn
a token.  For now, we only support burning from the sender account.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `burnFromAddress` | [string](#string) |  |  |






<a name="osmosis.tokenfactory.v1beta1.MsgBurnResponse"></a>

### MsgBurnResponse







<a name="osmosis.tokenfactory.v1beta1.MsgChangeAdmin"></a>

### MsgChangeAdmin
MsgChangeAdmin is the sdk.Msg type for allowing an admin account to reassign
adminship of a denom to a new account


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  |  |
| `denom` | [string](#string) |  |  |
| `new_admin` | [string](#string) |  |  |






<a name="osmosis.tokenfactory.v1beta1.MsgChangeAdminResponse"></a>

### MsgChangeAdminResponse
MsgChangeAdminResponse defines the response structure for an executed
MsgChangeAdmin message.






<a name="osmosis.tokenfactory.v1beta1.MsgCreateDenom"></a>

### MsgCreateDenom
MsgCreateDenom defines the message structure for the CreateDenom gRPC service
method. It allows an account to create a new denom. It requires a sender
address and a sub denomination. The (sender_address, sub_denomination) tuple
must be unique and cannot be re-used.

The resulting denom created is defined as
<factory/{creatorAddress}/{subdenom}>. The resulting denom's admin is
originally set to be the creator, but this can be changed later. The token
denom does not indicate the current admin.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  |  |
| `subdenom` | [string](#string) |  | subdenom can be up to 44 "alphanumeric" characters long. |






<a name="osmosis.tokenfactory.v1beta1.MsgCreateDenomResponse"></a>

### MsgCreateDenomResponse
MsgCreateDenomResponse is the return value of MsgCreateDenom
It returns the full string of the newly created denom


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `new_token_denom` | [string](#string) |  |  |






<a name="osmosis.tokenfactory.v1beta1.MsgForceTransfer"></a>

### MsgForceTransfer



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `transferFromAddress` | [string](#string) |  |  |
| `transferToAddress` | [string](#string) |  |  |






<a name="osmosis.tokenfactory.v1beta1.MsgForceTransferResponse"></a>

### MsgForceTransferResponse







<a name="osmosis.tokenfactory.v1beta1.MsgMint"></a>

### MsgMint
MsgMint is the sdk.Msg type for allowing an admin account to mint
more of a token.  For now, we only support minting to the sender account


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `mintToAddress` | [string](#string) |  |  |






<a name="osmosis.tokenfactory.v1beta1.MsgMintResponse"></a>

### MsgMintResponse







<a name="osmosis.tokenfactory.v1beta1.MsgSetDenomMetadata"></a>

### MsgSetDenomMetadata
MsgSetDenomMetadata is the sdk.Msg type for allowing an admin account to set
the denom's bank metadata


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  |  |
| `metadata` | [cosmos.bank.v1beta1.Metadata](#cosmos.bank.v1beta1.Metadata) |  |  |






<a name="osmosis.tokenfactory.v1beta1.MsgSetDenomMetadataResponse"></a>

### MsgSetDenomMetadataResponse
MsgSetDenomMetadataResponse defines the response structure for an executed
MsgSetDenomMetadata message.






<a name="osmosis.tokenfactory.v1beta1.MsgUpdateParams"></a>

### MsgUpdateParams
MsgUpdateParams is the Msg/UpdateParams request type.

Since: cosmos-sdk 0.47


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authority` | [string](#string) |  | authority is the address of the governance account. |
| `params` | [Params](#osmosis.tokenfactory.v1beta1.Params) |  | params defines the x/mint parameters to update.

NOTE: All parameters must be supplied. |






<a name="osmosis.tokenfactory.v1beta1.MsgUpdateParamsResponse"></a>

### MsgUpdateParamsResponse
MsgUpdateParamsResponse defines the response structure for executing a
MsgUpdateParams message.

Since: cosmos-sdk 0.47





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="osmosis.tokenfactory.v1beta1.Msg"></a>

### Msg
Msg defines the tokefactory module's gRPC message service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `CreateDenom` | [MsgCreateDenom](#osmosis.tokenfactory.v1beta1.MsgCreateDenom) | [MsgCreateDenomResponse](#osmosis.tokenfactory.v1beta1.MsgCreateDenomResponse) |  | |
| `Mint` | [MsgMint](#osmosis.tokenfactory.v1beta1.MsgMint) | [MsgMintResponse](#osmosis.tokenfactory.v1beta1.MsgMintResponse) |  | |
| `Burn` | [MsgBurn](#osmosis.tokenfactory.v1beta1.MsgBurn) | [MsgBurnResponse](#osmosis.tokenfactory.v1beta1.MsgBurnResponse) |  | |
| `ChangeAdmin` | [MsgChangeAdmin](#osmosis.tokenfactory.v1beta1.MsgChangeAdmin) | [MsgChangeAdminResponse](#osmosis.tokenfactory.v1beta1.MsgChangeAdminResponse) |  | |
| `SetDenomMetadata` | [MsgSetDenomMetadata](#osmosis.tokenfactory.v1beta1.MsgSetDenomMetadata) | [MsgSetDenomMetadataResponse](#osmosis.tokenfactory.v1beta1.MsgSetDenomMetadataResponse) |  | |
| `ForceTransfer` | [MsgForceTransfer](#osmosis.tokenfactory.v1beta1.MsgForceTransfer) | [MsgForceTransferResponse](#osmosis.tokenfactory.v1beta1.MsgForceTransferResponse) |  | |
| `UpdateParams` | [MsgUpdateParams](#osmosis.tokenfactory.v1beta1.MsgUpdateParams) | [MsgUpdateParamsResponse](#osmosis.tokenfactory.v1beta1.MsgUpdateParamsResponse) | UpdateParams defines a governance operation for updating the x/mint module parameters. The authority is hard-coded to the x/gov module account.

Since: cosmos-sdk 0.47 | |

 <!-- end services -->



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

