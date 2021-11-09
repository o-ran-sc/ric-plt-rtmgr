# Routing Manager

## Table of contents
* [Introduction](#introduction)
* [Release notes](#release-notes)
* [Prerequisites](#prerequisites)
* [Project folders structure](#project-folders-structure)
* [Installation guide](#installation-guide)
  * [Compiling code](#compiling-code)
  * [Building docker container](#building-docker-container)
  * [Installing Routing Manager](#installing-routing-manager)
  * [Testing and Troubleshoting](#testing-and-troubleshoting)
* [Upcoming changes](#upcoming-changes)
* [License](#license)

## Introduction
__Routing Manager__ is a basic platform service of RIC. It is responsible for distributing routing policies among the other platform components and xApps.

The routing manager has two ways to get the xapp details from xapp manager - httpGetter or httpRESTful.
In case of httpGetter, the implemented logic periodically queries the xApp Manager component for xApps' list.
Where in httpRESTful, starts a http server and creates a webhook subscription in xapp manager to update about changes in xapps and waits changed data to arrive on the REST http server.
Either ways, the xapp data received is stored and then processed to create routing policies and distributes them to all xApps.

The architecture consists of the following five well defined functions:
* NorthBound Interface (__NBI__): Maintains the communication channels towards RIC manager components 
* Routing Policy Engine (__RPE__): Provides the logic to calculate routing policies
* Shared Data Layer (__SDL__): Provides access to different kind persistent data stores
* SouthBound Interface (__SBI__): Maintains the communication channels towards RIC tenants and control components
* Control Logic (__RTMGR__): Controls the operation of above functions

Current implementation provides support for the followings:
* NBI:
  * __httpGet__: simple HTTP GET interface. Expects an URL where it gets the xApps' list in JSON format
  * __httRESTful__: provides REST API endpoints towards RIC manager components. Expects REST port and url where the HTTP service will be started to listen on.
* RPE:
  * __rmr__: creates routing policies formatted for RIC RMR
* SDL:
  * __file__: stores xApp data in container's local filesystem (or in a mountpoint)
  * (backlog) __sdl__: Shared Data Library to Redis database
* SBI:
  * __nngpub__: distributes RPE created policies via NNG Pub channel
  * __nngpipe__: distributes RPE created policies via NNG Pipeline channel

## Release notes
Check the separated `RELNOTES` file.

## Prerequisites
* Healthy kubernetes cluster (for Kubernetes testing)
* Access to the common docker registry (alternatively, you can set up your own private registry for testing: https://docs.docker.com/registry/deploying/)
* In case of non-Docker build: golang 11.1 at least, go-swagger (https://github.com/go-swagger/go-swagger, v0.19.0), glide (https://github.com/Masterminds/glide), XApp Manager spec file (available in ORAN: https://gerrit.o-ran-sc.org/r/admin/repos/ric-plt/appmgr under api folder)

## Project folder structure
* /api: contains Swagger spec files
* /manifest: contains deployment files (Kubernetes manifests, Helm chart)
* /cmd: contains go project's main file
* /pkg: contains go project's internal packages
* /test: contains CI/CD testing files (scripts, mocks, manifests)
* Dockerfile: contains main docker file
* container-tag.yaml: contains CI specific container tag information
* run_rtmgr.sh: shell script to run rtmgr (requires environment variables to be set)

## Installation guide

### Compiling code
#### Docker compile
The Dockerfile located in the project root folder does the following three things:
- As a first step, it creates a build container, fetches XApp Manager's spec file, generates rest api code from swagger spec and builds rtmgr.
- As a second step, it executes UTs on rtmgr source code.
- As a third step, it creates the final container from rtmgr binary (Ubuntu based).
For a docker build execute `docker build --tag=rtmgr-build:test .` in the project root directory (feel free to replace the name:tag with your own)

#### Compiling without docker
Compiling without Docker involves some manual steps before compiling directly with "go build".
The XApp manager's spec file must be fetched, then api generated with swagger. (these steps are included in the Dockerfile).
After the code is generated, glide can install the dependencies of rtmgr.
Make sure you set your GOPATH variable correctly (example: $HOME/go/src/routing-manager)
Code generation and building example (from project root folder):
```bash
git clone "https://gerrit.o-ran-sc.org/r/ric-plt/appmgr" && cp appmgr/api/appmgr_rest_api.yaml api/
swagger generate server -f api/routing_manager.yaml -t pkg/ --exclude-main -r LICENSE
swagger generate client -f api/appmgr_rest_api.yaml -t pkg/ -m appmgr_model -c appmgr_client -r LICENSE
glide install --strip-vendor
go build cmd/rtmgr.go
```

**NOTE:** before doing a docker build it is advised to remove any generated files and vendor packages:
```bash
# assuming that you stand in project root dir
rm -rf appmgr vendor pkg/appmgr_* pkg/models pkg/restapi
```

### Installing Routing Manager
#### Preparing environment
Tag the `rtmgr` container according to the project release and push it to a registry accessible from all minions of the Kubernetes cluster.
Edit the container image section of `rtmgr-dep.yaml` file according to the `rtmgr` image tag.

#### Deploying Routing Manager 
Issue the `kubectl create -f {manifest.yaml}` command in the following order:
  1. `manifests/namespace.yaml`: creates the `example` namespace for routing-manager resources
  2. `manifests/rtmgr/rtmgr-cfg.yaml`: creates default routes config file for routing-manager
  3. `manifests/rtmgr/rtmgr-dep.yaml`: instantiates the `rtmgr` deployment in the `example` namespace
  4. `manifests/rtmgr/rtmgr-svc.yaml`: creates the `rtmgr` service in `example` namespace

**NOTE:** The above manifest files will deploy routing manager with NBI as httpRESTful which would not succeed unless there is an xapp manager running at the defined xm-url. The solution is either to deploy a real XAPP manager before deploying routing-manager or start the mock xmgr as mentioned in [Testing](#testing-and-troubleshoting).

### Testing and Troubleshoting
### Testing with Kubernetes
Routing Manager's behaviour can be tested using the mocked xApp Manager, traffic generator xApp and receiver xApp.

  1. Checkout and compile both xApp receiver and xApp Tx generator of RIC admission control project: `https://gerrit.o-ran-sc.org/r/admin/repos/ric-app/admin`
  2. Copy the `adm-ctrl-xapp` binary to `./test/docker/xapp.build` folder furthermore copy all RMR related dinamycally linked library under `./test/docker/xapp.build/usr` folder. Issue `docker build ./test/docker/xapp.build` command. Tag the recently created docker image and push it to the common registry.
  3. Copy the `test-tx` binary to `./test/docker/xapp-tx.build` folder furthermore copy all RMR related dinamycally linked library under `./test/docker/xapp.build/usr` folder. Issue `docker build ./test/docker/xapp-tx.build` command.  Tag the recently created docker image and push it to the common registry.
  4. Enter the `./test/docker/xmgr.build` folder and issue `docker build .`.  Tag the recently created docker image and push it to the common registry.
  5. Modify the docker image version in each kuberbetes manifest files under `./test/kubernetes/` folder accordingly then issue the `kubectl create -f {manifest.yaml}` on each file.
  6. [Compile](#compiling-code) and [Install routing manager](#installing-routing-manager)
  7. Once the routing manager is started, it retrievs the initial xapp list from `xmgr` via HTTPGet additonaly it starts to listen on http://rtmgr:8888/v1/handles/xapp-handle endpoint and ready to receive xapp list updates.
  8. Edit the provided `test/data/xapp.json` file accordingly and issue the following curl command to update `rtmgr's` xapp list.
     ``` curl --header "Content-Type: application/json" --request POST --data '@./test/data/xapps.json' http://10.244.2.104:8888/v1/handles/xapp-handle ```

### Executing unit tests
For running unit tests, execute the following command:
  `go test ./pkg/nbi` (or any package - feel free to add your own parameters)
If you wish to execute the full UT set with coverage:
```bash
  mkdir -p unit-test
  go test ./pkg/sbi ./pkg/rpe ./pkg/nbi ./pkg/sdl -cover -race -coverprofile=$PWD/unit-test/c.out
  go tool cover -html=$PWD/unit-test/c.out -o $PWD/unit-test/coverage.htm
```

#### Command line arguments
Routing manager binary can be called with `-h` flag when it displays the available command line arguments and it's default value.

Example:

```bash
Usage of ./rtmgr:
  -configfile string
        Routing manager's configuration file path (default "/etc/rtmgrcfg.json")
  -filename string
        Absolute path of file where the route information to be stored (default "/db/rt.json")
  -loglevel string
        INFO | WARN | ERROR | DEBUG (default "INFO")
  -nbi string
        Northbound interface module to be used. Valid values are: 'httpGetter | httpRESTful' (default "httpGetter")
  -nbi-if string
        Base HTTP URL where routing manager will be listening on (default "http://localhost:8888")
  -rpe string
        Route Policy Engine to be used. Valid values are: 'rmrpush | rmrpub' (default "rmrpush")
  -sbi string
        Southbound interface module to be used. Valid values are: 'nngpush | nngpub' (default "nngpush")
  -sbi-if string
        IPv4 address of interface where Southbound socket to be opened (default "0.0.0.0")
  -sdl string
        Datastore enginge to be used. Valid values are: 'file' (default "file")
  -xm-url string
        HTTP URL where xApp Manager exposes the entire xApp List (default "http://localhost:3000/xapps")

```

For troubleshooting purpose the default logging level can be increased to `DEBUG`. (by hand launch it's set to INFO, kubernetes manifest has DEBUG set by default).

## Upcoming changes
[] Add unit tests


## License
This project is licensed under the Apache License, Version 2.0 - see the [LICENSE](LICENSE)

## Building arm64 rtmgr docker image

docker build -f Dockerfile-arm64 -t ric-plt-rtmgr:0.6.3 .
