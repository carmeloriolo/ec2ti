#!/bin/bash

set -eox pipefail

VERSION=$1
PLATFORMS="$PLATFORMS darwin/arm64 darwin/amd64"
PLATFORMS="$PLATFORMS windows/amd64 windows/386"
PLATFORMS="$PLATFORMS linux/amd64 linux/386"
PLATFORMS="$PLATFORMS linux/ppc64 linux/ppc64le"

mkdir -p artifacts

for PLATFORM in $PLATFORMS; do
  GOOS=${PLATFORM%/*}
  GOARCH=${PLATFORM#*/}
  APP="ec2ti"
  ARTIFACT="${GOOS}-${GOARCH}-ec2ti.tar.gz"
  echo "Building ${PLATFORM}"
  if [[ "${GOOS}" == "windows" ]]; then APP="ec2ti.exe"; fi
  GOOS=${GOOS} GOARCH=${GOARCH} go build -v -o ./bin/${APP} -ldflags="-s -w -X main.AppVersion=${VERSION}" ./cmd/main.go
  tar cvzf ${ARTIFACT} -C ./bin ${APP}
  if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    md5sum ${ARTIFACT} >> ./artifacts/checksum.txt
  elif [[ "$OSTYPE" == "darwin"* ]]; then
    md5 ${ARTIFACT} >> ./artifacts/checksum.txt
  fi
  mv ${ARTIFACT} ./artifacts
done

