package bandwidth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/cybercongress/cyberd/app/coin"
	"github.com/cybercongress/cyberd/types"
)

func NewBandwidthHandler(stakeKeeper stake.Keeper, accKeeper auth.AccountKeeper, bwKeeper AccountBandwidthKeeper) types.BandwidthHandler {

	return func(ctx sdk.Context, tx sdk.Tx) sdk.Error {

		account, sdkErr := getAccount(ctx, accKeeper, tx.(auth.StdTx))
		if sdkErr != nil {
			return sdkErr
		}

		accountBandwidth, err := bwKeeper.GetAccountBandwidth(account, ctx)
		if err != nil {
			return sdk.ErrInternal("Cannot process tx")
		}

		addressStake := accKeeper.GetAccount(ctx, account).GetCoins().AmountOf(coin.CBD)

		pool := stakeKeeper.GetPool(ctx)
		totalStake := pool.BondedTokens.RoundInt64() + pool.LooseTokens.RoundInt64()

		addressMaxBandwidth := ((addressStake.Int64() / totalStake) * MaxNetworkBandwidth) / 2

		// We should call this function instead of Recover() cause total stake could be changed since last update
		// and currently we can't intercept all AccountKeeper interactions.
		// This method calls Recover() under the hood, so everything should work fine.
		accountBandwidth.UpdateMax(addressMaxBandwidth, ctx.BlockHeight())
		bandwidthForTx := int64(1) // TODO: add calculations

		if !accountBandwidth.HasEnoughRemained(bandwidthForTx) {
			return sdk.ErrInternal("Not enough bandwidth to make transaction! ")
		}

		accountBandwidth.Consume(bandwidthForTx)
		bwKeeper.SetAccountBandwidth(ctx, accountBandwidth)

		return nil
	}
}

func getAccount(ctx sdk.Context, accKeeper auth.AccountKeeper, tx auth.StdTx) (sdk.AccAddress, sdk.Error) {

	if tx.GetMsgs() == nil || len(tx.GetMsgs()) == 0 {
		return nil, sdk.ErrInternal("Tx.GetMsgs() must return at least one message in list")
	}

	if err := tx.ValidateBasic(); err != nil {
		return nil, err
	}

	// signers acc [0] bandwidth will be consumed
	account := tx.GetSigners()[0]

	acc := accKeeper.GetAccount(ctx, account)
	if acc == nil {
		return nil, sdk.ErrUnknownAddress(account.String())
	}

	return account, nil
}
