#!/bin/bash
##############################################################################
#
#   Copyright (c) 2020 AT&T Intellectual Property.
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.
#
##############################################################################
set -eux

echo "--> rtmgr-build-ubuntu.sh"
curdir=`pwd`
RMRVERSION=3.6.5
wget --content-disposition https://packagecloud.io/o-ran-sc/staging/packages/debian/stretch/rmr_${RMRVERSION}_amd64.deb/download.deb && sudo dpkg -i rmr_${RMRVERSION}_amd64.deb && rm -rf rmr_${RMRVERSION}_amd64.deb
wget --content-disposition https://packagecloud.io/o-ran-sc/staging/packages/debian/stretch/rmr-dev_${RMRVERSION}_amd64.deb/download.deb && sudo dpkg -i rmr-dev_${RMRVERSION}_amd64.deb && rm -rf rmr-dev_${RMRVERSION}_amd64.deb

# required to find nng and rmr libs
export LD_LIBRARY_PATH=/usr/local/lib

# go installs tools like go-acc to $HOME/go/bin
# ubuntu minion path lacks go
export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin

# install the go coverage tool helper
go get -v github.com/ory/go-acc

mkdir -p /tmp/go/bin /tmp/go/src/routing-manager

wget --quiet https://github.com/go-swagger/go-swagger/releases/download/v0.19.0/swagger_linux_amd64 && mv swagger_linux_amd64 swagger && chmod +x swagger

mv swagger /tmp/go/bin

export GOPATH=/tmp/go

mkdir -p /tmp/go/src/routing-manager

git clone "https://gerrit.o-ran-sc.org/r/ric-plt/appmgr" \
  && cp appmgr/api/appmgr_rest_api.yaml api/ \
  && rm -rf appmgr

cp -r $curdir/* /tmp/go/src/routing-manager/.

cd /tmp/go/src/routing-manager

currnewdir=`pwd`
/tmp/go/bin/swagger generate server -f api/routing_manager.yaml -t pkg/ --exclude-main -r LICENSE
/tmp/go/bin/swagger generate client -f api/appmgr_rest_api.yaml -t pkg/ -m appmgr_model -c appmgr_client -r LICENSE
  
export GO111MODULE=on 
sudo ldconfig
go build ./cmd/rtmgr.go

export RMR_SEED_RT=/tmp/go/src/routing-manager/uta_rtg_ric.rt

cd $currnewdir/pkg/sbi
go-acc . -- -f "/go/src/routing-manager/manifests/rtmgr/rtmgr-cfg.yaml"

cd $currnewdir/pkg/rpe
go-acc $(go list ./...) -- -f "/go/src/routing-manager/manifests/rtmgr/rtmgr-cfg.yaml"

cd $currnewdir/pkg/sdl
go-acc $(go list ./...) -- -f "/go/src/routing-manager/manifests/rtmgr/rtmgr-cfg.yaml"

cd $currnewdir/pkg/nbi
go-acc $(go list ./...) -- -f "/go/src/routing-manager/manifests/rtmgr/rtmgr-cfg.yaml"

cd $currnewdir

cat $currnewdir/pkg/rpe/coverage.txt | grep -v atomic > coverage_tmp.txt
cat $currnewdir/pkg/sdl/coverage.txt | grep -v atomic >> coverage_tmp.txt
cat $currnewdir/pkg/nbi/coverage.txt | grep -v atomic >> coverage_tmp.txt
cp  $currnewdir/pkg/sbi/coverage.txt  coverage_tmp2.txt
cat coverage_tmp2.txt coverage_tmp.txt > $curdir/coverage.txt

sed -i -e 's/^routing-manager/./' $curdir/coverage.txt

echo "--> rtmgr-build-ubuntu.sh ends"
