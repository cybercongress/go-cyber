###########################################################################################
# Build cyber
###########################################################################################
FROM ubuntu:20.04  as build_stage_cuda

ENV GO_VERSION '1.17.2'
ENV GO_ARCH 'linux-amd64'
ENV GO_BIN_SHA 'f242a9db6a0ad1846de7b6d94d507915d14062660616a61ef7c808a76e4f1676'
ENV DAEMON_HOME /root/.cyber
ENV DAEMON_RESTART_AFTER_UPGRADE=true
ENV DAEMON_LOG_BUFFER_SIZE=1048
ENV UNSAFE_SKIP_BACKUP=true
ENV DAEMON_NAME cyber
ENV BUILD_DIR /build
ENV COSMWASM_VER "1.0.0-beta9"
ENV PATH /usr/local/go/bin:/root/.cargo/bin:/root/cargo/env:/root/.cyber/scripts:$PATH


# Install required dev tools to compile cyberd
###########################################################################################
RUN apt-get update && apt-get install -y --no-install-recommends wget git ca-certificates

# Install golang
###########################################################################################
RUN wget -O go.tgz https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    echo "${GO_BIN_SHA} *go.tgz" | sha256sum -c - && \
    tar -C /usr/local -xzf go.tgz &&\
    rm go.tgz

ENV PATH="/usr/local/go/bin:$PATH"
RUN apt-get -y install --no-install-recommends \
    make gcc g++ \
    wget \
    curl \
    git \
    nvidia-cuda-toolkit \
&& go version

# Create appropriate folders layout
###########################################################################################
 RUN mkdir -p /cyber/cosmovisor/genesis/bin
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
WORKDIR /
RUN curl https://sh.rustup.rs -sSf | sh -s -- -y \
 && wget --quiet https://github.com/CosmWasm/wasmvm/archive/v${COSMWASM_VER}.tar.gz -P /tmp \
 && tar xzf /tmp/v${COSMWASM_VER}.tar.gz -C $BUILD_DIR \
 && cd $BUILD_DIR/wasmvm-${COSMWASM_VER}/ && make build \
 && cp $BUILD_DIR/wasmvm-${COSMWASM_VER}/api/libwasmvm.so /usr/lib/ \
 && cp $BUILD_DIR/wasmvm-${COSMWASM_VER}/api/libwasmvm.dylib /usr/lib/

# Compile cuda kernel
###########################################################################################
COPY . /sources
WORKDIR /sources/x/rank/cuda
RUN make build
RUN cp ./build/libcbdrank.so /usr/lib/ && cp cbdrank.h /usr/lib/

# Compile cyberd for genesis version
###########################################################################################

WORKDIR /sources
RUN make build CUDA_ENABLED=true \
 && chmod +x ./build/cyber \
 && cp ./build/cyber /cyber/cosmovisor/genesis/bin/ \
 && cp ./build/cyber /cyber/ \ 
 && rm -rf ./build \
 && git reset --hard

###########################################################################################
# Create runtime cyber image
###########################################################################################
FROM ubuntu:20.04

ENV DAEMON_HOME /root/.cyber
ENV DAEMON_RESTART_AFTER_UPGRADE=true
ENV GAIA_HOME ${DAEMON_HOME}
ENV DAEMON_NAME cyber

# Install useful dev tools
###########################################################################################
RUN apt-get update && apt-get install -y --no-install-recommends wget curl ca-certificates 

# Download genesis file and links file from IPFS
###########################################################################################
RUN wget -O /genesis.json https://gateway.ipfs.cybernode.ai/ipfs/QmYubyVNfghD4xCrTFj26zBwrF9s5GJhi1TmxvrwmJCipr

WORKDIR /

# Copy compiled kernel and binaries for current bin version
###########################################################################################
COPY --from=build_stage_cuda /cyber /cyber

COPY --from=build_stage_cuda /cyber/cyber /usr/bin

COPY --from=build_stage_cuda /usr/bin/cosmovisor /usr/bin/cosmovisor

COPY --from=build_stage_cuda /usr/lib/cbdrank.h /usr/lib/cbdrank.h
COPY --from=build_stage_cuda /usr/lib/libcbdrank.so /usr/lib/libcbdrank.so

COPY --from=build_stage_cuda /usr/lib/libwasmvm.so /usr/lib/libwasmvm.so
COPY --from=build_stage_cuda /usr/lib/libwasmvm.dylib /usr/lib/libwasmvm.dylib

# Copy startup scripts
###########################################################################################

COPY start_script.sh start_script.sh
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x start_script.sh
RUN chmod +x /entrypoint.sh


#  Start
###############################################################################
EXPOSE 26656 26657 1317 9090
ENTRYPOINT ["/entrypoint.sh"]
CMD ["./start_script.sh"]
