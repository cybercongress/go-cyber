package v6

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/cybercongress/go-cyber/v6/app/keepers"
	"time"
)

func CreateV6UpgradeHandler(
	mm *module.Manager,
	cfg module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		before := time.Now()

		logger := ctx.Logger().With("upgrade", UpgradeName)

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

		logger.Info(fmt.Sprintf("pre migrate version map: %v", vm))
		versionMap, err := mm.RunMigrations(ctx, cfg, vm)
		if err != nil {
			return nil, err
		}
		logger.Info(fmt.Sprintf("post migrate version map: %v", versionMap))

		after := time.Now()
		ctx.Logger().Info("upgrade time", "duration ms", after.Sub(before).Milliseconds())

		return versionMap, err
	}
}
