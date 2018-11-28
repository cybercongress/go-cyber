package util

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/app"
	cbd "github.com/cybercongress/cyberd/app/types"
)

// build the sendTx msg
func BuildMsg(address sdk.AccAddress, fromCid cbd.Cid, toCid cbd.Cid) sdk.Msg {
	return app.NewMsgLink(address, fromCid, toCid)
}
