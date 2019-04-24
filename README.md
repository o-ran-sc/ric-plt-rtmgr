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

The implemented logic periodically queries the xApp Manager component for xApps' list. Stores the data then processes it to create routing policies and distributes them to all xApps.
The architecture consists of the following five well defined functions:
* NorthBound Interface (__NBI__): Maintains the communication channels towards RIC manager components 
* Routing Policy Engine (__RPE__): Provides the logic to calculate routing policies
* Shared Data Layer (__SDL__): Provides access to different kind persistent data stores
* SouthBound Interface (__SBI__): Maintains the communication channels towards RIC tenants and control components
* Control Logic (__RTMGR__): Controls the operatin of above functions

Current implementation provides support for the followings:
* NBI:
  * __httpGet__: simple HTTP GET interface. Expects an URL where it gets the xApps' list in JSON format
  * (WIP) __httRESTful__: provides REST API endpoints towards RIC manager components 
* RPE:
  * __rmr__: creates routing policies formatted for RIC RMR
* SDL:
  * __file__: stores xApp data in container's local filesystem (or in a mountpoint)
  * (backlog) __sdl__: Shared Data Library to Redis database
* SBI:
  * __nngpub__: distributes RPE created policies via NNG Pub channel
  * (WIP) __nngpipe__: distributes RPE created policies via NNG Pipeline channel

## Release notes
Check the separated `RELNOTES` file.

## Prerequisites
* Healthy kubernetes cluster
* Access to the common docker registry

## Project folder structure
* /api: contains Swagger source files
* /build: contains build tools (scripts, Dockerfiles, etc.)
* /manifest: contains deployment files (Kubernetes manifests, Helm chart)
* /cmd: contains go project's main file
* /pkg: contains go project's internal packages
* /test: contains CI/CD testing files (scripts, mocks, manifests)

## Installation guide

### Compiling code
Enter the project root and execute `./build.sh` script.
The build script has two main phases. First is the code compilation, where it creates a temporary container for downloading all dependencies then compiles the code. In the second phase it builds the production ready container and taggs it to `rtmgr:builder`

**NOTE:** The script puts a copy of the binary into the `./bin` folder for further use cases

### Installing Routing Manager
#### Preparing environment
Tag the `rtmgr` container according to the project release and push it to a registry accessible from all minions of the Kubernetes cluster.
Edit the container image section of `rtmgr-dep.yaml` file according to the `rtmgr` image tag.

#### Deploying Routing Manager 
Issue the `kubectl create -f {manifest.yaml}` command in the following order
  1. `manifests/namespace.yaml`: creates the `example` namespace for routing-manager resources
  2. `manifests/rtmgr/rtmgr-dep.yaml`: instantiates the `rtmgr` deployment in the `example` namespace
  3. `manifests/rtmgr/rtmgr-svc.yaml`: creates the `rtmgr` service in `example` namespace

### Testing and Troubleshoting
Routing Manager's behaviour can be tested using the mocked xApp Manager, traffic generator xApp and receiver xApp.

  1. Checkout and compile both xApp receiver and xApp Tx generator of RIC admission control project
  2. Copy the `adm-ctrl-xapp` binary to `./test/docker/xapp.build` folder furthermore copy all RMR related dinamycally linked library under `./test/docker/xapp.build/usr` folder. Issue `docker build ./test/docker/xapp.build` command. Tag the recently created docker image and push it to the common registry.
  3. Copy the `test-tx` binary to `./test/docker/xapp-tx.build` folder furthermore copy all RMR related dinamycally linked library under `./test/docker/xapp.build/usr` folder. Issue `docker build ./test/docker/xapp-tx.build` command.  Tag the recently created docker image and push it to the common registry.
  4. Enter the `./test/docker/xmgr.build` folder and issue `docker build .`.  Tag the recently created docker image and push it to the common registry.
  5. Modify the docker image version in each kuberbetes manifest files under `./test/kubernetes/` folder accordingly then issue the `kubectl create -f {manifest.yaml}` on each file.
  6. [Compile](#compiling-code) and [Install routing manager](#installing-routing-manager)

#### Command line arguments
Routing manager binary can be called with `-h` flag when it displays the available command line arguments and it's default value.

Example:

```bash
Usage of ./rtmgr:
  -filename string
        Absolute path of file where the route information to be stored (default "/db/rt.json")
  -loglevel string
        INFO | WARN | ERROR | DEBUG (default "INFO")
  -nbi string
        Northbound interface module to be used. Valid values are: 'httpGetter' (default "httpGetter")
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

For troubleshooting purpose the default logging level can be increased to `DEBUG`.

## Upcoming changes
[] Add RESTful NBI based on swagger api definition

[] Add unit tests

## License
This project is licensed under the Apache License, Version 2.0 - see the [LICENSE](LICENSE)

