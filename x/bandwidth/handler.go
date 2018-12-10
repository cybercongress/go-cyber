package bandwidth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cybercongress/cyberd/app/coin"
	"github.com/cybercongress/cyberd/x/bandwidth/types"
)

func NewBandwidthHandler(
	accKeeper auth.AccountKeeper, bwKeeper AccountBandwidthKeeper,
	msgCost types.MsgBandwidthCost, maxBandwidth types.MaxAccBandwidth,
) types.BandwidthHandler {

	return func(ctx sdk.Context, price float64, tx sdk.Tx) sdk.Error {

		account, sdkErr := getAccount(ctx, accKeeper, tx.(auth.StdTx))
		if sdkErr != nil {
			return sdkErr
		}

		accountBandwidth, err := bwKeeper.GetAccountBandwidth(account, ctx)
		if err != nil {
			return sdk.ErrInternal("Cannot process tx")
		}

		addressStake := accKeeper.GetAccount(ctx, account).GetCoins().AmountOf(coin.CBD)

		// We should call this function instead of Recover() cause total stake could be changed since last update
		// and currently we can't intercept all AccountKeeper interactions.
		// This method calls Recover() under the hood, so everything should work fine.
		accountBandwidth.UpdateMax(maxBandwidth(ctx, addressStake.Int64()), RecoveryPeriod, ctx.BlockHeight())

		bandwidthForTx := TxCost
		for _, msg := range tx.GetMsgs() {
			bandwidthForTx = bandwidthForTx + msgCost(msg)
		}

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
