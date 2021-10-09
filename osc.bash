#!/bin/bash
if [[ $TRAVIS_TAG =~ ^v[0-9]+\.[0-9]+\.[0-9]+ ]]; then
  echo "Release tag: $TRAVIS_TAG"
else
  exit 0
fi
echo "[general]" >> ~/.oscrc
echo "apiurl = https://api.opensuse.org" >> ~/.oscrc
echo "[https://api.opensuse.org]" >> ~/.oscrc
echo "user = $OSC_USER" >> ~/.oscrc
echo "pass = $OSC_PASS" >> ~/.oscrc

sudo apt-get update
sudo apt-get install osc -y
echo "osc install"

mkdir ../build
cd ../build
echo "osc checkout"
yes 1 | osc checkout $OSC_PATH
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
osc ci -m "$TRAVIS_TAG"
echo "osc done"
exit 0
