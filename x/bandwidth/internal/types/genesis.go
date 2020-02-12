package types

type GenesisState struct {
	Params Params `json:"Params" yaml:"Params"`
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
	if err := data.Params.Validate(); err != nil {
		return err
	}
	return nil
}
