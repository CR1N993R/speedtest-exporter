#!/bin/bash

build() {
  OS=$1
  ARCH=$2
  ARM_VERSION=$3
  if [ -n "$ARM_VERSION" ]; then
    ARM_VERSION_SUFFIX="v$ARM_VERSION"
  else
    ARM_VERSION_SUFFIX=""
  fi
  if [ "$OS" = "windows" ]; then
    SUFFIX=".exe"
  else
    SUFFIX=""
  fi
  echo "Building for $OS/$ARCH$ARM_VERSION_SUFFIX..."
  GOOS=$OS GOARCH=$ARCH GOARM=$ARM_VERSION go build  -ldflags "-s -w" -o "./build/speedtest-exporter${SUFFIX}"
  cd build
  tar -czvf "./speedtest-exporter-${OS}_${ARCH}${ARM_VERSION_SUFFIX}.tar.gz" "speedtest-exporter${SUFFIX}"
  rm "./speedtest-exporter${SUFFIX}"
  if [[ $? -eq 0 ]]; then
    echo "Build successful: ${OS}_${ARCH}${ARM_VERSION_SUFFIX}"
  else
    echo "Build failed for ${OS}_${ARCH}${ARM_VERSION_SUFFIX}"
  fi
  cd ..
}

collect_checksums() {
  OUTPUT_FILE="./build/checksums.txt"
  find "./build" -type f | while read -r file
  do
    sha256sum "$file" | awk '{print $1 "  " $2}' >> "$OUTPUT_FILE"
  done
  echo "Checksums written"
}

# darwin
build darwin arm64
build darwin amd64

# freebsd
build freebsd "386"
build freebsd amd64
build freebsd arm64
build freebsd arm "6"
build freebsd arm "7"

# linux
build linux "386"
build linux amd64
build linux arm64
build linux arm "5"
build linux arm "6"
build linux arm "7"
build linux mips
build linux mips64
build linux mips64le
build linux mipsle
build linux ppc64
build linux ppc64le
build linux riscv64
build linux s390x

# netbsd
build netbsd "386" ""
build netbsd amd64 ""
build netbsd arm64 ""
build netbsd arm "6"
build netbsd arm "7"

# openbsd
build openbsd "386" ""
build openbsd amd64 ""
build openbsd arm64 ""
build openbsd arm "7"

# windows
build windows "386" ""
build windows amd64 ""
build windows arm64 ""

collect_checksums
