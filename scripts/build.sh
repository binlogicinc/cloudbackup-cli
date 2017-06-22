#!/usr/bin/env bash
set -e

##
# NOTE: This simple script is intended to use it locally while developing.
##

getCurrTag() {
  echo `git describe --always --tags --abbrev=0 | tr -d "[v\r\n]"`
}

[ -e "./build" ] && \
  echo "Cleaning up old builds..." && \
  rm -rf "./build"

go build -ldflags "-X github.com/binlogicinc/cloudbackup-cli/cmd.version=${getCurrTag}" \
  -o="./build/cloudbackup-cli"
