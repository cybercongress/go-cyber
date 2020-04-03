package link

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	cbd "github.com/cybercongress/go-cyber/types"
	"github.com/cybercongress/go-cyber/x/bandwidth/exported"
	cyberlink "github.com/cybercongress/go-cyber/x/link/internal/types"

	"github.com/cosmos/cosmos-sdk/x/auth"
)

// NewHandler returns a handler for "link" type messages.
// cis  - cids index storage
// ils  - links storage
// as   - account storage
// imms - in-memory storage
func NewLinksHandler(
	cis CidNumberKeeper,
	ls IndexedKeeper,
	as auth.AccountKeeper,
	abk exported.BaseAccountBandwidthKeeper,
	meter exported.Meter,
) sdk.Handler {

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		linkMsg := msg.(cyberlink.Msg)
		accNumber := cbd.AccNumber(as.GetAccount(ctx, linkMsg.Address).GetAccountNumber())

		//todo: optimize validations
		for _, link := range linkMsg.Links {
			// if cid not exists it automatically means that this is new link
			fromCidNumber, exists := cis.GetCidNumber(ctx, link.From)
			if !exists {
				continue
			}
			toCidNumber, exists := cis.GetCidNumber(ctx, link.To)
			if !exists {
				continue
			}

			compactLink := cyberlink.NewLink(fromCidNumber, toCidNumber, accNumber)

			if ls.IsLinkExist(compactLink) {
				return nil, ErrCyberlinkExist
			}
		}

		// TODO later we will migrate to advanced karma ranking
		linkMsgCost := abk.GetParams(ctx).LinkMsgCost
		currentCreditPrice := meter.GetCurrentCreditPrice()
		linkCost := int64(float64(linkMsgCost) * currentCreditPrice)
		karma := linkCost*int64(len(linkMsg.Links))

		for _, link := range linkMsg.Links {
			fromCidNumber := cis.GetOrPutCidNumber(ctx, link.From)
			toCidNumber := cis.GetOrPutCidNumber(ctx, link.To)

			ls.PutLink(ctx, cyberlink.NewLink(fromCidNumber, toCidNumber, accNumber))

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					cyberlink.EventTypeCyberlink,
					sdk.NewAttribute(cyberlink.AttributeKeyObjectFrom, string(link.From)),
					sdk.NewAttribute(cyberlink.AttributeKeyObjectTo, string(link.To)),
				),
			)
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				cyberlink.EventTypeCyberlinkMeta,
				sdk.NewAttribute(cyberlink.AttributeKeySubject, linkMsg.Address.String()),
				sdk.NewAttribute(cyberlink.AttributeKeyKarma, strconv.FormatInt(karma, 10)),
			),
		)

		return &sdk.Result{Events: ctx.EventManager().Events()}, nil
	}
}
