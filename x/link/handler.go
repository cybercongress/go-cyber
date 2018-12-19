package link

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	cbd "github.com/cybercongress/cyberd/types"
	"github.com/cybercongress/cyberd/x/link/keeper"
	"github.com/cybercongress/cyberd/x/link/types"
	cbdlink "github.com/cybercongress/cyberd/x/link/types"
)

// NewHandler returns a handler for "link" type messages.
// cis  - cids index storage
// ils  - links storage
// as   - account storage
// imms - in-memory storage
func NewLinksHandler(cis keeper.CidNumberKeeper, ls keeper.LinkIndexedKeeper, as auth.AccountKeeper) sdk.Handler {

	getCidNumber := GetCidNumberFunc(cis)

	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {

		linkMsg := msg.(Msg)
		fromCidNumber := getCidNumber(ctx, linkMsg.From)
		toCidNumber := getCidNumber(ctx, linkMsg.To)
		accNumber := cbd.AccNumber(as.GetAccount(ctx, linkMsg.Address).GetAccountNumber())
		link := cbdlink.NewLink(fromCidNumber, toCidNumber, accNumber)

		if ls.IsLinkExist(ctx, link) {
			return sdk.Result{Code: cbd.CodeLinkAlreadyExist, Codespace: cbd.CodespaceCbd}
		}

		if !ctx.IsCheckTx() {
			ls.PutIntoIndex(link)
		}
		ls.PutLink(ctx, link)

		return sdk.Result{Code: cbd.CodeOK, Codespace: cbd.CodespaceCbd}
	}
}

func GetCidNumberFunc(cis keeper.CidNumberKeeper) func(sdk.Context, cbdlink.Cid) cbdlink.CidNumber {

	return func(ctx sdk.Context, cid types.Cid) types.CidNumber {

		index, exist := cis.GetCidNumber(ctx, cid)
		if !exist { // new cid
			index = cis.GetOrPutCidNumber(ctx, cid)
		}
		return index
	}
}
