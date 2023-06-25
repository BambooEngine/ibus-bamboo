#!/bin/bash
if [[ $GH_TAG =~ ^v[0-9]+\.[0-9]+\.[0-9]+ ]]; then
  echo "Release tag: $GH_TAG"
else
  echo "Release tag ($GH_TAG) is invalid, exiting"
  exit 0
fi
#set -euxo pipefail
echo "[general]" >> ~/.oscrc
echo "apiurl = https://api.opensuse.org" >> ~/.oscrc
echo "[https://api.opensuse.org]" >> ~/.oscrc
echo "user = $OSC_USER" >> ~/.oscrc
echo "pass = $OSC_PASS" >> ~/.oscrc
export DEBIAN_FRONTEND=noninteractive

mkdir /build && pushd /build
osc checkout $OSC_PATH
popd
rm -rf /build/$OSC_PATH/*
make build src DESTDIR=/build/$OSC_PATH
cd /build/$OSC_PATH
osc add *.spec *.changes *.tar.gz
osc addremove
osc st
osc ci -m "$GH_TAG"
exit 0
