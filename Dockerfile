# syntax=docker/dockerfile:experimental

# ============= Setting up base Stage ================
# Set required LUX_VERSION parameter in build image script
ARG LUX_VERSION

# ============= Compilation Stage ================
FROM golang:1.20.10-bullseye AS builder

WORKDIR /build

# Copy lux dependencies first (intermediate docker image caching)
# Copy node directory if present (for manual CI case, which uses local dependency)
COPY go.mod go.sum node* ./

# Download lux dependencies using go mod
RUN go mod download && go mod tidy -compat=1.20

# Copy the code into the container
COPY . .

# Pass in SUBNET_EVM_COMMIT as an arg to allow the build script to set this externally
ARG SUBNET_EVM_COMMIT
ARG CURRENT_BRANCH

RUN export SUBNET_EVM_COMMIT=$SUBNET_EVM_COMMIT && export CURRENT_BRANCH=$CURRENT_BRANCH && ./scripts/build.sh /build/srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy

# ============= Cleanup Stage ================
FROM SkyChains/node:$LUX_VERSION AS builtImage

# Copy the evm binary into the correct location in the container
COPY --from=builder /build/srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy /node/build/plugins/srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy
