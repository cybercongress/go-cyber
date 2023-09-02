package types

func NewGenesisState(params Params, routes []Route) *GenesisState {
	return &GenesisState{
		Params: params,
		Routes: routes,
	}
}

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

func ValidateGenesis(state GenesisState) error {
	return state.Params.Validate()
}
