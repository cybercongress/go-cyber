package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"

	"github.com/cybercongress/go-cyber/v4/x/bandwidth/types"
)

func InitGenesis(ctx sdk.Context, bm *BandwidthMeter, ak authkeeper.AccountKeeper, data types.GenesisState) {
	if err := bm.SetParams(ctx, data.Params); err != nil {
		panic(err)
	}

	bm.currentCreditPrice = bm.GetBandwidthPrice(ctx, data.Params.BasePrice)

	for _, address := range ak.GetAllAccounts(ctx) {
		accMaxBw := bm.GetAccountMaxBandwidth(ctx, address.GetAddress())
		bm.SetAccountBandwidth(ctx, types.NewGenesisNeuronBandwidth(address.GetAddress(), accMaxBw))
	}
}

func ExportGenesis(ctx sdk.Context, bm *BandwidthMeter) *types.GenesisState {
	return types.NewGenesisState(bm.GetParams(ctx))
}
