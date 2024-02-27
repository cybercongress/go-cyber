#include "rank.cu"
#include <stdint.h>



void test_getCompressedInLinksStartIndex() {

    uint32_t compressedInLinksCount [6] = { 0, 2, 0, 40, 13, 0 };
    uint64_t compressedInLinksStartIndex [6] = { };
    uint64_t size = get_links_start_index(6, compressedInLinksCount, compressedInLinksStartIndex);

    if (size != 55) {
        printf("getCompressedInLinksStartIndex() wrong composed in links size!\n");
    }

    uint64_t expected [6] = {0,0,2,2,42,55};
    if (std::equal(std::begin(expected), std::end(expected), std::begin(compressedInLinksStartIndex)))
        printf("getCompressedInLinksStartIndex() works as expected!\n");
    else {
        printf("getCompressedInLinksStartIndex() doesn't works :(\n");
        for (int i = sizeof(expected) / sizeof(expected[0])-1; i >= 0; i--)
            std::cout << compressedInLinksStartIndex[i] << ' ' << expected[i] << '\n';
    }
}

void test_getCompressedInLinksCount() {

    uint64_t cidsSize = 6;
    uint32_t inLinksCount [6] = { 0, 2, 0, 1, 3, 3 };
    uint64_t inLinksStartIndex [6] = { 0, 0, 2, 2, 3, 6 };
    uint64_t inLinksOuts [] = { 1, 1, 2, 2, 2, 1, 2, 1, 1};
    int outSize = (sizeof(inLinksOuts)/sizeof(*inLinksOuts));

    uint32_t *dev_inLinksCount;
    uint32_t *dev_compressedInLinksCount;
    uint64_t *dev_inLinksStartIndex;
    uint64_t *dev_inLinksOuts;

    cudaMalloc(&dev_inLinksCount, cidsSize*sizeof(uint32_t));
    cudaMalloc(&dev_compressedInLinksCount, cidsSize*sizeof(uint32_t));
    cudaMalloc(&dev_inLinksStartIndex, cidsSize*sizeof(uint64_t));
    cudaMalloc(&dev_inLinksOuts, outSize*sizeof(uint64_t));

    cudaMemcpy(dev_inLinksCount, inLinksCount, cidsSize*sizeof(uint32_t), cudaMemcpyHostToDevice);
    cudaMemcpy(dev_inLinksStartIndex, inLinksStartIndex, cidsSize*sizeof(uint64_t), cudaMemcpyHostToDevice);
    cudaMemcpy(dev_inLinksOuts, inLinksOuts, outSize*sizeof(uint64_t), cudaMemcpyHostToDevice);

    cudaDeviceSynchronize();
    get_compressed_in_links_count<<<2,3>>>(
        cidsSize,
        dev_inLinksStartIndex, dev_inLinksCount,
        dev_inLinksOuts, dev_compressedInLinksCount
    );
    cudaDeviceSynchronize();

    uint32_t actual[6] = {};
    cudaMemcpy(actual, dev_compressedInLinksCount, cidsSize*sizeof(uint32_t), cudaMemcpyDeviceToHost);

    uint64_t expected[6] = {0,1,0,1,2,2};
    if (std::equal(std::begin(expected), std::end(expected), std::begin(actual)))
        printf("getCompressedInLinksCount() works as expected!\n");
    else {
       printf("getCompressedInLinksCount() doesn't works :(\n");
       for (int i = sizeof(actual) / sizeof(actual[0])-1; i >= 0; i--)
           std::cout << actual[i] << ' ' << expected[i] << '\n';
    }
}

