#ifndef TYPES_HEADER_FILE
#define TYPES_HEADER_FILE

#include <stdint.h>

/*******************************************/
/* REPRESENTS LINK                         */
/*******************************************/
typedef struct {
    /* Index of opposite cid in cids array */
    uint64_t opposite_cid_index;
    /* Index of user stake in stakes array */
    uint64_t user_index;
} cid_link;

/*******************************************/
/* REPRESENTS CID                          */
/*******************************************/
typedef struct {
	uint64_t in_links_start_index;
	uint64_t in_links_count;

	uint64_t out_links_start_index;
	uint64_t out_links_count;
} cid;

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
