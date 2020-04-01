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


.. contents::
   :depth: 3
   :local:

Pre-requisites
--------------
    * Kubernetes v.1.16.0 or above
    * helm v2.12.3 or above
    * Appmgr pods should be deployed and in running state

Clone the ric-plt/dep git repository
------------------------------------

.. code:: bash

  git clone "https://gerrit.o-ran-sc.org/r/ric-plt/ric-dep"


Modify the deployment recipe
----------------------------
Locate for 'rtmgr' in the recipe file and edit the tag and repo.

.. code:: bash

  rtmgr:
    image:
      registry: "nexus3.o-ran-sc.org:10004/o-ran-sc"
      name: ric-plt-rtmgr
      tag: 0.5.3

Copy the ric-common helm charts for it/dep, configure the helm repo and start local helm server

.. code:: bash

   git clone "https://gerrit.o-ran-sc.org/r/it/dep"
   HELM_HOME=$(helm home)
   COMMON_CHART_VERSION=$(cat dep/ric-common/Common-Template/helm/ric-common/Chart.yaml | grep version | awk '{print $2}')
   helm package -d /tmp dep/ric-common/Common-Template/helm/ric-common
   cp /tmp/ric-common-$COMMON_CHART_VERSION.tgz $HELM_HOME/repository/local/
   helm repo index $HELM_HOME/repository/local/
   helm serve >& /dev/null &
   helm repo remove local
   helm repo add local http://127.0.0.1:8879/charts


At this stage, routing manager can be deployed.

.. code:: bash

  cd ric-dep/bin
  ./install -f ../RECIPE_EXAMPLE/PLATFORM/example_recipe.yaml -c rtmgr

Checking the Deployment Status
------------------------------

Now check the deployment status after a short wait. Results similar to the
output shown below indicate a complete and successful deployment. Check the
STATUS column from both kubectl outputs to ensure that all are either
"Completed" or "Running", and that none are "Error" or "ImagePullBackOff".

.. code:: bash

  #helm list | grep rtmgr
  r3-rtmgr                1               Wed Mar 25 08:34:22 2020        DEPLOYED        rtmgr-3.0.0             1.0          ricplt

  # kubectl get pods -n ricplt | grep rtmgr
  deployment-ricplt-rtmgr-6446b96b65-8mxzn           1/1     Running   0          46s

Checking Container Health
-------------------------

Check the health of the routing manager platform component by querying it
with the following command.

.. code:: bash

 #kubectl get pods -n ricplt -o wide | grep rtmgr
 deployment-ricplt-rtmgr-6446b96b65-8mxzn           1/1     Running   0          16m   10.244.0.17    master-node   <none>           <none>


 curl -v http://10.244.0.17:8080/ric/v1/health/alive
 *   Trying 10.244.0.17...
 * TCP_NODELAY set
 * Connected to 10.244.0.17 (10.244.0.17) port 8080 (#0)
 > GET /ric/v1/health/alive HTTP/1.1
 > Host: 10.244.0.17:8080
 > User-Agent: curl/7.58.0
 > Accept: */*
 >
 < HTTP/1.1 200 OK
 < Content-Type: application/json
 < Date: Wed, 25 Mar 2020 03:19:05 GMT
 < Content-Length: 0
 <  
 * Connection #0 to host 10.244.0.17 left intact

Undeploying Routing Manager
---------------------------

.. code:: bash

 #helm delete --purge r3-rtmgr
