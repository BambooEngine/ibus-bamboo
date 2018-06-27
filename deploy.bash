#!/bin/bash
echo "Check tag -------------------------------------------------------------------------------------------------------"
if [[ $TRAVIS_TAG =~ ^v[0-9]+\.[0-9]+\.[0-9]+ ]]; then
  echo "Release tag: $TRAVIS_TAG detected"
else
  echo "Release tag not found ($TRAVIS_TAG), do not deploy"
  exit 0
fi

echo "Check OSC ENV ---------------------------------------------------------------------------------------------------"
if [ -z "$OSC_USER" ] || [ -z "$OSC_PASS" ] || [ -z "$OSC_PATH" ]
then
  echo "OSC_USER|OSC_PASS|OSC_PATH is not defined, do not deploy"
  exit 0
fi

echo "Install OSC -----------------------------------------------------------------------------------------------------"
sudo apt-get update
sudo apt-get install -y osc
osc --version

echo "Make OSC config -------------------------------------------------------------------------------------------------"
echo "[general]" >> ~/.oscrc
echo "apiurl = https://api.opensuse.org" >> ~/.oscrc
echo "[https://api.opensuse.org]" >> ~/.oscrc
echo "user = $OSC_USER" >> ~/.oscrc
echo "pass = $OSC_PASS" >> ~/.oscrc

echo "OSC checkout ----------------------------------------------------------------------------------------------------"
mkdir ../obs
cd ../obs
osc checkout $OSC_PATH
cd $TRAVIS_BUILD_DIR

echo "Build new OSC source --------------------------------------------------------------------------------------------"
rm -rf ../obs/$OSC_PATH/*
make build src DESTDIR=../obs/$OSC_PATH
cd ../obs/$OSC_PATH

echo "OSC status ------------------------------------------------------------------------------------------------------"
osc addremove
osc st

echo "OSC commit ------------------------------------------------------------------------------------------------------"
osc ci -m "$TRAVIS_TAG" 
