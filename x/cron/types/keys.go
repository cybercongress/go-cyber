package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	ModuleName 		= "cron"
	StoreKey 		= ModuleName
  	RouterKey 		= ModuleName
	QuerierRoute 	= ModuleName

	ActionCronAddJob 			= "add_job"
	ActionCronRemoveJob 		= "remove_job"
	ActionCronChangeJobLabel 	= "change_label"
	ActionCronChangeJobCID 		= "change_cid"
	ActionCronChangeJobCallData = "change_call_data"
	ActionCronChangeJobGasPrice = "change_gas_price"
	ActionCronChangeJobPeriod   = "change_period"
	ActionCronChangeJobBlock    = "change_block"

	QueryParams		= "params"
	QueryJob		= "job"
	QueryJobStats   = "job_stats"
	QueryJobs		= "jobs"
	QueryJobsStats  = "jobs_stats"
)

var (
	JobKey          = []byte{0x00}
	JobStatsKey		= []byte{0x01}
)

func GetJobKey(contract, creator sdk.AccAddress, label string) []byte {
	key := append(append(contract.Bytes(), creator.Bytes()...), []byte(label)...)
	return append(JobKey, key...)
}

func GetJobStatsKey(contract, creator sdk.AccAddress, label string) []byte {
	key := append(append(contract.Bytes(), creator.Bytes()...), []byte(label)...)
	return append(JobStatsKey, key...)
}

