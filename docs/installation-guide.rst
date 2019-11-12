..
..  Copyright (c) 2019 AT&T Intellectual Property.
..  Copyright (c) 2019 Nokia.
..
..  Licensed under the Creative Commons Attribution 4.0 International
..  Public License (the "License"); you may not use this file except
..  in compliance with the License. You may obtain a copy of the License at
..
..    https://creativecommons.org/licenses/by/4.0/
..
..  Unless required by applicable law or agreed to in writing, documentation
..  distributed under the License is distributed on an "AS IS" BASIS,
..  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
..
..  See the License for the specific language governing permissions and
..  limitations under the License.
..



Installation Guide
==================

.. contents::
   :depth: 3
   :local:

Abstract
--------

This document describes how to install ric-plt/rtmgr, it's dependencies and required system resources.

Introduction
------------
Routing Manager is a basic platform service of RIC. It is responsible for distributing routing policies among the other platform components and xApps.

Pre-requisites
--------------
* Healthy kubernetes cluster (for Kubernetes testing)
* Access to the common docker registry (alternatively, you can set up your own private registry for testing: "https://docs.docker.com/registry/deploying/")
* In case of non-Docker build: 
    * golang 11.1 at least
    * go-swagger ("https://github.com/go-swagger/go-swagger", v0.19.0)
    * glide ("https://github.com/Masterminds/glide")
    * XApp Manager spec file (available in ORAN: "https://gerrit.o-ran-sc.org/r/admin/repos/ric-plt/appmgr" under api folder)

Software Installation and Deployment
------------------------------------
This section describes the installation of the ric-plt/rtmgr installation.

Docker Image
------------
The Dockerfile located in the project root folder does the following three things:

* As a first step, it creates a build container, fetches XApp Manager's spec file, generates rest api code from swagger spec and builds rtmgr.
* As a second step, it executes UTs on rtmgr source code.
* As a third step, it creates the final container from rtmgr binary (Ubuntu based).

For a docker build execute `docker build --tag=rtmgr-build:test .` in the project root directory (feel free to replace the name:tag with your own)

Linux Binary
------------
* Compiling without Docker involves some manual steps before compiling directly with "go build".
* The XApp manager's spec file must be fetched, then api generated with swagger. (these steps are included in the Dockerfile).
* After the code is generated, glide can install the dependencies of rtmgr.
* Make sure you set your GOPATH variable correctly (example: $HOME/go/src/routing-manager)
* Code generation and building example (from project root folder):

  * git clone https://gerrit.o-ran-sc.org/r/ric-plt/appmgr && cp appmgr/api/appmgr_rest_api.yaml api/
  * swagger generate server -f api/routing_manager.yaml -t pkg/ --exclude-main -r LICENSE
  * swagger generate client -f api/appmgr_rest_api.yaml -t pkg/ -m appmgr_model -c appmgr_client -r LICENSE
  * glide install --strip-vendor
  * go build cmd/rtmgr.go

NOTE: before doing a docker build it is advised to remove any generated files and vendor packages:
assuming that you stand in project root dir

	rm -rf appmgr vendor pkg/appmgr_* pkg/models pkg/restapi

Command line arguments
----------------------
Routing manager binary can be called with `-h` flag when it displays the available command line arguments and it's default value.

Example:

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

