# This Docker allows running the Hephaestus binary by providing the path to the configuration file.
#
# How to build the image:
# > docker build --tag desmoslabs/hephaestus .
#
# How to run the image:
# > docker run --volume /path/to/folder/containing/config:/data desmoslabs/hephaestus /data/config.yaml

FROM golang:1.18-alpine
ARG arch=x86_64

# Set up dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python3 ca-certificates build-base
RUN set -eux; apk add --no-cache $PACKAGES;

# Set working directory for the build
WORKDIR /code

# Add sources files
COPY . /code/

# See https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.1.1/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
RUN sha256sum /lib/libwasmvm_muslc.aarch64.a | grep 9ecb037336bd56076573dc18c26631a9d2099a7f2b40dc04b6cae31ffb4c8f9a

ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.1.1/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a
RUN sha256sum /lib/libwasmvm_muslc.x86_64.a | grep 6e4de7ba9bad4ae9679c7f9ecf7e283dd0160e71567c6a7be6ae47c81ebe7f32

# Copy the library you want to the final location that will be found by the linker flag `-lwasmvm_muslc`
RUN cp /lib/libwasmvm_muslc.${arch}.a /usr/local/lib/libwasmvm_muslc.a

# Set the entrypoint, so that the user can set the config using the CMD
RUN BUILD_TAGS=muslc GOOS=linux GOARCH=amd64 LINK_STATICALLY=true make build
RUN cp /code/build/hephaestus /usr/bin/hephaestus
ENTRYPOINT ["hephaestus"]
