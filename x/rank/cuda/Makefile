# NVCC is path to nvcc. Here it is assumed /usr/local/cuda is on one's PATH.

NVCC = nvcc

BUILD_DIR = build

NVCCFLAGS = -fmad=false --compiler-options '-fPIC -frounding-math -fsignaling-nans'

all: build

build: build_dir gpu

build_dir:
	mkdir -p $(BUILD_DIR)

gpu:
	$(NVCC) $(NVCCFLAGS) -o $(BUILD_DIR)/libcbdrank.so -shared rank.cu

clean:
	rm $(BUILD_DIR)/libcbdrank.so