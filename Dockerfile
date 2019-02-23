FROM nvidia/cuda:10.0-devel-ubuntu18.04 as build_stage

ENV GO_VERSION 1.11.5
ENV GO_ARCH 'linux-amd64'
ENV GO_BIN_SHA 'ff54aafedff961eb94792487e827515da683d61a5f9482f668008832631e5d25'


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

#  Install useful dev tools
###############################################################################
RUN apt-get update && apt-get install -y --no-install-recommends wget curl

#  Download genesis file and links file from IPFS
###############################################################################
# To slow using ipget, currently we use gateway
RUN wget -O /genesis.json https://ipfs.io/ipfs/QmQ88ZGztF7QvZc9nHnXZ161jn8WWDzfaAQyxPtAm8TDfh
RUN wget -O /links https://ipfs.io/ipfs/QmaXGhTQx3qBkcq5QBS6Se61bsqd8hjgivsRzGRKpzoqvE
RUN wget -O /config.toml https://ipfs.io/ipfs/QmWDdNUvVtzAgeQnAFcWohQe3DozKBVzhESnhzniwxqeQ3

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
