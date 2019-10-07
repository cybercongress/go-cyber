package types

type GenesisState struct {
	Params Params
}

func NewGenesisState(params Params) GenesisState {
	return GenesisState{
		Params: params,
	}
}

func DefaultGenesisState() GenesisState {
	return NewGenesisState(NewBaseDefaultParams())
}

func ValidateGenesis(data GenesisState) error {
	return data.Params.Validate()
}
