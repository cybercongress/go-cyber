package bandwidth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/x/link"
)

func MsgBandwidthCost(msg sdk.Msg) int64 {
	switch msg.(type) {
	case link.Msg:
		return LinkMsgCost
	default:
		return NonLinkMsgCost
	}
}
