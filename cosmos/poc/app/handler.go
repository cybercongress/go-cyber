package app

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ContentIdLinks struct {
	ContentID  string         `json:"cid"`
	LinkedCIDS map[string]int `json:"linkedCids"`
}

// NewHandler returns a handler for "link" type messages.
func NewHandler(keyLink *sdk.KVStoreKey) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		store := ctx.KVStore(keyLink)

		msgLink := msg.(MsgLink)

		cid1Bytes, err := json.Marshal(msgLink.ContentID1)

		if err != nil {
			return sdk.ErrInternal("Error serializing cid1").Result()
		}

		linksBytes := store.Get(cid1Bytes)
		var links ContentIdLinks

		if linksBytes == nil {
			links = ContentIdLinks{ContentID: msgLink.ContentID1, LinkedCIDS: map[string]int{}}
		} else {
			err := json.Unmarshal(linksBytes, &links)
			if err != nil {
				return sdk.ErrInternal("Error when deserializing links").Result()
			}
		}

		link := links.LinkedCIDS[msgLink.ContentID2]
		links.LinkedCIDS[msgLink.ContentID2] = link + 1

		linksBytes, err = json.Marshal(links)
		if err != nil {
			return sdk.ErrInternal("Links encoding error").Result()
		}

		store.Set(cid1Bytes, linksBytes)
		return sdk.Result{}
	}
}
