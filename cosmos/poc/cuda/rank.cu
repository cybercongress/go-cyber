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

/*****************************************************/
/* KERNEL: RUN SINGLE RANK ITERATION                 */
/*****************************************************/
/* All in links used here are compressed in links    */
/*****************************************************/
__global__
void run_rank_iteration(
    CompressedInLink *inLinks,                            /* all compressed in links */
    double *prevRank, double *rank,                       /* array index - cid index */
    uint64_t *inLinksStartIndex, uint32_t *inLinksCount,  /* array index - cid index */
    uint64_t rankSize,
    double innerProductOverSize, double defaultRank
) {

	int index = blockIdx.x * blockDim.x + threadIdx.x;
    uint64_t stride = blockDim.x * gridDim.x;

    for (uint64_t i = index; i < rankSize; i += stride) {
        double ksum = innerProductOverSize;
        for (uint64_t j = inLinksStartIndex[i]; j < inLinksStartIndex[i] + inLinksCount[i]; j++) {
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
/* HOST: FINDS MAXIMUM RANKS DIFFERENCE              */
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

/*****************************************************/
/* KERNEL: CALCULATE CID TOTAL OUTS STAKE            */
/*****************************************************/
__global__
void calculateCidTotalOutStake(
    uint64_t cidsSize,
    uint64_t *stakes,                                        /*array index - user index*/
    uint64_t *outLinksStartIndex, uint32_t *outLinksCount,   /*array index - cid index*/
    uint64_t *outLinksUsers,                                 /*all out links from all users*/
    /*returns*/ uint64_t *cidsTotalOutStakes                 /*array index - cid index*/
) {

	int index = blockIdx.x * blockDim.x + threadIdx.x;
    uint64_t stride = blockDim.x * gridDim.x;

    for (uint64_t i = index; i < cidsSize; i += stride) {
        uint64_t totalOutStake = 0;
        for (uint64_t j = outLinksStartIndex[i]; j < outLinksStartIndex[i] + outLinksCount[i]; j++) {
           totalOutStake += stakes[outLinksUsers[j]];
        }
        cidsTotalOutStakes[i] = totalOutStake;
    }
}

/*********************************************************/
/* KERNEL: CALCULATE COMPRESSED IN LINKS COUNT FOR CIDS  */
/*********************************************************/
__global__
void getCompressedInLinksCount(
    uint64_t cidsSize,
    uint64_t *inLinksStartIndex, uint32_t *inLinksCount,                    /*array index - cid index*/
    uint64_t *inLinksOuts,                                                  /*all incoming links from all users*/
    /*returns*/ uint32_t *compressedInLinksCount                            /*array index - cid index*/
) {

	int index = blockIdx.x * blockDim.x + threadIdx.x;
    uint64_t stride = blockDim.x * gridDim.x;

    for (uint64_t i = index; i < cidsSize; i += stride) {

        if(inLinksCount[i] == 0) {
            continue;
        }

        uint32_t compressedLinksCount = 0;
        for(uint64_t j = inLinksStartIndex[i]; j < inLinksStartIndex[i]+inLinksCount[i]; j++) {
            if(j == inLinksStartIndex[i] || inLinksOuts[j] != inLinksOuts[j-1]) {
                compressedLinksCount++;
            }
        }
        compressedInLinksCount[i] = compressedLinksCount;
    }
}

/*********************************************************/
/* DEVICE: USER TO DIVIDE TWO uint64                     */
/*********************************************************/
__device__
double ddiv_rz(uint64_t *a, uint64_t *b) {
    return __ddiv_rz(__ull2double_rz(*a), __ull2double_rz(*b));
}


/*********************************************************/
/* KERNEL: CALCULATE COMPRESSED IN LINKS                 */
/*********************************************************/
__global__
void getCompressedInLinks(
    uint64_t cidsSize,
    uint64_t *inLinksStartIndex, uint32_t *inLinksCount, uint64_t *cidsTotalOutStakes,   /*array index - cid index*/
    uint64_t *inLinksOuts, uint64_t *inLinksUsers,                                       /*all incoming links from all users*/
    uint64_t *stakes,                                                                    /*array index - user index*/
    uint64_t *compressedInLinksStartIndex, uint32_t *compressedInLinksCount,             /*array index - cid index*/
    /*returns*/ CompressedInLink *compressedInLinks                                      /*all incoming compressed links*/
) {

	int index = blockIdx.x * blockDim.x + threadIdx.x;
    uint64_t stride = blockDim.x * gridDim.x;

    for (uint64_t i = index; i < cidsSize; i += stride) {

        if(inLinksCount[i] == 0) {
            continue;
        }

        uint32_t compressedLinksIndex = compressedInLinksStartIndex[i];

        if(inLinksCount[i] == 1) {
            uint64_t oppositeCid = inLinksOuts[inLinksStartIndex[i]];
            uint64_t compressedLinkStake = stakes[inLinksUsers[inLinksStartIndex[i]]];
            double weight = ddiv_rz(&compressedLinkStake, &cidsTotalOutStakes[oppositeCid]);
            compressedInLinks[compressedLinksIndex] = CompressedInLink {oppositeCid, weight};
            continue;
        }

        uint64_t compressedLinkStake = 0;
        uint64_t lastLinkIndex = inLinksStartIndex[i] + inLinksCount[i] - 1;
        for(uint64_t j = inLinksStartIndex[i]; j < lastLinkIndex + 1; j++) {

            compressedLinkStake += stakes[inLinksUsers[j]];
            if(j == lastLinkIndex || inLinksOuts[j] != inLinksOuts[j+1]) {
                uint64_t oppositeCid = inLinksOuts[j];
                double weight = ddiv_rz(&compressedLinkStake, &cidsTotalOutStakes[oppositeCid]);
                compressedInLinks[compressedLinksIndex] = CompressedInLink {oppositeCid, weight};
                compressedLinksIndex++;
                compressedLinkStake=0;
            }
        }
    }
}

/************************************************************/
/* HOST: CALCULATE COMPRESSED IN LINKS START INDEXES        */
/************************************************************/
/* SEQUENTIAL LOGIC -> CALCULATE ON CPU                     */
/************************************************************/
__host__
void getCompressedInLinksStartIndex(
    uint64_t cidsSize,
    uint32_t *compressedInLinksCount,                   /*array index - cid index*/
    /*returns*/ uint64_t *compressedInLinksStartIndex   /*array index - cid index*/
) {

    uint64_t index = 0;
    for (uint64_t i = 0; i < cidsSize; i++) {
        compressedInLinksStartIndex[i] = index;
        index += compressedInLinksCount[i];
    }
}

extern "C" {

    void calculate_rank(
        uint64_t *stakes, uint64_t stakesSize, /* User stakes and corresponding array size */
        cid *cids, uint64_t cidsSize, /* Cids links */
        cid_link *inLinks, cid_link *outLinks /* Incoming and Outgoing cids links */
    ) {

        /*-------------------------------------------------------------------*/
        printf("Cuda !!!!!!!!!!!!!!!!!!\n");
        printf("Initializing device memory\n");

        uint64_t *cidsTotalOutStakes; // for each cid sum of all out links stake
        uint32_t *compressedInLinksCount; // for each cid count of compressed links

        cudaMalloc(&cidsTotalOutStakes, cidsSize*sizeof(uint64_t));
        cudaMalloc(&compressedInLinksCount, cidsSize*sizeof(uint32_t));
        //todo

        cudaFree(cidsTotalOutStakes);
        cudaFree(compressedInLinksCount);

        /*-------------------------------------------------------------------*/
        printf("Calculating rank\n");

        double *prevRank, *rank;
        cudaMalloc(&rank, cidsSize*sizeof(double));
        cudaMalloc(&prevRank, cidsSize*sizeof(double));

        int steps = 0;
        double change = TOLERANCE + 1;
        while(change > TOLERANCE) {
        	//run_rank_iteration()
        	//change = find_max_ranks_diff(prevrank, rank, cidsSize);
        	//prevrank = rank
        	steps++;
        	return;
        }

        cudaFree(rank);
        cudaFree(prevRank);
    }
};
