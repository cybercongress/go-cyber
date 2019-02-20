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

const int CUDA_THREAD_BLOCK_SIZE = 256;

/*****************************************************/
/* KERNEL: RUN SINGLE RANK ITERATION                 */
/*****************************************************/
/* All in links used here are compressed in links    */
/*****************************************************/
__global__
void run_rank_iteration(
    CompressedInLink *inLinks,                            /* all compressed in links */
    double *prevRank, double *rank, uint64_t rankSize,    /* array index - cid index */
    uint64_t *inLinksStartIndex, uint32_t *inLinksCount,  /* array index - cid index */
    double defaultRankWithCorrection // default rank + inner product correction
) {

	int index = blockIdx.x * blockDim.x + threadIdx.x;
    uint64_t stride = blockDim.x * gridDim.x;

    for (uint64_t i = index; i < rankSize; i += stride) {

        if(inLinksCount[i] == 0) {
            continue;
        }

        double ksum = 0;
        for (uint64_t j = inLinksStartIndex[i]; j < inLinksStartIndex[i] + inLinksCount[i]; j++) {
           ksum = prevRank[inLinks[j].fromIndex] * inLinks[j].weight + ksum;
           //ksum = __fmaf_rn(prevRank[inLinks[j].fromIndex], inLinks[j].weight, ksum);
        }
        rank[i] = ksum * DUMP_FACTOR + defaultRankWithCorrection;
        //rank[i] = __fmaf_rn(ksum, DUMP_FACTOR, defaultRankWithCorrection);
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
            compressedInLinksCount[i]=0;
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
__device__ __forceinline__
double ddiv_rn(uint64_t *a, uint64_t *b) {
    return __ddiv_rn(__ull2double_rn(*a), __ull2double_rn(*b));
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
            double weight = ddiv_rn(&compressedLinkStake, &cidsTotalOutStakes[oppositeCid]);
            compressedInLinks[compressedLinksIndex] = CompressedInLink {oppositeCid, weight};
            continue;
        }

        uint64_t compressedLinkStake = 0;
        uint64_t lastLinkIndex = inLinksStartIndex[i] + inLinksCount[i] - 1;
        for(uint64_t j = inLinksStartIndex[i]; j < lastLinkIndex + 1; j++) {

            compressedLinkStake += stakes[inLinksUsers[j]];
            if(j == lastLinkIndex || inLinksOuts[j] != inLinksOuts[j+1]) {
                uint64_t oppositeCid = inLinksOuts[j];
                double weight = ddiv_rn(&compressedLinkStake, &cidsTotalOutStakes[oppositeCid]);
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
/* RETURNS TOTAL COMPRESSED LINKS SIZE                      */
/************************************************************/
__host__
uint64_t getLinksStartIndex(
    uint64_t cidsSize,
    uint32_t *linksCount,                   /*array index - cid index*/
    /*returns*/ uint64_t *linksStartIndex   /*array index - cid index*/
) {

    uint64_t index = 0;
    for (uint64_t i = 0; i < cidsSize; i++) {
        linksStartIndex[i] = index;
        index += linksCount[i];
    }
    return index;
}

void swap(double* &a, double* &b){
  double *temp = a;
  a = b;
  b = temp;
}

extern "C" {

    void calculate_rank(
        uint64_t *stakes, uint64_t stakesSize,                    /* User stakes and corresponding array size */
        uint64_t cidsSize, uint64_t linksSize,                    /* Cids count */
        uint32_t *inLinksCount, uint32_t *outLinksCount,          /* array index - cid index*/
        uint64_t *inLinksOuts, uint64_t *inLinksUsers,            /*all incoming links from all users*/
        uint64_t *outLinksUsers,                                  /*all outgoing links from all users*/
        double *rank                                              /* array index - cid index*/
    ) {

        // setbuf(stdout, NULL);
        int CUDA_BLOCKS_NUMBER = (cidsSize + CUDA_THREAD_BLOCK_SIZE - 1) / CUDA_THREAD_BLOCK_SIZE;


        // STEP0: Calculate compressed in links start indexes
        /*-------------------------------------------------------------------*/
        // calculated on cpu
        uint64_t *inLinksStartIndex = (uint64_t*) malloc(cidsSize*sizeof(uint64_t));
        uint64_t *outLinksStartIndex = (uint64_t*) malloc(cidsSize*sizeof(uint64_t));
        getLinksStartIndex(cidsSize, inLinksCount, inLinksStartIndex);
        getLinksStartIndex(cidsSize, outLinksCount, outLinksStartIndex);


        // STEP1: Calculate for each cid total stake by out links
        /*-------------------------------------------------------------------*/
        uint64_t *d_outLinksStartIndex;
        uint32_t *d_outLinksCount;
        uint64_t *d_outLinksUsers;
        uint64_t *d_stakes;  // will be used to calculated links weights, should be freed before rank iterations
        uint64_t *d_cidsTotalOutStakes; // will be used to calculated links weights, should be freed before rank iterations

        cudaMalloc(&d_outLinksStartIndex, cidsSize*sizeof(uint64_t));
        cudaMalloc(&d_outLinksCount,      cidsSize*sizeof(uint32_t));
        cudaMalloc(&d_outLinksUsers,     linksSize*sizeof(uint64_t));
        cudaMalloc(&d_stakes,           stakesSize*sizeof(uint64_t));
        cudaMalloc(&d_cidsTotalOutStakes, cidsSize*sizeof(uint64_t));   //calculated

        cudaMemcpy(d_outLinksStartIndex, outLinksStartIndex, cidsSize*sizeof(uint64_t), cudaMemcpyHostToDevice);
        cudaMemcpy(d_outLinksCount,      outLinksCount,      cidsSize*sizeof(uint32_t), cudaMemcpyHostToDevice);
        cudaMemcpy(d_outLinksUsers,      outLinksUsers,     linksSize*sizeof(uint64_t), cudaMemcpyHostToDevice);
        cudaMemcpy(d_stakes,             stakes,           stakesSize*sizeof(uint64_t), cudaMemcpyHostToDevice);

        calculateCidTotalOutStake<<<CUDA_BLOCKS_NUMBER,CUDA_THREAD_BLOCK_SIZE>>>(
            cidsSize, d_stakes, d_outLinksStartIndex,
            d_outLinksCount, d_outLinksUsers, d_cidsTotalOutStakes
        );

        cudaFree(d_outLinksStartIndex);
        cudaFree(d_outLinksCount);
        cudaFree(d_outLinksUsers);
        /*-------------------------------------------------------------------*/



        // STEP2: Calculate compressed in links count
        /*-------------------------------------------------------------------*/
        uint64_t *d_inLinksStartIndex;
        uint32_t *d_inLinksCount;
        uint64_t *d_inLinksOuts;
        uint32_t *d_compressedInLinksCount;

        // free all before rank iterations
        cudaMalloc(&d_inLinksStartIndex,      cidsSize*sizeof(uint64_t));
        cudaMalloc(&d_inLinksCount,           cidsSize*sizeof(uint32_t));
        cudaMalloc(&d_inLinksOuts,           linksSize*sizeof(uint64_t));
        cudaMalloc(&d_compressedInLinksCount, cidsSize*sizeof(uint32_t));   //calculated

        cudaMemcpy(d_inLinksStartIndex, inLinksStartIndex, cidsSize*sizeof(uint64_t), cudaMemcpyHostToDevice);
        cudaMemcpy(d_inLinksCount,      inLinksCount,      cidsSize*sizeof(uint32_t), cudaMemcpyHostToDevice);
        cudaMemcpy(d_inLinksOuts,       inLinksOuts,      linksSize*sizeof(uint64_t), cudaMemcpyHostToDevice);

        getCompressedInLinksCount<<<CUDA_BLOCKS_NUMBER,CUDA_THREAD_BLOCK_SIZE>>>(
            cidsSize, d_inLinksStartIndex, d_inLinksCount, d_inLinksOuts, d_compressedInLinksCount
        );
        /*-------------------------------------------------------------------*/



        // STEP3: Calculate compressed in links start indexes
        /*-------------------------------------------------------------------*/
        uint32_t *compressedInLinksCount = (uint32_t*) malloc(cidsSize*sizeof(uint32_t));
        uint64_t *compressedInLinksStartIndex = (uint64_t*) malloc(cidsSize*sizeof(uint64_t));
        cudaMemcpy(compressedInLinksCount, d_compressedInLinksCount, cidsSize * sizeof(uint32_t), cudaMemcpyDeviceToHost);

        // calculated on cpu
        uint64_t compressedInLinksSize = getLinksStartIndex(
            cidsSize, compressedInLinksCount, compressedInLinksStartIndex
        );

        uint64_t *d_compressedInLinksStartIndex;
        cudaMalloc(&d_compressedInLinksStartIndex, cidsSize*sizeof(uint64_t));
        cudaMemcpy(d_compressedInLinksStartIndex, compressedInLinksStartIndex, cidsSize*sizeof(uint64_t), cudaMemcpyHostToDevice);
        free(compressedInLinksStartIndex);
        /*-------------------------------------------------------------------*/



        // STEP4: Calculate compressed in links
        /*-------------------------------------------------------------------*/
        uint64_t *d_inLinksUsers;
        CompressedInLink *d_compressedInLinks; //calculated

        cudaMalloc(&d_inLinksUsers,                   linksSize*sizeof(uint64_t));
        cudaMalloc(&d_compressedInLinks,  compressedInLinksSize*sizeof(CompressedInLink));
        cudaMemcpy(d_inLinksUsers, inLinksUsers,      linksSize*sizeof(uint64_t), cudaMemcpyHostToDevice);

        getCompressedInLinks<<<CUDA_BLOCKS_NUMBER,CUDA_THREAD_BLOCK_SIZE>>>(
            cidsSize,
            d_inLinksStartIndex, d_inLinksCount, d_cidsTotalOutStakes,
            d_inLinksOuts, d_inLinksUsers, d_stakes,
            d_compressedInLinksStartIndex, d_compressedInLinksCount,
            d_compressedInLinks
        );

        cudaFree(d_inLinksUsers);
        cudaFree(d_inLinksStartIndex);
        cudaFree(d_inLinksCount);
        cudaFree(d_inLinksOuts);
        cudaFree(d_stakes);
        cudaFree(d_cidsTotalOutStakes);
        /*-------------------------------------------------------------------*/



        // STEP5: Calculate dangling nodes rank, and default rank
        /*-------------------------------------------------------------------*/
        double defaultRank = (1.0 - DUMP_FACTOR) / cidsSize;
        uint64_t danglingNodesSize = 0;
        for(uint64_t i=0; i< cidsSize; i++){
            rank[i] = defaultRank;
            if(inLinksCount[i] == 0) {
                danglingNodesSize++;
            }
        }

        double innerProductOverSize = defaultRank * ((double) danglingNodesSize / (double)cidsSize);
        double defaultRankWithCorrection = (DUMP_FACTOR * innerProductOverSize) + defaultRank; //fma point
        /*-------------------------------------------------------------------*/




        // STEP6: Calculate rank
        /*-------------------------------------------------------------------*/
        double *d_rank, *d_prevRank;

        cudaMalloc(&d_rank, cidsSize*sizeof(double));
        cudaMalloc(&d_prevRank, cidsSize*sizeof(double));

        cudaMemcpy(d_rank,     rank, cidsSize*sizeof(double), cudaMemcpyHostToDevice);
        cudaMemcpy(d_prevRank, rank, cidsSize*sizeof(double), cudaMemcpyHostToDevice);

        int steps = 0;
        double change = TOLERANCE + 1.0;
        while(change > TOLERANCE) {
            swap(d_rank, d_prevRank);
            steps++;
        	run_rank_iteration<<<CUDA_BLOCKS_NUMBER,CUDA_THREAD_BLOCK_SIZE>>>(
                d_compressedInLinks,
                d_prevRank, d_rank, cidsSize,
                d_compressedInLinksStartIndex, d_compressedInLinksCount,
                defaultRankWithCorrection
        	);
        	change = find_max_ranks_diff(d_prevRank, d_rank, cidsSize);
        	cudaDeviceSynchronize();
        }

        cudaMemcpy(rank, d_rank, cidsSize * sizeof(double), cudaMemcpyDeviceToHost);
        /*-------------------------------------------------------------------*/


        cudaFree(d_rank);
        cudaFree(d_prevRank);
        cudaFree(d_compressedInLinksStartIndex);
        cudaFree(d_compressedInLinksCount);
        cudaFree(d_compressedInLinks);
    }
};
