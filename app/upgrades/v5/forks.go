package v5

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/go-cyber/v5/app/keepers"
)

func RunForkLogic(ctx sdk.Context, keepers *keepers.AppKeepers) {
	logger := ctx.Logger().With("upgrade", UpgradeName)

	logger.Info("Applying emergency hard fork for v5")

	liquidityParams := keepers.LiquidityKeeper.GetParams(ctx)
	liquidityParams.CircuitBreakerEnabled = true
	err := keepers.LiquidityKeeper.SetParams(ctx, liquidityParams)
	if err != nil {
		panic(err)
	}
	logger.Info("set liquidity circuit breaker enabled, disable dex")

	keepers.BankKeeper.SetSendEnabled(ctx, "millivolt", false)
	keepers.BankKeeper.SetSendEnabled(ctx, "milliampere", false)

	logger.Info("set bank send disabled for millivolt and amperes")
}
