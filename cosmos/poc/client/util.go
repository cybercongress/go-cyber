package client

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/cosmos/poc/app"
)

// build the sendTx msg
func BuildMsg(from sdk.AccAddress, cid1 string, cid2 string) sdk.Msg {
	return app.NewMsgLink(from, cid1, cid2)
}
