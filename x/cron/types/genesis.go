package types

func NewGenesisState(params Params, jobs []Job) *GenesisState {
	return &GenesisState{
		Params: 	 params,
		Jobs: 	 	 jobs,
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
