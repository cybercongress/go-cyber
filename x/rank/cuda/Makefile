# NVCC is path to nvcc. Here it is assumed /usr/local/cuda is on one's PATH.

NVCC = nvcc

BUILD_DIR = build

NVCCFLAGS = -fmad=false --compiler-options '-fPIC -frounding-math -fsignaling-nans'

all: build

build: build_dir gpu install_gpu

build_dir:
	mkdir -p $(BUILD_DIR)

build_gpu:
	$(NVCC) $(NVCCFLAGS) -o $(BUILD_DIR)/libcbdrank.so -shared ./cuda/rank.cu

install_gpu:
	sudo cp $(BUILD_DIR)/libcbdrank.so /usr/lib
	sudo cp ./cuda/cbdrank.h /usr/lib`

clean:
	rm $(BUILD_DIR)