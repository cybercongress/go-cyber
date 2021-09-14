package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	ModuleName 		= "cron"
	StoreKey 		= ModuleName
  	RouterKey 		= ModuleName
	QuerierRoute 	= ModuleName
)

var (
	JobKey               = []byte{0x00}
	JobStatsKey		     = []byte{0x01}
)

func GetJobKey(program sdk.AccAddress, label string) []byte {
	key := append(program.Bytes(), []byte(label)...)
	return append(JobKey, key...)
}

func GetJobStatsKey(program sdk.AccAddress, label string) []byte {
	key := append(program.Bytes(), []byte(label)...)
	return append(JobStatsKey, key...)
}

