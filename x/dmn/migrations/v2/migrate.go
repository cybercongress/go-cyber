package v2

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/v5/x/dmn/exported"
	"github.com/cybercongress/go-cyber/v5/x/dmn/types"
)

func Migrate(
	ctx sdk.Context,
	store sdk.KVStore,
	legacySubspace exported.Subspace,
	cdc codec.BinaryCodec,
) error {
	var currParams types.Params
	legacySubspace.GetParamSet(ctx, &currParams)

	if err := currParams.Validate(); err != nil {
		return err
	}

	bz := cdc.MustMarshal(&currParams)
	store.Set(types.ParamsKey, bz)

	return nil
}
