#include <stdint.h>
#include <stdio.h>
#include <thrust/transform.h>
#include <thrust/transform_reduce.h>
#include <thrust/device_vector.h>
#include <thrust/execution_policy.h>
#include <thrust/functional.h>
#include "types.h"

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
    double defaultRankWithCorrection,                     /* default rank + inner product correction */
    double dampingFactor
) {

	int index = blockIdx.x * blockDim.x + threadIdx.x;
    uint64_t stride = blockDim.x * gridDim.x;

    for (uint64_t i = index; i < rankSize; i += stride) {
        if(inLinksCount[i] == 0) { continue; }

        double ksum = 0;
        for (uint64_t j = inLinksStartIndex[i]; j < inLinksStartIndex[i] + inLinksCount[i]; j++) {
           if (inLinks[j].weight == 0) {continue;}
           ksum = prevRank[inLinks[j].fromIndex] * inLinks[j].weight + ksum;
        }
        rank[i] = ksum * dampingFactor + defaultRankWithCorrection;
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

/*******************************************************/
/* KERNEL: CALCULATE PARTICLE STAKE BY IN OR OUT LINKS */
/*******************************************************/
__global__
void get_particle_stake_by_links(
    uint64_t cidsSize,
    uint64_t *stakes,                                /*array index - user index*/
    uint64_t *linksStartIndex, uint32_t *linksCount, /*array index - cid index*/
    uint64_t *linksUsers,                            /*all links from all users*/
    /*returns*/ uint64_t *cidsTotalOutStakes         /*array index - cid index*/
) {

	int index = blockIdx.x * blockDim.x + threadIdx.x;
    uint64_t stride = blockDim.x * gridDim.x;

    for (uint64_t i = index; i < cidsSize; i += stride) {
        uint64_t totalOutStake = 0;
        for (uint64_t j = linksStartIndex[i]; j < linksStartIndex[i] + linksCount[i]; j++) {
           totalOutStake = __dadd_rn(totalOutStake, stakes[linksUsers[j]]);
        }
        cidsTotalOutStakes[i] = totalOutStake;
    }
}

/*********************************************************/
/* DEVICE: USER TO DIVIDE TWO uint64                     */
/*********************************************************/
__device__ __forceinline__
double ddiv_rn(uint64_t *a, uint64_t *b) {
    return __ddiv_rn(__ull2double_rn(*a), __ull2double_rn(*b));
}

/*****************************************************/
/* KERNEL: CALCULATE CYBERLINKS WEIGHTS BY STAKE     */
/*****************************************************/
__global__
void get_cyberlinks_weight_by_stake(
    uint64_t cidsSize,
    uint64_t *stakes,                                /*array index - user index*/
    uint64_t *linksStartIndex, uint32_t *linksCount, /*array index - cid index*/
    uint64_t *linksUsers,                            /*all out links from all users*/
    uint64_t *cidsTotalStakes,                       /*array index - cid index*/
    /*returns*/ double *cyberlinksLocalWeights       /*array index - links index*/
) {

	int index = blockIdx.x * blockDim.x + threadIdx.x;
    uint64_t stride = blockDim.x * gridDim.x;

    for (uint64_t i = index; i < cidsSize; i += stride) {
        uint64_t stake = cidsTotalStakes[i];
        for (uint64_t j = linksStartIndex[i]; j < linksStartIndex[i] + linksCount[i]; j++) {
            if (&stakes[linksUsers[j]] == 0 || &stake == 0) { continue; }
            double weight = ddiv_rn(&stakes[linksUsers[j]], &stake);
            if (isnan(weight)) { continue; }
            cyberlinksLocalWeights[j] = weight;
        }
    }
}

/*********************************************************/
/* KERNEL: MULTIPLY TWO ARRAYS                           */
/*********************************************************/
__global__
void multiply_arrays(
    uint64_t size,
    double   *a,
    double   *b,
    double   *output
) {
    int tx = blockIdx.x * blockDim.x + threadIdx.x;
    if (tx < size) output[tx] = __dmul_rn(a[tx], b[tx]);
}

/*************************************************************************/
/* KERNEL: CALCULATE PARTICLE TOTAL STAKE TRANSORMED WITH DAMPING FACTOR */
/*************************************************************************/
__global__
void get_stake_with_damping(
    uint64_t size,
    uint64_t *outStake,
    uint64_t *inStake,
    double   *swd,
    double   damping
) {
    int tx = blockIdx.x * blockDim.x + threadIdx.x;
    if (tx < size) swd[tx] = __dadd_rn(
        __dmul_rn(damping, __ull2double_rn(inStake[tx])),
        __dmul_rn(1-damping, __ull2double_rn(outStake[tx]))
    );
}

/******************************************************************************************/
/* KERNEL: CALCULATE SUM OF ADJACENT PARTICLES STAKE WITH DAMPING BY IN OR OUT CYBERLINKS */
/******************************************************************************************/
__global__
void sum_stake_with_damping_by_links(
    uint64_t cidsSize,
    uint64_t *linksStartIndex, uint32_t *linksCount, /*array index - cid index*/
    uint64_t *linksOuts, // linksIns                 /*all incoming or outgoing links from all users*/
    double *swd,                                     /*array index - cid index*/
    double damping,
    /*returns*/ double *sumswd                       /*array index - cid index*/
) {

	int index = blockIdx.x * blockDim.x + threadIdx.x;
    uint64_t stride = blockDim.x * gridDim.x;

    for (uint64_t i = index; i < cidsSize; i += stride) {
        for (uint64_t j = linksStartIndex[i]; j < linksStartIndex[i] + linksCount[i]; j++) {
            sumswd[i] = __dadd_rn(sumswd[i], __dmul_rn(damping, swd[linksOuts[j]]));
        }
    }
}

/******************************************************************/
/* KERNEL: CALCULATE ENTROPY BY IN OR OUT CYBERLINKS FOR PARTICLE */
/******************************************************************/
__global__
void calculate_entropy_by_links(
    uint64_t cidsSize,
    uint64_t *linksStartIndex, uint32_t *linksCount, /*array index - cid index*/
    uint64_t *linksOuts, // linksIns                 /*all incoming or outgoing links from all users*/
    double *swd,                                     /*array index - cid index*/
    double *d_sumswd,                                /*array index - cid index*/
    /*returns*/ double *entropy                      /*array index - cid index*/
) {

	int index = blockIdx.x * blockDim.x + threadIdx.x;
    uint64_t stride = blockDim.x * gridDim.x;

    for (uint64_t i = index; i < cidsSize; i += stride) {
        for (uint64_t j = linksStartIndex[i]; j < linksStartIndex[i] + linksCount[i]; j++) {
            if (swd[i] == 0 || d_sumswd[linksOuts[j]] == 0) { continue; }
            double weight = __ddiv_rn(swd[i], d_sumswd[linksOuts[j]]);
            if (isnan(weight)) { continue; }
            double logw = log2(weight);
            entropy[i] = __dadd_rn(entropy[i], fabs(__dmul_rn(weight, logw)));
        }
    }
}

/*********************************************************/
/* KERNEL: CALCULATE COMPRESSED IN LINKS COUNT FOR CIDS  */
/*********************************************************/
__global__
void get_compressed_in_links_count(
    uint64_t cidsSize,
    uint64_t *inLinksStartIndex, uint32_t *inLinksCount, /*array index - cid index*/
    uint64_t *inLinksOuts,                               /*all incoming links from all users*/
    /*returns*/ uint32_t *compressedInLinksCount         /*array index - cid index*/
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
/* KERNEL: CALCULATE COMPRESSED IN LINKS                 */
/*********************************************************/
__global__
void get_compressed_in_links(
    uint64_t cidsSize,
    uint64_t *inLinksStartIndex, uint32_t *inLinksCount, uint64_t *cidsTotalOutStakes, /*array index - cid index*/
    uint64_t *inLinksOuts, uint64_t *inLinksUsers,                                     /*all incoming links from all users*/
    uint64_t *stakes,                                                                  /*array index - user index*/
    uint64_t *compressedInLinksStartIndex, uint32_t *compressedInLinksCount,           /*array index - cid index*/
    /*returns*/ CompressedInLink *compressedInLinks                                    /*all incoming compressed links*/
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
            if (isnan(weight)) { weight = 0; }
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
                if (isnan(weight)) { weight = 0; }
                compressedInLinks[compressedLinksIndex] = CompressedInLink {oppositeCid, weight};
                compressedLinksIndex++;
                compressedLinkStake=0;
            }
        }
    }
}

