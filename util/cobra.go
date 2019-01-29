package util

import "github.com/spf13/cobra"

type CobraCmdFunc func(cmd *cobra.Command, args []string) error

// Returns a single function that calls each argument function in sequence
// RunE, PreRunE, PersistentPreRunE, etc. all have this same signature
func ConcatCobraCmdFuncs(fns ...CobraCmdFunc) CobraCmdFunc {
	return func(cmd *cobra.Command, args []string) error {
		for _, fn := range fns {
			if fn != nil {
				if err := fn(cmd, args); err != nil {
					return err
				}
			}
		}
		return nil
	}
}
