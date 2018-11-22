package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	. "github.com/cybercongress/cyberd/cosmos/poc/app/storage"
	cbd "github.com/cybercongress/cyberd/cosmos/poc/app/types"
)

// NewHandler returns a handler for "link" type messages.
// cis  - cids index storage
// ils  - links storage
// as   - account storage
// imms - in-memory storage
func NewLinksHandler(cis CidIndexStorage, ls LinksStorage, imms *InMemoryStorage, as auth.AccountKeeper) sdk.Handler {

	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {

		linkMsg := msg.(MsgLink)
		fromCidNumber := cis.GetOrPutCidNumber(ctx, linkMsg.From)
		toCidNumber := cis.GetOrPutCidNumber(ctx, linkMsg.To)
		accNumber := cbd.AccountNumber(as.GetAccount(ctx, linkMsg.Address).GetAccountNumber())
		link := cbd.NewLink(fromCidNumber, toCidNumber, accNumber)

		if ls.IsLinkExist(ctx, link) {
			return sdk.Result{Code: cbd.LinkAlreadyExistsCode()}
		}
		ls.AddLink(ctx, link)

		if !ctx.IsCheckTx() {
			imms.AddLink(link)
		}
		return sdk.Result{Code: sdk.ABCICodeOK}
	}

}
