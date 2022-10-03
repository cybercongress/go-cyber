package cyberbank

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/joinresistance/space-pussy/x/cyberbank/keeper"
	"github.com/joinresistance/space-pussy/x/cyberbank/types"
)

func EndBlocker(ctx sdk.Context, k *keeper.IndexedKeeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	k.UpdateAccountsStakeAmpere(ctx)
}
