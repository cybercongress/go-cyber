package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/stake"
)

type MaxAccBandwidth func(ctx sdk.Context, accStake int64) int64

func NewMaxAccBandwidth(stakeKeeper stake.Keeper, MaxNetworkBandwidth int64) MaxAccBandwidth {
	return func(ctx sdk.Context, accStake int64) int64 {
		pool := stakeKeeper.GetPool(ctx)
		totalStake := pool.BondedTokens.RoundInt64() + pool.LooseTokens.RoundInt64()

		return ((accStake / totalStake) * MaxNetworkBandwidth) / 2
	}
}

type MsgBandwidthCost func(msg sdk.Msg) int64

type BandwidthHandler func(ctx sdk.Context, price float64, tx sdk.Tx) (spent int64, err sdk.Error)
