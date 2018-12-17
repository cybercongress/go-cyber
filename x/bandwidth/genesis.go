package bandwidth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/app/genesis"
	"github.com/cybercongress/cyberd/x/bandwidth/types"
)

// Genesis accounts should contains fully restored bandwidth on block 0
func InitGenesis(ctx sdk.Context, bwHandler types.Handler, bwKeeper types.Keeper, accs []genesis.GenesisAccount) {

	for _, acc := range accs {
		accMaxBw := bwHandler.GetAccMaxBandwidth(ctx, acc.Address)
		bwKeeper.SetAccBandwidth(ctx, types.NewGenesisAccBandwidth(acc.Address, accMaxBw))
	}
}
