package main

import (
	"github.com/deep-foundation/deep-chain/cmd/iavltool/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		return
	}
}
