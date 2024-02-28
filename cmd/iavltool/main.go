package main

import (
	"github.com/cybercongress/go-cyber/v3/cmd/iavltool/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		return
	}
}
