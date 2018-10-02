package bank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cybercongress/cyberd/cosmos/poc/app/coin"
	"github.com/cybercongress/cyberd/cosmos/poc/app/storage"
	"reflect"
)

// NewHandler returns a handler for "bank" type messages.
func NewBankHandler(k Keeper, imms *storage.InMemoryStorage) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSend:
			return handleMsgSend(ctx, k, msg, imms)
		default:
			errMsg := "Unrecognized bank Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// NOTE: totalIn == totalOut should already have been checked
func handleMsgSend(ctx sdk.Context, k Keeper, msg MsgSend, imms *storage.InMemoryStorage) sdk.Result {

	tags, err := k.InputOutputCoins(ctx, msg.Inputs, msg.Outputs)
	if err != nil {
		return err.Result()
	}

	for _, input := range msg.Inputs {
		imms.UpdateStake(input.Address, -input.Coins.AmountOf(coin.CBD).Int64())
	}

	for _, output := range msg.Outputs {
		imms.UpdateStake(output.Address, output.Coins.AmountOf(coin.CBD).Int64())
	}

	return sdk.Result{Tags: tags}
}
