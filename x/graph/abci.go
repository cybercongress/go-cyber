package graph

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/go-cyber/x/graph/keeper"
	"github.com/cybercongress/go-cyber/x/graph/types"
)

func EndBlocker(ctx sdk.Context, gk *keeper.GraphKeeper, _ *keeper.IndexKeeper) {
	amountParticles := gk.GetCidsCount(ctx)
	amountLinks := gk.GetLinksCount(ctx)
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)
	defer telemetry.ModuleSetGauge(types.ModuleName, float32(amountLinks), "total_cyberlinks")
	defer telemetry.ModuleSetGauge(types.ModuleName, float32(amountParticles), "total_particles")

	gk.UpdateMemNeudegs(ctx)
	// ik.MergeContextLinks(ctx)
}
