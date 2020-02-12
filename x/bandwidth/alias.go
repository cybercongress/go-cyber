package bandwidth

import (
	"github.com/cybercongress/cyberd/x/bandwidth/internal/keeper"
	"github.com/cybercongress/cyberd/x/bandwidth/internal/types"
)

const (
	ModuleName        		= types.ModuleName
	DefaultParamspace 		= types.DefaultParamspace
	StoreKey          		= types.StoreKey
	QuerierRoute            = types.QuerierRoute
	QueryParameters         = types.QueryParameters
	QueryDesirableBandwidth = types.QueryDesirableBandwidth
	QueryMaxBlockBandwidth  = types.QueryMaxBlockBandwidth
	QueryRecoveryPeriod     = types.QueryRecoveryPeriod
	QueryAdjustPricePeriod  = types.QueryAdjustPricePeriod
	QueryBaseCreditPrice    = types.QueryBaseCreditPrice
	QueryTxCost             = types.QueryTxCost
	QueryLinkMsgCost        = types.QueryLinkMsgCost
	QueryNonLinkMsgCost     = types.QueryNonLinkMsgCost
)

var (
	// functions aliases
	NewAccountBandwidthKeeper    = keeper.NewAccountBandwidthKeeper
	NewBlockSpentBandwidthKeeper = keeper.NewBlockSpentBandwidthKeeper
	NewQuerier          = keeper.NewQuerier
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
	ParamKeyTable       = types.ParamKeyTable
	NewParams			= types.NewParams
	DefaultParams       = types.DefaultParams

	// variable aliases
	ModuleCdc             = types.ModuleCdc
	KeyTxCost             = types.KeyTxCost
	KeyLinkMsgCost 		  = types.KeyLinkMsgCost
	KeyNonLinkMsgCost     = types.KeyNonLinkMsgCost
	KeyRecoveryPeriod     = types.KeyRecoveryPeriod
	KeyAdjustPricePeriod  = types.KeyAdjustPricePeriod
	KeyBaseCreditPrice    = types.KeyBaseCreditPrice
	KeyDesirableBandwidth = types.KeyDesirableBandwidth
	KeyMaxBlockBandwidth  = types.KeyMaxBlockBandwidth
)

type (
	AccountBandwidthKeeper    = keeper.BaseAccountBandwidthKeeper
	BlockSpentBandwidthKeeper = keeper.BaseBlockSpentBandwidthKeeper

	Meter            = types.BandwidthMeter
	AccountBandwidth = types.AcсountBandwidth
	GenesisState     = types.GenesisState
	Params           = types.Params
)
