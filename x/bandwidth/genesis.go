package bandwidth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/cybercongress/go-cyber/x/bandwidth/types"
)

// Genesis accounts should contains fully restored bandwidth on block 0
func InitGenesis(ctx sdk.Context, bm *BandwidthMeter, ak auth.AccountKeeper, data GenesisState) {
	bm.SetParams(ctx, data.Params)
	for _, address := range ak.GetAllAccounts(ctx) {
		accMaxBw := bm.GetAccountMaxBandwidth(ctx, address.GetAddress())
		bm.SetAccountBandwidth(ctx, types.NewGenesisAccountBandwidth(address.GetAddress(), accMaxBw))
	}
}

func ExportGenesis(ctx sdk.Context, bm *BandwidthMeter) GenesisState {
	return NewGenesisState(bm.GetParams(ctx))
}
