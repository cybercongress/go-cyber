package staking

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var _ module.AppModule = AppModule{}

type AppModule struct {
	staking.AppModule
	sk stakingkeeper.Keeper
	bk bankkeeper.Keeper
	ak authkeeper.AccountKeeper
}

func NewAppModule(
	cdc codec.Codec,
	stakingKeeper stakingkeeper.Keeper,
	accountKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
) AppModule {
	return AppModule{
		AppModule: staking.NewAppModule(cdc, stakingKeeper, accountKeeper, bankKeeper),
		sk:        stakingKeeper,
		bk:        bankKeeper,
		ak:        accountKeeper,
	}
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	stakingtypes.RegisterMsgServer(cfg.MsgServer(), NewMsgServerImpl(am.sk, am.bk))
	querier := stakingkeeper.Querier{Keeper: am.sk}
	stakingtypes.RegisterQueryServer(cfg.QueryServer(), querier)

	m := stakingkeeper.NewMigrator(am.sk)
	_ = cfg.RegisterMigration(stakingtypes.ModuleName, 1, m.Migrate1to2)
}
