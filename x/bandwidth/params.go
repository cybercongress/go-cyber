package bandwidth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cybercongress/cyberd/x/bandwidth/types"
)

const (
	// DefaultParamspace default name for parameter store
	DefaultParamspace = types.ModuleName
)

// ParamTable for bandwidth module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func (bk *BaseAccBandwidthKeeper) GetParams(ctx sdk.Context) (params Params) {
	bk.paramSpace.GetParamSet(ctx, &params)
	return params
}

// set the params
func (bk *BaseAccBandwidthKeeper) SetParams(ctx sdk.Context, params Params) {
	bk.paramSpace.SetParamSet(ctx, &params)
}
