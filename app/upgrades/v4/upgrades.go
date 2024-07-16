package v3

import (
	"fmt"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	generaltypes "github.com/cybercongress/go-cyber/v4/types"
	tokenfactorytypes "github.com/cybercongress/go-cyber/v4/x/tokenfactory/types"
	"time"

	liquiditytypes "github.com/cybercongress/go-cyber/v4/x/liquidity/types"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"

	bandwidthtypes "github.com/cybercongress/go-cyber/v4/x/bandwidth/types"
	dmntypes "github.com/cybercongress/go-cyber/v4/x/dmn/types"
	gridtypes "github.com/cybercongress/go-cyber/v4/x/grid/types"
	ranktypes "github.com/cybercongress/go-cyber/v4/x/rank/types"
	resourcestypes "github.com/cybercongress/go-cyber/v4/x/resources/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/cybercongress/go-cyber/v4/app/keepers"
)

const NewDenomCreationGasConsume uint64 = 2_000_000

func CreateV4UpgradeHandler(
	mm *module.Manager,
	cfg module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		before := time.Now()

		logger := ctx.Logger().With("upgrade", UpgradeName)

		for _, subspace := range keepers.ParamsKeeper.GetSubspaces() {
			subspace := subspace

			var keyTable paramstypes.KeyTable
			switch subspace.Name() {
			case authtypes.ModuleName:
				keyTable = authtypes.ParamKeyTable() //nolint:staticcheck
			case banktypes.ModuleName:
				keyTable = banktypes.ParamKeyTable() //nolint:staticcheck
			case stakingtypes.ModuleName:
				keyTable = stakingtypes.ParamKeyTable() //nolint:staticcheck

			case minttypes.ModuleName:
				keyTable = minttypes.ParamKeyTable() //nolint:staticcheck
			case distrtypes.ModuleName:
				keyTable = distrtypes.ParamKeyTable() //nolint:staticcheck
			case slashingtypes.ModuleName:
				keyTable = slashingtypes.ParamKeyTable() //nolint:staticcheck
			case govtypes.ModuleName:
				keyTable = govv1.ParamKeyTable() //nolint:staticcheck
			case crisistypes.ModuleName:
				keyTable = crisistypes.ParamKeyTable() //nolint:staticcheck

			// ibc types
			case ibctransfertypes.ModuleName:
				keyTable = ibctransfertypes.ParamKeyTable()

			// wasm
			case wasmtypes.ModuleName:
				keyTable = wasmtypes.ParamKeyTable() //nolint:staticcheck

			// cyber modules
			case bandwidthtypes.ModuleName:
				keyTable = bandwidthtypes.ParamKeyTable() //nolint:staticcheck
			case dmntypes.ModuleName:
				keyTable = dmntypes.ParamKeyTable() //nolint:staticcheck
			case gridtypes.ModuleName:
				keyTable = gridtypes.ParamKeyTable() //nolint:staticcheck
			case ranktypes.ModuleName:
				keyTable = ranktypes.ParamKeyTable() //nolint:staticcheck
			case resourcestypes.ModuleName:
				keyTable = resourcestypes.ParamKeyTable() //nolint:staticcheck
			case liquiditytypes.ModuleName:
				keyTable = liquiditytypes.ParamKeyTable()
			}
			if !subspace.HasKeyTable() {
				subspace.WithKeyTable(keyTable)
			}
		}

		// Migrate Tendermint consensus parameters from x/params module to a deprecated x/consensus module.
		// The old params module is required to still be imported in your app.go in order to handle this migration.
		baseAppLegacySS := keepers.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable())
		baseapp.MigrateParams(ctx, baseAppLegacySS, &keepers.ConsensusParamsKeeper)

		// Run migrations
		logger.Info(fmt.Sprintf("pre migrate version map: %v", vm))
		versionMap, err := mm.RunMigrations(ctx, cfg, vm)
		if err != nil {
			return nil, err
		}
		logger.Info(fmt.Sprintf("post migrate version map: %v", versionMap))

		// TODO check ibc-go state after migration
		// https://github.com/cosmos/ibc-go/blob/v7.1.0/docs/migrations/v7-to-v7_1.md
		// explicitly update the IBC 02-client params, adding the localhost client type
		params := keepers.IBCKeeper.ClientKeeper.GetParams(ctx)
		params.AllowedClients = append(params.AllowedClients, exported.Localhost)
		keepers.IBCKeeper.ClientKeeper.SetParams(ctx, params)

		logger.Info("set ibc params")

		newTokenFactoryParams := tokenfactorytypes.Params{
			DenomCreationFee:        sdk.NewCoins(sdk.NewCoin(generaltypes.CYB, sdk.NewInt(10*generaltypes.Giga))),
			DenomCreationGasConsume: NewDenomCreationGasConsume,
		}
		if err := keepers.TokenFactoryKeeper.SetParams(ctx, newTokenFactoryParams); err != nil {
			return nil, err
		}
		logger.Info("set tokenfactory params")

		after := time.Now()

		ctx.Logger().Info("migration time", "duration_ms", after.Sub(before).Milliseconds())

		return versionMap, err
	}
}
