package rank

type GenesisState struct{}

func NewGenesisState() GenesisState {
	return GenesisState{}
}

func ValidateGenesis(data GenesisState) error {
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func ExportGenesis() GenesisState {
	return GenesisState{}
}

func InitGenesis() {}
