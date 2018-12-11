FROM nvidia/cuda:10.0-devel-ubuntu18.04 as build_stage

ENV GO_VERSION 1.11.2
ENV GO_ARCH 'linux-amd64'
ENV GO_BIN_SHA '1dfe664fa3d8ad714bbd15a36627992effd150ddabd7523931f077b3926d736d'


#  Install required dev tools to compile cyberd
###############################################################################
RUN apt-get update && apt-get install -y --no-install-recommends wget git


#  Install golang
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
RUN nvcc -fmad=false -shared -o libcbdrank.so rank.cu --compiler-options '-fPIC -frounding-math -fsignaling-nans' && \
    cp libcbdrank.so /usr/lib/ && cp cbdrank.h /usr/lib/


# Compile cyberd
###############################################################################
WORKDIR /sources
RUN go build -tags cuda -o daemon ./cyberd
RUN go build -o cli ./cyberdcli
RUN go build -o daemon_proxy ./proxy


###############################################################################
# Create base image
###############################################################################
FROM nvidia/cuda:10.0-runtime-ubuntu18.04

ENV GO_VERSION 1.11.2
ENV GO_ARCH 'linux-amd64'
ENV GO_BIN_SHA '1dfe664fa3d8ad714bbd15a36627992effd150ddabd7523931f077b3926d736d'


#  Install required dev tools to install go
###############################################################################
RUN apt-get update && apt-get install -y --no-install-recommends wget curl


#  Install golang
###############################################################################
RUN url="https://golang.org/dl/go${GO_VERSION}.${GO_ARCH}.tar.gz" && \
	wget -O go.tgz "$url" && \
	echo "${GO_BIN_SHA} *go.tgz" | sha256sum -c - && \
	tar -C /usr/local -xzf go.tgz &&\
	rm go.tgz


#  Copy compiled kernel and binaries
###############################################################################
COPY --from=build_stage /sources/daemon /usr/bin/cyberd
COPY --from=build_stage /sources/cli /usr/bin/cyberdcli
COPY --from=build_stage /sources/daemon_proxy /usr/bin/cyberdproxy

COPY --from=build_stage /usr/lib/cbdrank.h /usr/lib/cbdrank.h
COPY --from=build_stage /usr/lib/libcbdrank.so /usr/lib/libcbdrank.so


#  Copy configs and startup scripts
###############################################################################
COPY ./testnet/genesis.json /genesis.json
COPY ./testnet/config.toml /config.toml

COPY start_script.sh start_script.sh
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x start_script.sh
RUN chmod +x /entrypoint.sh


#  Start
###############################################################################
EXPOSE 26656 26657 26660
ENTRYPOINT ["/entrypoint.sh"]
CMD ["./start_script.sh"]
