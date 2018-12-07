package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	. "github.com/cybercongress/cyberd/app/storage"
	cbd "github.com/cybercongress/cyberd/app/types"
)

// NewHandler returns a handler for "link" type messages.
// cis  - cids index storage
// ils  - links storage
// as   - account storage
// imms - in-memory storage
func NewLinksHandler(cis CidIndexStorage, ls LinksStorage, imms *InMemoryStorage, as auth.AccountKeeper) sdk.Handler {

	getCidNumber := GetCidNumberFunc(cis, imms)

	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {

		linkMsg := msg.(MsgLink)
		fromCidNumber := getCidNumber(ctx, linkMsg.From)
		toCidNumber := getCidNumber(ctx, linkMsg.To)
		accNumber := cbd.AccountNumber(as.GetAccount(ctx, linkMsg.Address).GetAccountNumber())
		link := cbd.NewLink(fromCidNumber, toCidNumber, accNumber)

		if ls.IsLinkExist(ctx, link) {
			return sdk.Result{Code: cbd.CodeLinkAlreadyExist, Codespace: cbd.CodespaceCbd}
		}

		if !ctx.IsCheckTx() {
			imms.AddLink(link)
		}
		ls.AddLink(ctx, link)

		return sdk.Result{Code: cbd.CodeOK, Codespace: cbd.CodespaceCbd}
	}
}

func GetCidNumberFunc(cis CidIndexStorage, imms *InMemoryStorage) func(sdk.Context, cbd.Cid) cbd.CidNumber {

	return func(ctx sdk.Context, cid cbd.Cid) cbd.CidNumber {

		index, exist := imms.GetCidIndex(cid)
		if !exist { // new cid
			index = cis.GetOrPutCidNumber(ctx, cid)
			if !ctx.IsCheckTx() {
				imms.AddCid(cid, index)
			}
		}
		return index
	}
}
