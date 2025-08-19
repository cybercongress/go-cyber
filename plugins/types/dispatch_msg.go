package types

import (
	"encoding/json"
	"errors"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	dmnkeeper "github.com/cybercongress/go-cyber/v6/x/dmn/keeper"
	dmntypes "github.com/cybercongress/go-cyber/v6/x/dmn/types"
	graphkeeper "github.com/cybercongress/go-cyber/v6/x/graph/keeper"
	graphtypes "github.com/cybercongress/go-cyber/v6/x/graph/types"
	gridkeeper "github.com/cybercongress/go-cyber/v6/x/grid/keeper"
	gridtypes "github.com/cybercongress/go-cyber/v6/x/grid/types"
	resourceskeeper "github.com/cybercongress/go-cyber/v6/x/resources/keeper"
	resourcestypes "github.com/cybercongress/go-cyber/v6/x/resources/types"
	tokenfactorykeeper "github.com/cybercongress/go-cyber/v6/x/tokenfactory/keeper"
	tokenfactorytypes "github.com/cybercongress/go-cyber/v6/x/tokenfactory/types"

	errorsmod "cosmossdk.io/errors"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ wasmkeeper.Messenger = (*CustomMessenger)(nil)

var ErrHandleMsg = errors.New("error handle message")

type ModuleMessenger interface {
	HandleMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg CyberMsg) ([]sdk.Event, [][]byte, error)
}

type CustomMessenger struct {
	wrapped            wasmkeeper.Messenger
	moduleMessengers   map[string]ModuleMessenger
	graphKeeper        *graphkeeper.GraphKeeper
	dmnKeeper          *dmnkeeper.Keeper
	gridKeeper         *gridkeeper.Keeper
	resourcesKeeper    *resourceskeeper.Keeper
	bankKeeper         *bankkeeper.Keeper
	tokenFactoryKeeper *tokenfactorykeeper.Keeper
}

func CustomMessageDecorator(
	moduleMessengers map[string]ModuleMessenger,
	graph *graphkeeper.GraphKeeper,
	dmn *dmnkeeper.Keeper,
	grid *gridkeeper.Keeper,
	resources *resourceskeeper.Keeper,
	bank *bankkeeper.Keeper,
	tokenFactory *tokenfactorykeeper.Keeper,
) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessenger{
			wrapped:            old,
			moduleMessengers:   moduleMessengers,
			graphKeeper:        graph,
			dmnKeeper:          dmn,
			gridKeeper:         grid,
			resourcesKeeper:    resources,
			bankKeeper:         bank,
			tokenFactoryKeeper: tokenFactory,
		}
	}
}

func (m *CustomMessenger) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) ([]sdk.Event, [][]byte, error) {
	// Return early if msg.Custom is nil
	if msg.Custom == nil {
		return m.wrapped.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
	}

	var contractMsg CyberMsg
	if err := json.Unmarshal(msg.Custom, &contractMsg); err != nil {
		ctx.Logger().Debug("json.Unmarshal: failed to decode incoming custom cosmos message",
			"from_address", contractAddr.String(),
			"message", string(msg.Custom),
			"error", err,
		)
		return nil, nil, errorsmod.Wrap(err, "cyber msg error")
	}

	switch {
	case contractMsg.Cyberlink != nil:
		return m.moduleMessengers[graphtypes.ModuleName].HandleMsg(ctx, contractAddr, contractIBCPortID, contractMsg)
	case contractMsg.Investmint != nil:
		return m.moduleMessengers[resourcestypes.ModuleName].HandleMsg(ctx, contractAddr, contractIBCPortID, contractMsg)
	case contractMsg.CreateEnergyRoute != nil:
		return m.moduleMessengers[gridtypes.ModuleName].HandleMsg(ctx, contractAddr, contractIBCPortID, contractMsg)
	case contractMsg.EditEnergyRoute != nil:
		return m.moduleMessengers[gridtypes.ModuleName].HandleMsg(ctx, contractAddr, contractIBCPortID, contractMsg)
	case contractMsg.EditEnergyRouteName != nil:
		return m.moduleMessengers[gridtypes.ModuleName].HandleMsg(ctx, contractAddr, contractIBCPortID, contractMsg)
	case contractMsg.DeleteEnergyRoute != nil:
		return m.moduleMessengers[gridtypes.ModuleName].HandleMsg(ctx, contractAddr, contractIBCPortID, contractMsg)

	case contractMsg.CreateThought != nil:
		return m.moduleMessengers[dmntypes.ModuleName].HandleMsg(ctx, contractAddr, contractIBCPortID, contractMsg)
	case contractMsg.ForgetThought != nil:
		return m.moduleMessengers[dmntypes.ModuleName].HandleMsg(ctx, contractAddr, contractIBCPortID, contractMsg)
	case contractMsg.ChangeThoughtInput != nil:
		return m.moduleMessengers[dmntypes.ModuleName].HandleMsg(ctx, contractAddr, contractIBCPortID, contractMsg)
	case contractMsg.ChangeThoughtPeriod != nil:
		return m.moduleMessengers[dmntypes.ModuleName].HandleMsg(ctx, contractAddr, contractIBCPortID, contractMsg)
	case contractMsg.ChangeThoughtBlock != nil:
		return m.moduleMessengers[dmntypes.ModuleName].HandleMsg(ctx, contractAddr, contractIBCPortID, contractMsg)
	case contractMsg.ChangeThoughtGasPrice != nil:
		return m.moduleMessengers[dmntypes.ModuleName].HandleMsg(ctx, contractAddr, contractIBCPortID, contractMsg)
	case contractMsg.ChangeThoughtParticle != nil:
		return m.moduleMessengers[dmntypes.ModuleName].HandleMsg(ctx, contractAddr, contractIBCPortID, contractMsg)
	case contractMsg.ChangeThoughtName != nil:
		return m.moduleMessengers[dmntypes.ModuleName].HandleMsg(ctx, contractAddr, contractIBCPortID, contractMsg)
	case contractMsg.TokenFactory != nil:
		return m.moduleMessengers[tokenfactorytypes.ModuleName].HandleMsg(ctx, contractAddr, contractIBCPortID, contractMsg)
	default:
		// If no handler could handle the message, return an error
		return nil, nil, errorsmod.Wrap(sdkerrors.ErrUnknownRequest, "unknown cyber message type")
	}

	// Iterate over the module message handlers and dispatch to the appropriate one
	//for _, handler := range m.moduleMessengers {
	//	event, resp, err := handler.HandleMsg(ctx, contractAddr, contractIBCPortID, contractMsg)
	//	if err != nil {
	//		if err == ErrHandleMsg {
	//			// This handler cannot handle the message, try the next one
	//			continue
	//		}
	//		// Some other error occurred, return it
	//		return nil, nil, err
	//	}
	//	// Message was handled successfully, return the result
	//	return event, resp, nil
	//}
	//
	//// If no handler could handle the message, return an error
	//return nil, nil, errorsmod.Wrap(sdkerrors.ErrUnknownRequest, "unknown cyber message type")

}
