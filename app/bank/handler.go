package bank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	. "github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cybercongress/cyberd/app/coin"
	"github.com/cybercongress/cyberd/app/storage"
	cbd "github.com/cybercongress/cyberd/app/types"
	"github.com/cybercongress/cyberd/x/bandwidth"
	bw "github.com/cybercongress/cyberd/x/bandwidth/types"
	"reflect"
)

// NewHandler returns a handler for "bank" type messages.
func NewBankHandler(
	k Keeper, imms *storage.InMemoryStorage, am auth.AccountKeeper,
	maxBandwidth bw.MaxAccBandwidth, abk bandwidth.AccountBandwidthKeeper,
) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSend:
			return handleMsgSend(ctx, k, msg, imms, am, maxBandwidth, abk)
		default:
			errMsg := "Unrecognized bank Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// NOTE: totalIn == totalOut should already have been checked
func handleMsgSend(
	ctx sdk.Context, k Keeper, msg MsgSend, imms *storage.InMemoryStorage, am auth.AccountKeeper,
	maxBandwidth bw.MaxAccBandwidth, abk bandwidth.AccountBandwidthKeeper,
) sdk.Result {

	tags, err := k.InputOutputCoins(ctx, msg.Inputs, msg.Outputs)

	if err != nil {
		return err.Result()
	}

	if !ctx.IsCheckTx() {
		for _, input := range msg.Inputs {
			acc := am.GetAccount(ctx, input.Address)
			imms.UpdateStake(cbd.AccountNumber(acc.GetAccountNumber()), -input.Coins.AmountOf(coin.CBD).Int64())

			updateBandwidth(ctx, acc.GetAddress(), maxBandwidth(ctx, acc.GetCoins().AmountOf(coin.CBD).Int64()), abk)
		}

		for _, output := range msg.Outputs {
			acc := am.GetAccount(ctx, output.Address)
			imms.UpdateStake(cbd.AccountNumber(acc.GetAccountNumber()), output.Coins.AmountOf(coin.CBD).Int64())

			updateBandwidth(ctx, acc.GetAddress(), maxBandwidth(ctx, acc.GetCoins().AmountOf(coin.CBD).Int64()), abk)
		}
	}

	return sdk.Result{Tags: tags}
}

func updateBandwidth(ctx sdk.Context, addr sdk.AccAddress, maxBandwidth int64, abk bandwidth.AccountBandwidthKeeper) {
	accBandwidth, _ := abk.GetAccountBandwidth(addr, ctx)
	accBandwidth.UpdateMax(maxBandwidth, ctx.BlockHeight(), bandwidth.RecoveryPeriod)
	abk.SetAccountBandwidth(ctx, accBandwidth)
}
