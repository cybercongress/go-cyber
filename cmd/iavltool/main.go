package main

import (
	"fmt"

	"github.com/cybercongress/go-cyber/cmd/iavltool/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Println("error executing the iavl tool")
		panic(err)
	}
}
