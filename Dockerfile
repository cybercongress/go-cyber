###########################################################################################
# Build cyber
###########################################################################################
FROM nvidia/cuda:10.2-base as build_stage_cuda

ENV GO_VERSION 1.16.2
ENV GO_ARCH 'linux-amd64'
ENV GO_BIN_SHA '542e936b19542e62679766194364f45141fde55169db2d8d01046555ca9eb4b8'
ENV DAEMON_HOME /root/.cyber
ENV DAEMON_RESTART_AFTER_UPGRADE=on
ENV GAIA_HOME ${DAEMON_HOME}
ENV DAEMON_NAME cyber
ENV BUILD_DIR /build
ENV COSMWASM_VER "0.13.0"
ENV PATH /usr/local/go/bin:/root/.cargo/bin:/root/.cyber/scripts:$PATH


# Install required dev tools to compile cyberd
###########################################################################################
RUN apt-get update && apt-get install -y --no-install-recommends wget git

# Install golang
###########################################################################################
RUN url="https://golang.org/dl/go${GO_VERSION}.${GO_ARCH}.tar.gz" && \
    wget -O go.tgz "$url" && \
    echo "${GO_BIN_SHA} *go.tgz" | sha256sum -c - && \
    tar -C /usr/local -xzf go.tgz &&\
    rm go.tgz

ENV PATH="/usr/local/go/bin:$PATH"
RUN apt-get -y install --no-install-recommends \
    ca-certificates \
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
 RUN git clone https://github.com/cosmos/cosmos-sdk.git $BUILD_DIR \
 && cd $BUILD_DIR/cosmovisor/ \
 && make cosmovisor \
 && cp cosmovisor /usr/bin/cosmovisor \
 && chmod +x /usr/bin/cosmovisor \
 && rm -fR $BUILD_DIR/* && rm -fR $BUILD_DIR/.*[a-z]


# Compile cuda kernel
###########################################################################################
COPY . /sources
WORKDIR /sources/x/rank/cuda
RUN make build
RUN cp ./build/libcbdrank.so /usr/lib/ && cp cbdrank.h /usr/lib/

###########################################################################################
# Build wasmvm
###########################################################################################

FROM rustlang/rust:nightly as build_stage_rust
RUN wget --quiet https://github.com/CosmWasm/wasmvm/archive/v${COSMWASM_VER}.tar.gz -P /tmp \
 && tar xzf /tmp/v${COSMWASM_VER}.tar.gz -C $BUILD_DIR \
 && cd $BUILD_DIR/wasmvm-${COSMWASM_VER}/ && make build \
 && cp $BUILD_DIR/wasmvm-${COSMWASM_VER}/api/libgo_cosmwasm.so /usr/lib/ \
 && cp $BUILD_DIR/wasmvm-${COSMWASM_VER}/api/libgo_cosmwasm.dylib /usr/lib/

# Compile cyberd for genesis version
###########################################################################################

ADD https://github.com/CosmWasm/wasmvm/releases/download/v0.13.0/libwasmvm_muslc.a /lib/libwasmvm_muslc.a
RUN sha256sum /lib/libwasmvm_muslc.a | grep 39dc389cc6b556280cbeaebeda2b62cf884993137b83f90d1398ac47d09d3900

WORKDIR /sources
# TODO: Update brach to master before merge\relaese
RUN git checkout bostrom-dev
RUN make build
COPY ./build/cyber /cyber/cosmovisor/genesis/bin



###########################################################################################
# Create runtime cyber image
###########################################################################################
FROM nvidia/cuda:10.2-base

ENV DAEMON_HOME /root/.cyber
ENV DAEMON_RESTART_AFTER_UPGRADE=on
ENV GAIA_HOME ${DAEMON_HOME}
ENV DAEMON_NAME cyber

# Install useful dev tools
###########################################################################################
RUN apt-get update && apt-get install -y --no-install-recommends wget curl

# Download genesis file and links file from IPFS
###########################################################################################
# To slow using ipget, currently we use gateway
# PUT needed CID_OF_GENESIS and CID_OF_CONFIG here
RUN wget -O /genesis.json https://ipfs.io/ipfs/QmYS6UdU1VFqkYjQDahxiESntk8pFCZkwMANqf1a8ZNDfy

WORKDIR /

# Copy compiled kernel and binaries for current bin version
###########################################################################################
COPY --from=build_stage_cuda /cyber /cyber

COPY --from=build_stage_cuda /usr/bin/cosmovisor /usr/local/bin/cosmovisor

COPY --from=build_stage_cuda /usr/lib/cbdrank.h /usr/lib/cbdrank.h
COPY --from=build_stage_cuda /usr/lib/libcbdrank.so /usr/lib/libcbdrank.so

COPY --from=build_stage_rust /usr/lib/libgo_cosmwasm.so /usr/lib/libgo_cosmwasm.so
COPY --from=build_stage_rust /usr/lib/libgo_cosmwasm.dylib /usr/lib/libgo_cosmwasm.dylib

# Copy startup scripts
###########################################################################################

COPY start_script.sh /start_script.sh
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /start_script.sh
RUN chmod +x /entrypoint.sh

# Cleanup for runtime container
###########################################################################################
#RUN cd / && rm -fR $BUILD_DIR $HOME/.rustup $HOME/.cargo $HOME/.cache $HOME/go 
#&& apt-get remove -y --auto-remove \
#ca-certificates \
#make gcc g++ \
#curl \
#git \
#nvidia-cuda-toolkit \
#; \
#rm -rf /var/lib/apt/lists/*; \
#rm -fR /tmp/*;

# Start
###########################################################################################
EXPOSE 26656 26657 1317
ENTRYPOINT ["/entrypoint.sh"]
CMD ["./start_script.sh"]
