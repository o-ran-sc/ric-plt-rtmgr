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
#
#   This source code is part of the near-RT RIC (RAN Intelligent Controller)
#   platform project (RICP).
#==================================================================================

# The CI system creates and publishes the rtmgr Docker image
# from the last step in this multi-stage build and applies 
# a Docker tag from the string in file container-tag.yaml

#FROM golang:1.12.1 as rtmgrbuild
FROM nexus3.o-ran-sc.org:10004/o-ran-sc/bldr-ubuntu18-c-go:8-u18.04 as rtmgrbuild

# Install RMr shared library
ARG RMRVERSION=4.0.5
RUN wget --content-disposition https://packagecloud.io/o-ran-sc/release/packages/debian/stretch/rmr_${RMRVERSION}_amd64.deb/download.deb && dpkg -i rmr_${RMRVERSION}_amd64.deb && rm -rf rmr_${RMRVERSION}_amd64.deb
# Install RMr development header files
RUN wget --content-disposition https://packagecloud.io/o-ran-sc/release/packages/debian/stretch/rmr-dev_${RMRVERSION}_amd64.deb/download.deb && dpkg -i rmr-dev_${RMRVERSION}_amd64.deb && rm -rf rmr-dev_${RMRVERSION}_amd64.deb

ENV GOLANG_VERSION 1.12.1
RUN wget --quiet https://dl.google.com/go/go$GOLANG_VERSION.linux-amd64.tar.gz \
     && tar xvzf go$GOLANG_VERSION.linux-amd64.tar.gz -C /usr/local 
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH /go

RUN mkdir -p /go/bin
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

RUN /go/bin/swagger generate server -f api/routing_manager.yaml -t pkg/ --exclude-main -r LICENSE
RUN /go/bin/swagger generate client -f api/appmgr_rest_api.yaml -t pkg/ -m appmgr_model -c appmgr_client -r LICENSE

ENV GO111MODULE=on 
ENV GOPATH ""
COPY go.sum go.sum
COPY go.mod go.mod
COPY pkg pkg
COPY cmd cmd
COPY run_rtmgr.sh /run_rtmgr.sh
RUN mkdir manifests
COPY manifests/ /go/src/routing-manager/manifests
COPY "uta_rtg_ric.rt" /
ENV GOPATH /go

ENV GOBIN /go/bin
RUN go install ./cmd/rtmgr.go

# UT intermediate container
#FROM rtmgrbuild as rtmgrut
#RUN ldconfig
#ENV RMR_SEED_RT "/uta_rtg_ric.rt"
#RUN go test ./pkg/sbi ./pkg/rpe ./pkg/nbi ./pkg/sdl -f "/go/src/routing-manager/manifests/rtmgr/rtmgr-cfg.yaml" -cover -race

# Final, executable container
FROM ubuntu:18.04
COPY --from=rtmgrbuild /go/bin/rtmgr /
COPY --from=rtmgrbuild /run_rtmgr.sh /
COPY --from=rtmgrbuild /usr/local/include /usr/local/include
COPY --from=rtmgrbuild /usr/local/lib /usr/local/lib
COPY "uta_rtg_ric.rt" /
RUN ldconfig
RUN apt-get update && apt-get install -y iputils-ping net-tools curl tcpdump
RUN mkdir /db && touch /db/rt.json && chmod 777 /db/rt.json
RUN chmod 755 /run_rtmgr.sh
CMD /run_rtmgr.sh
