package bank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkbank "github.com/cosmos/cosmos-sdk/x/bank"
)

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data sdkbank.GenesisState) {
	keeper.SetSendEnabled(ctx, data.SendEnabled)
}
