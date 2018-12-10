package util

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	cbd "github.com/cybercongress/cyberd/app/types"
	"github.com/cybercongress/cyberd/x/link"
)

// build the sendTx msg
func BuildMsg(address sdk.AccAddress, fromCid cbd.Cid, toCid cbd.Cid) sdk.Msg {
	return link.NewMsgLink(address, fromCid, toCid)
}
