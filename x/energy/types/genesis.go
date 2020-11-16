package types

type GenesisState struct {
	Params Params `json:"params" yaml:"params"`
	Routes Routes `json:"routes" yaml:"routes"`
}

func NewGenesisState(params Params, routes []Route) GenesisState {
	return GenesisState{
		Params: 	 params,
		Routes: 	 routes,
	}
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

func ValidateGenesis(data GenesisState) error {
	// TODO add validation
	return nil
}
