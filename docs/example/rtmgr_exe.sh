#!/bin/bash
#set -x

curdir=`pwd`
RMRVERSION=3.6.0
wget --content-disposition https://packagecloud.io/o-ran-sc/staging/packages/debian/stretch/rmr_${RMRVERSION}_amd64.deb/download.deb && dpkg -i rmr_${RMRVERSION}_amd64.deb && rm -rf rmr_${RMRVERSION}_amd64.deb
wget --content-disposition https://packagecloud.io/o-ran-sc/staging/packages/debian/stretch/rmr-dev_${RMRVERSION}_amd64.deb/download.deb && dpkg -i rmr-dev_${RMRVERSION}_amd64.deb && rm -rf rmr-dev_${RMRVERSION}_amd64.deb

wget http://launchpadlibrarian.net/463891089/libnng1_1.2.6-1_amd64.deb && dpkg -i libnng1_1.2.6-1_amd64.deb && rm -rf libnng1_1.2.6-1_amd64.deb 

mkdir -p /tmp/go/bin

wget --quiet https://github.com/go-swagger/go-swagger/releases/download/v0.19.0/swagger_linux_amd64 && mv swagger_linux_amd64 swagger && chmod +x swagger

mv swagger /tmp/go/bin

export GOPATH=/tmp/go
export GO111MODULE=on

mkdir -p /tmp/go/src/routing-manager

git clone "https://gerrit.o-ran-sc.org/r/ric-plt/appmgr" \
  && cp appmgr/api/appmgr_rest_api.yaml api/ \
  && rm -rf appmgr

cp -r $PWD/../../rtmgr/* /tmp/go/src/routing-manager/. 

cd /tmp/go/src/routing-manager
/tmp/go/bin/swagger generate server -f api/routing_manager.yaml -t pkg/ --exclude-main -r LICENSE
/tmp/go/bin/swagger generate client -f api/appmgr_rest_api.yaml -t pkg/ -m appmgr_model -c appmgr_client -r LICENSE
  
export GO111MODULE=on 
ldconfig
go build ./cmd/rtmgr.go

cp -f rtmgr $curdir/.

mkdir -p /db && touch /db/rt.json && chmod 777 /db/rt.json
