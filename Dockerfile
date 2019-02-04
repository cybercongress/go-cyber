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
RUN go build -tags cuda -o cyberd ./daemon
RUN go build -o cyberdcli ./cli


###############################################################################
# Create base image
###############################################################################
FROM nvidia/cuda:10.0-runtime-ubuntu18.04

ENV GO_VERSION 1.11.2
ENV GO_ARCH 'linux-amd64'
ENV GO_BIN_SHA '1dfe664fa3d8ad714bbd15a36627992effd150ddabd7523931f077b3926d736d'

ENV IPGET_VERSION 0.3.0
ENV IPGET_ARCH 'linux-amd64'

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

#  Download genesis file and links file from IPFS
###############################################################################
RUN url="https://dist.ipfs.io/ipget/v${IPGET_VERSION}/ipget_v${IPGET_VERSION}_${IPGET_ARCH}.tar.gz" && \
    wget -O ipget.tgz "$url" && \
    tar -xzf ipget.tgz && \
    rm ipget.tgz

WORKDIR /ipget
RUN ./ipget QmSFTpNaXB3FhB4EWjsrUydupZXTL8Z44c2j18o5CGnN5h -o /genesis.json
RUN ./ipget QmepwmLe7vQcK2W6WmvfEk46de3cJ4Jp6jXRXNhuR2AfJ9 -o /links
RUN ./ipget QmVKoU5E4MWkSDmkVqGRwZwDzY9Jgfi4jv16SwiMQTufmy -o /config.toml

WORKDIR /

#  Copy compiled kernel and binaries
###############################################################################
COPY --from=build_stage /sources/cyberd /usr/bin/cyberd
COPY --from=build_stage /sources/cyberdcli /usr/bin/cyberdcli

COPY --from=build_stage /usr/lib/cbdrank.h /usr/lib/cbdrank.h
COPY --from=build_stage /usr/lib/libcbdrank.so /usr/lib/libcbdrank.so

#  Copy startup scripts
###############################################################################

COPY start_script.sh start_script.sh
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x start_script.sh
RUN chmod +x /entrypoint.sh


#  Start
###############################################################################
EXPOSE 26656 26657 1317
ENTRYPOINT ["/entrypoint.sh"]
CMD ["./start_script.sh"]
