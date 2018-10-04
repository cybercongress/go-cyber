package util

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/cosmos/poc/app"
	"github.com/cybercongress/cyberd/cosmos/poc/app/storage"
)

// build the sendTx msg
func BuildMsg(address sdk.AccAddress, fromCid storage.Cid, toCid storage.Cid) sdk.Msg {
	return app.NewMsgLink(address, fromCid, toCid)
}
