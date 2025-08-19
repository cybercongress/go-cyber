package v6

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/cybercongress/go-cyber/v6/app/keepers"
	resourcestypes "github.com/cybercongress/go-cyber/v6/x/resources/types"
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

		updated := 0
		totalBurned := sdk.NewCoins()
		for _, acc := range keepers.AccountKeeper.GetAllAccounts(ctx) {
			switch v := acc.(type) {
			case *vestingtypes.PeriodicVestingAccount:
				// compute unvested (still vesting) coins at upgrade time before conversion
				unvested := v.GetVestingCoins(ctx.BlockTime())
				// convert to base account to make all balances spendable for burning
				bacc := authtypes.NewBaseAccount(acc.GetAddress(), acc.GetPubKey(), acc.GetAccountNumber(), acc.GetSequence())
				keepers.AccountKeeper.SetAccount(ctx, bacc)
				updated++

				if unvested.IsAllPositive() {
					addr := acc.GetAddress()
					balances := keepers.BankKeeper.GetAllBalances(ctx, addr)
					// prepare coins to burn: minimum of unvested and current balance per denom
					coinsToBurn := sdk.NewCoins()
					for _, c := range unvested {
						balAmt := balances.AmountOf(c.Denom)
						if balAmt.IsPositive() {
							amt := sdk.MinInt(c.Amount, balAmt)
							if amt.IsPositive() {
								coinsToBurn = coinsToBurn.Add(sdk.NewCoin(c.Denom, amt))
							}
						}
					}
					if coinsToBurn.IsAllPositive() {
						if err := keepers.BankKeeper.SendCoinsFromAccountToModule(ctx, addr, resourcestypes.ResourcesName, coinsToBurn); err != nil {
							logger.Error("failed to move coins for burning", "addr", addr.String(), "coins", coinsToBurn.String(), "err", err)
						} else {
							if err := keepers.BankKeeper.BurnCoins(ctx, resourcestypes.ResourcesName, coinsToBurn); err != nil {
								logger.Error("failed to burn coins", "addr", addr.String(), "coins", coinsToBurn.String(), "err", err)
							} else {
								totalBurned = totalBurned.Add(coinsToBurn...)
							}
						}
					}
			}
		}
		logger.Info("vesting cleanup completed", "accounts_updated", updated, "total_burned", totalBurned.String())

		return versionMap, err
	}
}
