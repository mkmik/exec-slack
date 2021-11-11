#!/bin/bash

set -e
set -x

# default is 256, which can cause tests to fail
echo "Setting ulimit..."
ulimit -n 5000
ulimit -a

CLONE_DIR=/tmp/iox
CARGO_PATH=$HOME/.cargo/bin/cargo 

# Required to use `arm64` system protoc 
export PROTOC=/opt/homebrew/bin/protoc 
export PROTOC_INCLUDE=/opt/homebrew/include

echo "Checking build dependencies..."
command -v brew >/dev/null || (echo "Homebrew needs to be installed manually." && exit 1)
command -v $CARGO_PATH >/dev/null || (echo "Rust needs to be installed manually." && exit 1)
command -v clang >/dev/null || (echo "Installing clang via Homebrew..." && brew update && brew install clang)
command -v protoc >/dev/null || (echo "Installing protoc via Homebrew..." && brew update && brew install protobuf)


[ -d $CLONE_DIR ] && echo "Cleaning $CLONE_DIR" && rm -rf $CLONE_DIR
echo "Creating $CLONE_DIR"
mkdir -p $CLONE_DIR

echo "Cloning IOx..."
cd $CLONE_DIR
git clone https://github.com/influxdata/influxdb_iox.git .

echo "Running tests...\n"
$CARGO_PATH test --workspace
