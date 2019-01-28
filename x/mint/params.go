package mint

import "fmt"

type Params struct {
	TokensPerBlock int64 `json:"tokens_per_block"`
}

func DefaultParams() Params {
	return Params{
		TokensPerBlock: 634195840, // 200% for 10^16 tokens in case 1 seconds blocks
	}
}

func validateParams(params Params) error {
	if params.TokensPerBlock < 0 {
		return fmt.Errorf("mint parameter TokensPerBlock should be not negative. Got %v ", params.TokensPerBlock)
	}
	return nil
}
