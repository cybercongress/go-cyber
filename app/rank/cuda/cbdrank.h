#include <stdint.h>

void calculate_rank(
    uint64_t *stakes, uint64_t stakesSize,                    /* User stakes and corresponding array size */
    uint64_t cidsSize, uint64_t linksSize,                    /* Cids count */
    uint32_t *inLinksCount, uint32_t *outLinksCount,          /* array index - cid index*/
    uint64_t *inLinksOuts, uint64_t *inLinksUsers,            /*all incoming links from all users*/
    uint64_t *outLinksUsers,                                  /*all outgoing links from all users*/
    double *rank                                              /* array index - cid index*/
);