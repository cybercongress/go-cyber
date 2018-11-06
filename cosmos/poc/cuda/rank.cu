#include <stdint.h>
#include <stdio.h>
#include "types.h"

extern "C" {

    void calculate_rank(
        uint64_t *stakes, uint64_t stakesSize, /* User stakes and corresponding array size */
        cid *cids, uint64_t cidsSize, /* Cids links */
        cid_link *inLinks, cid_link *outLinks /* Incoming and Outgoing cids links */
    ) {

        for(int i = 0; i < stakesSize; i++) {
            printf("%lu \n", stakes[i]);
        }
    }
};
