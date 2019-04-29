RELEASE ?= $(shell git describe --tags || echo -n)
SHELL := /usr/bin/env bash

test:
	go test -v ./...

install:
	dep ensure

build:
	test -n "${RELEASE}" && \
	rm -fr release/${RELEASE} && \
	env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o release/${RELEASE}/tcpwait_${RELEASE}_darwin_amd64 && \
	upx --ultra-brute release/${RELEASE}/tcpwait_${RELEASE}_darwin_amd64 && \
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o release/${RELEASE}/tcpwait_${RELEASE}_linux_amd64 && \
	upx --ultra-brute release/${RELEASE}/tcpwait_${RELEASE}_linux_amd64 && \
	env GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o release/${RELEASE}/tcpwait_${RELEASE}_linux_arm64 && \
	upx --ultra-brute release/${RELEASE}/tcpwait_${RELEASE}_linux_arm64 && \
	env GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o release/${RELEASE}/tcpwait_${RELEASE}_windows_amd64 && \
	upx --ultra-brute release/${RELEASE}/tcpwait_${RELEASE}_windows_amd64 && \
	dgstore 'release/**/*'

doctoc:
	command -v doctoc &>/dev/null && doctoc README.md || { >&2 echo "Error: install doctoc with \`npm install -g doctoc\`"; exit 1; }
