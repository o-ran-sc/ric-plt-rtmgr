#! /bin/sh -e
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
#	Mnemonic:	run_test-tx.sh
#	Abstract:	Runs the TX transmitter xApp with proper arguments
#	Date:		19 March 2019
#
NAME=${NAME}
PORT=${PORT}
RATE=${RATE}
MESSAGE_TYPE=0
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib

exec /test-tx -n $NAME -p $PORT -r $RATE -m $MESSAGE_TYPE
