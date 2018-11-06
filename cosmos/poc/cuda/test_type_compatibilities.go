package main

import (
	"fmt"
	"math/rand"
	"unsafe"
)

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -lrank
#cgo LDFLAGS: -lcudart
#include "rank.h"
*/
import "C"

type CidLink struct {
	oppositeCidIndex uint64
	userIndex        uint64
}

type Cid struct {
	inLinksStartIndex  uint64
	inLinksCount       uint64
	outLinksStartIndex uint64
	outLinksCount      uint64
}

func main() {

	/* --- Stakes ------------------------------------ */
	stakes := make([]uint64, 2)

	for i := range stakes {
		stakes[i] = uint64(rand.Intn(30))
	}

	cStakes := (*C.ulong)(&stakes[0])
	cStakesLen := C.ulong(len(stakes))

	/* --- Cids ------------------------------------   */
	cids := make([]Cid, 10)

	cCids := (**C.cid)(unsafe.Pointer(&cids[0]))
	cCidsLen := C.ulong(len(cids))

	/* --- Cids Links -------------------------------  */
	inLinks := make([]CidLink, 100)
	outLinks := make([]CidLink, 100)

	cInLinks := (**C.cid_link)(unsafe.Pointer(&inLinks[0]))
	cOutLinks := (**C.cid_link)(unsafe.Pointer(&outLinks[0]))

	/* --- Run Computation --------------------------  */
	fmt.Printf("Invoking cuda library...\n")
	C.calculate_rank(
		cStakes, cStakesLen,
		cCids, cCidsLen,
		cInLinks, cOutLinks,
	)
}
