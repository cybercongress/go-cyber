package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	ModuleName   = "dmn"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

var (
	ThoughtKey      = []byte{0x00}
	ThoughtStatsKey = []byte{0x01}
	ParamsKey       = []byte{0x02}
)

func GetThoughtKey(program sdk.AccAddress, name string) []byte {
	key := append(program.Bytes(), []byte(name)...)
	return append(ThoughtKey, key...)
}

func GetThoughtStatsKey(program sdk.AccAddress, name string) []byte {
	key := append(program.Bytes(), []byte(name)...)
	return append(ThoughtStatsKey, key...)
}
