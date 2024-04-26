package staking

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"
	"github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

var _ module.AppModule = AppModule{}

type AppModule struct {
	staking.AppModule
	sk *keeper.Keeper
	bk bankkeeper.Keeper
	ak authkeeper.AccountKeeper

	// legacySubspace is used solely for migration of x/params managed parameters
	legacySubspace exported.Subspace
}

func NewAppModule(
	cdc codec.Codec,
	stakingKeeper *keeper.Keeper,
	accountKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
	ls exported.Subspace,
) AppModule {
	return AppModule{
		AppModule:      staking.NewAppModule(cdc, stakingKeeper, accountKeeper, bankKeeper, ls),
		sk:             stakingKeeper,
		bk:             bankKeeper,
		ak:             accountKeeper,
		legacySubspace: ls,
	}
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), NewMsgServerImpl(am.sk, am.bk))
	querier := keeper.Querier{Keeper: am.sk}
	types.RegisterQueryServer(cfg.QueryServer(), querier)

	m := keeper.NewMigrator(am.sk, am.legacySubspace)
	if err := cfg.RegisterMigration(types.ModuleName, 1, m.Migrate1to2); err != nil {
		panic(fmt.Sprintf("failed to migrate x/%s from version 1 to 2: %v", types.ModuleName, err))
	}
	if err := cfg.RegisterMigration(types.ModuleName, 2, m.Migrate2to3); err != nil {
		panic(fmt.Sprintf("failed to migrate x/%s from version 2 to 3: %v", types.ModuleName, err))
	}
	if err := cfg.RegisterMigration(types.ModuleName, 3, m.Migrate3to4); err != nil {
		panic(fmt.Sprintf("failed to migrate x/%s from version 3 to 4: %v", types.ModuleName, err))
	}
}
