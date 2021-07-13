#include <stdint.h>

void calculate_rank(
    uint64_t *stakes, uint64_t stakesSize,                    /* user stakes and corresponding array size */
    uint64_t cidsSize, uint64_t linksSize,                    /* cids and links array size */
    uint32_t *inLinksCount, uint32_t *outLinksCount,          /* array index - cid index*/
    uint64_t *inLinksOuts, uint64_t *outLinksIns,
    uint64_t *inLinksUsers,                                   /*all incoming links from all users*/
    uint64_t *outLinksUsers,                                  /*all outgoing links from all users*/
    double dampingFactor,                                     /* value of damping factor*/
    double tolerance,                                         /* value of needed tolerance */
    double *rank,                                             /* array index - cid index*/
    double *entropy,                                          /* array index - cid index*/
    double *luminosity,                                       /* array index - cid index*/
    double *karma                                             /* array index - cid index*/
);