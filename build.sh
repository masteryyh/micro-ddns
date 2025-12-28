#!/bin/bash

set -e

get_commit_hash() {
    local hash
    hash=$(git rev-parse HEAD)

    local dirty
    dirty=$(git status --porcelain)
    if [ -n "${dirty}" ]; then
        local hash_short
        hash_short=$(echo "$hash" | cut -c 1-6)
        echo "${hash_short}-dirty"
    else
        echo "$hash"
    fi
}

get_build_time() {
    if grep -q "alpine" /etc/os-release; then
        date -Iseconds
    else
        date --iso=seconds
    fi
}

compile() {
    local os=$1
    local arch=$2
    local ldflags=$3

    echo "Compiling for arch $arch... of OS $os"
    local binary="bin/micro-ddns-${os}-${arch}"
    if [ "$os" == "windows" ]; then
        binary="${binary}.exe"
    fi
    GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build -ldflags="${ldflags}" -o "$binary" cmd/main.go
}

if [ -z "${VERSION}" ]; then
    VERSION="$(get_commit_hash)"
fi

if [ -z "${BUILD_TIME}" ]; then
    BUILD_TIME="$(get_build_time)"
fi

if [ -z "${GO_VERSION}" ]; then
    GO_VERSION="$(go version | awk '{print $3}')"
fi

if [ -z "${PLATFORMS}" ]; then
    PLATFORMS="$(go tool dist list | grep 'linux\|windows\|freebsd\|darwin' | tr '\n' ',' | sed 's/,$//')"
fi

LDFLAGS="-X 'github.com/masteryyh/micro-ddns/internal/version.Version=${VERSION}'"
LDFLAGS="${LDFLAGS} -X 'github.com/masteryyh/micro-ddns/internal/version.BuildTime=${BUILD_TIME}'"
LDFLAGS="${LDFLAGS} -X 'github.com/masteryyh/micro-ddns/internal/version.GoVersion=${GO_VERSION}'"
LDFLAGS="${LDFLAGS} -X 'github.com/masteryyh/micro-ddns/internal/version.CommitHash=$(get_commit_hash)'"

IFS=',' read -r -a platforms <<< "$PLATFORMS"
for platform in "${platforms[@]}"
do
    IFS='/' read -r -a os_arch <<< "$platform"
    compile "${os_arch[0]}" "${os_arch[1]}" "$LDFLAGS" &
done
wait
