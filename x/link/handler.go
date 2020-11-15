package link

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	сtypes "github.com/cybercongress/go-cyber/types"
	"github.com/cybercongress/go-cyber/x/link/types"

	"github.com/cosmos/cosmos-sdk/x/auth"
)

func NewHandler(
	gk GraphKeeper,
	ik *IndexKeeper,
	as auth.AccountKeeper,
) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case MsgCyberlink:
			return HandleMsgCyberlink(ctx, gk, ik, as, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized message type: %T", msg)
		}
	}
}

func HandleMsgCyberlink(
	ctx sdk.Context,
	gk GraphKeeper,
	ik *IndexKeeper,
	as auth.AccountKeeper,
	linkMsg MsgCyberlink,
) (*sdk.Result, error) {
	ctx = ctx.WithEventManager(sdk.NewEventManager())

	var accNumber сtypes.AccNumber
	acc := as.GetAccount(ctx, linkMsg.Address); if (acc != nil) {
		accNumber = сtypes.AccNumber(acc.GetAccountNumber())
	} else { // case when cyberlinks transferred over ikp and given account not exist
		accNumber = сtypes.AccNumber(as.NewAccountWithAddress(ctx, linkMsg.Address).GetAccountNumber())
	}

	for _, link := range linkMsg.Links {
		// if cid not exists it automatically means that this is new link
		fromCidNumber, exists := gk.GetCidNumber(ctx, link.From); if !exists { continue	}
		toCidNumber, exists := gk.GetCidNumber(ctx, link.To); if !exists { continue	}

		compactLink := NewLink(fromCidNumber, toCidNumber, accNumber)

		if ik.IsLinkExist(compactLink) {
			return nil, types.ErrCyberlinkExist
		}
	}

	for _, link := range linkMsg.Links {
		fromCidNumber := gk.GetOrPutCidNumber(ctx, link.From)
		toCidNumber := gk.GetOrPutCidNumber(ctx, link.To)

		gk.PutLink(ctx, NewLink(fromCidNumber, toCidNumber, accNumber))
		ik.PutLink(ctx, NewLink(fromCidNumber, toCidNumber, accNumber))

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCyberlink,
				sdk.NewAttribute(types.AttributeKeyObjectFrom, string(link.From)),
				sdk.NewAttribute(types.AttributeKeyObjectTo, string(link.To)),
			),
		)
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}


//func NewLinksHandler(
//	cis CidNumberKeeper,
//	ls IndexedKeeper,
//	as auth.AccountKeeper,
//) sdk.Handler {
//
//	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
//		ctx = ctx.WithEventManager(sdk.NewEventManager())
//
//		linkMsg := msg.(cyberlink.MsgCyberlink)
//		accNumber := cbd.AccNumber(as.GetAccount(ctx, linkMsg.Address).GetAccountNumber())
//
//		//todo: optimize validations
//		for _, link := range linkMsg.Links {
//			// if cid not exists it automatically means that this is new link
//			fromCidNumber, exists := cis.GetCidNumber(ctx, link.From)
//			if !exists {
//				continue
//			}
//			toCidNumber, exists := cis.GetCidNumber(ctx, link.To)
//			if !exists {
//				continue
//			}
//
//			compactLink := cyberlink.NewLink(fromCidNumber, toCidNumber, accNumber)
//
//			if ls.IsLinkExist(compactLink) {
//				return nil, ErrCyberlinkExist
//			}
//		}
//
//		for _, link := range linkMsg.Links {
//			fromCidNumber := cis.GetOrPutCidNumber(ctx, link.From)
//			toCidNumber := cis.GetOrPutCidNumber(ctx, link.To)
//
//			ls.PutLink(ctx, cyberlink.NewLink(fromCidNumber, toCidNumber, accNumber))
//
//			ctx.EventManager().EmitEvent(
//				sdk.NewEvent(
//					cyberlink.EventTypeCyberlink,
//					sdk.NewAttribute(cyberlink.AttributeKeyObjectFrom, string(link.From)),
//					sdk.NewAttribute(cyberlink.AttributeKeyObjectTo, string(link.To)),
//				),
//			)
//		}
//
//		return &sdk.Result{Events: ctx.EventManager().Events()}, nil
//	}
//}
