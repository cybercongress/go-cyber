package main

import (
	"encoding/binary"
	"fmt"
	"github.com/cybercongress/cyberd/merkle"
)

func main() {

	t := merkle.NewTree()

	data := make([]byte, 8)

	for i := 0 ; i < 31 ; i++ {
		binary.LittleEndian.PutUint64(data, uint64(i))
		t.Push(data)
	}

	proofs := t.GetIndexProofs(18)

	for _, proof := range proofs {
		fmt.Println(proof)
	}
}
