# syntax=docker/dockerfile:1
# check=skip=InvalidDefaultArgInFrom

FROM golangci/golangci-lint:v2.11.4-alpine@sha256:72bcd68512b4e27540dd3a778a1b7afd45759d8145cfb3c089f1d7af53e718e9 AS golangci-lint-bin

FROM registry.suse.com/bci/golang:1.25.7 AS builder
ARG MK_HOST_ARCH
ENV ARCH=$MK_HOST_ARCH
RUN zypper in -y bash git gcc docker vim less file curl wget ca-certificates trousers-devel
COPY --from=golangci-lint-bin /usr/bin/golangci-lint /usr/local/bin/golangci-lint
ENV HOME=/go/src/github.com/harvester/rancherd

# ---- base ----
FROM builder AS base
WORKDIR /go/src/github.com/harvester/rancherd
COPY . .

# ---- build ----
FROM base AS build
ARG MK_REPO_ID
RUN --mount=type=cache,target=/go/pkg/mod,id=rancherd-go-mod-${MK_REPO_ID} \
    --mount=type=cache,target=/go/src/github.com/harvester/rancherd/.cache/go-build,id=rancherd-go-build-${MK_REPO_ID} \
    ./scripts/build

FROM scratch AS build-output
COPY --from=build /go/src/github.com/harvester/rancherd/bin/ /bin/

# ---- validate ----
FROM base AS validate
ARG MK_REPO_ID
RUN --mount=type=cache,target=/go/pkg/mod,id=rancherd-go-mod-${MK_REPO_ID} \
    --mount=type=cache,target=/go/src/github.com/harvester/rancherd/.cache/go-build,id=rancherd-go-build-${MK_REPO_ID} \
    ./scripts/validate

# ---- test ----
FROM base AS test
ARG MK_REPO_ID
RUN --mount=type=cache,target=/go/pkg/mod,id=rancherd-go-mod-${MK_REPO_ID} \
    --mount=type=cache,target=/go/src/github.com/harvester/rancherd/.cache/go-build,id=rancherd-go-build-${MK_REPO_ID} \
    ./scripts/test
