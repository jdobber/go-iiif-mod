#!/bin/bash

pushd () {
    command pushd "$@" > /dev/null
}

popd () {
    command popd "$@" > /dev/null
}

DEPDIR=".deps"
ARCHIVE="https://github.com/chai2010/webp/archive/v1.1.0.tar.gz"

mkdir -p $DEPDIR 
pushd $DEPDIR

echo "Downloading and extracting " $ARCHIVE
curl -s -L $ARCHIVE -o - | tar -xz -C .

echo "Patching sources..."
pushd webp-1.1.0
patch -p1 < ../../webp.patch
popd 

popd
echo "Modifying go.mod"
go mod edit -replace github.com/chai2010/webp=.deps/webp-1.1.0

echo "DONE"
echo -e "Now build the test program:\n\ngo build -o bin/test main.go\n"
echo -e "Then run:\n\n./bin/test -config ./config.json -prefix ./fixtures/ -uri colorize.webp/0,0,1024,1024/1024,/0/default.webp\n"