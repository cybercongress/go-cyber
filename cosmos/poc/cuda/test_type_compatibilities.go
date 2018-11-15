package main

import (
	"fmt"
)

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -lrank
#cgo LDFLAGS: -lcudart
#include "rank.h"
*/
import "C"

func main() {

	/* --- Init network ------------------------------- */
	stakes := []uint64{3, 1, 2}

	inLinksCount := []uint32{0, 0, 1, 5, 4, 0, 1, 0}
	inLinksStartIndex := []uint64{0, 0, 0, 1, 6, 10, 10, 11}
	outLinksCount := []uint32{2, 2, 1, 1, 3, 1, 0, 1}
	outLinksStartIndex := []uint64{0, 2, 4, 5, 6, 9, 10, 10}

	inLinksOuts := []uint64{7, 1, 4, 4, 4, 2, 5, 0, 0, 1, 3}
	inLinksUsers := []uint64{0, 2, 0, 1, 2, 0, 1, 1, 2, 1, 1}
	outLinksUsers := []uint64{1, 2, 1, 2, 0, 1, 0, 1, 2, 1, 0}

	/* --- Convert to C ------------------------------- */
	cStakesSize := C.ulong(len(stakes))
	cCidsSize := C.ulong(len(inLinksStartIndex))
	cLinksSize := C.ulong(len(inLinksOuts))

	cStakes := (*C.ulong)(&stakes[0])

	cInLinksStartIndex := (*C.ulong)(&inLinksStartIndex[0])
	cInLinksCount := (*C.uint)(&inLinksCount[0])

	cOutLinksStartIndex := (*C.ulong)(&outLinksStartIndex[0])
	cOutLinksCount := (*C.uint)(&outLinksCount[0])

	cInLinksOuts := (*C.ulong)(&inLinksOuts[0])
	cInLinksUsers := (*C.ulong)(&inLinksUsers[0])
	cOutLinksUsers := (*C.ulong)(&outLinksUsers[0])

	/* --- Init rank ---------------------------------- */
	rank := make([]float64, len(inLinksStartIndex))
	cRank := (*C.double)(&rank[0])
	/* --- Run Computation ---------------------------- */
	fmt.Printf("Invoking cuda library...\n")
	C.calculate_rank(
		cStakes, cStakesSize, cCidsSize, cLinksSize,
		cInLinksStartIndex, cInLinksCount,
		cOutLinksStartIndex, cOutLinksCount,
		cInLinksOuts, cInLinksUsers, cOutLinksUsers,
		cRank,
	)

	for c, r := range rank {
		fmt.Printf("%v -> %v\n", c, r)
	}
}
