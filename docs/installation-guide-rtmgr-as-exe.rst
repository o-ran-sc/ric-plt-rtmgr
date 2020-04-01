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

Pre-requisites
--------------
    * Ubuntu machine
    * golang 1.12.1 minimum

Compiling Routing Manager
-------------------------
Clone the ric-plt/dep git repository.

.. code:: bash

  git clone "https://gerrit.o-ran-sc.org/r/ric-plt/rtmgr"

Execute this shell script which will give you the rtmgr as executable

.. code:: bash

  cd rtmgr/example
  ./rtmgr_exe.sh

Run rtmgr by passing the config file as parameter. Note that the rtmgr may abort after sometime as it needs appmgr to be running. This can be tweaked by modifying the actual code. As this would be needed only for actual Development, the sam eis not being mentioned here.

.. code:: bash

  ./rtmgr -f rtmgr-config.yaml

