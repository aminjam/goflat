#!/bin/bash

build() {
    local commit=$(git rev-parse HEAD)
    local dirty=$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)
    local message=$(git describe --tags)
    local ldflag="-X main.GitCommit=${commit}${dirty} -X main.GitDescribe=${message}"

    local os="darwin freebsd linux windows"
    local arch="386 amd64 arm"
    if [ "${DEV}" == "dev" ]; then
        local os=$(go env GOOS)
        local arch=$(go env GOARCH)
    fi

    # Build!
    echo "==> Building ${MAIN_PACKAGE}...(${arch}) for (${os})"
    gox \
        -os="${os}" \
        -arch="${arch}" \
        -ldflags="${ldflag}"  \
        -output="$DIR/pkg/{{.OS}}_{{.Arch}}/${BINARY_NAME}" \
        ${MAIN_PACKAGE}
}

copy_to_path() {
    # Copy our OS/Arch to the bin/ directory
    local binary_file="$DIR/pkg/$(go env GOOS)_$(go env GOARCH)"
    for F in $(find ${binary_file} -mindepth 1 -maxdepth 1 -type f); do
        cp ${F} "${GOPATH}/bin/"
    done
}

package() {
    echo "==> Packaging..."
    for PLATFORM in $(find $DIR/pkg -mindepth 1 -maxdepth 1 -type d); do
        local osarch=$(basename ${PLATFORM})
        echo "--> ${osarch}"
        pushd $PLATFORM >/dev/null 2>&1
        zip ../${osarch}.zip ./*
        popd >/dev/null 2>&1
    done
}

main() {
    set +e
    set -o pipefail

    SOURCE="${BASH_SOURCE[0]}"
    while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
    DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"
    rm -rf $DIR/pkg
    mkdir -p $DIR/pkg

    #parse args
    BINARY_NAME=$1
    MAIN_PACKAGE=$2
    DEV=$3

    build
    copy_to_path

    if [ "$DEV" != "dev" ]; then
        package
    fi
}

main "$@"

