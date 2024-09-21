#!/bin/bash

set -e

get_commit_hash() {
    local hash=$(git rev-parse HEAD)

    local dirty=$(git status --porcelain)
    if [ -n "${dirty}" ]; then
        local hash_short=$(echo $hash | cut -c 1-6)
        echo "${hash_short}-dirty"
    fi
    echo $hash
}

get_build_time() {
    cat /etc/os-release | grep -q "alpine"
    if [ $? -eq 0 ]; then
        date -Iseconds
    else
        date --iso=seconds
    fi
}

compile() {
    local goarch=$1
    local ldflags=$2

    echo "Compiling for arch $goarch..."
    GOOS=linux GOARCH=$goarch go build -ldflags="${ldflags}" -o "bin/micro-ddns" cmd/main.go
}

if [ -z "${VERSION}" ]; then
    VERSION="0.0.1-$(get_commit_hash)"
fi

if [ -z "${BUILD_TIME}" ]; then
    BUILD_TIME="$(get_build_time)"
fi

if [ -z "${GO_VERSION}" ]; then
    GO_VERSION="$(go version | awk '{print $3}')"
fi

if [ -z "${OUTPUT_PATH}" ]; then
    OUTPUT_PATH="bin/micro-ddns"
fi

if [ -z "${TARGETARCH}" ]; then
    TARGETARCH="amd64"
    IFS="," read -ra ARCHS <<< $TARGETARCH
fi

LDFLAGS="-X 'github.com/masteryyh/micro-ddns/internal/version.Version=${VERSION}'"
LDFLAGS="${LDFLAGS} -X 'github.com/masteryyh/micro-ddns/internal/version.BuildTime=${BUILD_TIME}'"
LDFLAGS="${LDFLAGS} -X 'github.com/masteryyh/micro-ddns/internal/version.GoVersion=${GO_VERSION}'"
LDFLAGS="${LDFLAGS} -X 'github.com/masteryyh/micro-ddns/internal/version.CommitHash=$(get_commit_hash)'"

for arch in "${ARCHS[@]}"
do
    compile "${arch}" "${LDFLAGS}"
done


