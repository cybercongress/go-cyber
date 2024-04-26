package main

import (
	"github.com/cybercongress/go-cyber/v4/cmd/iavltool/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		return
	}
}
