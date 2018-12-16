package bandwidth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cybercongress/cyberd/x/bandwidth/types"
)

func NewBandwidthHandler(
	accKeeper auth.AccountKeeper, bwKeeper AccountBandwidthKeeper, msgCost types.MsgBandwidthCost,
) types.BandwidthHandler {

	return func(ctx sdk.Context, price float64, tx sdk.Tx) (int64, sdk.Error) {

		account, sdkErr := getAccount(ctx, accKeeper, tx.(auth.StdTx))
		if sdkErr != nil {
			return 0, sdkErr
		}

		// We should call this function cause total stake could be changed since last update
		// and currently we can't intercept all AccountKeeper interactions.
		// This method calls bandwidth.`Recover()` under the hood, so everything should work fine.
		accountBandwidth := bwKeeper.GetCurrentAccBandwidth(ctx, account)

		bandwidthForTx := TxCost
		for _, msg := range tx.GetMsgs() {
			bandwidthForTx = bandwidthForTx + msgCost(msg)
		}

		if !accountBandwidth.HasEnoughRemained(int64(float64(bandwidthForTx) * price)) {
			return 0, sdk.ErrInternal("Not enough bandwidth to make transaction! ")
		}

		accountBandwidth.Consume(int64(float64(bandwidthForTx) * price))
		bwKeeper.SetAccBandwidth(ctx, accountBandwidth)

		return bandwidthForTx, nil
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
