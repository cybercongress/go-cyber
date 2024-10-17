package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/v5/x/resources/exported"
	v2 "github.com/cybercongress/go-cyber/v5/x/resources/migrations/v2"
)

// Migrator is a struct for handling in-place state migrations.
type Migrator struct {
	keeper         Keeper
	legacySubspace exported.Subspace
}

func NewMigrator(k Keeper, ss exported.Subspace) Migrator {
	return Migrator{
		keeper:         k,
		legacySubspace: ss,
	}
}

func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return v2.Migrate(ctx, ctx.KVStore(m.keeper.storeKey), m.legacySubspace, m.keeper.cdc)
}
