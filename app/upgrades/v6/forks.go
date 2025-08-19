package v6

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/go-cyber/v6/app/keepers"
)

func RunForkLogic(ctx sdk.Context, keepers *keepers.AppKeepers) {
	logger := ctx.Logger().With("upgrade", UpgradeName)

	logger.Info("Applying upgrade v6")

	liquidityParams := keepers.LiquidityKeeper.GetParams(ctx)
	liquidityParams.CircuitBreakerEnabled = false
	err := keepers.LiquidityKeeper.SetParams(ctx, liquidityParams)
	if err != nil {
		panic(err)
	}
	logger.Info("set liquidity circuit breaker disabled, enable dex")

	keepers.BankKeeper.SetSendEnabled(ctx, "millivolt", true)
	keepers.BankKeeper.SetSendEnabled(ctx, "milliampere", true)

	logger.Info("set bank send enabled for millivolt and amperes")
}
