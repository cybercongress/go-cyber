package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"

	pluginstypes "github.com/cybercongress/go-cyber/v4/plugins/types"
	bandwidthkeeper "github.com/cybercongress/go-cyber/v4/x/bandwidth/keeper"
	cyberbankkeeper "github.com/cybercongress/go-cyber/v4/x/cyberbank/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/v4/x/graph/keeper"
)

type Messenger struct {
	gk *keeper.GraphKeeper
	ik *keeper.IndexKeeper
	ak *authkeeper.AccountKeeper
	bk *cyberbankkeeper.IndexedKeeper
	bm *bandwidthkeeper.BandwidthMeter
}

func NewMessenger(
	gk *keeper.GraphKeeper,
	ik *keeper.IndexKeeper,
	ak *authkeeper.AccountKeeper,
	bk *cyberbankkeeper.IndexedKeeper,
	bm *bandwidthkeeper.BandwidthMeter,
) *Messenger {
	return &Messenger{
		gk,
		ik,
		ak,
		bk,
		bm,
	}
}

func (m *Messenger) HandleMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg pluginstypes.CyberMsg) ([]sdk.Event, [][]byte, error) {
	switch {
	case msg.Cyberlink != nil:
		if msg.Cyberlink.Neuron != contractAddr.String() {
			return nil, nil, wasmvmtypes.InvalidRequest{Err: "cyberlink wrong neuron"}
		}

		msgServer := keeper.NewMsgServerImpl(m.gk, m.ik, *m.ak, m.bk, m.bm)

		if err := msg.Cyberlink.ValidateBasic(); err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed validating msg")
		}

		res, err := msgServer.Cyberlink(
			sdk.WrapSDKContext(ctx),
			msg.Cyberlink,
		)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "cyberlink msg")
		}

		responseBytes, err := json.Marshal(*res)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed to serialize cyberlink response")
		}

		resp := [][]byte{responseBytes}

		return nil, resp, nil
	default:
		return nil, nil, pluginstypes.ErrHandleMsg
	}
}

type Querier struct {
	*keeper.GraphKeeper
}

func NewWasmQuerier(keeper *keeper.GraphKeeper) *Querier {
	return &Querier{keeper}
}

func (querier *Querier) HandleQuery(ctx sdk.Context, query pluginstypes.CyberQuery) ([]byte, error) {
	switch {
	case query.GraphStats != nil:
		res, err := querier.GraphKeeper.GraphStats(ctx, query.GraphStats)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to get graph stats")
		}

		responseBytes, err := json.Marshal(res)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to serialize graph stats response")
		}
		return responseBytes, nil

	default:
		return nil, pluginstypes.ErrHandleQuery
	}
}
