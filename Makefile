ROOT := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

DOCKER_BUILDKIT := 1
export DOCKER_BUILDKIT

ifdef CI
  BOLD  :=
  CYAN  :=
  RESET :=
else
  BOLD  := \033[1m
  CYAN  := \033[36m
  RESET := \033[0m
endif

BANNER = @printf "$(BOLD)$(CYAN)[target: $@]$(RESET)\n"

MK_REPO_ID := $(shell echo -n "$(ROOT)$$(cat /etc/machine-id 2>/dev/null)" | sha256sum | cut -c1-8)
export MK_REPO_ID

MK_HOST_ARCH := $(shell uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')
export MK_HOST_ARCH

MK_DOCKER_PROGRESS ?= plain
export MK_DOCKER_PROGRESS


DOCKER_BUILD = docker build \
    --progress=$(MK_DOCKER_PROGRESS) \
    --build-arg MK_REPO_ID \
    --build-arg MK_HOST_ARCH \
    -f $(ROOT)/Dockerfile $(ROOT)

# ---- Pre-generate version env for container builds (no .git needed inside Docker) ----
# Also handles git worktree checkouts where .git is a pointer file to an external directory.
gen-version-env:
	@bash $(ROOT)/scripts/version > /dev/null

# ---- build ----
build: gen-version-env | $(ROOT)/bin
	$(BANNER)
	$(DOCKER_BUILD) --target build-output \
	    --output type=local,dest=$(ROOT)

$(ROOT)/bin:
	@mkdir -p $@

# ---- validate ----
validate: gen-version-env
	$(BANNER)
	$(DOCKER_BUILD) --target validate

# ---- test ----
test: gen-version-env
	$(BANNER)
	$(DOCKER_BUILD) --target test

# ---- package (host-side: copies built binaries to dist/artifacts) ----
package: build
	$(BANNER)
	mkdir -p $(ROOT)/dist/artifacts
	cp $(ROOT)/bin/rancherd-* $(ROOT)/dist/artifacts/

# ---- build-sha-file (host-side: sha256 checksums for built binaries) ----
build-sha-file: build
	$(BANNER)
	mkdir -p $(ROOT)/dist/artifacts
	cd $(ROOT)/bin && sha256sum rancherd-amd64 > $(ROOT)/dist/artifacts/sha256sum-amd64.txt
	cd $(ROOT)/bin && sha256sum rancherd-arm64 > $(ROOT)/dist/artifacts/sha256sum-arm64.txt

# ---- ci ----
ci: build test validate package build-sha-file

# ---- clean ----
clean:
	$(BANNER)
	@rm -rf $(ROOT)/bin $(ROOT)/dist

# ---- clean-all ----
clean-all: clean

.DEFAULT_GOAL := ci

.PHONY: build validate test package build-sha-file ci clean clean-all gen-version-env
