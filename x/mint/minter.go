package mint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

var (
	// default paramspace for mint params
	DefaultParamspace = "cbdmint"
	// params store for inflation params
	ParamStoreKeyParams = []byte("params")
)

// mint new CYB tokens for every block
type Minter struct {
	fck         auth.FeeCollectionKeeper
	stakeKeeper *keeper.Keeper
	paramSpace  params.Subspace
}

// get inflation params from the global param store
func (m Minter) GetParams(ctx sdk.Context) (p Params) {
	m.paramSpace.Get(ctx, ParamStoreKeyParams, &p)
	return
}

// set inflation params from the global param store
func (m Minter) SetParams(ctx sdk.Context, p Params) {
	m.paramSpace.Set(ctx, ParamStoreKeyParams, &p)
}

func NewMinter(fck auth.FeeCollectionKeeper, stakeKeeper *keeper.Keeper, paramSpace params.Subspace) Minter {
	return Minter{
		fck:         fck,
		stakeKeeper: stakeKeeper,
		paramSpace:  paramSpace.WithKeyTable(ParamTypeTable()),
	}
}

func ParamTypeTable() params.KeyTable {
	return params.NewKeyTable(
		ParamStoreKeyParams, Params{},
	)
}