/************************************************************/
/* HOST: CALCULATE KARMA                                    */
/************************************************************/
/* SEQUENTIAL LOGIC -> CALCULATE ON CPU                     */
/* RETURNS KARMA FOR ALL ACCOUNTS                           */
/************************************************************/
__host__
void calculate_karma(
    uint64_t cidsSize,
    uint64_t *outLinksStartIndex, uint32_t *outLinksCount,
    uint64_t *outLinksUsers,
    double   *cyberlinksLocalWeights,
    double   *light,
    /*returns*/ double *karma
) {
    for (uint64_t i = 0; i < cidsSize; i++) {
        for (uint64_t j = outLinksStartIndex[i]; j < outLinksStartIndex[i] + outLinksCount[i]; j++) {
            if (isnan(cyberlinksLocalWeights[j])) { continue; } // !
            karma[outLinksUsers[j]] = light[i]*cyberlinksLocalWeights[j] + karma[outLinksUsers[j]];
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
uint64_t get_links_start_index(
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

void printSize(size_t usageOffset) {
	size_t free = 0, total = 0;
	cudaMemGetInfo(&free, &total);
	fprintf(stderr, "-[GPU]: Free: %.2fMB\tUsed: %.2fMB\n", free / 1048576.0f, (total - usageOffset - free) / 1048576.0f);
}

extern "C" {

    void calculate_rank(
        uint64_t *stakes, uint64_t stakesSize,                    /* User stakes and corresponding array size */
        uint64_t cidsSize, uint64_t linksSize,                    /* Cids count */
        uint32_t *inLinksCount, uint32_t *outLinksCount,          /* array index - cid index*/
        uint64_t *inLinksOuts,
        uint64_t *outLinksIns,
        uint64_t *inLinksUsers,                                   /*all incoming links from all users*/
        uint64_t *outLinksUsers,                                  /*all outgoing links from all users*/
        double dampingFactor,                                     /* value of damping factor*/
        double tolerance,                                         /* value of needed tolerance */
        double *rank,                                             /* array index - cid index*/
        double *entropy,                                          /* array index - cid index*/
        double *light,                                            /* array index - cid index*/
        double *karma                                             /* array index - account index*/
    ) {

        // setbuf(stdout, NULL);
        int CUDA_BLOCKS_NUMBER = (cidsSize + CUDA_THREAD_BLOCK_SIZE - 1) / CUDA_THREAD_BLOCK_SIZE;

        // size_t freeStart = 0, totalStart = 0, usageOffset = 0;
        // cudaMemGetInfo(&freeStart, &totalStart);
        // usageOffset = totalStart - freeStart;
        // fprintf(stderr, "[GPU]: Usage Offset: %.2fMB\n", usageOffset / 1048576.0f);

        // STEP0: Calculate compressed in/out links start indexes
        /*-------------------------------------------------------------------*/
        // calculated on CPU
        // printf("STEP0: Calculate compressed in/out links start indexes\n");

        uint64_t *inLinksStartIndex = (uint64_t*) malloc(cidsSize*sizeof(uint64_t));
        uint64_t *outLinksStartIndex = (uint64_t*) malloc(cidsSize*sizeof(uint64_t));
        get_links_start_index(cidsSize, inLinksCount, inLinksStartIndex);
        get_links_start_index(cidsSize, outLinksCount, outLinksStartIndex);

        // printSize(usageOffset);
        /*-------------------------------------------------------------------*/


        // STEP1.1: Calculate for each particle stake by OUT cyberlinks
        /*-------------------------------------------------------------------*/
        // printf("STEP1.1: Calculate for each particle stake by OUT cyberlinks\n");

        uint64_t *d_outLinksStartIndex;
        uint32_t *d_outLinksCount;
        uint64_t *d_outLinksUsers;
        uint64_t *d_stakes;             // will be used to calculated links weights, should be freed before rank iterations
        uint64_t *d_cidsTotalOutStakes; // will be used to calculated links weights, should be freed before rank iterations

        cudaMalloc(&d_outLinksStartIndex, cidsSize*sizeof(uint64_t));
        cudaMalloc(&d_outLinksCount,      cidsSize*sizeof(uint32_t));
        cudaMalloc(&d_outLinksUsers,     linksSize*sizeof(uint64_t));
        cudaMalloc(&d_stakes,           stakesSize*sizeof(uint64_t));
        cudaMalloc(&d_cidsTotalOutStakes, cidsSize*sizeof(uint64_t)); //calculated

        cudaMemcpy(d_outLinksStartIndex, outLinksStartIndex, cidsSize*sizeof(uint64_t), cudaMemcpyHostToDevice);
        cudaMemcpy(d_outLinksCount,      outLinksCount,      cidsSize*sizeof(uint32_t), cudaMemcpyHostToDevice);
        cudaMemcpy(d_outLinksUsers,      outLinksUsers,     linksSize*sizeof(uint64_t), cudaMemcpyHostToDevice);
        cudaMemcpy(d_stakes,             stakes,           stakesSize*sizeof(uint64_t), cudaMemcpyHostToDevice);

        get_particle_stake_by_links<<<CUDA_BLOCKS_NUMBER,CUDA_THREAD_BLOCK_SIZE>>>(
            cidsSize, d_stakes, d_outLinksStartIndex,
            d_outLinksCount, d_outLinksUsers, d_cidsTotalOutStakes
        );

        // printSize(usageOffset);
        /*-------------------------------------------------------------------*/


        // STEP1.2: Calculate for each particle total stake by IN links
        /*-------------------------------------------------------------------*/
        // printf("STEP1.2: Calculate for each particle stake by IN links\n");

        uint64_t *d_inLinksStartIndex;
        uint32_t *d_inLinksCount;
        uint64_t *d_inLinksUsers;
        uint64_t *d_cidsTotalInStakes; // will be used to calculated links weights, should be freed before rank iterations

        cudaMalloc(&d_inLinksStartIndex, cidsSize*sizeof(uint64_t));
        cudaMalloc(&d_inLinksCount,      cidsSize*sizeof(uint32_t));
        cudaMalloc(&d_inLinksUsers,      linksSize*sizeof(uint64_t));
        cudaMalloc(&d_cidsTotalInStakes, cidsSize*sizeof(uint64_t));   //calculated

        cudaMemcpy(d_inLinksStartIndex, inLinksStartIndex, cidsSize*sizeof(uint64_t), cudaMemcpyHostToDevice);
        cudaMemcpy(d_inLinksCount,      inLinksCount,      cidsSize*sizeof(uint32_t), cudaMemcpyHostToDevice);
        cudaMemcpy(d_inLinksUsers,      inLinksUsers,      linksSize*sizeof(uint64_t), cudaMemcpyHostToDevice);

        get_particle_stake_by_links<<<CUDA_BLOCKS_NUMBER,CUDA_THREAD_BLOCK_SIZE>>>(
            cidsSize, d_stakes, d_inLinksStartIndex,
            d_inLinksCount, d_inLinksUsers, d_cidsTotalInStakes
        );

        // printSize(usageOffset);
       /*-------------------------------------------------------------------*/


        // STEP1.3: Calculate Stake With Damping
        /*-------------------------------------------------------------------*/
        // printf("STEP1.3: Calculate Stake With Damping\n");

        double *d_swd;
        cudaMalloc(&d_swd, cidsSize*sizeof(double));
        cudaMemcpy(d_swd, entropy, cidsSize*sizeof(double), cudaMemcpyHostToDevice);

        get_stake_with_damping<<<CUDA_BLOCKS_NUMBER,CUDA_THREAD_BLOCK_SIZE>>>(
            cidsSize, d_cidsTotalOutStakes, d_cidsTotalInStakes, d_swd, dampingFactor);
        cudaFree(d_cidsTotalInStakes);
        // printSize(usageOffset);
        /*-------------------------------------------------------------------*/


        // STEP1.4: Calculate Local weights
        /*-------------------------------------------------------------------*/
        // printf("STEP1.4: Calculate Local weights\n");

        // local weight for future karma for contributed light
        double *d_cyberlinksLocalWeights;
        cudaMalloc(&d_cyberlinksLocalWeights, linksSize*sizeof(double));

        get_cyberlinks_weight_by_stake<<<CUDA_BLOCKS_NUMBER,CUDA_THREAD_BLOCK_SIZE>>>(
            cidsSize, d_stakes, d_outLinksStartIndex,
            d_outLinksCount, d_outLinksUsers, d_cidsTotalOutStakes, d_cyberlinksLocalWeights
        );
        cudaFree(d_outLinksUsers);
        // printSize(usageOffset);
        /*-------------------------------------------------------------------*/


        // STEP2: Calculate compressed in links count
        /*-------------------------------------------------------------------*/
        // printf("STEP2: Calculate compressed in links count\n");

        uint64_t *d_inLinksOuts;
        uint32_t *d_compressedInLinksCount;

        cudaMalloc(&d_inLinksOuts,           linksSize*sizeof(uint64_t));
        cudaMalloc(&d_compressedInLinksCount, cidsSize*sizeof(uint32_t));   //calculated
        cudaMemcpy(d_inLinksOuts,       inLinksOuts,      linksSize*sizeof(uint64_t), cudaMemcpyHostToDevice);

        get_compressed_in_links_count<<<CUDA_BLOCKS_NUMBER,CUDA_THREAD_BLOCK_SIZE>>>(
            cidsSize, d_inLinksStartIndex, d_inLinksCount, d_inLinksOuts, d_compressedInLinksCount
        );

        // printSize(usageOffset);
        /*-------------------------------------------------------------------*/


        // STEP3: Calculate world entropy
        /*-------------------------------------------------------------------*/
        // printf("STEP3: Calculate world entropy\n");

        double *d_sumswd;
        cudaMalloc(&d_sumswd, cidsSize*sizeof(double));
        cudaMemcpy(d_sumswd, entropy, cidsSize*sizeof(double), cudaMemcpyHostToDevice);

        sum_stake_with_damping_by_links<<<CUDA_BLOCKS_NUMBER,CUDA_THREAD_BLOCK_SIZE>>>(
            cidsSize, d_inLinksStartIndex,
            d_inLinksCount, d_inLinksOuts, d_swd, dampingFactor, d_sumswd);

        uint64_t *d_outLinksIns;
        cudaMalloc(&d_outLinksIns, linksSize*sizeof(uint64_t));
        cudaMemcpy(d_outLinksIns, outLinksIns, linksSize*sizeof(uint64_t), cudaMemcpyHostToDevice);

        sum_stake_with_damping_by_links<<<CUDA_BLOCKS_NUMBER,CUDA_THREAD_BLOCK_SIZE>>>(
            cidsSize, d_outLinksStartIndex,
            d_outLinksCount, d_outLinksIns, d_swd, 1-dampingFactor, d_sumswd);

        // calculate entropy by in/out links
        double *d_entropy;
        cudaMalloc(&d_entropy, cidsSize*sizeof(double));
        cudaMemcpy(d_entropy, entropy, cidsSize*sizeof(double), cudaMemcpyHostToDevice);

        calculate_entropy_by_links<<<CUDA_BLOCKS_NUMBER,CUDA_THREAD_BLOCK_SIZE>>>(
            cidsSize, d_inLinksStartIndex,
            d_inLinksCount, d_inLinksOuts, d_swd, d_sumswd, d_entropy);

        calculate_entropy_by_links<<<CUDA_BLOCKS_NUMBER,CUDA_THREAD_BLOCK_SIZE>>>(
        cidsSize, d_outLinksStartIndex,
        d_outLinksCount, d_outLinksIns, d_swd, d_sumswd, d_entropy);

        cudaFree(d_swd);
        cudaFree(d_sumswd);
        cudaFree(d_outLinksIns);
        cudaFree(d_outLinksStartIndex);
        cudaFree(d_outLinksCount);


        // STEP4: Calculate compressed in links start indexes
        /*-------------------------------------------------------------------*/
        // printf("STEP4: Calculate compressed in links start indexes\n");

        uint32_t *compressedInLinksCount = (uint32_t*) malloc(cidsSize*sizeof(uint32_t));
        uint64_t *compressedInLinksStartIndex = (uint64_t*) malloc(cidsSize*sizeof(uint64_t));
        cudaMemcpy(compressedInLinksCount, d_compressedInLinksCount, cidsSize * sizeof(uint32_t), cudaMemcpyDeviceToHost);

        // calculated on CPU
        uint64_t compressedInLinksSize = get_links_start_index(
            cidsSize, compressedInLinksCount, compressedInLinksStartIndex
        );

        uint64_t *d_compressedInLinksStartIndex;
        cudaMalloc(&d_compressedInLinksStartIndex, cidsSize*sizeof(uint64_t));
        cudaMemcpy(d_compressedInLinksStartIndex, compressedInLinksStartIndex, cidsSize*sizeof(uint64_t), cudaMemcpyHostToDevice);
        free(compressedInLinksStartIndex);

        // printSize(usageOffset);
        /*-------------------------------------------------------------------*/


        // STEP5: Calculate compressed in links
        /*-------------------------------------------------------------------*/
        // printf("STEP5: Calculate compressed in links\n");

        CompressedInLink *d_compressedInLinks; //calculated
        cudaMalloc(&d_compressedInLinks,  compressedInLinksSize*sizeof(CompressedInLink));

        get_compressed_in_links<<<CUDA_BLOCKS_NUMBER,CUDA_THREAD_BLOCK_SIZE>>>(
            cidsSize,
            d_inLinksStartIndex, d_inLinksCount, d_cidsTotalOutStakes,
            d_inLinksOuts, d_inLinksUsers, d_stakes,
            d_compressedInLinksStartIndex, d_compressedInLinksCount,
            d_compressedInLinks
        );

        cudaFree(d_inLinksStartIndex);
        cudaFree(d_inLinksCount);
        cudaFree(d_inLinksUsers);
        cudaFree(d_inLinksOuts);
        cudaFree(d_stakes);
        cudaFree(d_cidsTotalOutStakes);

        // printSize(usageOffset);
        /*-------------------------------------------------------------------*/


        // STEP6: Calculate dangling nodes rank, and default rank
        /*-------------------------------------------------------------------*/
        // printf("STEP6: Calculate dangling nodes rank, and default rank\n");

        double defaultRank = (1.0 - dampingFactor) / cidsSize;
        uint64_t danglingNodesSize = 0;
        for(uint64_t i=0; i< cidsSize; i++){
            rank[i] = defaultRank;
            if(inLinksCount[i] == 0) {
                danglingNodesSize++;
            }
        }

        double innerProductOverSize = defaultRank * ((double) danglingNodesSize / (double)cidsSize);
        double defaultRankWithCorrection = (dampingFactor * innerProductOverSize) + defaultRank; //fma point

        // printSize(usageOffset);
        /*-------------------------------------------------------------------*/


        // STEP7: Calculate Rank
        /*-------------------------------------------------------------------*/
        // printf("STEP7: Calculate Rank\n");

        double *d_rank, *d_prevRank;
        cudaMalloc(&d_rank,     cidsSize*sizeof(double));
        cudaMalloc(&d_prevRank, cidsSize*sizeof(double));
        cudaMemcpy(d_rank,     rank, cidsSize*sizeof(double), cudaMemcpyHostToDevice);
        cudaMemcpy(d_prevRank, rank, cidsSize*sizeof(double), cudaMemcpyHostToDevice);

        int steps = 0;
        double change = tolerance + 1.0;
        while(change > tolerance) {
            swap(d_rank, d_prevRank);
            steps++;
        	run_rank_iteration<<<CUDA_BLOCKS_NUMBER,CUDA_THREAD_BLOCK_SIZE>>>(
                d_compressedInLinks,
                d_prevRank, d_rank, cidsSize,
                d_compressedInLinksStartIndex, d_compressedInLinksCount,
                defaultRankWithCorrection, dampingFactor
        	);
        	change = find_max_ranks_diff(d_prevRank, d_rank, cidsSize);
        	cudaDeviceSynchronize();
        }

        cudaMemcpy(rank, d_rank, cidsSize * sizeof(double), cudaMemcpyDeviceToHost);

        cudaFree(d_prevRank);
        cudaFree(d_compressedInLinksStartIndex);
        cudaFree(d_compressedInLinksCount);
        cudaFree(d_compressedInLinks);

        // printSize(usageOffset);
        /*-------------------------------------------------------------------*/


        // STEP8: Calculate Light
        /*-------------------------------------------------------------------*/
        // printf("STEP8: Calculate Light\n");

        double *d_light;
        cudaMalloc(&d_light, cidsSize*sizeof(double));
        cudaMemcpy(d_light, light, cidsSize*sizeof(double), cudaMemcpyHostToDevice);

        multiply_arrays<<<CUDA_BLOCKS_NUMBER,CUDA_THREAD_BLOCK_SIZE>>>(
            cidsSize, d_rank, d_entropy, d_light
        );

        cudaMemcpy(light, d_light, cidsSize * sizeof(double), cudaMemcpyDeviceToHost);
        cudaMemcpy(entropy, d_entropy, cidsSize * sizeof(double), cudaMemcpyDeviceToHost);

        cudaFree(d_entropy);
        cudaFree(d_rank);
        cudaFree(d_light);

        // printSize(usageOffset);
        /*-------------------------------------------------------------------*/


        // STEP9: Calculate Karma
        /*-------------------------------------------------------------------*/
        // printf("STEP9: Calculate Karma\n");
        // calculated on CPU
        double *cyberlinksLocalWeights = (double*) malloc(linksSize*sizeof(double));
        cudaMemcpy(cyberlinksLocalWeights, d_cyberlinksLocalWeights, linksSize*sizeof(double), cudaMemcpyDeviceToHost);
        cudaFree(d_cyberlinksLocalWeights);

        calculate_karma(
            cidsSize,
            outLinksStartIndex,
            outLinksCount,
            outLinksUsers,
            cyberlinksLocalWeights,
            light,
            karma
        );
        free(cyberlinksLocalWeights);

        // printSize(usageOffset);
        free(inLinksStartIndex);
        free(outLinksStartIndex);
        free(compressedInLinksCount);
    }
};