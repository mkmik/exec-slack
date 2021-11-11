#!/bin/bash

set -e
set -x

CLONE_DIR=/tmp/iox

echo "Checking build dependencies..."
command -v clang >/dev/null || (echo "Installing clang..." && brew update && brew install clang)
command -v cargo >/dev/null || (echo "Installing rustup..." && brew update && brew install rustup)

[ -d $CLONE_DIR ] && echo "Cleaning $CLONE_DIR" && rm -rf $CLONE_DIR
echo "Creating $CLONE_DIR"
mkdir -p $CLONE_DIR

echo "Cloning IOx..."
cd $CLONE_DIR
git clone https://github.com/influxdata/influxdb_iox.git .

echo "Running tests...\n"
cargo test --workspace
