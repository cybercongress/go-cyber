package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"

	ctypes "github.com/cybercongress/go-cyber/types"
	cyberbankkeeper "github.com/cybercongress/go-cyber/x/cyberbank/keeper"

	//sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cybercongress/go-cyber/x/graph/types"
)

type msgServer struct {
	GraphKeeper
	*IndexKeeper
	authkeeper.AccountKeeper
	*cyberbankkeeper.IndexedKeeper
}

// NewMsgServerImpl returns an implementation of the stored MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(
	gk GraphKeeper,
	ik *IndexKeeper,
	ak authkeeper.AccountKeeper,
	bk *cyberbankkeeper.IndexedKeeper,
) types.MsgServer {
	return &msgServer{
		gk, ik, ak, bk,
	}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) Cyberlink(goCtx context.Context, msg *types.MsgCyberlink) (*types.MsgCyberlinkResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var accNumber ctypes.AccNumber
	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, err
	}
	acc := k.GetAccount(ctx, addr); if (acc != nil) {
		accNumber = ctypes.AccNumber(acc.GetAccountNumber())
	} else {
		return nil, types.ErrInvalidAccount
	}

	if ampers, ok := k.GetTotalStakesAmper()[uint64(accNumber)]; ok {
		if ampers == 0 {
			return nil, types.ErrZeroPower
		}
	} else {
		return nil, types.ErrZeroPower
	}

	for _, link := range msg.Links {
		// if cid not exists it automatically means that this is new link
		fromCidNumber, exists := k.GetCidNumber(ctx, types.Cid(link.From)); if !exists { continue }
		toCidNumber, exists := k.GetCidNumber(ctx, types.Cid(link.To)); if !exists { continue }

		compactLink := types.NewLink(fromCidNumber, toCidNumber, accNumber)

		if k.IndexKeeper.IsLinkExist(compactLink) {
			return nil, types.ErrCyberlinkExist
		}
	}

	for _, link := range msg.Links {
		fromCidNumber := k.GetOrPutCidNumber(ctx, types.Cid(link.From))
		toCidNumber := k.GetOrPutCidNumber(ctx, types.Cid(link.To))

		k.GraphKeeper.SaveLink(ctx, types.NewLink(fromCidNumber, toCidNumber, accNumber))
		k.IndexKeeper.PutLink(ctx, types.NewLink(fromCidNumber, toCidNumber, accNumber))

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCyberlink,
				sdk.NewAttribute(types.AttributeKeyObjectFrom, string(link.From)),
				sdk.NewAttribute(types.AttributeKeyObjectTo, string(link.To)),
			),
		)
	}

	return &types.MsgCyberlinkResponse{}, nil
}