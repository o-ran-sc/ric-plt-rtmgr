#!/bin/sh -e
#
#==================================================================================
#   Copyright (c) 2019 AT&T Intellectual Property.
#   Copyright (c) 2019 Nokia
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
#==================================================================================
#
#
#	Mnemonic:	build.sh
#	Abstract:	Compiles the rtmgr source and builds the docker container
#	Date:		12 March 2019
#

echo 'Creating compiler container'
docker build --no-cache --tag=rtmgr_compiler:0.1 build/binary/

echo 'Running rtmgr compiler'
docker run --rm --name=rtmgr_compiler -v ${PWD}:/opt/ rtmgr_compiler:0.1

echo 'Cleaning up compiler container'
docker rmi -f rtmgr_compiler:0.1

echo 'rtmgr binary successfully built!'

echo 'Creating rtmgr container'
cp ${PWD}/bin/* ${PWD}/build/container/
docker build --no-cache --tag=rtmgr:builder build/container/
