package client

import (
	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/client"
	"github.com/cybercongress/cyberd/claim/context"
	"github.com/cybercongress/cyberd/util"
	"github.com/pkg/errors"
	"net"
	"net/http"
)

const (
	token    = "CBD"
	maxValue = "9000000000"
)

func ClaimHandlerFn(ctx context.ClaimContext) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		address, err := util.GetSingleParamValue("address", r)
		if err != nil {
			util.HandleError(err, w)
			return
		}

		amount, err := util.GetSingleParamValue("amount", r)
		if err != nil {
			amount = maxValue
		}

		claimTo, err := types.AccAddressFromBech32(address)
		if err != nil {
			util.HandleError(err, w)
			return
		}

		// Claim if address doesn't exists, error otherwise
		if err = ctx.CliContext.EnsureAccountExistsFromAddr(claimTo); err != nil {

			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				util.HandleError(err, w)
				return
			}

			err = ctx.IncrementIp(ip)
			if err != nil {
				util.HandleError(err, w)
				return
			}

			coins, _ := sdk.ParseCoins(amount + token)
			msg := client.CreateMsg(ctx.ClaimFrom, claimTo, coins)

			ctx.Mtx.Lock()
			defer ctx.Mtx.Unlock()

			txBldr, err := ctx.TxBuilder()
			if err != nil {
				util.HandleError(err, w)
				return
			}

			txBytes, err := txBldr.BuildAndSign(ctx.Name, ctx.Passphrase, []sdk.Msg{msg})
			if err != nil {
				util.HandleError(err, w)
				return
			}

			result, err := ctx.CliContext.BroadcastTxSync(txBytes)
			if err != nil {
				util.HandleError(err, w)
				return
			}

			resultJson, err := ctx.Codec.MarshalJSON(result)
			if err != nil {
				util.HandleError(err, w)
				return
			}
			*ctx.Sequence++
			w.Write(resultJson)
		} else {
			util.HandleError(errors.New("Account already has tokens"), w)
			return
		}
	}
}