void test_calculateCidTotalOutStake() {

    int cidsSize = 6;
    int linksSize = 9;
    int usersSize = 3;

    uint32_t outLinksCount [6] = { 0, 2, 0, 1, 3, 3 };
    uint64_t outLinksStartIndex [6] = { 0, 0, 2, 2, 3, 6 };
    uint64_t outLinksUsers [9] = { 1, 0, 2, 0, 2, 1, 2, 1, 0};
    uint64_t stakes [3] = { 1, 2, 3};

    uint32_t *dev_outLinksCount;
    uint64_t *dev_outLinksStartIndex;
    uint64_t *dev_outLinksUsers;
    uint64_t *dev_stakes;
    uint64_t *dev_cidsTotalOutStakes;

    cudaMalloc(&dev_outLinksCount, cidsSize*sizeof(uint32_t));
    cudaMalloc(&dev_outLinksStartIndex, cidsSize*sizeof(uint64_t));
    cudaMalloc(&dev_outLinksUsers, linksSize*sizeof(uint64_t));
    cudaMalloc(&dev_stakes, usersSize*sizeof(uint64_t));
    cudaMalloc(&dev_cidsTotalOutStakes, cidsSize*sizeof(uint64_t));

    cudaMemcpy(dev_outLinksCount, outLinksCount, cidsSize*sizeof(uint32_t), cudaMemcpyHostToDevice);
    cudaMemcpy(dev_outLinksStartIndex, outLinksStartIndex, cidsSize*sizeof(uint64_t), cudaMemcpyHostToDevice);
    cudaMemcpy(dev_outLinksUsers, outLinksUsers, linksSize*sizeof(uint64_t), cudaMemcpyHostToDevice);
    cudaMemcpy(dev_stakes, stakes, usersSize*sizeof(uint64_t), cudaMemcpyHostToDevice);

    cudaDeviceSynchronize();

    get_particle_stake_by_links<<<2,3>>>(
        cidsSize, dev_stakes,
        dev_outLinksStartIndex, dev_outLinksCount,
        dev_outLinksUsers, dev_cidsTotalOutStakes
    );
    cudaDeviceSynchronize();

    uint64_t actual[6] = {};
    cudaMemcpy(actual, dev_cidsTotalOutStakes, cidsSize*sizeof(uint64_t), cudaMemcpyDeviceToHost);

    uint64_t expected[6] = {0,3,0,3,6,6};
    if (std::equal(std::begin(expected), std::end(expected), std::begin(actual)))
        printf("calculateCidTotalOutStake() works as expected!\n");
    else {
       printf("calculateCidTotalOutStake() doesn't works :(\n");
       for (int i = sizeof(actual) / sizeof(actual[0])-1; i >= 0; i--)
           std::cout << actual[i] << ' ' << expected[i] << '\n';
    }
}

void test_find_max_ranks_diff() {

    double prevRank [6] = { -1.324, 32.1, 0.001, 2.231, -3.22, -0.02 };
    double newRank [6] = {1.3242, 32.22, 0.032, 2.231, -3.232, 0.02 };

    double *dev_prevRank;
    double *dev_newRank;
    cudaMalloc(&dev_prevRank, 6*sizeof(double));
    cudaMalloc(&dev_newRank, 6*sizeof(double));
    cudaMemcpy(dev_prevRank, prevRank, 6*sizeof(double), cudaMemcpyHostToDevice);
    cudaMemcpy(dev_newRank, newRank, 6*sizeof(double), cudaMemcpyHostToDevice);

    double maxDiff = find_max_ranks_diff(dev_prevRank, dev_newRank, 6);
    if (maxDiff == 2.6482)
        printf("find_max_ranks_diff() works as expected!\n");
    else {
       printf("find_max_ranks_diff() doesn't works :(\n");
       std::cout << maxDiff << ' ' << 2.6482 << '\n';
    }
}

