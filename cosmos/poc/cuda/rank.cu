#include <stdint.h>
#include <stdio.h>
#include "types.h"


/******************************************/
/* CELL STRUCT LEADING TO ARRAY OF STRUCT */
/******************************************/
typedef struct {
    /* Index of opposite cid in cids array */
    uint64_t fromIndex;
    /* Index of user stake in stakes array */
    uint32_t weight;
} InLink;


/*****************************************************/
/* KERNEL: RUN SINGLE RANK ITERATION                         */
/*****************************************************/
/* For all given arrays, array index = cidId         */
/* Except: *inLinks, that represent 1d array of all  */
/*   i->j links with corresponding weights           */
/*****************************************************/
__global__
void run_rank_iteration(
    InLink *inLinks,
    uint64_t *prevRank,
    uint64_t *rank,
    uint64_t *inLinksStartIndex,
    uint32_t *inLinksCount,
    uint64_t rankSize,
    uint64_t innerProductOverSize,
    uint64_t defaultRank
) {

	int index = blockIdx.x * blockDim.x + threadIdx.x;
    int stride = blockDim.x * gridDim.x;

    for (int i = index; i < rankSize; i += stride) {
        rank[i] = innerProductOverSize;
        for (int j = 0; j < inLinksCount[i]; j++) {
           rank[i] += prevRank[inLinks[j].fromIndex] * inLinks[j].weight;
        }
        rank[i] = rank[i] / 20 * 17 + defaultRank;
    }
}

/******************************************/
/* CELL STRUCT LEADING TO ARRAY OF STRUCT */
/******************************************/
extern "C" {

    void calculate_rank(
        uint64_t *stakes, uint64_t stakesSize, /* User stakes and corresponding array size */
        cid *cids, uint64_t cidsSize, /* Cids links */
        cid_link *inLinks, cid_link *outLinks /* Incoming and Outgoing cids links */
    ) {


    }
};
