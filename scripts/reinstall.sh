#!/usr/bin/env bash
set -exuo pipefail

mkdir -p bin
rm -f bin/*

TARGET_OS=${1:-darwin}

source ./scripts/install.sh $TARGET_OS