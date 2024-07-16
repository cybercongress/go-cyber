package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/v4/x/liquidity/exported"
	v3 "github.com/cybercongress/go-cyber/v4/x/liquidity/migrations/v3"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper         Keeper
	legacySubspace exported.Subspace
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper, ss exported.Subspace) Migrator {
	return Migrator{
		keeper:         keeper,
		legacySubspace: ss,
	}
}

// Migrate2to3 migrates from version 2 to 3.
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	return v3.Migrate(ctx, ctx.KVStore(m.keeper.storeKey), m.legacySubspace, m.keeper.cdc)
}
