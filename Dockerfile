###########################################################################################
# Build cyber
###########################################################################################
FROM ubuntu:20.04

ENV GO_VERSION '1.19.13'
ENV GO_ARCH 'linux-amd64'
ENV GO_BIN_SHA '4643d4c29c55f53fa0349367d7f1bb5ca554ea6ef528c146825b0f8464e2e668'
ENV DEBIAN_FRONTEND=noninteractive 
ENV DAEMON_HOME /root/.cyber
ENV DAEMON_RESTART_AFTER_UPGRADE=true
ENV DAEMON_ALLOW_DOWNLOAD_BINARIES=true
ENV DAEMON_LOG_BUFFER_SIZE=1048
ENV UNSAFE_SKIP_BACKUP=true
ENV DAEMON_NAME cyber
ENV BUILD_DIR /build
ENV PATH /usr/local/go/bin:/root/.cargo/bin:/root/cargo/env:/root/.cyber/scripts:$PATH
ENV CUDA_VER '11.4.4-1'
ENV PATH="/usr/local/go/bin:/usr/local/cuda/bin:$PATH"


# Install go and required deps
###########################################################################################
RUN apt-get update && apt-get install -y --no-install-recommends wget ca-certificates \
&& wget -O go.tgz https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz \
&& echo "${GO_BIN_SHA} *go.tgz" | sha256sum -c - \
&& tar -C /usr/local -xzf go.tgz \
&& rm go.tgz \
&& go version 


COPY . /sources
WORKDIR /sources

# Install CUDA, build tools and compile cyber
###########################################################################################
RUN apt-get update && apt-get -y install --no-install-recommends \
    make gcc g++ \
    curl \
    gnupg \
    git \
    software-properties-common \
#    nvidia-cuda-toolkit \
# Install cuda selected version instead nvidia-cuda-toolkit
&& wget https://developer.download.nvidia.com/compute/cuda/repos/ubuntu2004/x86_64/cuda-ubuntu2004.pin \
&& mv cuda-ubuntu2004.pin /etc/apt/preferences.d/cuda-repository-pin-600 \
&& apt-key adv --fetch-keys https://developer.download.nvidia.com/compute/cuda/repos/ubuntu2004/x86_64/3bf863cc.pub \
&& add-apt-repository "deb https://developer.download.nvidia.com/compute/cuda/repos/ubuntu2004/x86_64/ /" \
&& apt-get update \
&& apt-get install cuda=${CUDA_VER} -y --no-install-recommends \
&& mkdir -p /cyber/cosmovisor/genesis/bin \
&& mkdir -p /cyber/cosmovisor/upgrades/cyberfrey/bin \
&& mkdir -p /cyber/cosmovisor/upgrades/v3/bin \
# Compile cyber for genesis version
###########################################################################################
&& git checkout v0.2.0 \
&& cd /sources/x/rank/cuda \
&& make build \
&& cd /sources \
&& make build CUDA_ENABLED=true \
&& cp ./build/cyber /cyber/cosmovisor/genesis/bin/ \
&& rm -rf ./build \
 # Compile cyber for cyberfrey version
###########################################################################################
&& git checkout v2.0.0 \
&& cd /sources/x/rank/cuda \
&& make build \
&& cd  /sources \
&& make build CUDA_ENABLED=true \
&& cp ./build/cyber /cyber/cosmovisor/upgrades/cyberfrey/bin/ \
&& rm -rf ./build \
 # Compile cyber for v3 version
###########################################################################################
&& git checkout v3.0.0 \
&& cd /sources/x/rank/cuda \
&& make build \
&& cd  /sources \
&& make build CUDA_ENABLED=true \
&& cp ./build/cyber /cyber/cosmovisor/upgrades/v3/bin/ \
&& rm -rf ./build \
# Cleanup 
###########################################################################################
&& apt-get purge -y git \
    make \
    cuda \
    gcc g++ \
    curl \
    gnupg \
    python3.8 \
&& go clean --cache -i \
&& apt-get remove --purge '^nvidia-.*' -y \
&& apt-get autoremove -y \
&& apt-get clean 

# Install cosmovisor
###########################################################################################
 RUN wget -O cosmovisor.tgz https://github.com/cosmos/cosmos-sdk/releases/download/cosmovisor%2Fv1.1.0/cosmovisor-v1.1.0-linux-amd64.tar.gz \
 && tar -xzf cosmovisor.tgz \
 && cp cosmovisor /usr/bin/cosmovisor \
 && chmod +x /usr/bin/cosmovisor \
 && rm cosmovisor.tgz && rm -fR $BUILD_DIR/* && rm -fR $BUILD_DIR/.*[a-z]

# Copy startup scripts and genesis
###########################################################################################
WORKDIR /
COPY start_script.sh start_script.sh
COPY entrypoint.sh /entrypoint.sh
RUN wget -O /genesis.json https://gateway.ipfs.cybernode.ai/ipfs/QmYubyVNfghD4xCrTFj26zBwrF9s5GJhi1TmxvrwmJCipr \
&& chmod +x start_script.sh \
&& chmod +x /entrypoint.sh 

#  Start
###############################################################################
EXPOSE 26656 26657 1317 9090 26660
ENTRYPOINT ["/entrypoint.sh"]
CMD ["./start_script.sh"]
