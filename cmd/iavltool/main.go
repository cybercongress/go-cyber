package main

import (
	"github.com/cybercongress/go-cyber/v2/cmd/iavltool/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		return
	}
}
