#!/bin/bash
set -euxo pipefail
echo "[general]" >> ~/.oscrc
echo "apiurl = https://api.opensuse.org" >> ~/.oscrc
echo "[https://api.opensuse.org]" >> ~/.oscrc
echo "user = $OSC_USER" >> ~/.oscrc
echo "pass = $OSC_PASS" >> ~/.oscrc
export DEBIAN_FRONTEND=noninteractive

mkdir ../build && cd ../build
echo "osc checkout $OSC_PATH"
yes 2>/dev/null | osc checkout $OSC_PATH
cd $TRAVIS_BUILD_DIR
rm -rf ../build/$OSC_PATH/*
echo "osc build"
make build src DESTDIR=../build/$OSC_PATH
cd ../build/$OSC_PATH
osc add *.spec *.changes *.tar.gz
echo "osc addremove"
osc addremove
echo "osc st"
osc st
echo "osc commit"
echo "$TRAVIS_TAG"
yes 2>/dev/null | osc ci -m "$TRAVIS_TAG"
echo "osc done"
exit 0
