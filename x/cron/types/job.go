package types

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewJob(
	program string,
	trigger Trigger,
	load Load,
	label string,
	particle string,
) Job {
	return Job{
		Program:  program,
		Trigger:  trigger,
		Load:     load,
		Label:    label,
		Particle: particle,
	}
}

type Jobs []Job

type JobsStats []JobStats

func (js Jobs) Sort() {
	sort.Sort(js)
}

func (js Jobs) Len() int { return len(js) }

func (js Jobs) Less(i, j int) bool {
	return js[j].Load.GasPrice.IsLT(js[i].Load.GasPrice)
	//return js[i].Load.GasPrice.Amount.GT(js[j].Load.GasPrice.Amount)
}

func (js Jobs) Swap(i, j int) { js[i], js[j] = js[j], js[i] }

//______________________________________________________________________


func NewTrigger (period, block uint64) Trigger {
	return Trigger{
		Period: period,
		Block: block,
	}
}

func NewStats (
	program, label string,
	calls, fees, gas, block uint64,
) JobStats {
	return JobStats{
		Program:   program,
		Label:     label,
		Calls:     calls,
		Fees:      fees,
		Gas:       gas,
		LastBlock: block,
	}
}

func NewLoad (calldata string, gasprice sdk.Coin) Load {
	return Load{
		CallData: calldata,
		GasPrice: gasprice,
	}
}
