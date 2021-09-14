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
	Program sdk.AccAddress
	Label    string
}

func NewQueryJobParams(program sdk.AccAddress, label string) QueryJobParams {
	return QueryJobParams{program, label}
}