void test_getCompressedInLinks() {

    int cidsSize = 8;
    int linksSize = 11;
    int compressedLinksSize = 8;
    int usersSize = 3;

    uint32_t inLinksCount [8] =           {0,0,1,5,4,0,1,0};
    uint32_t compressedInLinksCount [8] = {0,0,1,3,3,0,1,0};
    uint64_t inLinksStartIndex [8] =                {0,0,0,1,6,10,10,11};
    uint64_t compressedInLinksStartIndex [8] =      {0,0,0,1,4,7,7,8};
    uint64_t cidsTotalOutStakes [8] =    {3,3,3,1,6,1,0,3};
    uint64_t inLinksOuts [11]  = {7,1,4,4,4,2,5,0,0,1,3};
    uint64_t inLinksUsers [11] = {0,2,0,1,2,0,1,1,2,1,1};
    uint64_t stakes [3] = {3,1,2};

    uint64_t *dev_inLinksStartIndex;
    uint32_t *dev_inLinksCount;
    uint64_t *dev_cidsTotalOutStakes;
    uint64_t *dev_inLinksOuts;
    uint64_t *dev_inLinksUsers;
    uint64_t *dev_stakes;
    uint64_t *dev_compressedInLinksStartIndex;
    uint32_t *dev_compressedInLinksCount;
    CompressedInLink *dev_compressedInLinks;

    cudaMalloc(&dev_inLinksStartIndex, cidsSize*sizeof(uint64_t));
    cudaMalloc(&dev_inLinksCount, cidsSize*sizeof(uint32_t));
    cudaMalloc(&dev_cidsTotalOutStakes, cidsSize*sizeof(uint64_t));
    cudaMalloc(&dev_inLinksOuts, linksSize*sizeof(uint64_t));
    cudaMalloc(&dev_inLinksUsers, linksSize*sizeof(uint64_t));
    cudaMalloc(&dev_stakes, usersSize*sizeof(uint64_t));
    cudaMalloc(&dev_compressedInLinksStartIndex, cidsSize*sizeof(uint64_t));
    cudaMalloc(&dev_compressedInLinksCount, cidsSize*sizeof(uint32_t));
    cudaMalloc(&dev_compressedInLinks, compressedLinksSize*sizeof(CompressedInLink));

    cudaMemcpy(dev_inLinksStartIndex, inLinksStartIndex, cidsSize*sizeof(uint64_t), cudaMemcpyHostToDevice);
    cudaMemcpy(dev_inLinksCount, inLinksCount, cidsSize*sizeof(uint32_t), cudaMemcpyHostToDevice);
    cudaMemcpy(dev_cidsTotalOutStakes, cidsTotalOutStakes, cidsSize*sizeof(uint64_t), cudaMemcpyHostToDevice);
    cudaMemcpy(dev_inLinksOuts, inLinksOuts, linksSize*sizeof(uint64_t), cudaMemcpyHostToDevice);
    cudaMemcpy(dev_inLinksUsers, inLinksUsers, linksSize*sizeof(uint64_t), cudaMemcpyHostToDevice);
    cudaMemcpy(dev_stakes, stakes, usersSize*sizeof(uint64_t), cudaMemcpyHostToDevice);
    cudaMemcpy(dev_compressedInLinksStartIndex, compressedInLinksStartIndex, cidsSize*sizeof(uint64_t), cudaMemcpyHostToDevice);
    cudaMemcpy(dev_compressedInLinksCount, compressedInLinksCount, cidsSize*sizeof(uint32_t), cudaMemcpyHostToDevice);

    cudaDeviceSynchronize();
    get_compressed_in_links<<<4,2>>>(
        cidsSize,
        dev_inLinksStartIndex, dev_inLinksCount, dev_cidsTotalOutStakes,
        dev_inLinksOuts, dev_inLinksUsers,
        dev_stakes,
        dev_compressedInLinksStartIndex, compressedInLinksCount,
        dev_compressedInLinks
    );
    cudaDeviceSynchronize();

    CompressedInLink actual[8] = {};
    cudaMemcpy(actual, dev_compressedInLinks, compressedLinksSize*sizeof(CompressedInLink), cudaMemcpyDeviceToHost);

    CompressedInLink expected[8] = {
        {7,1.0},{1,0.666667},{4,1},{2,1},{5,1},{0,1},{1,0.333333},{3,1}
    };

    printf("calculateCidTotalOutStake() output\n");
    for (int i = sizeof(actual) / sizeof(actual[0])-1; i >= 0; i--) {
       std::cout << actual[i].fromIndex <<'_'<< actual[i].weight << "   ";
       std::cout << expected[i].fromIndex <<'_'<< expected[i].weight << '\n';
    }
}

// To run use `nvcc test_rank.cu -o test && ./test && rm test` command.
int main(void) {
    printf("Start testing !!!!!!!!!!!!!!!!!!\n");
    test_getCompressedInLinksStartIndex();
    test_getCompressedInLinksCount();
    test_calculateCidTotalOutStake();
    test_find_max_ranks_diff();
    test_getCompressedInLinks();
}