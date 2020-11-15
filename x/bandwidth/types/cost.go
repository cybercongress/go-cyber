package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/x/link"
)

func MsgBandwidthCosts(ctx sdk.Context, params Params, msg sdk.Msg) uint64 {
	switch msg.(type) {
	case link.MsgCyberlink:
		linkMsg := msg.(link.MsgCyberlink)
		return uint64(len(linkMsg.Links)) * params.LinkMsgCost
	default:
		return params.NonLinkMsgCost
	}
}
