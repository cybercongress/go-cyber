package bank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	. "github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cybercongress/cyberd/app/coin"
	"github.com/cybercongress/cyberd/app/storage"
	cbd "github.com/cybercongress/cyberd/app/types"
	"reflect"
)

// NewHandler returns a handler for "bank" type messages.
func NewBankHandler(k Keeper, imms *storage.InMemoryStorage, am auth.AccountKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSend:
			return handleMsgSend(ctx, k, msg, imms, am)
		default:
			errMsg := "Unrecognized bank Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// NOTE: totalIn == totalOut should already have been checked
func handleMsgSend(
	ctx sdk.Context, k Keeper, msg MsgSend, imms *storage.InMemoryStorage, am auth.AccountKeeper,
) sdk.Result {

	tags, err := k.InputOutputCoins(ctx, msg.Inputs, msg.Outputs)

	if err != nil {
		return err.Result()
	}

	if !ctx.IsCheckTx() {
		for _, input := range msg.Inputs {
			accNumber := am.GetAccount(ctx, input.Address).GetAccountNumber()
			imms.UpdateStake(cbd.AccountNumber(accNumber), -input.Coins.AmountOf(coin.CBD).Int64())
		}

		for _, output := range msg.Outputs {
			accNumber := am.GetAccount(ctx, output.Address).GetAccountNumber()
			imms.UpdateStake(cbd.AccountNumber(accNumber), output.Coins.AmountOf(coin.CBD).Int64())
		}
	}
	return sdk.Result{Tags: tags}
}
