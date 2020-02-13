package types

type GenesisState struct {
	Params Params `json:"params" yaml:"params"`
}

func NewGenesisState(params Params) GenesisState {
	return GenesisState{
		Params: params,
	}
}

func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultParams())
}

func ValidateGenesis(data GenesisState) error {
	return data.Params.Validate()
}
