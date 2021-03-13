package types

func NewGenesisState(params Params, routes []Route) *GenesisState {
	return &GenesisState{
		Params: 	 params,
		Routes: 	 routes,
	}
}

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

func ValidateGenesis(_ GenesisState) error {
	// TODO add validation
	return nil
}
