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
func NewLinksHandler(cis CidIndexStorage, ls LinksStorage, imms *InMemoryStorage) sdk.Handler {

	getCidNumber := GetCidNumberFunc(cis, imms)

	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {

		link := msg.(MsgLink)

		linkedCids := LinkedCids{
			FromCid: getCidNumber(ctx, link.CidFrom),
			ToCid:   getCidNumber(ctx, link.CidTo),
			Creator: AccountNumber(link.Address.String()),
		}

		ls.AddLink(ctx, linkedCids)
		imms.AddLink(linkedCids)
		return sdk.Result{}
	}

}

func GetCidNumberFunc(cis CidIndexStorage, imms *InMemoryStorage) func(sdk.Context, Cid) CidNumber {

	return func(ctx sdk.Context, cid Cid) CidNumber {

		index, exist := imms.GetCidIndex(cid)
		if !exist { // new cid
			index = cis.GetOrPutCidIndex(ctx, cid)
			imms.AddCid(cid, index)
		}
		return index
	}
}
