package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	QueryParams		= "params"
	QueryJob		= "job"
	QueryJobStats   = "job_stats"
	QueryJobs		= "jobs"
	QueryJobsStats  = "jobs_stats"
)

type QueryJobParams struct {
	Creator  sdk.AccAddress
	Contract sdk.AccAddress
	Label    string
}

func NewQueryJobParams(creator, contract sdk.AccAddress, label string) QueryJobParams {
	return QueryJobParams{creator, contract, label}
}