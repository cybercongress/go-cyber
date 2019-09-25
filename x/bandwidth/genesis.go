package bandwidth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/x/bandwidth/types"
)

type GenesisState struct{}

func NewGenesisState() GenesisState {
	return GenesisState{}
}

func ValidateGenesis(data GenesisState) error {
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func ExportGenesis() GenesisState {
	return GenesisState{}
}

// Genesis accounts should contains fully restored bandwidth on block 0
func InitGenesis(
	ctx sdk.Context, handler types.BandwidthMeter,
	keeper types.Keeper, addresses []sdk.AccAddress) {

	for _, address := range addresses {
		accMaxBw := handler.GetAccMaxBandwidth(ctx, address)
		keeper.SetAccBandwidth(ctx, types.NewGenesisAccBandwidth(address, accMaxBw))
	}
}
