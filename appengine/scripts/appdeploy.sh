#!/usr/bin/env bash
# Deploy the web app to Google App Engine.

name=consul
package=github.com/hashicorp/consul

OLD_GOPATH=${GOPATH}
export GOPATH=/tmp/gopath_${name}

echo GOROOT=${GOROOT}
echo GOPATH=${GOPATH}
echo OLD_GOPATH=${OLD_GOPATH}

echo ========= ${name} ========= "rm -rf ${GOPATH}/src/*"
rm -rf ${GOPATH}/src/*

echo ========= ${name} ========= "mkdir -p ${GOPATH}/src/${package}/"
mkdir -p ${GOPATH}/src/${package}

echo ========= ${name} ========= "cp -r ${OLD_GOPATH}/src/${package}/vendor/* ${GOPATH}/src/"
cp -r ${OLD_GOPATH}/src/${package}/vendor/* ${GOPATH}/src/

echo ========= ${name} ========= "cp -r ${OLD_GOPATH}/src/${package}/* ${GOPATH}/src/${package}/"
cp -r ${OLD_GOPATH}/src/${package}/* ${GOPATH}/src/${package}/

echo ========= ${name} ========= "rm -rf ${GOPATH}/src/${package}/vendor/"
rm -rf ${GOPATH}/src/${package}/vendor/

echo ========= ${name} ========= "./scripts/go-get.sh"
./scripts/go-get.sh

# Set GCLOUD_PROJECT in .env
export GCLOUD_PROJECT=
. ../.env

# Deploy to Google Cloud.
echo ========= ${name} ========= gcloud config set project ${GCLOUD_PROJECT}
gcloud config set project ${GCLOUD_PROJECT}

echo ========= ${name} ========= gcloud config list project
gcloud config list project

echo ========= ${name} ========= gcloud app deploy app.yaml --quiet --verbosity=info
gcloud app deploy app.yaml --quiet --verbosity=info

echo ========= ${name} Deployed! =========

export GOPATH=${OLD_GOPATH}
echo GOPATH=${GOPATH}

# sudo apt install golang
# sudo apt install google-cloud-sdk-app-engine-go
# echo 'export GOROOT=$HOME/go' >>~/.zshrc
# echo 'export PATH=$PATH:$GOROOT/bin' >>~/.zshrc
# echo 'export GOPATH=$GOROOT/bin' >>~/.zshrc
# go get -u google.golang.org/appengine
