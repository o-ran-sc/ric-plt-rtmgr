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

Developer-Guide
===============

.. contents::
   :depth: 3
   :local:

Overview
--------
Routing Manager is a basic platform service of RIC. It is responsible for distributing routing policies among the other platform components and xApps.

The routing manager has two ways to get the xapp details from xapp manager - httpGetter or httpRESTful.
In case of httpGetter, the implemented logic periodically queries the xApp Manager component for xApps' list.
Where in httpRESTful, starts a http server and creates a webhook subscription in xapp manager to update about changes in xapps and waits changed data to arrive on the REST http server.
Either ways, the xapp data received is stored and then processed to create routing policies and distributes them to all xApps.

Architecture
------------
The architecture consists of the following five well defined functions:

* NorthBound Interface (__NBI__): Maintains the communication channels towards RIC manager components
* Routing Policy Engine (__RPE__): Provides the logic to calculate routing policies
* Shared Data Layer (__SDL__): Provides access to different kind persistent data stores
* SouthBound Interface (__SBI__): Maintains the communication channels towards RIC tenants and control components
* Control Logic (__RTMGR__): Controls the operation of above functions


Installing Routing Manager
--------------------------
* Tag the `rtmgr` container according to the project release and push it to a registry accessible from all minions of the Kubernetes cluster.
* Edit the container image section of `rtmgr-dep.yaml` file according to the `rtmgr` image tag.

Deploying Routing Manager
-------------------------
* Issue the `kubectl create -f {manifest.yaml}` command in the following order:
   1. `manifests/namespace.yaml`: creates the `example` namespace for routing-manager resources
   2. `manifests/rtmgr/rtmgr-cfg.yaml`: creates default routes config file for routing-manager
   3. `manifests/rtmgr/rtmgr-dep.yaml`: instantiates the `rtmgr` deployment in the `example` namespace
   4. `manifests/rtmgr/rtmgr-svc.yaml`: creates the `rtmgr` service in `example` namespace

NOTE: The above manifest files will deploy routing manager with NBI as httpRESTful which would not succeed unless there is an xapp manager running at the defined xm-url. The solution is either to deploy a real XAPP manager before deploying routing-manager or start the mock xmgr as mentioned in [Testing](#testing-and-troubleshoting).

Testing and Troubleshooting
---------------------------
Testing with Kubernetes
-----------------------
Routing Manager's behaviour can be tested using the mocked xApp Manager, traffic generator xApp and receiver xApp.

* Checkout and compile both xApp receiver and xApp Tx generator of RIC admission control project:
  `https://gerrit.o-ran-sc.org/r/admin/repos/ric-app/admin`

* Copy the `adm-ctrl-xapp` binary to `./test/docker/xapp.build` folder furthermore copy all RMR related dinamycally linked library under `./test/docker/xapp.build/usr` folder. Issue `docker build ./test/docker/xapp.build` command. Tag the recently created docker image and push it to the common registry.

* Copy the `test-tx` binary to `./test/docker/xapp-tx.build` folder furthermore copy all RMR related dinamycally linked library under `./test/docker/xapp.build/usr` folder. Issue `docker build ./test/docker/xapp-tx.build` command.  Tag the recently created docker image and push it to the common registry.

* Enter the `./test/docker/xmgr.build` folder and issue `docker build .`.  Tag the recently created docker image and push it to the common registry.

* Modify the docker image version in each kuberbetes manifest files under `./test/kubernetes/` folder accordingly then issue the `kubectl create -f {manifest.yaml}` on each file.

Once the routing manager is started, it retrievs the initial xapp list from `xmgr` via HTTPGet additonaly it starts to listen on http://rtmgr:8888/v1/handles/xapp-handle endpoint and ready to receive xapp list updates.

* Edit the provided `test/data/xapp.json` file accordingly and issue the following curl command to update `rtmgr's` xapp list.

  `curl --header "Content-Type: application/json" --request POST --data '@./test/data/xapps.json' http://10.244.2.104:8888/v1/handles/xapp-handle`

Executing unit tests
--------------------
For running unit tests, execute the following command:
   `go test ./pkg/nbi` (or any package - feel free to add your own parameters)

If you wish to execute the full UT set with coverage:

   mkdir -p unit-test

   go test ./pkg/sbi ./pkg/rpe ./pkg/nbi ./pkg/sdl -cover -race -coverprofile=$PWD/unit-test/c.out

   go tool cover -html=$PWD/unit-test/c.out -o $PWD/unit-test/coverage.html


For troubleshooting purpose the default logging level can be increased to `DEBUG`. (by hand launch it's set to INFO, kubernetes manifest has DEBUG set by default).
