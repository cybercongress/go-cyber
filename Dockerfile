###########################################################################################
# Build cyber
###########################################################################################
FROM ubuntu:20.04  as build_stage_cuda

ENV GO_VERSION '1.17.8'
ENV GO_ARCH 'linux-amd64'
ENV GO_BIN_SHA '980e65a863377e69fd9b67df9d8395fd8e93858e7a24c9f55803421e453f4f99'
ENV DEBIAN_FRONTEND=noninteractive 
ENV DAEMON_HOME /root/.cyber
ENV DAEMON_RESTART_AFTER_UPGRADE=true
ENV DAEMON_LOG_BUFFER_SIZE=1048
ENV UNSAFE_SKIP_BACKUP=true
ENV DAEMON_NAME cyber
ENV BUILD_DIR /build
ENV COSMWASM_VER "1.0.0-beta9"
ENV PATH /usr/local/go/bin:/root/.cargo/bin:/root/cargo/env:/root/.cyber/scripts:$PATH
#ENV CUDA_VER '11.6.1-1'


# Install required dev tools to compile cyber
###########################################################################################
RUN apt-get update && apt-get install -y --no-install-recommends wget git ca-certificates

# Install golang
###########################################################################################
RUN wget -O go.tgz https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    echo "${GO_BIN_SHA} *go.tgz" | sha256sum -c - && \
    tar -C /usr/local -xzf go.tgz &&\
    rm go.tgz

ENV PATH="/usr/local/go/bin:/usr/local/cuda/bin:$PATH"
RUN apt-get -y install --no-install-recommends \
    make gcc g++ \
    wget \
    curl \
    git \
#    gnupg \
#    software-properties-common \
    nvidia-cuda-toolkit  \
&& go version

# Install requested CUDA version
###########################################################################################
#RUN wget https://developer.download.nvidia.com/compute/cuda/repos/ubuntu2004/x86_64/cuda-ubuntu2004.pin \
#&& mv cuda-ubuntu2004.pin /etc/apt/preferences.d/cuda-repository-pin-600 \
#&& apt-key adv --fetch-keys https://developer.download.nvidia.com/compute/cuda/repos/ubuntu2004/x86_64/7fa2af80.pub \
#&& add-apt-repository "deb https://developer.download.nvidia.com/compute/cuda/repos/ubuntu2004/x86_64/ /" \
#&& apt-get update \
#&& apt-get install cuda=${CUDA_VER} -y 

# Create appropriate folders layout
###########################################################################################
 RUN mkdir -p /cyber/cosmovisor/genesis/bin \
 && mkdir -p /cyber/cosmovisor/upgrades/upgrade-1/bin
# Compile cosmovisor
###########################################################################################
 RUN git clone --depth 1 https://github.com/cosmos/cosmos-sdk.git $BUILD_DIR/ \
 && cd $BUILD_DIR/cosmovisor/ \
 && make cosmovisor \
 && cp cosmovisor /usr/bin/cosmovisor \
 && chmod +x /usr/bin/cosmovisor \
 && rm -fR $BUILD_DIR/* && rm -fR $BUILD_DIR/.*[a-z]

###########################################################################################
# Build wasmvm
###########################################################################################
#WORKDIR /
#RUN curl https://sh.rustup.rs -sSf | sh -s -- -y \
# && wget --quiet https://github.com/CosmWasm/wasmvm/archive/v${COSMWASM_VER}.tar.gz -P /tmp \
# && tar xzf /tmp/v${COSMWASM_VER}.tar.gz -C $BUILD_DIR \
# && cd $BUILD_DIR/wasmvm-${COSMWASM_VER}/ && make build \
# && cp $BUILD_DIR/wasmvm-${COSMWASM_VER}/api/libwasmvm.so /usr/lib/ \
# && cp $BUILD_DIR/wasmvm-${COSMWASM_VER}/api/libwasmvm.dylib /usr/lib/

# Compile cuda kernel
###########################################################################################
COPY . /sources
WORKDIR /sources/x/rank/cuda
RUN make build
RUN cp ./build/libcbdrank.so /usr/lib/ && cp cbdrank.h /usr/lib/

# Compile cyber for genesis version
###########################################################################################

WORKDIR /sources
RUN git checkout v0.2.0 \
 && make build CUDA_ENABLED=true \
 && chmod +x ./build/cyber \
 && cp ./build/cyber /cyber/cosmovisor/genesis/bin/ \
 && cp ./build/cyber /usr/local/bin \ 
 && rm -rf ./build \
 && git reset --hard


 # Compile cyber for genesis version
###########################################################################################

WORKDIR /sources
RUN git checkout upgrade-1 \
 && make build CUDA_ENABLED=true \
 && chmod +x ./build/cyber \
 && cp ./build/cyber /cyber/cosmovisor/upgrades/upgrade-1/bin/ \
 && rm -rf ./build \
 && git reset --hard

###########################################################################################
# Create runtime cyber image
###########################################################################################
#FROM ubuntu:20.04

ENV DAEMON_HOME /root/.cyber
ENV DAEMON_RESTART_AFTER_UPGRADE=true
ENV DAEMON_NAME cyber
ENV DAEMON_ALLOW_DOWNLOAD_BINARIES=false
ENV DAEMON_LOG_BUFFER_SIZE=812
ENV UNSAFE_SKIP_BACKUP=true

# Install useful dev tools
###########################################################################################
#RUN apt-get update && apt-get install -y --no-install-recommends wget curl ca-certificates 

# Download genesis file and links file from IPFS
###########################################################################################
RUN wget -O /genesis.json https://gateway.ipfs.cybernode.ai/ipfs/QmYe81dBfxgYsVhX1mX4uiLQyu3jx2kJvR6CgytDnFgKzc

WORKDIR /

# Copy compiled kernel and binaries for current bin version
###########################################################################################
#COPY --from=build_stage_cuda /cyber /cyber

#COPY /cyber/cyber /usr/bin

#COPY --from=build_stage_cuda /usr/bin/cosmovisor /usr/bin/cosmovisor

#COPY --from=build_stage_cuda /usr/lib/cbdrank.h /usr/lib/cbdrank.h
#COPY --from=build_stage_cuda /usr/lib/libcbdrank.so /usr/lib/libcbdrank.so

#COPY --from=build_stage_cuda /usr/lib/libwasmvm.so /usr/lib/libwasmvm.so
#COPY --from=build_stage_cuda /usr/lib/libwasmvm.dylib /usr/lib/libwasmvm.dylib

# Copy startup scripts
###########################################################################################

COPY start_script.sh start_script.sh
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x start_script.sh
RUN chmod +x /entrypoint.sh
RUN go version
RUN cyber version


#  Start
###############################################################################
EXPOSE 26656 26657 1317 9090
ENTRYPOINT ["/entrypoint.sh"]
CMD ["./start_script.sh"]
