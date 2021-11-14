#!/bin/bash

VERSION=$1

PLATFORMS="$PLATFORMS windows/amd64 windows/386"
PLATFORMS="$PLATFORMS linux/amd64 linux/386"
PLATFORMS="$PLATFORMS linux/ppc64 linux/ppc64le"
PLATFORMS="darwin/arm64 darwin/amd64"


for PLATFORM in $PLATFORMS; do
  mkdir -p built
  GOOS=${PLATFORM%/*}
  GOARCH=${PLATFORM#*/}
  APP="ec2ti"
  echo "Building ${PLATFORM}"
  if [[ "${GOOS}" == "windows" ]]; then APP="ec2ti.exe"; fi
  GOOS=${GOOS} GOARCH=${GOARCH} go build -v -o ./bin/${APP} -ldflags="-X github.com/carmeloriolo/ec2ti/cmd/main.AppVersion=${VERSION}" ./cmd/main.go
  cd ./bin
  tar cvzf ${GOOS}-${GOARCH}-ec2ti.tar.gz ${APP}
  mv ${GOOS}-${GOARCH}-ec2ti.tar.gz ../built/
  cd -
done

