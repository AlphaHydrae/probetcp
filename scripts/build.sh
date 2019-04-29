#!/usr/bin/env bash
set -e

bold=$(tput bold)
normal=$(tput sgr0)
versions="darwin_amd64 linux_amd64 linux_arm64 windows_amd64"

test -z "$RELEASE" && RELEASE="$(git describe --tags || echo -n)"

if test -z "$RELEASE"; then
  >&2 echo No Git tag found
  exit 1
fi

rm -fr release/${RELEASE}

echo "\n${bold}○ Building binaries...${normal}"

for version in $versions; do
  os="$(echo $version | cut -d _ -f 1)"
  arch="$(echo $version | cut -d _ -f 2)"
  env GOOS=$os GOARCH=$arch go build -ldflags="-s -w" -o release/${RELEASE}/tcpwait_${RELEASE}_${os}_${arch} &
done

wait

for version in $versions; do
  file="release/${RELEASE}/tcpwait_${RELEASE}_${version}"
  echo "\n${bold}○ Compressing ${file}...${normal}\n"
  upx --ultra-brute "$file"
done

echo "\n${bold}○ Calculating digests...${normal}\n"
dgstore 'release/**/*'

echo
