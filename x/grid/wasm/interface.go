package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	pluginstypes "github.com/cybercongress/go-cyber/v5/plugins/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/v5/x/grid/keeper"
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
	case msg.CreateEnergyRoute != nil:
		if msg.CreateEnergyRoute.Source != contractAddr.String() {
			return nil, nil, wasmvmtypes.InvalidRequest{Err: "create route wrong source"}
		}

		msgServer := keeper.NewMsgServerImpl(*m.keeper)

		if err := msg.CreateEnergyRoute.ValidateBasic(); err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed validating msg")
		}

		res, err := msgServer.CreateRoute(
			sdk.WrapSDKContext(ctx),
			msg.CreateEnergyRoute,
		)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "create energy route msg")
		}

		responseBytes, err := json.Marshal(*res)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed to serialize create energy route response")
		}

		resp := [][]byte{responseBytes}

		return nil, resp, nil
	case msg.EditEnergyRoute != nil:
		if msg.EditEnergyRoute.Source != contractAddr.String() {
			return nil, nil, wasmvmtypes.InvalidRequest{Err: "create route wrong source"}
		}

		msgServer := keeper.NewMsgServerImpl(*m.keeper)

		if err := msg.EditEnergyRoute.ValidateBasic(); err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed validating msg")
		}

		res, err := msgServer.EditRoute(
			sdk.WrapSDKContext(ctx),
			msg.EditEnergyRoute,
		)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "edit energy route msg")
		}

		responseBytes, err := json.Marshal(*res)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed to serialize edit energy route response")
		}

		resp := [][]byte{responseBytes}

		return nil, resp, nil
	case msg.EditEnergyRouteName != nil:
		if msg.EditEnergyRouteName.Source != contractAddr.String() {
			return nil, nil, wasmvmtypes.InvalidRequest{Err: "edit energy route wrong source"}
		}

		msgServer := keeper.NewMsgServerImpl(*m.keeper)

		if err := msg.EditEnergyRouteName.ValidateBasic(); err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed validating msg")
		}

		res, err := msgServer.EditRouteName(
			sdk.WrapSDKContext(ctx),
			msg.EditEnergyRouteName,
		)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "edit energy route name msg")
		}

		responseBytes, err := json.Marshal(*res)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed to serialize edit energy route name response")
		}

		resp := [][]byte{responseBytes}

		return nil, resp, nil
	case msg.DeleteEnergyRoute != nil:
		if msg.DeleteEnergyRoute.Source != contractAddr.String() {
			return nil, nil, wasmvmtypes.InvalidRequest{Err: "delete route wrong source"}
		}

		msgServer := keeper.NewMsgServerImpl(*m.keeper)

		if err := msg.DeleteEnergyRoute.ValidateBasic(); err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed validating msg")
		}

		res, err := msgServer.DeleteRoute(
			sdk.WrapSDKContext(ctx),
			msg.DeleteEnergyRoute,
		)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "delete energy route name msg")
		}

		responseBytes, err := json.Marshal(*res)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed to serialize delete energy route name response")
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
	case query.SourceRoutes != nil:
		res, err := querier.Keeper.SourceRoutes(ctx, query.SourceRoutes)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to get source routes")
		}

		responseBytes, err := json.Marshal(res)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to serialize source routes response")
		}
		return responseBytes, nil

	case query.SourceRoutedEnergy != nil:
		res, err := querier.Keeper.SourceRoutedEnergy(ctx, query.SourceRoutedEnergy)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to get source routed energy")
		}

		responseBytes, err := json.Marshal(res)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to serialize source routed energy response")
		}
		return responseBytes, nil

	case query.DestinationRoutedEnergy != nil:
		res, err := querier.Keeper.DestinationRoutedEnergy(ctx, query.DestinationRoutedEnergy)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to get destination routed energy")
		}

		responseBytes, err := json.Marshal(res)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to serialize destination routed energy response")
		}
		return responseBytes, nil

	case query.Route != nil:
		res, err := querier.Keeper.Route(ctx, query.Route)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to get route")
		}

		responseBytes, err := json.Marshal(res)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to serialize route response")
		}
		return responseBytes, nil

	default:
		return nil, pluginstypes.ErrHandleQuery
	}
}
