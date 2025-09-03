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
		// -- reset all vesting accounts to base accounts

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

		type burnEntry struct {
			addr   string
			denom  string
			amount string
		}

		entries := []burnEntry{
			{addr: "bostrom1fvh29xlzzatu390p3lpe7uhdsl7dpffcnlxdsq", denom: "milliampere", amount: "26932298723"},
			{addr: "bostrom1fvh29xlzzatu390p3lpe7uhdsl7dpffcnlxdsq", denom: "millivolt", amount: "2602715242"},
			{addr: "bostrom1vt38jz5g2hrszegvs929y6zp59tgkwlsgjkprq", denom: "milliampere", amount: "2882337639"},
			{addr: "bostrom1k09afxhedzgrkn60rrnpa0u66yv3vz7yxj7ja9", denom: "milliampere", amount: "519384871"},
			{addr: "bostrom1srklpacfmq889tedn9gpxmlxulcxp3vc6dfxdd", denom: "milliampere", amount: "279939249"},
			{addr: "bostrom1nmr4flrrzrka3lanfwxklunsapr92fqq5rfrd3", denom: "milliampere", amount: "1400231678"},
			{addr: "bostrom1qkw8dtdansmvayk2jxfkpel9q4d20rvfhwcf4k", denom: "millivolt", amount: "801420796"},
			{addr: "bostrom1qkw8dtdansmvayk2jxfkpel9q4d20rvfhwcf4k", denom: "milliampere", amount: "3701530803"},
			{addr: "bostrom1fj8d32l98legw3l23ag5gl6ulc64akfrpp8jzw", denom: "milliampere", amount: "1819889604"},
			{addr: "bostrom1latgea9pkd8qspxag6am6568h6sgcsywjxskr9", denom: "millivolt", amount: "2262900956"},
			{addr: "bostrom1latgea9pkd8qspxag6am6568h6sgcsywjxskr9", denom: "milliampere", amount: "26039233057"},
			{addr: "bostrom1f88q956xu3k6lvt8d60u3kj8gkk45nmqzt4h7y", denom: "millivolt", amount: "263762923"},
			{addr: "bostrom1qqays2hpyw2zy3qhvvvhp8w953hqc8la9w4u7n", denom: "milliampere", amount: "351433460"},
			{addr: "bostrom1qe0mlgcumms5qsvervcpsvwgltc8vl4r234lc5", denom: "millivolt", amount: "919947132573"},
			{addr: "bostrom1qkw8dtdansmvayk2jxfkpel9q4d20rvfhwcf4k", denom: "millivolt", amount: "810823925"},
			{addr: "bostrom1qkw8dtdansmvayk2jxfkpel9q4d20rvfhwcf4k", denom: "milliampere", amount: "3607499517"},
			{addr: "bostrom1ymprf45c44rp9k0g2r84w2tjhsq7kalv98rgpt", denom: "milliampere", amount: "23309976454"},
			{addr: "bostrom1ymprf45c44rp9k0g2r84w2tjhsq7kalv98rgpt", denom: "millivolt", amount: "1945783863"},
			{addr: "bostrom125uwvkh0je7ug5dsv0pgrquef89k33v0rywexu", denom: "milliampere", amount: "7803967398"},
			{addr: "bostrom1k09afxhedzgrkn60rrnpa0u66yv3vz7yxj7ja9", denom: "milliampere", amount: "551211514"},
			{addr: "bostrom1k09afxhedzgrkn60rrnpa0u66yv3vz7yxj7ja9", denom: "milliampere", amount: "2726353018"},
			{addr: "bostrom1k09afxhedzgrkn60rrnpa0u66yv3vz7yxj7ja9", denom: "millivolt", amount: "508210244"},
			{addr: "bostrom1v2rqwg7zvj5rrllgp6wu5h3q8ertfetsngznvy", denom: "milliampere", amount: "10523379733"},
			{addr: "bostrom1k09afxhedzgrkn60rrnpa0u66yv3vz7yxj7ja9", denom: "milliampere", amount: "3980209389"},
			{addr: "bostrom1zswrrnduv7w3vqqxll2vr2u7u67677egl46k8z", denom: "millivolt", amount: "343855512"},
			{addr: "bostrom1qqays2hpyw2zy3qhvvvhp8w953hqc8la9w4u7n", denom: "millivolt", amount: "5341580419"},
			{addr: "bostrom1f5warat4vc0q98k7ygys4saka8u04rfxpmthvl", denom: "milliampere", amount: "59911163882"},
			{addr: "bostrom1f5warat4vc0q98k7ygys4saka8u04rfxpmthvl", denom: "millivolt", amount: "5121347544"},
			{addr: "bostrom1vt38jz5g2hrszegvs929y6zp59tgkwlsgjkprq", denom: "milliampere", amount: "413761980"},
			{addr: "bostrom1v2rqwg7zvj5rrllgp6wu5h3q8ertfetsngznvy", denom: "milliampere", amount: "351922703"},
			{addr: "bostrom1k09afxhedzgrkn60rrnpa0u66yv3vz7yxj7ja9", denom: "milliampere", amount: "4553479199"},
			{addr: "bostrom1jpkh0h7npgp3cq3cy2xt25tus7k7kfey7s8cl0", denom: "milliampere", amount: "275137941"},
			{addr: "bostrom1qkw8dtdansmvayk2jxfkpel9q4d20rvfhwcf4k", denom: "millivolt", amount: "788685037"},
			{addr: "bostrom1qkw8dtdansmvayk2jxfkpel9q4d20rvfhwcf4k", denom: "milliampere", amount: "3828888401"},
			{addr: "bostrom1v2rqwg7zvj5rrllgp6wu5h3q8ertfetsngznvy", denom: "milliampere", amount: "455964649"},
			{addr: "bostrom1k09afxhedzgrkn60rrnpa0u66yv3vz7yxj7ja9", denom: "millivolt", amount: "712020894"},
			{addr: "bostrom1k09afxhedzgrkn60rrnpa0u66yv3vz7yxj7ja9", denom: "millivolt", amount: "393178603"},
			{addr: "bostrom1k09afxhedzgrkn60rrnpa0u66yv3vz7yxj7ja9", denom: "milliampere", amount: "9613298490"},
			{addr: "bostrom1v2rqwg7zvj5rrllgp6wu5h3q8ertfetsngznvy", denom: "milliampere", amount: "494072354"},
			{addr: "bostrom1vt38jz5g2hrszegvs929y6zp59tgkwlsgjkprq", denom: "millivolt", amount: "288237732"},
			{addr: "bostrom1vt38jz5g2hrszegvs929y6zp59tgkwlsgjkprq", denom: "millivolt", amount: "303654781"},
			{addr: "bostrom1v2rqwg7zvj5rrllgp6wu5h3q8ertfetsngznvy", denom: "millivolt", amount: "260954257"},
			{addr: "bostrom1k09afxhedzgrkn60rrnpa0u66yv3vz7yxj7ja9", denom: "millivolt", amount: "1286096483"},
			{addr: "bostrom15dme26c0jv6gcxzed4yc4q8d7rlypev0t0vt5c", denom: "millivolt", amount: "261126476"},
			{addr: "bostrom1k09afxhedzgrkn60rrnpa0u66yv3vz7yxj7ja9", denom: "millivolt", amount: "302091183"},
			{addr: "bostrom1pxgpxr2xty0r77sfgah9zc3vtatrw56y329n2l", denom: "millivolt", amount: "327958125"},
			{addr: "bostrom1vt38jz5g2hrszegvs929y6zp59tgkwlsgjkprq", denom: "milliampere", amount: "413761980"},
			{addr: "bostrom15dme26c0jv6gcxzed4yc4q8d7rlypev0t0vt5c", denom: "milliampere", amount: "315286008"},
			{addr: "bostrom1k09afxhedzgrkn60rrnpa0u66yv3vz7yxj7ja9", denom: "millivolt", amount: "614199696"},
			{addr: "bostrom1vt38jz5g2hrszegvs929y6zp59tgkwlsgjkprq", denom: "millivolt", amount: "303654781"},
			{addr: "bostrom1vt38jz5g2hrszegvs929y6zp59tgkwlsgjkprq", denom: "milliampere", amount: "2188331236"},
			{addr: "bostrom1v2rqwg7zvj5rrllgp6wu5h3q8ertfetsngznvy", denom: "milliampere", amount: "265056074"},
			{addr: "bostrom1k09afxhedzgrkn60rrnpa0u66yv3vz7yxj7ja9", denom: "milliampere", amount: "3076363293"},
			{addr: "bostrom15dme26c0jv6gcxzed4yc4q8d7rlypev0t0vt5c", denom: "milliampere", amount: "2601393556"},
			{addr: "bostrom15dme26c0jv6gcxzed4yc4q8d7rlypev0t0vt5c", denom: "milliampere", amount: "1509398189"},
			{addr: "bostrom1vt38jz5g2hrszegvs929y6zp59tgkwlsgjkprq", denom: "millivolt", amount: "537330942"},
			{addr: "bostrom1k09afxhedzgrkn60rrnpa0u66yv3vz7yxj7ja9", denom: "millivolt", amount: "469031783"},
			{addr: "bostrom1v2rqwg7zvj5rrllgp6wu5h3q8ertfetsngznvy", denom: "millivolt", amount: "353620879"},
			{addr: "bostrom1vt38jz5g2hrszegvs929y6zp59tgkwlsgjkprq", denom: "millivolt", amount: "390355260"},
			{addr: "bostrom18zrtny9ghuyy9qhpvgtd0jjped9eujzpcc3lzk", denom: "millivolt", amount: "388804500"},
			{addr: "bostrom1v2rqwg7zvj5rrllgp6wu5h3q8ertfetsngznvy", denom: "milliampere", amount: "354330221"},
			{addr: "bostrom15dme26c0jv6gcxzed4yc4q8d7rlypev0t0vt5c", denom: "milliampere", amount: "14594313658"},
			{addr: "bostrom1qe0mlgcumms5qsvervcpsvwgltc8vl4r234lc5", denom: "milliampere", amount: "152062237261"},
			{addr: "bostrom1k09afxhedzgrkn60rrnpa0u66yv3vz7yxj7ja9", denom: "milliampere", amount: "12738135126"},
			{addr: "bostrom1k09afxhedzgrkn60rrnpa0u66yv3vz7yxj7ja9", denom: "millivolt", amount: "2231730860"},
			{addr: "bostrom1vt38jz5g2hrszegvs929y6zp59tgkwlsgjkprq", denom: "milliampere", amount: "5409149881"},
			{addr: "bostrom1vt38jz5g2hrszegvs929y6zp59tgkwlsgjkprq", denom: "millivolt", amount: "545310484"},
			{addr: "bostrom1vt38jz5g2hrszegvs929y6zp59tgkwlsgjkprq", denom: "milliampere", amount: "288735209"},
			{addr: "bostrom15dme26c0jv6gcxzed4yc4q8d7rlypev0t0vt5c", denom: "milliampere", amount: "771149592"},
			{addr: "bostrom14qtapaem0rcj2h9gp5r7q32mq9sh5cs3wlvsn7", denom: "millivolt", amount: "2429029471"},
			{addr: "bostrom1qe0mlgcumms5qsvervcpsvwgltc8vl4r234lc5", denom: "millivolt", amount: "818108814658"},
			{addr: "bostrom1qe0mlgcumms5qsvervcpsvwgltc8vl4r234lc5", denom: "millivolt", amount: "7128428104"},
			{addr: "bostrom1fty3ql39e7m4gy98z0jv68vdst0wksy932vxku", denom: "millivolt", amount: "2085965578"},
			{addr: "bostrom1ngymhfhty76w8r6e8dlyvj68g5zfm09qjmejts", denom: "millivolt", amount: "2370208604"},
			{addr: "bostrom1ngymhfhty76w8r6e8dlyvj68g5zfm09qjmejts", denom: "milliampere", amount: "24649814272"},
			{addr: "bostrom1a3nm66sw7ejp454azqem2d88g5m7lc6dwf9tsa", denom: "millivolt", amount: "1511899342"},
			{addr: "bostrom1a3nm66sw7ejp454azqem2d88g5m7lc6dwf9tsa", denom: "millivolt", amount: "726616669"},
			{addr: "bostrom1a3nm66sw7ejp454azqem2d88g5m7lc6dwf9tsa", denom: "millivolt", amount: "596926589"},
			{addr: "bostrom1a3nm66sw7ejp454azqem2d88g5m7lc6dwf9tsa", denom: "millivolt", amount: "524364059"},
			{addr: "bostrom1a3nm66sw7ejp454azqem2d88g5m7lc6dwf9tsa", denom: "milliampere", amount: "3387538353"},
			{addr: "bostrom1a3nm66sw7ejp454azqem2d88g5m7lc6dwf9tsa", denom: "milliampere", amount: "2539415149"},
			{addr: "bostrom1a3nm66sw7ejp454azqem2d88g5m7lc6dwf9tsa", denom: "milliampere", amount: "2877122095"},
			{addr: "bostrom1vt38jz5g2hrszegvs929y6zp59tgkwlsgjkprq", denom: "millivolt", amount: "6823201994"},
			{addr: "bostrom15dme26c0jv6gcxzed4yc4q8d7rlypev0t0vt5c", denom: "milliampere", amount: "18783348485"},
			{addr: "bostrom1gvy4vk6v2rzpd9my487l358uajec8n0vn5c9us", denom: "milliampere", amount: "55387606103"},
			{addr: "bostrom1gvy4vk6v2rzpd9my487l358uajec8n0vn5c9us", denom: "millivolt", amount: "8690051483"},
			{addr: "bostrom10r6r8crpu0c8t4cumlpcp6uzc8ngut3j09flpy", denom: "milliampere", amount: "43311697895"},
			{addr: "bostrom10r6r8crpu0c8t4cumlpcp6uzc8ngut3j09flpy", denom: "millivolt", amount: "7033800128"},
			{addr: "bostrom10r6r8crpu0c8t4cumlpcp6uzc8ngut3j09flpy", denom: "milliampere", amount: "8821722561"},
			{addr: "bostrom1tkz3uvmgglns04c524quqx26aunj3h0k00q68f", denom: "milliampere", amount: "267406320"},
			{addr: "bostrom1k09afxhedzgrkn60rrnpa0u66yv3vz7yxj7ja9", denom: "millivolt", amount: "2774293087"},
			{addr: "bostrom1vt38jz5g2hrszegvs929y6zp59tgkwlsgjkprq", denom: "milliampere", amount: "5400859561"},
			{addr: "bostrom1c5g8xw3dfcxypw65j7aa4j29hetr4wy48acpcx", denom: "milliampere", amount: "12662635612"},
			{addr: "bostrom15dme26c0jv6gcxzed4yc4q8d7rlypev0t0vt5c", denom: "milliampere", amount: "18797239514"},
			{addr: "bostrom15dme26c0jv6gcxzed4yc4q8d7rlypev0t0vt5c", denom: "milliampere", amount: "11070841292"},
			{addr: "bostrom1ggda2n35v7zkwucufzp25mp34p0adp0x0w7wgc", denom: "millivolt", amount: "432000000"},
			{addr: "bostrom1ggda2n35v7zkwucufzp25mp34p0adp0x0w7wgc", denom: "milliampere", amount: "4320000000"},
			{addr: "bostrom1k09afxhedzgrkn60rrnpa0u66yv3vz7yxj7ja9", denom: "millivolt", amount: "1397957106"},
			{addr: "bostrom1ay267fakkrgfy9lf2m7wsj8uez2dgylhtkdf9k", denom: "milliampere", amount: "11162598095"},
			{addr: "bostrom1qqays2hpyw2zy3qhvvvhp8w953hqc8la9w4u7n", denom: "milliampere", amount: "18749017000"},
		}

		explosionBurned := sdk.NewCoins()
		processedAccounts := 0
		for _, e := range entries {
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
