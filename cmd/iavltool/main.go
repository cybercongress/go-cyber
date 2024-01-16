package main

import (
	"github.com/cybercongress/go-cyber/cmd/iavltool/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		return
	}
}
