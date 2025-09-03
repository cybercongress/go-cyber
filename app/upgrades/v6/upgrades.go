package v6

import (
	"fmt"
	"time"

	bandwidthtypes "github.com/cybercongress/go-cyber/v6/x/bandwidth/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
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

		before = time.Now()
		// -- burn unvested coins for accounts with periodic vesting accounts
		// -- delete unvested future vesting periods, set end time to current block time

		updatedVestingAccounts := 0
		totalUnvestedBurned := sdk.NewCoins()
		for _, acc := range keepers.AccountKeeper.GetAllAccounts(ctx) {
			switch v := acc.(type) {
			case *vestingtypes.PeriodicVestingAccount:

				// compute unvested (still vesting) coins at upgrade time before any conversion
				unvested := v.GetVestingCoins(ctx.BlockTime())

				// Trim all not-yet-finished periods so remaining schedule contains only fully completed periods.
				// Also set EndTime to current block time and align OriginalVesting with kept periods.
				if unvested.IsAllPositive() {
					elapsed := ctx.BlockTime().Unix() - v.StartTime
					if elapsed < 0 {
						elapsed = 0
					}
					cumLength := int64(0)
					keptPeriods := vestingtypes.Periods{}
					keptOriginal := sdk.NewCoins()
					for _, p := range v.VestingPeriods {
						cumLength += p.Length
						if cumLength <= elapsed {
							keptPeriods = append(keptPeriods, p)
							keptOriginal = keptOriginal.Add(p.Amount...)
						} else {
							break
						}
					}
					v.VestingPeriods = keptPeriods
					v.OriginalVesting = keptOriginal
					v.EndTime = v.StartTime + cumLength // or ctx.BlockTime().Unix()
					keepers.AccountKeeper.SetAccount(ctx, v)

					updatedVestingAccounts++

					// After unlocking, burn the unvested resources, limited by available balances.
					addr := acc.GetAddress()
					balances := keepers.BankKeeper.GetAllBalances(ctx, addr)
					coinsToBurn := sdk.NewCoins()
					for _, c := range unvested {
						if c.Denom == "millivolt" || c.Denom == "milliampere" {
							balAmt := balances.AmountOf(c.Denom)
							if balAmt.IsPositive() {
								amt := sdk.MinInt(c.Amount, balAmt)
								if amt.IsPositive() {
									coinsToBurn = coinsToBurn.Add(sdk.NewCoin(c.Denom, amt))
								}
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
								totalUnvestedBurned = totalUnvestedBurned.Add(coinsToBurn...)
							}
						}
					}
				}
			}
			logger.Info("vesting cleanup completed", "accounts updated", updatedVestingAccounts, "total unvested burned", totalUnvestedBurned.String())

		}

		// manually burn minted coins for accounts which made investmints during HFR break

		explosionBurned := sdk.NewCoins()
		processedAccounts := 0
		for _, e := range hfrEntries {
			addr, err := sdk.AccAddressFromBech32(e.addr)
			if err != nil {
				logger.Error("invalid address", "addr", e.addr, "err", err)
				continue
			}
			amt, ok := sdk.NewIntFromString(e.amount)
			if !ok || !amt.IsPositive() {
				continue
			}
			balanceAmt := keepers.BankKeeper.GetBalance(ctx, addr, e.denom).Amount
			toBurnAmt := sdk.MinInt(balanceAmt, amt)
			if !toBurnAmt.IsPositive() {
				continue
			}
			coin := sdk.NewCoin(e.denom, toBurnAmt)
			if err := keepers.BankKeeper.SendCoinsFromAccountToModule(ctx, addr, resourcestypes.ResourcesName, sdk.NewCoins(coin)); err != nil {
				logger.Error("failed to move coins for burning", "addr", addr.String(), "coin", coin.String(), "err", err)
				continue
			}
			if err := keepers.BankKeeper.BurnCoins(ctx, resourcestypes.ResourcesName, sdk.NewCoins(coin)); err != nil {
				logger.Error("failed to burn coins", "addr", addr.String(), "coin", coin.String(), "err", err)
				continue
			}

			explosionBurned = explosionBurned.Add(coin)
			processedAccounts++
		}

		logger.Info("burn completed", "accounts processed", processedAccounts, "total explosion burned", explosionBurned.String())
		after = time.Now()
		ctx.Logger().Info("balances fixed", "duration ms", after.Sub(before).Milliseconds())

		before = time.Now()
		for _, acc := range keepers.AccountKeeper.GetAllAccounts(ctx) {
			keepers.BandwidthMeter.SetZeroAccountBandwidth(ctx, acc.GetAddress())
		}
		params := keepers.BandwidthMeter.GetParams(ctx)
		err = keepers.BandwidthMeter.SetParams(ctx, bandwidthtypes.Params{
			BasePrice:         sdk.OneDec(),
			RecoveryPeriod:    params.RecoveryPeriod,
			AdjustPricePeriod: params.AdjustPricePeriod,
			BaseLoad:          sdk.NewDecWithPrec(2, 2),
			MaxBlockBandwidth: params.MaxBlockBandwidth,
		})
		if err != nil {
			return nil, err
		}
		after = time.Now()

		millivoltSupply := keepers.BankKeeper.GetSupply(ctx, "millivolt")
		keepers.BandwidthMeter.SetDesirableBandwidth(ctx, millivoltSupply.Amount.Uint64())

		ctx.Logger().Info("set zero bandwidth for all accounts", "duration ms", after.Sub(before).Milliseconds())

		return versionMap, err
	}
}
