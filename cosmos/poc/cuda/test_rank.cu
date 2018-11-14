#include "rank.cu"
#include <stdint.h>



void test_getCompressedInLinksStartIndex() {

    uint32_t compressedInLinksCount [6] = { 0, 2, 0, 40, 13, 0 };
    uint64_t compressedInLinksStartIndex [6] = { };
    getCompressedInLinksStartIndex(6, compressedInLinksCount, compressedInLinksStartIndex);

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

    cudaMallocManaged(&dev_inLinksCount, cidsSize*sizeof(uint32_t));
    cudaMallocManaged(&dev_compressedInLinksCount, cidsSize*sizeof(uint32_t));
    cudaMallocManaged(&dev_inLinksStartIndex, cidsSize*sizeof(uint64_t));
    cudaMallocManaged(&dev_inLinksOuts, outSize*sizeof(uint64_t));

    cudaMemcpy(dev_inLinksCount, inLinksCount, cidsSize*sizeof(uint32_t), cudaMemcpyHostToDevice);
    cudaMemcpy(dev_inLinksStartIndex, inLinksStartIndex, cidsSize*sizeof(uint64_t), cudaMemcpyHostToDevice);
    cudaMemcpy(dev_inLinksOuts, inLinksOuts, outSize*sizeof(uint64_t), cudaMemcpyHostToDevice);

    cudaDeviceSynchronize();
    getCompressedInLinksCount<<<1,6>>>(
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

// To run use `nvcc test_rank.cu -o test && ./test` command.
int main(void) {
    printf("Start testing !!!!!!!!!!!!!!!!!!\n");
    test_getCompressedInLinksStartIndex();
    test_getCompressedInLinksCount();
}