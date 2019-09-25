package bandwidth

import (
	"github.com/cybercongress/cyberd/x/bandwidth/types"
)

const (
	DefaultParamspace = types.DefaultParamspace
	ModuleName        = types.ModuleName
	StoreKey          = types.StoreKey
	RouterKey         = types.RouterKey
)

type (
	Keeper                    = types.Keeper
	BlockSpentBandwidthKeeper = types.BlockSpentBandwidthKeeper
	AcсBandwidth              = types.AcсBandwidth
	Params                    = types.Params
)

var (
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
	ParamKeyTable = types.ParamKeyTable
)
