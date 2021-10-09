package types

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewThought(
	program string,
	trigger Trigger,
	load Load,
	name string,
	particle string,
) Thought {
	return Thought{
		Program:  program,
		Trigger:  trigger,
		Load:     load,
		Name :    name,
		Particle: particle,
	}
}

type Thoughts []Thought

type ThoughtsStats []ThoughtStats

func (js Thoughts) Sort() {
	sort.Sort(js)
}

func (js Thoughts) Len() int { return len(js) }

func (js Thoughts) Less(i, j int) bool {
	return js[j].Load.GasPrice.IsLT(js[i].Load.GasPrice)
	//return js[i].Load.GasPrice.Amount.GT(js[j].Load.GasPrice.Amount)
}

func (js Thoughts) Swap(i, j int) { js[i], js[j] = js[j], js[i] }

//______________________________________________________________________


func NewTrigger (period, block uint64) Trigger {
	return Trigger{
		Period: period,
		Block: block,
	}
}

func NewStats (
	program, name string,
	calls, fees, gas, block uint64,
) ThoughtStats {
	return ThoughtStats{
		Program:   program,
		Name:      name,
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
