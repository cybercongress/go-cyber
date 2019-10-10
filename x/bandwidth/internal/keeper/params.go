package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/cybercongress/cyberd/x/bandwidth/internal/types"

	"strconv"
)

// ParamTable for bandwidth module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

func (bk *BaseAccBandwidthKeeper) GetParams(ctx sdk.Context) (params types.Params) {
	bk.paramSpace.GetParamSet(ctx, &params)
	return params
}

// set the params
func (bk *BaseAccBandwidthKeeper) SetParams(ctx sdk.Context, params types.Params) {
	params.SlidingWindowSize = params.RecoveryPeriod
	params.ShouldBeSpentPerSlidingWindow = strconv.Itoa(int(params.DesirableNetworkBandwidthForRecoveryPeriod))
	bk.paramSpace.SetParamSet(ctx, &params)
}
