ARG MK_GOLANGCI_LINT_IMAGE
ARG MK_PACKAGE_BASE registry.suse.com/bci/bci-base:16.0
FROM ${MK_GOLANGCI_LINT_IMAGE} AS golangci-lint

FROM golang:1.25.11-bookworm AS buildenv
ENV GOTOOLCHAIN=auto

COPY --from=golangci-lint /usr/bin/golangci-lint /usr/local/bin/golangci-lint
RUN --mount=type=cache,target=/var/lib/apt/lists apt-get update -qq \
 && apt-get install -y --no-install-recommends \
  gzip \
  tar

# ---- base ----
FROM buildenv AS base
ARG MK_REPO
ARG MK_REPO_ID
WORKDIR /go/src/${MK_REPO}
# to exclude some files, add them in .dockerignore
COPY . .

# ---- build ----
FROM base AS build
ARG MK_REPO
ARG MK_REPO_ID
ARG VERSION v0.0.0-dev
RUN --mount=type=cache,target=/go/pkg/mod,id=harvester-go-mod-${MK_REPO_ID} \
    --mount=type=cache,target=/go/src/${MK_REPO}/.cache/go-build,id=harvester-go-build-${MK_REPO_ID} \
    <<EOF
#!/bin/bash -e

mkdir -p bin
[ "$(uname)" != "Darwin" ] && LINKFLAGS="-extldflags -static -s"
GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.VERSION=$VERSION $LINKFLAGS" -o bin/docker-machine-driver-harvester-amd64
GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "-X main.VERSION=$VERSION $LINKFLAGS" -o bin/docker-machine-driver-harvester-arm64
EOF

# ---- package ----
FROM build AS package
RUN <<EOF
#!/bin/bash -e
mkdir -p dist/artifacts
for arch in amd64 arm64
do
  cp  bin/docker-machine-driver-harvester-${arch} bin/docker-machine-driver-harvester
  tar -zcvf dist/artifacts/docker-machine-driver-harvester-${arch}.tar.gz -C bin docker-machine-driver-harvester
done
EOF

# ---- test ----
FROM base AS test
ARG MK_REPO
ARG MK_REPO_ID
RUN --mount=type=cache,target=/go/pkg/mod,id=harvester-go-mod-${MK_REPO_ID} \
    --mount=type=cache,target=/go/src/${MK_REPO}/.cache/go-build,id=harvester-go-build-${MK_REPO_ID} \
    --mount=type=secret,id=codecov_token_${MK_REPO_ID},env=CODECOV_TOKEN \
    <<EOF
#!/bin/bash -e

mkdir -p cover

echo Running tests:
go test \
  -v \
  -cover \
  -coverprofile=cover/coverage.out \
  -tags=test \
  ./harvester

go tool cover \
  -html=cover/coverage.out \
  -o cover/coverage.html
EOF

# ---- validate ----
FROM base AS validate
ARG MK_REPO
ARG MK_REPO_ID
ARG VERSION v0.0.0-dev
RUN --mount=type=cache,target=/go/pkg/mod,id=harvester-go-mod-${MK_REPO_ID} \
    --mount=type=cache,target=/go/src/${MK_REPO}/.cache/go-build,id=harvester-go-build-${MK_REPO_ID} \
    <<EOF
#!/bin/bash -e

echo Running validation

PACKAGES="$(go list ./...)"

echo Running validation: golangci-lint
golangci-lint run --timeout 5m

echo Running validation: go fmt
test -z "$(go fmt ${PACKAGES} | tee /dev/stderr)"

echo "Running dirty check"

if echo "$VERSION" | grep dirty ; then
    echo "Git is dirty"
    git status
    git diff
    exit 1
fi

echo "All clean"
EOF

# ---- build output ----
FROM scratch AS build-output
ARG MK_REPO
COPY --from=build /go/src/${MK_REPO}/bin/ /bin/

# ---- test output ----
FROM scratch AS test-output
ARG MK_REPO
COPY --from=test /go/src/${MK_REPO}/cover/ /cover/

# ---- package output ----
FROM scratch AS package-output
ARG MK_REPO
COPY --from=package /go/src/${MK_REPO}/dist/ /dist/
