package keeper

import (
	"context"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	ctypes "github.com/cybercongress/go-cyber/types"
	bandwidthkeeper "github.com/cybercongress/go-cyber/x/bandwidth/keeper"
	bandwidthtypes "github.com/cybercongress/go-cyber/x/bandwidth/types"
	cyberbankkeeper "github.com/cybercongress/go-cyber/x/cyberbank/keeper"

	// sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cybercongress/go-cyber/x/graph/types"
)

type msgServer struct {
	*GraphKeeper
	*IndexKeeper
	authkeeper.AccountKeeper
	*cyberbankkeeper.IndexedKeeper
	*bandwidthkeeper.BandwidthMeter
}

// NewMsgServerImpl returns an implementation of the stored MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(
	gk *GraphKeeper,
	ik *IndexKeeper,
	ak authkeeper.AccountKeeper,
	bk *cyberbankkeeper.IndexedKeeper,
	bm *bandwidthkeeper.BandwidthMeter,
) types.MsgServer {
	return &msgServer{
		gk,
		ik,
		ak,
		bk,
		bm,
	}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) Cyberlink(goCtx context.Context, msg *types.MsgCyberlink) (*types.MsgCyberlinkResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var accNumber ctypes.AccNumber
	addr, err := sdk.AccAddressFromBech32(msg.Neuron)
	if err != nil {
		return nil, err
	}
	acc := k.GetAccount(ctx, addr)
	if acc != nil {
		accNumber = ctypes.AccNumber(acc.GetAccountNumber())
	} else {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Invalid neuron address")
	}

	// TODO move to ante and contract case below
	if ampers, ok := k.GetTotalStakesAmpere()[uint64(accNumber)]; ok {
		if ampers == 0 {
			return nil, types.ErrZeroPower
		}
	} else {
		return nil, types.ErrZeroPower
	}

	// case when programs and autonomous programs cyberlink
	if k.GetAccount(ctx, addr).GetPubKey() == nil {
		cost := uint64(k.GetCurrentCreditPrice().MulInt64(int64(len(msg.Links) * 1000)).TruncateInt64())
		accountBandwidth := k.GetCurrentAccountBandwidth(ctx, addr)

		currentBlockSpentBandwidth := k.GetCurrentBlockSpentBandwidth(ctx)
		maxBlockBandwidth := k.GetMaxBlockBandwidth(ctx)

		if !accountBandwidth.HasEnoughRemained(cost) {
			return nil, bandwidthtypes.ErrNotEnoughBandwidth
		} else if (cost + currentBlockSpentBandwidth) > maxBlockBandwidth {
			return nil, bandwidthtypes.ErrExceededMaxBlockBandwidth
		} else {
			err = k.ConsumeAccountBandwidth(ctx, accountBandwidth, cost)
			if err != nil {
				return nil, bandwidthtypes.ErrNotEnoughBandwidth
			}
			k.AddToBlockBandwidth(cost)
		}
	}

	for _, link := range msg.Links {
		// if cid not exists it automatically means that this is new link
		fromCidNumber, exists := k.GetCidNumber(ctx, types.Cid(link.From))
		if !exists {
			continue
		}
		toCidNumber, exists := k.GetCidNumber(ctx, types.Cid(link.To))
		if !exists {
			continue
		}

		compactLink := types.NewLink(fromCidNumber, toCidNumber, accNumber)

		if k.IndexKeeper.IsLinkExist(compactLink) {
			return nil, types.ErrCyberlinkExist
		}

		if k.IndexKeeper.IsLinkExistInCache(ctx, compactLink) {
			return nil, types.ErrCyberlinkExist
		}
	}

	for _, link := range msg.Links {
		fromCidNumber := k.GetOrPutCidNumber(ctx, types.Cid(link.From))
		toCidNumber := k.GetOrPutCidNumber(ctx, types.Cid(link.To))

		k.GraphKeeper.SaveLink(ctx, types.NewLink(fromCidNumber, toCidNumber, accNumber))
		k.GraphKeeper.IncrementNeudeg(ctx, uint64(accNumber))
		k.IndexKeeper.PutLink(ctx, types.NewLink(fromCidNumber, toCidNumber, accNumber))

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCyberlink,
				sdk.NewAttribute(types.AttributeKeyParticleFrom, link.From),
				sdk.NewAttribute(types.AttributeKeyParticleTo, link.To),
			),
		)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCyberlink,
			sdk.NewAttribute(types.AttributeKeyNeuron, msg.Neuron),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Neuron),
		),
	})

	return &types.MsgCyberlinkResponse{}, nil
}
