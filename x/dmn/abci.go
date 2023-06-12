package dmn

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/go-cyber/x/dmn/keeper"
	"github.com/cybercongress/go-cyber/x/dmn/types"
)

func BeginBlock(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	k.ExecuteThoughtsQueue(ctx)

	// TODO add block event
}
