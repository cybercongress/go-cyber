package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/cybercongress/cyberd/cosmos/poc/app/storage"
)

// NewHandler returns a handler for "link" type messages.
// cis  - cids index storage
// ils  - incoming links storage
// ols  - outgoing links storage
// imms - in-memory storage
func NewLinksHandler(cis CidIndexStorage, ils LinksStorage, ols LinksStorage, imms *InMemoryStorage) sdk.Handler {

	getCidIndex := func(ctx sdk.Context, cid Cid) CidNumber {

		index, exist := imms.GetCidIndex(cid)
		if !exist { // new cid
			index = cis.GetOrPutCidIndex(ctx, cid)
		}
		return index
	}

	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {

		link := msg.(MsgLink)

		linkedCids := LinkedCids{
			FromCid: getCidIndex(ctx, link.CidFrom),
			ToCid:   getCidIndex(ctx, link.CidTo),
			Creator: AccountNumber(link.Address.String()),
		}

		ils.AddLink(ctx, linkedCids)
		imms.AddLink(linkedCids)
		return sdk.Result{}
	}

}
