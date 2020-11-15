package rpc

import (
	"github.com/cosmos/cosmos-sdk/types"
	rpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"

	"github.com/cybercongress/go-cyber/x/link"
)

func IsLinkExist(ctx *rpctypes.Context, from string, to string, address string) (bool, error) {

	if len(address) == 0 {
		return cyberdApp.IsLinkExist(link.Cid(from), link.Cid(to), nil), nil
	}

	accAddress, err := types.AccAddressFromBech32(address)
	if err != nil {
		return false, err
	}

	return cyberdApp.IsLinkExist(link.Cid(from), link.Cid(to), accAddress), nil
}
