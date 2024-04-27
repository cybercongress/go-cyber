package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	pluginstypes "github.com/cybercongress/go-cyber/v4/plugins/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/v4/x/dmn/keeper"
)

type Messenger struct {
	keeper *keeper.Keeper
}

func NewMessenger(
	keeper *keeper.Keeper,
) *Messenger {
	return &Messenger{
		keeper: keeper,
	}
}

func (m *Messenger) HandleMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg pluginstypes.CyberMsg) ([]sdk.Event, [][]byte, error) {
	switch {
	case msg.CreateThought != nil:
		if msg.CreateThought.Program != contractAddr.String() {
			return nil, nil, wasmvmtypes.InvalidRequest{Err: "create thought wrong program"}
		}

		msgServer := keeper.NewMsgServerImpl(*m.keeper)

		if err := msg.CreateThought.ValidateBasic(); err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed validating msg")
		}

		res, err := msgServer.CreateThought(
			sdk.WrapSDKContext(ctx),
			msg.CreateThought,
		)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "create thought msg")
		}

		responseBytes, err := json.Marshal(*res)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed to serialize create thought response")
		}

		resp := [][]byte{responseBytes}

		return nil, resp, nil
	case msg.ForgetThought != nil:
		if msg.ForgetThought.Program != contractAddr.String() {
			return nil, nil, wasmvmtypes.InvalidRequest{Err: "forget thought wrong program"}
		}

		msgServer := keeper.NewMsgServerImpl(*m.keeper)

		if err := msg.ForgetThought.ValidateBasic(); err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed validating msg")
		}

		res, err := msgServer.ForgetThought(
			sdk.WrapSDKContext(ctx),
			msg.ForgetThought,
		)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "forget thought msg")
		}

		responseBytes, err := json.Marshal(*res)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed to serialize forget thought response")
		}

		resp := [][]byte{responseBytes}

		return nil, resp, nil
	case msg.ChangeThoughtInput != nil:
		if msg.ChangeThoughtInput.Program != contractAddr.String() {
			return nil, nil, wasmvmtypes.InvalidRequest{Err: "change thought wrong program"}
		}

		msgServer := keeper.NewMsgServerImpl(*m.keeper)

		if err := msg.ChangeThoughtInput.ValidateBasic(); err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed validating msg")
		}

		res, err := msgServer.ChangeThoughtInput(
			sdk.WrapSDKContext(ctx),
			msg.ChangeThoughtInput,
		)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "change thought input msg")
		}

		responseBytes, err := json.Marshal(*res)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed to serialize change thought input response")
		}

		resp := [][]byte{responseBytes}

		return nil, resp, nil
	case msg.ChangeThoughtPeriod != nil:
		if msg.ChangeThoughtPeriod.Program != contractAddr.String() {
			return nil, nil, wasmvmtypes.InvalidRequest{Err: "change thought wrong program"}
		}

		msgServer := keeper.NewMsgServerImpl(*m.keeper)

		if err := msg.ChangeThoughtPeriod.ValidateBasic(); err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed validating msg")
		}

		res, err := msgServer.ChangeThoughtPeriod(
			sdk.WrapSDKContext(ctx),
			msg.ChangeThoughtPeriod,
		)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "change thought period msg")
		}

		responseBytes, err := json.Marshal(*res)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed to serialize change thought period response")
		}

		resp := [][]byte{responseBytes}

		return nil, resp, nil

	case msg.ChangeThoughtBlock != nil:
		if msg.ChangeThoughtBlock.Program != contractAddr.String() {
			return nil, nil, wasmvmtypes.InvalidRequest{Err: "change thought wrong program"}
		}

		msgServer := keeper.NewMsgServerImpl(*m.keeper)

		if err := msg.ChangeThoughtBlock.ValidateBasic(); err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed validating msg")
		}

		res, err := msgServer.ChangeThoughtBlock(
			sdk.WrapSDKContext(ctx),
			msg.ChangeThoughtBlock,
		)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "change thought block msg")
		}

		responseBytes, err := json.Marshal(*res)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed to serialize change thought block response")
		}

		resp := [][]byte{responseBytes}

		return nil, resp, nil

	case msg.ChangeThoughtGasPrice != nil:
		if msg.ChangeThoughtGasPrice.Program != contractAddr.String() {
			return nil, nil, wasmvmtypes.InvalidRequest{Err: "change thought wrong program"}
		}

		msgServer := keeper.NewMsgServerImpl(*m.keeper)

		if err := msg.ChangeThoughtGasPrice.ValidateBasic(); err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed validating msg")
		}

		res, err := msgServer.ChangeThoughtGasPrice(
			sdk.WrapSDKContext(ctx),
			msg.ChangeThoughtGasPrice,
		)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "change thought gas price msg")
		}

		responseBytes, err := json.Marshal(*res)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed to serialize change thought gas price response")
		}

		resp := [][]byte{responseBytes}

		return nil, resp, nil

	case msg.ChangeThoughtParticle != nil:
		if msg.ChangeThoughtParticle.Program != contractAddr.String() {
			return nil, nil, wasmvmtypes.InvalidRequest{Err: "change thought wrong program"}
		}

		msgServer := keeper.NewMsgServerImpl(*m.keeper)

		if err := msg.ChangeThoughtParticle.ValidateBasic(); err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed validating msg")
		}

		res, err := msgServer.ChangeThoughtParticle(
			sdk.WrapSDKContext(ctx),
			msg.ChangeThoughtParticle,
		)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "change thought particle msg")
		}

		responseBytes, err := json.Marshal(*res)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed to serialize change thought particle response")
		}

		resp := [][]byte{responseBytes}

		return nil, resp, nil

	case msg.ChangeThoughtName != nil:
		if msg.ChangeThoughtName.Program != contractAddr.String() {
			return nil, nil, wasmvmtypes.InvalidRequest{Err: "change thought wrong program"}
		}

		msgServer := keeper.NewMsgServerImpl(*m.keeper)

		if err := msg.ChangeThoughtName.ValidateBasic(); err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed validating msg")
		}

		res, err := msgServer.ChangeThoughtName(
			sdk.WrapSDKContext(ctx),
			msg.ChangeThoughtName,
		)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "change thought name msg")
		}

		responseBytes, err := json.Marshal(*res)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed to serialize change thought name response")
		}

		resp := [][]byte{responseBytes}

		return nil, resp, nil
	default:
		return nil, nil, pluginstypes.ErrHandleMsg
	}
}

type Querier struct {
	*keeper.Keeper
}

func NewWasmQuerier(keeper *keeper.Keeper) *Querier {
	return &Querier{keeper}
}

func (querier *Querier) HandleQuery(ctx sdk.Context, query pluginstypes.CyberQuery) ([]byte, error) {
	switch {
	case query.Thought != nil:
		res, err := querier.Keeper.Thought(ctx, query.Thought)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to get dmn thought")
		}

		responseBytes, err := json.Marshal(res)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to serialize dmn thought response")
		}
		return responseBytes, nil
	case query.ThoughtStats != nil:
		res, err := querier.Keeper.ThoughtStats(ctx, query.ThoughtStats)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to get dmn thought stats")
		}

		responseBytes, err := json.Marshal(res)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to serialize dmn thought stats response")
		}
		return responseBytes, nil
	case query.ThoughtsFees != nil:
		res, err := querier.Keeper.ThoughtsFees(ctx, query.ThoughtsFees)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to get dmn thoughts fees")
		}

		responseBytes, err := json.Marshal(res)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to serialize dmn thoughts fees response")
		}
		return responseBytes, nil
	default:
		return nil, pluginstypes.ErrHandleQuery
	}
}
