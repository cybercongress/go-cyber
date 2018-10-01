package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ContentIdLinks struct {
	ContentID  string         `json:"cid"`
	LinkedCIDS map[string]int `json:"linkedCids"`
}

// NewHandler returns a handler for "link" type messages.
// cis - cids index storage
// ils - incoming links storage
// ols - outgoing links storage
func NewLinksHandler(cis CidIndexStorage, ils LinksStorage, ols LinksStorage) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {

		link := msg.(MsgLink)

		cis.GetOrPutCidIndex(ctx, link.ContentID1)
		cis.GetOrPutCidIndex(ctx, link.ContentID2)

		ils.AddLink(ctx, link.Address, link.ContentID2, link.ContentID1)
		ols.AddLink(ctx, link.Address, link.ContentID1, link.ContentID2)
		return sdk.Result{}
	}
}
