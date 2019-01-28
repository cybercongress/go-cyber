package mint

import sdk "github.com/cosmos/cosmos-sdk/types"

type GenesisState struct {
	Params Params `json:"params"`
}

func InitGenesis(ctx sdk.Context, minter Minter, data GenesisState) {
	minter.SetParams(ctx, data.Params)
}

func ExportGenesis(ctx sdk.Context, minter Minter) GenesisState {
	params := minter.GetParams(ctx)
	return GenesisState{
		Params: params,
	}
}

func ValidateGenesis(data GenesisState) error {
	err := validateParams(data.Params)
	if err != nil {
		return err
	}
	return nil
}
