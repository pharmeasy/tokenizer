#!/bin/bash

set -e

# Check if Package variable set and not empty, otherwise set to default value
if [ -z "$PKG_SRC" ]; then
  echo "Package not set"
  exit 0
fi

# Check if AppName variable set and not empty, otherwise set to default value
if [ -z "$APP_NAME" ]; then
  echo "App Name not set"
  exit 0
fi

# Get rid of existing binaries
rm -f build/$APP_NAME*

# Check if VERSION variable set and not empty, otherwise set to default value
if [ -z "$VERSION" ]; then
  echo "Version not set"
  exit 0
fi

echo "Building application version $VERSION"

echo "Building default binary"
CGO_ENABLED=1 go build -ldflags "-s -w" -ldflags "-X $PKG_SRC/cmd.version=${VERSION}" -o "build/$APP_NAME" $PKG_SRC

if [ ! -z "${BUILD_ONLY_DEFAULT}" ]; then
  echo "Only default binary was requested to build"
  exit 0
fi

# Build amd64 binaries
OS_PLATFORM_ARG=(linux darwin)
OS_ARCH_ARG=(amd64)
for OS in ${OS_PLATFORM_ARG[@]}; do
  for ARCH in ${OS_ARCH_ARG[@]}; do
    echo "Building binary for $OS/$ARCH..."
    GOARCH=$ARCH GOOS=$OS CGO_ENABLED=1 go build -ldflags "-s -w" -ldflags "-X $PKG_SRC/cmd.version=${VERSION}" -o "build/$APP_NAME-$OS-$ARCH" $PKG_SRC
  done
done
