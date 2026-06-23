ROOT := $(realpath $(dir $(realpath $(firstword $(MAKEFILE_LIST)))))
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

MK_HOST_ARCH := $(shell uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')
export MK_HOST_ARCH

MK_SYSTEM_ID := $(strip $(shell \
		if [ -s /etc/machine-id ]; then \
				cat /etc/machine-id 2>/dev/null; \
		elif command -v hostname >/dev/null 2>&1; then \
				hostname 2>/dev/null; \
		else \
				echo -n "unknown"; \
		fi))

MK_REPO             := github.com/harvester/docker-machine-driver-harvester
MK_REPO_ID          := $(shell printf '%s' "$(ROOT)$(MK_SYSTEM_ID)" | sha256sum | cut -c1-8)
MK_VERSION 					:= $(shell git describe --tags --always --dirty)
MK_CODECOV_TOKEN    ?=
MK_DOCKER_PROGRESS  ?= plain

MK_CODECOV_SECRET_ARG  := --secret id=codecov_token_$(MK_REPO_ID),env=MK_CODECOV_TOKEN --no-cache-filter=test
MK_GOLANGCI_LINT_IMAGE := golangci/golangci-lint:v2.12.2-alpine@sha256:91b27804074a0bacea298707f016911e60cf0cdbc6c7bf5ccacb5f0606d18d60
MK_PACKAGE_BASE        := registry.suse.com/bci/bci-base:16.0

DOCKER_BUILD := \
	docker build \
		--progress=$(MK_DOCKER_PROGRESS) \
		--build-arg MK_REPO=$(MK_REPO) \
		--build-arg MK_REPO_ID=$(MK_REPO_ID) \
		--build-arg MK_HOST_ARCH=$(MK_HOST_ARCH) \
		--build-arg PROVIDER_VERSION=$(MK_PROVIDER_VERSION) \
		--build-arg MK_GOLANGCI_LINT_IMAGE=$(MK_GOLANGCI_LINT_IMAGE) \
		--build-arg MK_PACKAGE_BASE=$(MK_PACKAGE_BASE) \
		--build-arg VERSION=$(MK_VERSION) \
		-f $(ROOT)/Dockerfile $(ROOT)

.PHONY: ci validate build test package

# ---- Directories ----
$(ROOT)/bin:
	@mkdir -p $@

# ---- Validate with static analysis ----
validate:
	$(BANNER)
	$(DOCKER_BUILD) --target validate

# ---- Compile harvester-terraform-provider binaries ----
build: $(ROOT)/bin
	$(BANNER)
	$(DOCKER_BUILD) \
		--target build-output \
		--output type=local,dest=.

# ---- Test ----
test: validate
	$(BANNER)
	$(DOCKER_BUILD) \
		$(if $(MK_CODECOV_TOKEN),$(MK_CODECOV_SECRET_ARG)) \
		--target test-output \
		--output type=local,dest=.

# ---- Package harvester-terraform-provider image ----
package: build
	$(BANNER)
	$(DOCKER_BUILD) \
		--target package-output \
		--output type=local,dest=.

ci: validate build test package
	$(BANNER)

.DEFAULT_GOAL := default
default: build test package
