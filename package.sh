#!/bin/bash

set -e

create_package() {
    local os=$1
    local arch=$2
    local version=$3

    echo "Packing for OS $os and arch $arch..."
    local path="output/micro-ddns-$os-$arch"
    if [ ! -d "$path" ]; then
        mkdir -p "$path"
    fi

    local binary="bin/micro-ddns-${os}-${arch}"
    if [ "$os" == "windows" ]; then
        binary="${binary}.exe"
    fi

    cp "$binary" "$path/micro-ddns"
    if [ "$os" == "windows" ]; then
        mv "$path/micro-ddns" "$path/micro-ddns.exe"
    fi

    if [ ! -d "$path/example" ]; then
        mkdir -p "$path/example"
    fi
    cp "example/homelab.yaml" "$path/example/config.yaml"

    if [ "$os" == "linux" ]; then
        cp -r "init" "$path"
    fi

    local tar_name="output/micro-ddns-${os}-${arch}-$version.tar.gz"
    tar -czf "$tar_name" -C "$path" .
    sha256sum "$tar_name" > "$tar_name.sha256"
}

if [ -z "${VERSION}" ]; then
    echo "VERSION is required"
    exit 1
fi

if [ -z "${PLATFORMS}" ]; then
    PLATFORMS="$(go tool dist list | grep 'linux\|windows\|freebsd\|darwin' | tr '\n' ',' | sed 's/,$//')"
fi

IFS=',' read -r -a platforms <<< "$PLATFORMS"
for platform in "${platforms[@]}"
do
    IFS='/' read -r -a os_arch <<< "$platform"
    create_package "${os_arch[0]}" "${os_arch[1]}" "$VERSION" &
done
wait

rm -rf output/*/
