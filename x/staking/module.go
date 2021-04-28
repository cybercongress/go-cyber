
package staking

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var (
	_ module.AppModule = AppModule{}
)

type AppModule struct {
	staking.AppModule
	sk stakingkeeper.Keeper
	bk bankkeeper.Keeper
	ak authkeeper.AccountKeeper
}

func NewAppModule(
	cdc codec.Marshaler,
	stakingKeeper stakingkeeper.Keeper,
	accountKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
) AppModule {
	return AppModule{
		AppModule:     staking.NewAppModule(cdc, stakingKeeper, accountKeeper, bankKeeper),
		sk: stakingKeeper,
		bk: bankKeeper,
		ak: accountKeeper,
	}
}

func NewHandler(
	sk stakingkeeper.Keeper,
	bk bankkeeper.Keeper,
) sdk.Handler {
	return WrapStakingHandler(sk, bk)
}

func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(stakingtypes.RouterKey, NewHandler(am.sk, am.bk))
}