package bandwidth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/x/link"
)

func MsgBandwidthCosts(msg sdk.Msg) int64 {
	switch msg.(type) {
	case link.Msg:
		linkMsg := msg.(link.Msg)
		return int64(len(linkMsg.Links)) * LinkMsgCost
	default:
		return NonLinkMsgCost
	}
}
