package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cybercongress/go-cyber/v4/x/clock/types"
)

func TestParamsValidate(t *testing.T) {
	testCases := []struct {
		name    string
		params  types.Params
		success bool
	}{
		{
			"Success - Default",
			types.DefaultParams(),
			true,
		},
		{
			"Success - Meets min Gas",
			types.NewParams(100_000),
			true,
		},
		{
			"Success - Meets min Gas",
			types.NewParams(500_000),
			true,
		},
		{
			"Fail - Not Enough Gas",
			types.NewParams(1),
			false,
		},
		{
			"Fail - Not Enough Gas",
			types.NewParams(100),
			false,
		},
		{
			"Fail - Not Enough Gas",
			types.NewParams(1_000),
			false,
		},
		{
			"Fail - Not Enough Gas",
			types.NewParams(10_000),
			false,
		},
	}

	for _, tc := range testCases {
		err := tc.params.Validate()

		if tc.success {
			require.NoError(t, err, tc.name)
		} else {
			require.Error(t, err, tc.name)
		}
	}
}
