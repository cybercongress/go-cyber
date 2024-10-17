package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	pluginstypes "github.com/cybercongress/go-cyber/v5/plugins/types"
	"github.com/cybercongress/go-cyber/v5/x/resources/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
	case msg.Investmint != nil:
		if msg.Investmint.Neuron != contractAddr.String() {
			return nil, nil, wasmvmtypes.InvalidRequest{Err: "investmint wrong neuron"}
		}

		msgServer := keeper.NewMsgServerImpl(*m.keeper)

		if err := msg.Investmint.ValidateBasic(); err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed validating msg")
		}

		res, err := msgServer.Investmint(
			sdk.WrapSDKContext(ctx),
			msg.Investmint,
		)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "investmint msg")
		}

		responseBytes, err := json.Marshal(*res)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed to serialize investmint response")
		}

		resp := [][]byte{responseBytes}

		return nil, resp, nil
	default:
		return nil, nil, pluginstypes.ErrHandleMsg
	}
}
