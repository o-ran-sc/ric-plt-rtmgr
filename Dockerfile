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

# The CI system creates and publishes the rtmgr Docker image
# from the last step in this multi-stage build and applies 
# a Docker tag from the string in file container-tag.yaml

FROM golang:1.12 as rtmgrbuild
ENV GOPATH /go
RUN apt-get update \
    && apt-get install -y golang-glide git wget

RUN cd /go/bin \
    && wget --quiet https://github.com/go-swagger/go-swagger/releases/download/v0.19.0/swagger_linux_amd64 \
    && mv swagger_linux_amd64 swagger \
    && chmod +x swagger

WORKDIR /go/src/routing-manager
COPY api/ /go/src/routing-manager/api
COPY LICENSE LICENSE
RUN mkdir pkg

RUN git clone "https://gerrit.o-ran-sc.org/r/ric-plt/appmgr" \
    && cp appmgr/api/appmgr_rest_api.yaml api/ \
    && rm -rf appmgr

RUN swagger generate server -f api/routing_manager.yaml -t pkg/ --exclude-main -r LICENSE
RUN swagger generate client -f api/appmgr_rest_api.yaml -t pkg/ -m appmgr_model -c appmgr_client -r LICENSE

COPY glide.lock glide.lock
COPY glide.yaml glide.yaml

RUN glide install --strip-vendor

COPY pkg pkg
COPY cmd cmd
COPY run_rtmgr.sh /run_rtmgr.sh

ENV GOBIN /go/bin
RUN go install ./cmd/rtmgr.go

# UT intermediate container
FROM rtmgrbuild as rtmgrut
RUN go test ./pkg/sbi ./pkg/rpe ./pkg/nbi ./pkg/sdl -cover -race

# Final, executable container
FROM ubuntu:16.04
COPY --from=rtmgrbuild /go/bin/rtmgr /
COPY --from=rtmgrbuild /run_rtmgr.sh /
RUN mkdir /db && touch /db/rt.json && chmod 777 /db/rt.json
RUN chmod 755 /run_rtmgr.sh
CMD /run_rtmgr.sh
