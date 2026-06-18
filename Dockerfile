FROM golangci/golangci-lint:v2.12.2-alpine@sha256:91b27804074a0bacea298707f016911e60cf0cdbc6c7bf5ccacb5f0606d18d60 AS golangci-lint-bin

FROM registry.suse.com/bci/golang:1.26.0 AS builder
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
