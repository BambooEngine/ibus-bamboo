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

mkdir ../build
cd ../build
osc checkout $OSC_PATH
cd $TRAVIS_BUILD_DIR
rm -rf ../build/$OSC_PATH/*
make build src DESTDIR=../build/$OSC_PATH
cd ../build/$OSC_PATH
osc add *.spec *.changes *.tar.gz
osc addremove
osc st
osc ci -m "$TRAVIS_TAG"

