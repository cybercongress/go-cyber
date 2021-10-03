package graph

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/x/graph/keeper"
	"github.com/cybercongress/go-cyber/x/graph/types"
)

func EndBlocker(ctx sdk.Context, gk *keeper.GraphKeeper, ik *keeper.IndexKeeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	gk.UpdateNeudegs(ctx)
	ik.MergeContextLinks(ctx)
}