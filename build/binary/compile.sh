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
#	Mnemonic:	compile.sh
#	Abstract:	Compiles the rtmgr source
#	Date:		19 March 2019
#
glide install --strip-vendor
p="UT"
if [ "$p" = "$1" ]
then
  echo "Starting Unit Tests..."
  mkdir -p $PWD/unit-test
  go test ./pkg/sbi ./pkg/rpe ./pkg/nbi ./pkg/sdl -cover -race -coverprofile=$PWD/unit-test/c.out
  go tool cover -html=$PWD/unit-test/c.out -o $PWD/unit-test/coverage.html
else
  echo "Compiling..."
  go build -o ./bin/rtmgr cmd/rtmgr.go
fi


