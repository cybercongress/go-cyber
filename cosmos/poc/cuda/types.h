#ifndef TYPES_HEADER_FILE
#define TYPES_HEADER_FILE

#include <stdint.h>

/* Go -> CGO passing types */
typedef struct {
    /* Index of opposite cid in cids array */
    uint64_t opposite_cid_index;
    /* Index of user stake in stakes array */
    uint64_t user_index;
} cid_link;

typedef struct {
	uint64_t in_links_start_index;
	uint64_t in_links_count;

	uint64_t out_links_start_index;
	uint64_t out_links_count;
} cid;

#endif
