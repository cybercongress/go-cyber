package bandwidth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/go-cyber/x/link"
)

func MsgBandwidthCosts(ctx sdk.Context, params Params, msg sdk.Msg) int64 {
	switch msg.(type) {
	case link.Msg:
		linkMsg := msg.(link.Msg)
		return int64(len(linkMsg.Links)) * params.LinkMsgCost
	default:
		return params.NonLinkMsgCost
	}
}
