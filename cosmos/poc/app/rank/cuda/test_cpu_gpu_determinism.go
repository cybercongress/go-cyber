package main

import (
	"fmt"
	cpurank "github.com/cybercongress/cyberd/cosmos/poc/app/rank"
	. "github.com/cybercongress/cyberd/cosmos/poc/app/storage"
	"github.com/tendermint/tendermint/libs/log"
)

/*
#cgo CFLAGS: -I/usr/lib/
#cgo LDFLAGS: -lcbdrank -lcudart
#include "cbdrank.h"
*/
import "C"

func main() {

	/* --- Init network ------------------------------- */
	stakes := []uint64{3, 1, 2}

	inLinksCount := []uint32{0, 0, 1, 5, 4, 0, 1, 0}
	inLinksStartIndex := []uint64{0, 0, 0, 1, 6, 10, 10, 11}
	outLinksCount := []uint32{2, 2, 1, 1, 3, 1, 0, 1}

	inLinksOuts := []uint64{7, 1, 4, 4, 4, 2, 5, 0, 0, 1, 3}
	inLinksUsers := []uint64{0, 2, 0, 1, 2, 0, 1, 1, 2, 1, 1}
	outLinksUsers := []uint64{1, 2, 1, 2, 0, 1, 0, 1, 2, 1, 0}

	/* --- Convert to C ------------------------------- */
	cStakesSize := C.ulong(len(stakes))
	cCidsSize := C.ulong(len(inLinksStartIndex))
	cLinksSize := C.ulong(len(inLinksOuts))

	cStakes := (*C.ulong)(&stakes[0])

	cInLinksCount := (*C.uint)(&inLinksCount[0])
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
		cInLinksCount, cOutLinksCount,
		cInLinksOuts, cInLinksUsers, cOutLinksUsers,
		cRank,
	)

	fmt.Printf("Rank calculated on raw gpu...\n")
	for c, r := range rank {
		fmt.Printf("%v -> %v\n", c, r)
	}

	m := InMemoryStorage{}
	m.Empty()
	for i := 0; i < 8; i++ {
		m.AddCid(Cid(i), CidNumber(i))
	}
	m.UpdateStakeByNumber(AccountNumber(0), 3)
	m.UpdateStakeByNumber(AccountNumber(1), 1)
	m.UpdateStakeByNumber(AccountNumber(2), 2)

	m.AddLink(LinkedCids{FromCid: CidNumber(0), ToCid: CidNumber(4), Creator: AccountNumber(1)})
	m.AddLink(LinkedCids{FromCid: CidNumber(0), ToCid: CidNumber(4), Creator: AccountNumber(2)})
	m.AddLink(LinkedCids{FromCid: CidNumber(4), ToCid: CidNumber(3), Creator: AccountNumber(1)})
	m.AddLink(LinkedCids{FromCid: CidNumber(7), ToCid: CidNumber(2), Creator: AccountNumber(0)})
	m.AddLink(LinkedCids{FromCid: CidNumber(1), ToCid: CidNumber(3), Creator: AccountNumber(2)})
	m.AddLink(LinkedCids{FromCid: CidNumber(2), ToCid: CidNumber(3), Creator: AccountNumber(0)})
	m.AddLink(LinkedCids{FromCid: CidNumber(3), ToCid: CidNumber(6), Creator: AccountNumber(1)})
	m.AddLink(LinkedCids{FromCid: CidNumber(1), ToCid: CidNumber(4), Creator: AccountNumber(1)})
	m.AddLink(LinkedCids{FromCid: CidNumber(4), ToCid: CidNumber(3), Creator: AccountNumber(0)})
	m.AddLink(LinkedCids{FromCid: CidNumber(4), ToCid: CidNumber(3), Creator: AccountNumber(2)})
	m.AddLink(LinkedCids{FromCid: CidNumber(5), ToCid: CidNumber(4), Creator: AccountNumber(1)})

	rank, _ = cpurank.CalculateRank(&m, cpurank.CPU, log.NewNopLogger())

	fmt.Printf("Rank calculated on cpu...\n")
	for c, r := range rank {
		fmt.Printf("%v -> %v\n", c, r)
	}

	rank, _ = cpurank.CalculateRank(&m, cpurank.GPU, log.NewNopLogger())

	fmt.Printf("Rank calculated on gpu via cyberd ...\n")
	for c, r := range rank {
		fmt.Printf("%v -> %v\n", c, r)
	}
}
