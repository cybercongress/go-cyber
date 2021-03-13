
package staking

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	//"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

var (
	_ module.AppModule = AppModule{}
)

type AppModule struct {
	staking.AppModule
	authkeeper.AccountKeeper
	stakingkeeper.Keeper
}

func NewAppModule(
	cdc codec.Marshaler,
	stakingKeeper stakingkeeper.Keeper,
	accKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
) AppModule {
	return AppModule{
		AppModule:  staking.NewAppModule(cdc, stakingKeeper, accKeeper, bankKeeper),
		AccountKeeper: accKeeper,
		Keeper: stakingKeeper,
	}
}

func NewHandler(ak authkeeper.AccountKeeper, sk stakingkeeper.Keeper) sdk.Handler {
	//stakingHandler := staking.NewHandler(sk)
	return WrapStakingHandler(ak, sk)
	//return wrappedHandler
}

func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(stakingtypes.RouterKey, NewHandler(am.AccountKeeper, am.Keeper))
}