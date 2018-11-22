#ifndef TYPES_HEADER_FILE
#define TYPES_HEADER_FILE

#include <stdint.h>

/*******************************************/
/* REPRESENTS INCOMING LINK WITH IT WEIGHT */
/*******************************************/
/* Finds maximum rank difference for single element  */
/*                                                   */
/*****************************************************/
typedef struct {
    /* Index of opposite cid in cids array */
    uint64_t fromIndex;
    /* Index of user stake in stakes array */
    double weight;
} CompressedInLink;

#endif
