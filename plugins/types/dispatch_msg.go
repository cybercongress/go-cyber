package types

import (
	"encoding/json"
	"errors"
	dmnkeeper "github.com/cybercongress/go-cyber/v4/x/dmn/keeper"
	graphkeeper "github.com/cybercongress/go-cyber/v4/x/graph/keeper"
	gridkeeper "github.com/cybercongress/go-cyber/v4/x/grid/keeper"
	resourceskeeper "github.com/cybercongress/go-cyber/v4/x/resources/keeper"

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
	wrapped          wasmkeeper.Messenger
	moduleMessengers []ModuleMessenger
	graphKeeper      *graphkeeper.GraphKeeper
	dmnKeeper        *dmnkeeper.Keeper
	gridKeeper       *gridkeeper.Keeper
	resourcesKeeper  *resourceskeeper.Keeper
}

func CustomMessageDecorator(
	moduleMessengers []ModuleMessenger,
	graph *graphkeeper.GraphKeeper,
	dmn *dmnkeeper.Keeper,
	grid *gridkeeper.Keeper,
	resources *resourceskeeper.Keeper,
) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessenger{
			wrapped:          old,
			moduleMessengers: moduleMessengers,
			graphKeeper:      graph,
			dmnKeeper:        dmn,
			gridKeeper:       grid,
			resourcesKeeper:  resources,
		}
	}
}

func (m *CustomMessenger) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) ([]sdk.Event, [][]byte, error) {
	if msg.Custom != nil {
		var contractMsg CyberMsg
		if err := json.Unmarshal(msg.Custom, &contractMsg); err != nil {
			return nil, nil, errorsmod.Wrap(err, "cyber msg error")
		}

		// Iterate over the module message handlers and dispatch to the appropriate one
		for _, handler := range m.moduleMessengers {
			event, resp, err := handler.HandleMsg(ctx, contractAddr, contractIBCPortID, contractMsg)
			if err != nil {
				if err == ErrHandleMsg {
					// This handler cannot handle the message, try the next one
					continue
				}
				// Some other error occurred, return it
				return nil, nil, err
			}
			// Message was handled successfully, return the result
			return event, resp, nil
		}

		// If no handler could handle the message, return an error
		return nil, nil, errorsmod.Wrap(sdkerrors.ErrUnknownRequest, "unknow cyber message type")
	}
	return m.wrapped.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
}
