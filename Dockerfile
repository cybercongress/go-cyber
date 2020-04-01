###############################################################################
# Build cyber
###############################################################################
FROM nvidia/cuda:10.0-devel-ubuntu18.04 as build_stage_cuda

ENV GO_VERSION 1.13.1
ENV GO_ARCH 'linux-amd64'
ENV GO_BIN_SHA '94f874037b82ea5353f4061e543681a0e79657f787437974214629af8407d124'

# Install required dev tools to compile cyberd
###############################################################################
RUN apt-get update && apt-get install -y --no-install-recommends wget git

# Install golang
###############################################################################
RUN url="https://golang.org/dl/go${GO_VERSION}.${GO_ARCH}.tar.gz" && \
	wget -O go.tgz "$url" && \
	echo "${GO_BIN_SHA} *go.tgz" | sha256sum -c - && \
	tar -C /usr/local -xzf go.tgz &&\
	rm go.tgz

ENV PATH="/usr/local/go/bin:$PATH"
RUN go version && nvcc --version

# Compile cuda kernel
###############################################################################
COPY . /sources
WORKDIR /sources/x/rank/cuda
RUN make build
RUN cp ./build/libcbdrank.so /usr/lib/ && cp cbdrank.h /usr/lib/

# Compile cyberd
###############################################################################
WORKDIR /sources
RUN make build



###############################################################################
# Build go-cosmwasm
###############################################################################
FROM rustlang/rust:nightly as build_stage_rust

# Install build dependencies
###############################################################################
RUN apt-get update
RUN apt install -y clang gcc g++ zlib1g-dev libmpc-dev libmpfr-dev libgmp-dev
RUN apt install -y build-essential cmake git

# Install repository
###############################################################################
RUN git clone https://github.com/confio/go-cosmwasm sources

# Compile go-cosmwasm
###############################################################################
WORKDIR /sources
RUN make build-rust-release



###############################################################################
# Create runtime cyber image
###############################################################################
FROM nvidia/cuda:10.0-runtime-ubuntu18.04

# Install useful dev tools
###############################################################################
RUN apt-get update && apt-get install -y --no-install-recommends wget curl

# Download genesis file and links file from IPFS
###############################################################################
# To slow using ipget, currently we use gateway
# PUT needed CID_OF_GENESIS and CID_OF_CONFIG here
RUN wget -O /genesis.json https://ipfs.io/ipfs/<CID_OF_GENESIS>
RUN wget -O /config.toml https://ipfs.io/ipfs/<CID_OF_CONFIG>

WORKDIR /

# Copy compiled kernel and binaries
###############################################################################
COPY --from=build_stage_cuda /sources/build/cyberd /usr/bin/cyberd
COPY --from=build_stage_cuda /sources/build/cyberdcli /usr/bin/cyberdcli

COPY --from=build_stage_cuda /usr/lib/cbdrank.h /usr/lib/cbdrank.h
COPY --from=build_stage_cuda /usr/lib/libcbdrank.so /usr/lib/libcbdrank.so

COPY --from=build_stage_rust /sources/api/libgo_cosmwasm.so /usr/lib/libgo_cosmwasm.so

# Copy startup scripts
###############################################################################

COPY start_script.sh start_script.sh
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x start_script.sh
RUN chmod +x /entrypoint.sh

# Start
###############################################################################
EXPOSE 26656 26657
ENTRYPOINT ["/entrypoint.sh"]
CMD ["./start_script.sh"]
