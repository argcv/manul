#!/usr/bin/env bash
# env TEST_MODE=true ./build.sh
TEST_MODE=${TEST_MODE:-false}
# Exit Once failed
set -Eeuxo pipefail

if [ "$TEST_MODE" = true ]; then

pushd log
go test -v
popd # log

fi # WITH_TEST

go get ./cmd/...


echo "Build Release"

PLATFORM="$(uname -s | tr 'A-Z' 'a-z')"

#function is_linux() {
#  [[ "${PLATFORM}" == "linux" ]]
#}

#function is_macos() {
#  [[ "${PLATFORM}" == "darwin" ]]
#}

export CGO_ENABLED=0
#GOOS=linux
export GOOS=${PLATFORM}

export BUILD_DATE=$(date '+%Y%m%d%H%M%S%Z')
export BUILD_LDFLAGS="-X github.com/argcv/manul/version.GitHash=$(git rev-parse HEAD | cut -c1-8) "
export BUILD_LDFLAGS="${BUILD_LDFLAGS} -X github.com/argcv/manul/version.BuildDate=\"${BUILD_DATE}\" "
export BUILD_LDFLAGS="${BUILD_LDFLAGS} \"-extldflags='-static'\""

go build -a -ldflags="$BUILD_LDFLAGS" ./cmd/manul

# Not in use: entrypoint

#env GOOS=linux go build -a -ldflags="$BUILD_LDFLAGS" ./cmd/manul-entrypoint
#
#if [[ "${GOOS}" != 'linux' ]]; then
#    go build -o manul-entrypoint-$GOOS -a -ldflags="$BUILD_LDFLAGS" ./cmd/manul-entrypoint
#fi