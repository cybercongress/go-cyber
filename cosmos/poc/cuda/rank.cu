#include <stdint.h>
#include <stdio.h>
#include <thrust/transform.h>
#include <thrust/transform_reduce.h>
#include <thrust/device_vector.h>
#include <thrust/execution_policy.h>
#include <thrust/functional.h>
#include "types.h"

const double DUMP_FACTOR = 0.85;
const double TOLERANCE = 1e-3;

/*******************************************/
/* REPRESENTS INCOMING LINK WITH IT WEIGHT */
/*******************************************/
typedef struct {
    /* Index of opposite cid in cids array */
    uint64_t fromIndex;
    /* Index of user stake in stakes array */
    double weight;
} InLink;


/*****************************************************/
/* KERNEL: RUN SINGLE RANK ITERATION                 */
/*****************************************************/
/* For all given arrays, array index = cidId         */
/* Except: *inLinks, that represent 1D array of all  */
/*   links with corresponding weights                */
/*****************************************************/
__global__
void run_rank_iteration(
    InLink *inLinks,
    double *prevRank,
    double *rank,
    uint64_t *inLinksStartIndex,
    uint32_t *inLinksCount,
    uint64_t rankSize,
    double innerProductOverSize,
    double defaultRank
) {

	int index = blockIdx.x * blockDim.x + threadIdx.x;
    uint64_t stride = blockDim.x * gridDim.x;

    for (uint64_t i = index; i < rankSize; i += stride) {
        double ksum = innerProductOverSize;
        for (uint64_t j = 0; j < inLinksCount[i]; j++) {
           // ksum = prevRank[inLinks[j].fromIndex] * inLinks[j].weight + ksum
           ksum = __fmaf_rz(prevRank[inLinks[j].fromIndex], inLinks[j].weight, ksum);
        }
        // rank[i] = ksum * DUMP_FACTOR + defaultRank
        rank[i] = __fmaf_rz(ksum, DUMP_FACTOR, defaultRank); // ksum * DUMP_FACTOR + defaultRank
    }
}


/*****************************************************/
/* KERNEL: DOUBLE ABS FUNCTOR                        */
/*****************************************************/
/* Return absolute value for double                  */
/*****************************************************/
struct absolute_value {
  __device__ double operator()(const double &x) const {
    return x < 0.0 ? -x : x;
  }
};


/*****************************************************/
/* KERNEL: FINDS MAXIMUM RANKS DIFFERENCE            */
/*****************************************************/
/* Finds maximum rank difference for single element  */
/*                                                   */
/*****************************************************/
double find_max_ranks_diff(double *prevRank, double *newRank, uint64_t rankSize) {

    thrust::device_vector<double> ranksDiff(rankSize);
    thrust::device_ptr<double> newRankBegin(newRank);
    thrust::device_ptr<double> prevRankBegin(prevRank);
    thrust::device_ptr<double> prevRankEnd(prevRank + rankSize);
    thrust::transform(thrust::device,
        prevRankBegin, prevRankEnd, newRankBegin, ranksDiff.begin(), thrust::minus<double>()
    );

    return thrust::transform_reduce(thrust::device,
        ranksDiff.begin(), ranksDiff.end(), absolute_value(), 0.0, thrust::maximum<double>()
    );
}

extern "C" {

    void calculate_rank(
        uint64_t *stakes, uint64_t stakesSize, /* User stakes and corresponding array size */
        cid *cids, uint64_t cidsSize, /* Cids links */
        cid_link *inLinks, cid_link *outLinks /* Incoming and Outgoing cids links */
    ) {

        printf("Cuda !!!!!!!!!!!!!!!!!!\n");

        double *prevRank, *rank;
        cudaMalloc(&rank, cidsSize*sizeof(double));
        cudaMalloc(&prevRank, cidsSize*sizeof(double));

        int steps = 0;
        double change = TOLERANCE + 1;
        while(change > TOLERANCE) {
        	//run_rank_iteration()
        	//change = calculateChange(prevrank, rank)
        	//prevrank = rank
        	steps++;
        	return;
        }

        cudaFree(rank);
        cudaFree(prevRank);
    }
};
